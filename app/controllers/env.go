package controllers

/*
 ============================================================================
 Name        : env.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 执行环境
 ============================================================================
*/
import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/xipfs/ipfsadmin/app/entity"
	"github.com/xipfs/ipfsadmin/app/service"
)

type EnvController struct {
	BaseController
}

func (this *EnvController) List() {
	envList, _ := service.EnvService.GetEnvList()
	this.Data["pageTitle"] = "发布环境配置"
	this.Data["envList"] = envList
	this.display()
}

func (this *EnvController) Add() {

	if this.isPost() {
		env := new(entity.Env)
		env.Name = this.GetString("name")
		env.SendMail, _ = this.GetInt("send_mail")
		env.MailTplId, _ = this.GetInt("mail_tpl_id")
		env.MailTo = this.GetString("mail_to")
		env.MailCc = this.GetString("mail_cc")

		if env.Name == "" {
			this.showMsg("环境名称", MSG_ERR)
		}

		serverIds := this.GetStrings("serverIds")
		if len(serverIds) < 1 {
			this.showMsg("请选择服务器", MSG_ERR)
		}

		if env.SendMail > 0 {
			if env.MailTplId == 0 {
				this.showMsg("请选择邮件模板", MSG_ERR)
			}
		}

		env.ServerList = make([]entity.Server, 0, len(serverIds))
		for _, v := range serverIds {
			if sid, _ := strconv.Atoi(v); sid > 0 {
				if sv, err := service.ServerService.GetServer(sid); err == nil {
					env.ServerList = append(env.ServerList, *sv)
				} else {
					this.showMsg("服务器ID不存在: "+v, MSG_ERR)
				}
			}
		}
		if err := service.EnvService.AddEnv(env); err != nil {
			this.checkError(err)
		}

		this.redirect(beego.URLFor("EnvController.List"))
	}

	this.Data["serverList"], _ = service.ServerService.GetServerList(1, -1)
	this.Data["mailTplList"], _ = service.MailService.GetMailTplList()
	this.Data["pageTitle"] = "添加发布环境"
	this.display()
}

func (this *EnvController) Edit() {
	id, _ := this.GetInt("id")

	env, err := service.EnvService.GetEnv(id)
	this.checkError(err)

	if this.isPost() {
		env.Name = this.GetString("name")
		env.SendMail, _ = this.GetInt("send_mail")
		env.MailTplId, _ = this.GetInt("mail_tpl_id")
		env.MailTo = this.GetString("mail_to")
		env.MailCc = this.GetString("mail_cc")

		if env.Name == "" {
			this.showMsg("环境名称不能为空。", MSG_ERR)
		}

		serverIds := this.GetStrings("serverIds")
		if len(serverIds) < 1 {
			this.showMsg("请选择服务器", MSG_ERR)
		}

		if env.SendMail > 0 {
			if env.MailTplId == 0 {
				this.showMsg("请选择邮件模板", MSG_ERR)
			}
		}

		env.ServerList = make([]entity.Server, 0, len(serverIds))
		for _, v := range serverIds {
			if sid, _ := strconv.Atoi(v); sid > 0 {
				if sv, err := service.ServerService.GetServer(sid); err == nil {
					env.ServerList = append(env.ServerList, *sv)
				} else {
					this.showMsg("服务器ID不存在: "+v, MSG_ERR)
				}
			}
		}

		service.EnvService.SaveEnv(env)

		this.redirect(beego.URLFor("EnvController.List"))
	}

	serverList, _ := service.ServerService.GetServerList(1, -1)

	serverIds := make([]int, 0, len(env.ServerList))
	for _, v := range env.ServerList {
		serverIds = append(serverIds, v.Id)
	}

	jsonData, err := json.Marshal(serverIds)
	this.checkError(err)
	mailTplList, _ := service.MailService.GetMailTplList()

	this.Data["serverList"] = serverList
	this.Data["mailTplList"] = mailTplList
	this.Data["serverIds"] = string(jsonData)
	this.Data["env"] = env
	this.Data["pageTitle"] = "编辑发布环境"
	this.display()
}

func (this *EnvController) Del() {
	id, _ := this.GetInt("id")
	service.EnvService.DeleteEnv(id)
	this.redirect(beego.URLFor("EnvController.List"))
}
