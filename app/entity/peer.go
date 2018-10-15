package entity

import (
	"time"
)

/*
 ============================================================================
 Name        : peer.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 节点
 ============================================================================
*/

type Peer struct {
	Id         int
	PeerId     string    `orm:"size(32)"` // 节点ID
	RepoSize   int64     // 使用大小
	StorageMax int64     // 存储空间
	NumObjects int       // 对象数量
	Status     int       // 状态
	UpdateTime time.Time `orm:"auto_now;type(datetime)"` // 更新时间
	CreateTime time.Time `orm:"auto_now_add;type(date)"` // 创建时间
}
