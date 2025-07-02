// Copyright(C) 2020-2025 PHCP Technologies. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"

	"template/adapter"
	"template/pkg/injector"

	"gitlab.com/phcp-common/library-golang/app"
	"gitlab.com/phcp-common/library-golang/env"
	libGin "gitlab.com/phcp-common/library-golang/gin"
	"gitlab.com/phcp-common/library-golang/httpserver"
	"gitlab.com/phcp-common/library-golang/log"

	"github.com/gin-gonic/gin"
)

func main() {
	// step 1: initial config file
	if err := env.InitEnv("config/app.toml"); err != nil {
		log.Errorf("initial environment config file failed: %s", err.Error())
		os.Exit(1)
	}
	log.Info("initial environment config file successful.")

	// step 2: inject all implements
	injector.InjectServices()

	// step 3: initial gin and mount controller/pprof
	router := libGin.InitGin()
	adapter.Mount(router)

	// step 4: start http server
	go func(r *gin.Engine) {
		if err := httpserver.Startup(r, env.Env().String("http.server.port")); err != nil {
			log.Errorf("startup http server failed: %s", err.Error())
			os.Exit(1)
		}
	}(router)

	// step 5: Log for application start successful
	log.Infof("%s start successful, version is %s, environment is %s.",
		env.Env().String("app.name"),
		env.Env().String("app.version"),
		env.Env().String("app.env.value"))

	// step 6: waiting for graceful exit
	app.WatingForExitSignal()
}
