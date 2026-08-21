# common-template-golang

**语言:** [English](README.md) | 中文

一个基于 `common-library-golang` 构建的通用 Go 语言开发模板。

## 1. 代码结构

- `adapter/`
  对外交互的适配层代码，例如 RESTful API 定义、命令行定义等。

- `config/`
  配置文件（`app.toml`）以及 SQLite 数据库结构定义（`schema_sqlite.sql`）。

- `data`
  SQLite 数据文件。

- `docs/`
  由 swagger 生成的文档目录。

- `domain/`
  领域模型的定义与实现。

- `infra/`
  实现外部依赖的基础设施代码，例如数据库访问对象（DAO）等。

- `pkg/`
  功能性代码包，例如工具函数、依赖注入函数、DTO 定义等。

- `service/`
  服务层代码，实现主要的业务流程。

- `CLAUDE.md`
  给 Claude Code（以及其他 AI 编程助手）看的指引文档：项目概览、构建/测试命令，以及一些改动代码前值得了解的非显而易见的实现细节。

- `Dockerfile`
  两阶段构建：第一阶段编译应用，第二阶段在精简的基础镜像上运行。

- `go.mod`、`go.sum`
  Go 模块定义与依赖锁定文件。

- `main.go`
  主程序，代码入口。

- `README.md`
  项目文档。

本模板遵循标准的 Go 项目布局，便于以最佳实践快速启动新的 Go 模块。

## 2. `main.go` 如何组装各个组件

`main.go` 本身并不做任何具体的初始化工作——它只是声明了一个有序的组件（**component**）列表，交给 `bootstrap.New()...Run()` 去执行，这是 `common-library-golang` 上所有服务共用的启动/关闭编排器。每个组件都有一个 `Init()`（启动时按注册顺序调用）和一个 `Close()`（关闭时按**逆序**调用）。理解这条链路，是理解一个基于本模板生成的新服务启动时到底做了什么的最快方式。

```go
bootstrap.New().
    Add(envComp.Component("config/app.toml")). // 第1个 - env
    Add(logComp.Component()).                  // 第2个 - log
    AddParallel(dbComp.Component()).
    PreReady(initServices).
    Add(ginComp.Component(func(r *gin.Engine) {
        router = r
        adapter.Mount(r)
    })).
    Add(httpComp.Component(func() http.Handler { return router })).
    PostReady(func() { log.Infof("%s start successfully, ...", ...) }).
    Run()
```

逐步说明：

1. **`envComp.Component("config/app.toml")`** —— 必须永远是第一个 `Add()` 调用。把 `config/app.toml` 加载进 koanf 配置单例（`env.Env()`），后续所有组件都会在各自的 `Init()` 里读取这个配置。`Close()` 是空操作。

2. **`logComp.Component()`** —— 必须永远是第二个 `Add()` 调用。读取 `[log]` 配置段（`log.level`，以及可选的 `log.file.*` 用于基于文件的日志轮转），初始化进程级别的日志器。因为它是最后被关闭的（LIFO 顺序），所以关闭过程中的日志消息不会丢失。

3. **`AddParallel(dbComp.Component())`** —— 打开数据库连接。本模板使用的是 `dbsqlx/sqlite/component`，它从 `config/app.toml` 读取 `db.sqlite.path`，并把连接注册为 `dbsqlx.Default()`。`AddParallel` 允许多个相互独立的组件（比如第二个数据存储、一个缓存客户端）在同一阶段并发初始化——这里只有一个，但实际项目里可以在同一阶段加入更多。把这一行换成 `dbsqlx/postgres/component`，就能把模板切换成使用 Postgres。

4. **`PreReady(initServices)`** —— 一个自定义的钩子函数（定义在这个文件里，不是库里的组件），在 env/log/db 都就绪之后、HTTP 服务器开始接受流量之前运行一次。这里就是做依赖注入的地方：在 `dbsqlx.Default()` 之上构造 DAO，包装成 service，再通过 `adapter.Svcs` 发布给适配层。往模板里加一个新的实体时，就按照这里 `service.NewUserService(dao.NewUserDao(dbsqlx.Default()))` 的方式扩展这个函数即可。`PreReady` 这一步没有 `Close()`——它只是纯粹的初始化逻辑，不是一个组件。

5. **`ginComp.Component(func(r *gin.Engine) {...})`** —— 创建 `*gin.Engine`（CORS 根据 `app.env.value` 的值，从 `cors.allow.origins.dev`/`.prod` 中选取配置），随后立即调用传入的闭包。这个闭包把 router 赋值给 `main()` 顶部声明的 `var router *gin.Engine`，并调用 `adapter.Mount(r)` 注册所有 REST 路由。`Close()` 是空操作——Gin 引擎本身不持有任何资源，真正的监听器由下一个组件负责。

6. **`httpComp.Component(func() http.Handler { return router })`** —— 在 `http.server.port` 上启动真正的 HTTP 监听，处理该闭包返回的 `http.Handler`。这个闭包是在 `Init()` 时才被惰性求值的（而不是在 `Add()` 调用时），这也是为什么即使 `router` 只在第 5 步的闭包里才被赋值也依然安全——等到这个组件的 `Init()` 运行时，第 5 步早已完成。`Close()` 会在 `httpserver.DefaultShutdownTimeout` 时间内优雅地关闭服务器。

7. **`PostReady(func() {...})`** —— 在上面所有组件都初始化成功之后、进程阻塞等待关闭信号之前，运行一次。本模板在这里打印应用名称/版本/环境/本机 IP 的日志；可以把任何"只需要在服务完全就绪后执行一次"的逻辑放在这里（比如注册到服务发现系统）。

8. **`Run()`** —— 按顺序执行以上所有步骤，然后调用 `shutdown.Wait()` 阻塞，直到收到一个操作系统信号（或者内部触发的 `shutdown.Trigger()`，比如 HTTP 服务器意外挂掉的情况）。关闭时，每个组件的 `Close()` 按**逆序（LIFO）**执行：HTTP 服务器 → Gin（空操作） → SQLite → log → env（空操作）——这样可以让正在处理中的请求有机会完成、数据库连接干净地关闭，并且最后写入的一定是关闭日志。

如果 `Init()` 在任何一步失败（包括 `PreReady` 函数返回非 nil 的 error），`bootstrap` 会记录这次失败，把已经启动的组件按 LIFO 顺序全部回滚，然后以退出码 1 结束进程——这样就不会出现一个新服务"启动到一半"的中间状态。

## 3. 如何运行？

**1. 更新依赖**
- go mod tidy

**2. 生成 swagger 文件**
- swag init

访问 http://localhost:8001/swagger/index.html

**3. 测试**
- go test ./... -cover

**4. 构建**
- go build

**5. 运行**
- ./template
