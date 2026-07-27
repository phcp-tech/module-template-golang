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

package dao

import (
	"context"

	"template/domain/model"
	"template/pkg/dto"

	"github.com/phcp-tech/common-library-golang/dbsqlx"
	libDto "github.com/phcp-tech/common-library-golang/dto"
	"github.com/vinovest/sqlx"
)

// compile-time interface check
var _ IUserDao = (*UserDao)(nil)

type UserDao struct {
	db *sqlx.DB
}

func NewUserDao(db *sqlx.DB) IUserDao {
	return &UserDao{db: db}
}

// Get User list in data access layer
func (d *UserDao) GetList(listPara *dto.UserListPara) (libDto.DataListResp, error) {
	var liststr string = `SELECT id, username, nickname, email, kind, status FROM users WHERE 1 = 1 `
	var totalstr string = `SELECT COUNT(*) FROM users WHERE 1 = 1 `
	var sqlstr string
	var args []any

	pagestr := dbsqlx.SortSql(&listPara.PageParameter) + dbsqlx.PageSql(&listPara.PageParameter)

	ctx := context.Background()
	var users []model.User
	var list libDto.DataListResp
	list.List = users

	countQuery := d.db.Rebind(totalstr + sqlstr)
	if err := d.db.GetContext(ctx, &list.Total, countQuery, args...); err != nil {
		return list, err
	}
	if list.Total == 0 {
		return list, nil
	}

	listQuery := d.db.Rebind(liststr + sqlstr + pagestr)
	if err := d.db.SelectContext(ctx, &users, listQuery, args...); err != nil {
		return list, err
	}
	list.List = users
	return list, nil
}
