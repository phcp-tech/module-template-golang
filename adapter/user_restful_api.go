// Copyright(C) 2020-2026 PHCP Technologies. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package adapter

import (
	"net/http"
	"strconv"

	"template/pkg/dto"
	"template/pkg/metrics"

	"github.com/phcp-tech/common-library-golang/dbsqlx"
	libDto "github.com/phcp-tech/common-library-golang/dto"
	"github.com/phcp-tech/common-library-golang/errors"
	"github.com/phcp-tech/common-library-golang/health"
	"github.com/phcp-tech/common-library-golang/version"

	"github.com/gin-gonic/gin"
)

// Mount all RESTful APIs
func Mount(router *gin.Engine) {
	MountUser(router)
	MountSwagger(router)
}

// MountUser
func MountUser(router *gin.Engine) *gin.Engine {
	// Need JWT authorization.
	r1 := router.Group("/usrapi/v1/users")
	r1.GET("/list", getUserList)

	// Don not need JWT authorization.
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, libDto.ResponseMessage{
			Code: errors.API_CODE_SUCCESS,
			Data: "OK"})
	})
	router.GET("/usrapi/v1/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, libDto.ResponseMessage{
			Code: errors.API_CODE_SUCCESS,
			Data: version.Get()})
	})
	router.GET("/usrapi/v1/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, libDto.ResponseMessage{
			Code: errors.API_CODE_SUCCESS,
			Data: health.Check(c.Request.Context(), dbsqlx.HealthChecker())})
	})
	router.GET("/usrapi/v1/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, libDto.ResponseMessage{
			Code: errors.API_CODE_SUCCESS,
			Data: metrics.GetMetrics()})
	})

	return router
}

// getUserList godoc
// @Summary Query users list
// @Schemes
// @Description
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserListResp
// @Failure 500 {object} dto.ResponseMessage
// @Router /users/list [get]
func getUserList(c *gin.Context) {
	var listPara dto.UserListPara
	// paginate parameters
	listPara.Page, _ = strconv.Atoi(c.Query("page"))
	listPara.Limit, _ = strconv.Atoi(c.Query("limit"))
	listPara.Sort = c.Query("sort")
	listPara.Direction = c.Query("direction")

	// get user list from service
	if user, err := Svcs.UserService.GetList(&listPara); err == nil {
		c.JSON(http.StatusOK, libDto.ResponseMessage{
			Code: errors.API_CODE_SUCCESS,
			Data: user})
	} else {
		c.JSON(http.StatusInternalServerError, libDto.ResponseMessage{Code: http.StatusInternalServerError, Message: err.Error()})
	}
}
