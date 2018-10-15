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

	"github.com/xipfs/ipfsadmin/app/service"
)

type ConfigController struct {
	BaseController
}

type ConfigJson struct {
	Version    string
	DataReport *DataReport
}

type DataReport struct {
	ReportUrls          []string
	BatchReportNum      int
	MaxCacheFileSize    int
	DisableEventActions []string
	DataGathers         []DataGathers
}

type DataGathers struct {
	Key         string
	EventAction string
	Param       map[string]string
}

// 获取状态
func (this *ConfigController) Get() {
	var r ConfigJson
	configs, err := service.ConfigService.GetAllConfig()
	this.checkError(err)
	fmt.Println(configs)
	for _, v := range configs {
		if v.Key == "version1" {
			err := json.Unmarshal([]byte(v.Value), &r)
			if err != nil {
				fmt.Printf("err was %v", err)
			}
		}
	}
	this.jsonResult(r)
}
