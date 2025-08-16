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

package adapter

import (
	"net/http"
	"strconv"
	"strings"

	"template/pkg/dto"
	"template/pkg/injector"

	"gitlab.com/phcp-common/library-golang/code"
	libDto "gitlab.com/phcp-common/library-golang/dto"
	"gitlab.com/phcp-common/library-golang/log"
	"gitlab.com/phcp-common/library-golang/util"

	"github.com/gin-gonic/gin"
)

// Mount all RESTful APIs
func Mount(router *gin.Engine) {
	MountUser(router)
	MountSwagger(router)
	log.Info("mount all RESTful APIs successful.")
}

// MountUser
func MountUser(router *gin.Engine) *gin.Engine {
	// Need JWT authorization.
	r1 := router.Group("/usrapi/v1/users")
	r1.GET("/list", getUserList)

	return router
}

// getUserList godoc
// @Summary Query users list
// @Schemes
// @Description
// @Tags User
// @Accept json
// @Produce json
// @Param  kind query string false "user kind, default is all" Enums(Reader, Author, Admin)
// @Success 200 {object} dto.UserListResp
// @Failure 500 {object} dto.ResponseMessage
// @Router /users/list [get]
func getUserList(c *gin.Context) {
	var listPara dto.UserListPara
	//分页等可选参数
	listPara.Page, _ = strconv.Atoi(c.Query("page"))
	listPara.Limit, _ = strconv.Atoi(c.Query("limit"))
	listPara.Sort = c.Query("sort")
	listPara.Direction = c.Query("direction")
	listPara.Kind = strings.TrimSpace(c.Query("kind"))

	// validate list parameters
	if err := util.Validator().Struct(&listPara); err != nil {
		c.JSON(http.StatusBadRequest, libDto.ResponseMessage{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	// get user list from service
	if user, err := injector.UserServiceImpl.GetList(&listPara); err == nil {
		c.JSON(http.StatusOK, libDto.ResponseMessage{
			Code: code.API_CODE_SUCCESS,
			Data: user})
	} else {
		c.JSON(http.StatusInternalServerError, libDto.ResponseMessage{Code: http.StatusInternalServerError, Message: err.Error()})
	}
}
