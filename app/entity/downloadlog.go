package entity

import (
	"time"
)

/*
 ============================================================================
 Name        : download_log.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 下载日志
 ============================================================================
*/

type DownloadLog struct {
	Id           int
	PeerId       string    //节点ID
	Name         string    //资源名称
	DownloadTime int       //下载耗时
	Size         int       //文件大小
	CreateTime   time.Time `orm:"auto_now_add;type(datetime)"` // 下载时间
}
