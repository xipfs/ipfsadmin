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
	"encoding/json"

	"github.com/astaxie/beego/logs"
	"github.com/xipfs/ipfsadmin/app/service"
)

type ConfigController struct {
	BaseController
}

type ConfigParam struct {
	PeerId  string `json:"peer_id"`
	TimeStr string `json:"timestr"`
}

// 获取状态
func (this *ConfigController) Get() {
	req := this.Ctx.Request
	addr := req.RemoteAddr
	p := &ConfigParam{}
	requestBody := this.GetRequestBody()
	json.Unmarshal(requestBody, p)
	config, _ := service.ConfigService.GetConfig(1)
	logs.Info("config_get:{ip:%s,pid:%s,config:%s,timestr:%s}", addr, p.PeerId, config.Value, p.TimeStr)
	this.Ctx.WriteString(config.Value)
}
