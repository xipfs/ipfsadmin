package controllers

/*
 ============================================================================
 Name        : server.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 服务器
 ============================================================================
*/

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/xipfs/ipfsadmin/app/entity"
	"github.com/xipfs/ipfsadmin/app/libs"
	"github.com/xipfs/ipfsadmin/app/service"
)

type ServerController struct {
	BaseController
}

// 列表
func (this *ServerController) List() {
	page, _ := strconv.Atoi(this.GetString("page"))
	if page < 1 {
		page = 1
	}
	count, _ := service.ServerService.GetTotal(service.SERVER_TYPE_NORMAL)
	//this.checkError(err)
	serverList, _ := service.ServerService.GetServerList(page, this.pageSize)
	//this.checkError(err)

	this.Data["count"] = count
	this.Data["list"] = serverList
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("ServerController.List"), true).ToString()
	this.Data["pageTitle"] = "服务器列表"
	this.display()
}

// 添加
func (this *ServerController) Add() {
	if this.isPost() {
		server := &entity.Server{}
		server.TypeId = service.SERVER_TYPE_NORMAL
		server.Ip = this.GetString("server_ip")
		server.Area = this.GetString("area")
		server.SshPort, _ = this.GetInt("ssh_port")
		server.SshUser = this.GetString("ssh_user")
		server.SshPwd = this.GetString("ssh_pwd")
		server.SshKey = this.GetString("ssh_key")
		server.WorkDir = this.GetString("work_dir")
		server.Description = this.GetString("description")
		err := this.validServer(server)
		this.checkError(err)
		err = service.ServerService.AddServer(server)
		this.checkError(err)
		this.redirect(beego.URLFor("ServerController.List"))
	}

	this.Data["pageTitle"] = "添加服务器"
	this.display()
}

// 编辑
func (this *ServerController) Edit() {
	id, _ := this.GetInt("id")
	server, err := service.ServerService.GetServer(id, service.SERVER_TYPE_NORMAL)
	this.checkError(err)

	if this.isPost() {
		server.Ip = this.GetString("server_ip")
		server.Area = this.GetString("area")
		server.SshPort, _ = this.GetInt("ssh_port")
		server.SshUser = this.GetString("ssh_user")
		server.SshPwd = this.GetString("ssh_pwd")
		server.SshKey = this.GetString("ssh_key")
		server.WorkDir = this.GetString("work_dir")
		server.Description = this.GetString("description")
		err := this.validServer(server)
		this.checkError(err)
		err = service.ServerService.UpdateServer(server)
		this.checkError(err)
		this.redirect(beego.URLFor("ServerController.List"))
	}

	this.Data["pageTitle"] = "编辑服务器"
	this.Data["server"] = server
	this.display()
}

// 删除
func (this *ServerController) Del() {
	id, _ := this.GetInt("id")

	_, err := service.ServerService.GetServer(id, service.SERVER_TYPE_NORMAL)
	this.checkError(err)

	err = service.ServerService.DeleteServer(id)
	this.checkError(err)
	this.redirect(beego.URLFor("ServerController.List"))
}

// 资源列表
func (this *ServerController) Resources() {
	id, _ := this.GetInt("id")
	server, _ := service.ServerService.GetServer(id, service.SERVER_TYPE_NORMAL)
	//this.checkError(err)
	envList, _ := service.EnvService.GetEnvListByServerId(id)
	//this.checkError(err)

	result := make(map[int]map[string]interface{})
	for _, env := range envList {
		if _, ok := result[env.ResourceId]; !ok {
			resource, err := service.ResourceService.GetResource(env.ResourceId)
			if err != nil {
				continue
			}
			row := make(map[string]interface{})
			row["resourceId"] = resource.Id
			row["resourceName"] = resource.Name
			row["envName"] = env.Name
			result[env.ResourceId] = row
		} else {
			result[env.ResourceId]["envName"] = result[env.ResourceId]["envName"].(string) + ", " + env.Name
		}
	}

	this.Data["list"] = result
	this.Data["server"] = server
	this.Data["pageTitle"] = server.Ip + " 下的资源列表"
	this.display()
}

func (this *ServerController) validServer(server *entity.Server) error {
	valid := validation.Validation{}
	valid.Required(server.Ip, "ip").Message("请输入服务器IP")
	valid.Range(server.SshPort, 1, 65535, "ssh_port").Message("SSH端口无效")
	valid.Required(server.SshUser, "ssh_user").Message("SSH用户名不能为空")
	valid.Required(server.SshPwd, "ssh_pwd").Message("SSH密码不能为空")
	valid.IP(server.Ip, "ip").Message("服务器IP无效")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}

	addr := fmt.Sprintf("%s:%d", server.Ip, server.SshPort)
	serv := libs.NewServerConn(addr, server.SshUser, server.SshPwd)

	if err := serv.TryConnect(); err != nil {
		return errors.New("无法连接到服务器: " + err.Error())
	}
	serv.Close()

	return nil
}
