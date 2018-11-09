package entity

import (
	"time"
)

/*
 ============================================================================
 Name        : task.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 任务
 ============================================================================
*/
type Task struct {
	Id             int
	ResourceId     int       `orm:"index"`                       // 资源id
	Message        string    `orm:"type(text)"`                  // 版本说明
	UserId         int       `orm:"index"`                       // 创建人ID
	UserName       string    `orm:"size(20)"`                    // 创建人名称
	BuildStatus    int       `orm:"default(0)"`                  // 构建状态
	ChangeLogs     string    `orm:"type(text)"`                  // 修改日志列表
	ChangeFiles    string    `orm:"type(text)"`                  // 修改文件列表
	FileName       string    `orm:"size(200)"`                   // 文件名
	PubEnvId       int       `orm:"default(0)"`                  // 发布环境ID
	PubStatus      int       `orm:"default(0)"`                  // 发布状态：1 正在发布，2 添加到本地节点，3 目标服务器执行命令，-2 添加本地节点失败，-3 目标服务器执行失败
	PubTime        time.Time `orm:"null;type(datetime)"`         // 发布时间
	ErrorMsg       string    `orm:"type(text)"`                  // 错误消息
	PubLog         string    `orm:"type(text)"`                  // 发布日志
	CreateTime     time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	UpdateTime     time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	ResourceInfo   *Resource `orm:"-"`                           // 资源信息
	EnvInfo        *Env      `orm:"-"`                           // 发布环境
	UploadFileName string    `orm:"size(32)"`                    //上传文件标识
}
