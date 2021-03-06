package controllers

/*
 ============================================================================
 Name        : role.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 角色
 ============================================================================
*/

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/xipfs/ipfsadmin/app/entity"
	"github.com/xipfs/ipfsadmin/app/service"
)

type RoleController struct {
	BaseController
}

func (this *RoleController) List() {
	roleList, err := service.RoleService.GetAllRoles()
	this.checkError(err)
	for k, role := range roleList {
		roleList[k].UserList, _ = service.UserService.GetUserListByRoleId(role.Id)
	}
	this.Data["pageTitle"] = "角色管理"
	this.Data["list"] = roleList
	this.display()
}

func (this *RoleController) Add() {
	if this.isPost() {
		role := &entity.Role{}
		role.RoleName = this.GetString("role_name")
		role.Description = this.GetString("description")
		if role.RoleName == "" {
			this.showMsg("角色名不能为空", MSG_ERR)
		}
		err := service.RoleService.AddRole(role)
		this.checkError(err)
		this.redirect(beego.URLFor("RoleController.List"))
	}
	this.Data["pageTitle"] = "创建角色"
	this.display()
}

func (this *RoleController) Edit() {
	id, _ := this.GetInt("id")
	role, err := service.RoleService.GetRole(id)
	this.checkError(err)

	if this.isPost() {
		role.RoleName = this.GetString("role_name")
		role.Description = this.GetString("description")
		err := service.RoleService.UpdateRole(role, "RoleName", "Description")
		this.checkError(err)
		this.redirect(beego.URLFor("RoleController.List"))
	}

	this.Data["pageTitle"] = "编辑角色"
	this.Data["role"] = role
	this.display()
}

func (this *RoleController) Del() {
	id, _ := this.GetInt("id")

	err := service.RoleService.DeleteRole(id)
	this.checkError(err)

	this.redirect(beego.URLFor("RoleController.List"))
}

func (this *RoleController) Perm() {
	id, _ := this.GetInt("id")
	role, err := service.RoleService.GetRole(id)
	this.checkError(err)

	if this.isPost() {
		pids := this.GetStrings("pids")
		perms := this.GetStrings("perms")
		if len(pids) == 0 {
			role.ResourceIds = ""
		} else {
			role.ResourceIds = strings.Join(pids, ",")
		}
		err := service.RoleService.UpdateRole(role, "ResourceIds")
		this.checkError(err)
		err = service.RoleService.SetPerm(role.Id, perms)
		this.checkError(err)
		this.redirect(beego.URLFor("RoleController.List"))
	}

	resourceList, _ := service.ResourceService.GetAllResource()
	permList := service.SystemService.GetPermList()

	chkmap := make(map[string]string)
	for _, v := range role.PermList {
		chkmap[v.Key] = "checked"
	}
	if role.ResourceIds != "" {
		pids := strings.Split(role.ResourceIds, ",")
		for _, pid := range pids {
			chkmap[pid] = "checked"
		}
	}

	this.Data["pageTitle"] = "编辑权限"
	this.Data["permList"] = permList
	this.Data["resourceList"] = resourceList
	this.Data["role"] = role
	this.Data["chkmap"] = chkmap
	this.display()
}
