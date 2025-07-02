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

package dao

import (
	"template/domain/model"
	"template/pkg/dto"

	libDto "gitlab.com/phcp-common/library-golang/dto"
)

type UserDao struct{}

func NewUserDao() IUserDao {
	return &UserDao{}
}

// Get User list in data access layer
func (d *UserDao) GetList(listPara *dto.UserListPara) (libDto.DataListResp, error) {
	var users []model.User
	user1 := model.User{
		Id:       1,
		Username: "Tom",
		Nickname: "Tommy",
		Email:    "tom@gmail.com",
		Kind:     "Reader",
		Status:   1}
	user2 := model.User{
		Id:       2,
		Username: "Jerry",
		Nickname: "Jerry",
		Email:    "jerry@hotmail.com",
		Kind:     "Author",
		Status:   1}
	users = append(users, user1, user2)

	var list libDto.DataListResp
	list.List = users
	list.Total = 2

	return list, nil
}
