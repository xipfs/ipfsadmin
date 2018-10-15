package controllers

/*
 ============================================================================
 Name        : api.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : api
 ============================================================================
*/

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/xipfs/ipfsadmin/app/libs"
	"github.com/xipfs/ipfsadmin/app/service"
)

type ApiController struct {
	BaseController
}

// API 列表
func (this *ApiController) List() {
	page, _ := strconv.Atoi(this.GetString("page"))
	if page < 1 {
		page = 1
	}

	count, _ := service.ApiService.GetTotal()
	apis, _ := service.ApiService.GetList(page, this.pageSize)

	this.Data["pageTitle"] = "API 管理"
	this.Data["count"] = count
	this.Data["list"] = apis
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("ApiController.List"), true).ToString()
	this.display()
}
