package controllers

/*
 ============================================================================
 Name        : config.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 配置
 ============================================================================
*/

import (
	"github.com/astaxie/beego/logs"
	"github.com/xipfs/ipfsadmin/app/service"
)

type ConfigController struct {
	BaseController
}

// 获取状态
func (this *ConfigController) Get() {
	req := this.Ctx.Request
	addr := req.RemoteAddr
	peer_id := this.GetString("peer_id")
	timestr := this.GetString("timestr")
	config, _ := service.ConfigService.GetConfig(1)
	logs.Info("{ip:%s,pid:%s,config:%s,timestr:%s}", addr, peer_id, config.Value, timestr)
	this.Ctx.WriteString(config.Value)
}
