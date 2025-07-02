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

//go:build !noswagger

package adapter

import (
	"template/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

// mount swagger adapter
func MountSwagger(router *gin.Engine) *gin.Engine {
	docs.SwaggerInfo.Title = "User RESTful API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/usrapi/v1"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	router.GET("/usrapi/v1/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	return router
}
