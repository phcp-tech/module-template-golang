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

package service

import (
	"template/infra/dao"
	"template/pkg/dto"

	libDto "gitlab.com/phcp-common/library-golang/dto"
)

type UserService struct {
	dao dao.IUserDao
}

func NewUserService(dao dao.IUserDao) *UserService {
	return &UserService{
		dao: dao,
	}
}

// Get user list in Service layer
func (s *UserService) GetList(listPara *dto.UserListPara) (libDto.DataListResp, error) {
	return s.dao.GetList(listPara)
}
