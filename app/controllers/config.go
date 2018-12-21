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
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/xipfs/ipfsadmin/app/entity"
	"github.com/xipfs/ipfsadmin/app/service"
)

type ConfigController struct {
	BaseController
}

type ConfigParam struct {
	PeerId        string `json:"peer_id"`
	TimeStr       string `json:"timestr"`
	ExtParams     string `json:"extParams"`
	DynamicParams string `json:"dynamicParams"`
}

// 获取状态
func (this *ConfigController) Get() {
	req := this.Ctx.Request
	addr := req.RemoteAddr
	p := &ConfigParam{}
	requestBody := this.GetRequestBody()
	json.Unmarshal(requestBody, p)
	config := &entity.Config{}

	if strings.Contains(p.ExtParams, "android") {
		if strings.Contains(p.ExtParams, "WIFI") {
			config, _ = service.ConfigService.GetConfig(1)
		} else {
			config, _ = service.ConfigService.GetConfig(2)
		}
	} else {
		config, _ = service.ConfigService.GetConfig(1)
	}

	if strings.Contains(p.DynamicParams, "WIFI") {
		config, _ = service.ConfigService.GetConfig(1)
	} else if strings.Contains(p.DynamicParams, "MOBILE") {
		config, _ = service.ConfigService.GetConfig(2)
	} else {

	}

	logs.Info("config_get:{ip:%s,pid:%s,config:%s,timestr:%s}", addr, p.PeerId, config.Value, p.TimeStr)
	this.Ctx.WriteString(config.Value)
}
