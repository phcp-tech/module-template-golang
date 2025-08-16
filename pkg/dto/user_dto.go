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

package dto

import (
	"template/domain/model"

	libDto "gitlab.com/phcp-common/library-golang/dto"
)

// Define User parameters for query list.
type UserListPara struct {
	libDto.PageParameter
	model.User
}

// those dtos only for swagger api doc, maybe repeated in other files
// common dto
type ResponseMessage struct {
	Code    int    `json:"code"`              // 0: success; others: failed
	Message string `json:"message,omitempty"` // omit in success response
	Data    any    `json:"data,omitempty"`    // omit in failed response
}

// those dtos only for swagger api doc, maybe repeated in other files
// special dto
type UserListResp struct {
	Code int `json:"code"` // 0: success; others: failed
	Data struct {
		Total int          `json:"total"`
		List  []model.User `json:"list"`
	} `json:"data"`
}
