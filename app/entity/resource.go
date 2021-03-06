package entity

import (
	"time"
)

/*
 ============================================================================
 Name        : resource.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 资源
 ============================================================================
*/

type Resource struct {
	Id             int
	Name           string    `orm:"size(100)"`                   // 资源名称
	Domain         string    `orm:"size(100)"`                   // 资源标识
	Version        string    `orm:"size(20)"`                    // 最后发布版本
	VersionTime    time.Time `orm:"type(datetime)"`              // 最后发版时间
	Hash           string    `orm:"size(64)"`                    // Hash
	MD5            string    `orm:"size(32);column(md5)"`        // md5
	RepoUrl        string    `orm:"size(100)"`                   // 地址
	Status         int       `orm:"default(0)"`                  // 初始化状态
	ErrorMsg       string    `orm:"type(text)"`                  // 错误消息
	AgentId        int       `orm:"default(0)"`                  // 跳板机ID
	IgnoreList     string    `orm:"type(text)"`                  // 忽略文件列表
	BeforeShell    string    `orm:"type(text)"`                  // 发布前要执行的shell脚本
	AfterShell     string    `orm:"type(text)"`                  // 发布后要执行的shell脚本
	CreateVerfile  int       `orm:"default(0)"`                  // 是否生成版本号文件
	VerfilePath    string    `orm:"size(50)"`                    // 版本号文件目录
	TaskReview     int       `orm:"default(0)"`                  // 发布单是否需要经过审批
	CreateTime     time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	UpdateTime     time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	UploadFileName string    `orm:"size(32)"`                    // 上传标志
}
