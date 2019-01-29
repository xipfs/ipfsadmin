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
	"fmt"
	"strings"
	"time"

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
		if strings.Contains(p.DynamicParams, "network:") {
			if strings.Contains(p.DynamicParams, "WIFI") {
				config, _ = service.ConfigService.GetConfig(1)
			} else {
				config, _ = service.ConfigService.GetConfig(2)
			}
		} else {
			if strings.Contains(p.ExtParams, "WIFI") {
				config, _ = service.ConfigService.GetConfig(1)
			} else {
				config, _ = service.ConfigService.GetConfig(2)
			}
		}
	} else {
		config, _ = service.ConfigService.GetConfig(1)
	}
	msg := fmt.Sprintf("config_get:{ip:%s,pid:%s,config:%d,timestr:%s}", addr, p.PeerId, config.Id, p.TimeStr)
	time := time.Now().Format("2006-01-02 15:04:05")
	msg = time + "\u0003" + addr + "\u0003" + msg
	logs.Info(msg)
	if flag {
		SendToKafka(msg, "test")
	} else {
		err := InitKafka()
		if err != nil {

		} else {
			SendToKafka(msg, "test")
		}

	}
	this.Ctx.WriteString(config.Value)
}
