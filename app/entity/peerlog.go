package entity

import (
	"time"
)

/*
 ============================================================================
 Name        : peer_log.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 节点日志
 ============================================================================
*/

type PeerLog struct {
	Id            int
	PeerCount     int       // 节点数量
	DownloadCount int       // 下载量
	CreateTime    time.Time `orm:"type(datetime)"`
	Goarch        string    `orm:"size(16)"`
	Goos          string    `orm:"size(16)"`
	Mac           string    `orm:"size(16)"`
	PeerId        string    `orm:"size(64)"`
	EventAction   string    `orm:"size(16)"`
}
