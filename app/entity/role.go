package entity

import (
	"time"
)

/*
 ============================================================================
 Name        : role.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 角色
 ============================================================================
*/

type Role struct {
	Id          int
	RoleName    string    `orm:"size(20)"`                    // 角色名称
	ResourceIds string    `orm:"size(1000)"`                  // 项目权限
	Description string    `orm:"size(200)"`                   // 说明
	CreateTime  time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	UpdateTime  time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	PermList    []Perm    `orm:"-"`                           // 权限列表
	UserList    []User    `orm:"-"`                           // 用户列表
}

// 角色权限
type RolePerm struct {
	RoleId int    // 角色id
	Perm   string `orm:"size(50)"` // 权限
}
