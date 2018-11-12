package service

/*
 ============================================================================
 Name        : action.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 动作
 ============================================================================
*/

import (
	"fmt"

	"github.com/xipfs/ipfsadmin/app/entity"
)

// 系统动态
type actionService struct{}

func (this *actionService) table() string {
	return tableName("action")
}

// 添加记录
func (this *actionService) Add(action, actor, objectType string, objectId int, extra string) bool {
	act := new(entity.Action)
	act.Action = action
	act.Actor = actor
	act.ObjectType = objectType
	act.ObjectId = objectId
	act.Extra = extra
	o.Insert(act)
	return true
}

// 登录动态
func (this *actionService) Login(userName string, userId int, ip string) {
	this.Add("login", userName, "user", userId, userName+" 从 "+ip+" 登录 ！")
}

// 退出登录
func (this *actionService) Logout(userName string, userId int, ip string) {
	this.Add("logout", userName, "user", userId, userName+" 从 "+ip+" 退出 ！")
}

// 更新个人信息
func (this *actionService) UpdateProfile(userName string, userId int) {
	this.Add("update_profile", userName, "user", userId, userName+"更新个人信息！")
}

// 获取动态列表
func (this *actionService) GetList(page, pageSize int) ([]entity.Action, int64) {
	var (
		list  []entity.Action
		count int64
	)
	query := o.QueryTable(this.table())
	count, _ = query.Count()
	num, err := query.OrderBy("-create_time").Offset((page - 1) * pageSize).Limit(pageSize).All(&list)
	if num > 0 && err == nil {
		for i := 0; i < int(num); i++ {
			this.format(&list[i])
		}
	}
	return list, count

}

// 格式化
func (this *actionService) format(action *entity.Action) {
	action.Message = fmt.Sprintf("<b>%s</b>", action.Extra)
}
