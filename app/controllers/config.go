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

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xipfs/ipfsadmin/app/entity"
	"github.com/xipfs/ipfsadmin/app/service"
)

type ConfigController struct {
	BaseController
}

type ConfigParam struct {
	PeerId        string `json:"peer_id"`
	TimeStr       string `json:"local_timestr"`
	Version       string `json:"version"`
	Goos          string `json:"goos"`
	Goarch        string `json:"goarch"`
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
	time := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("%s\u0003%s\u0003%s\u0003%s\u0001%s\u0001%s\u0001%s\u0001%s\u0001%s\u0001%s\u0001%s\u0001%s\u0001%s\u0001%s\u0001%d\u0002", time, addr, p.Version, "config_get", p.Version, p.Version, p.PeerId, p.Goarch, "", "", p.TimeStr, "", p.ExtParams, p.DynamicParams, config.Id)
	if p.Version == "3.0" {
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
	}
	this.Ctx.WriteString(config.Value)
}

// 模板列表
func (this *ConfigController) List() {
	list, _ := service.ConfigService.GetAllConfig()
	this.Data["pageTitle"] = "配置"
	this.Data["list"] = list
	this.display()
}

// 添加配置
func (this *ConfigController) Add() {
	if this.isPost() {
		key := this.GetString("key")
		value := this.GetString("value")

		if key == "" || value == "" {
			this.showMsg("配置 key、value 不能为空", MSG_ERR)
		}

		config := new(entity.Config)
		config.Key = key
		config.Value = value
		err := service.ConfigService.AddConfig(config)
		this.checkError(err)
		this.redirect(beego.URLFor("ConfigController.List"))
	}
	this.Data["pageTitle"] = "添加配置"
	this.display()
}

// 编辑模板
func (this *ConfigController) Edit() {
	id, _ := this.GetInt("id")
	config, err := service.ConfigService.GetConfig(id)
	this.checkError(err)

	if this.isPost() {
		key := this.GetString("key")
		value := this.GetString("value")

		if key == "" || value == "" {
			this.showMsg("配置 key、value 不能为空", MSG_ERR)
		}

		config.Key = key
		config.Value = value
		err := service.ConfigService.SaveConfig(config)
		this.checkError(err)

		this.redirect(beego.URLFor("ConfigController.List"))
	}

	this.Data["pageTitle"] = "修改配置"
	this.Data["config"] = config
	this.display()
}

// 删除模板
func (this *ConfigController) Del() {
	id, _ := this.GetInt("id")

	err := service.ConfigService.DelConfig(id)
	this.checkError(err)

	this.redirect(beego.URLFor("ConfigController.List"))
}
