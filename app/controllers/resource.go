package controllers

/*
 ============================================================================
 Name        : resource.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 资源
 ============================================================================
*/

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/xipfs/ipfsadmin/app/entity"
	"github.com/xipfs/ipfsadmin/app/libs"
	"github.com/xipfs/ipfsadmin/app/service"
)

type ResourceController struct {
	BaseController
}

type Resp struct {
	Status int    `json:"status"`
	Length int    `json:"length"`
	Datas  []Data `json:"data"`
}
type Data struct {
	Pn  string `json:"pn"`
	Url string `json:"url"`
}

// 资源列表
func (this *ResourceController) List() {
	page, _ := strconv.Atoi(this.GetString("page"))
	status, _ := this.GetInt("status")
	if page < 1 {
		page = 1
	}
	filter := make([]interface{}, 0, 6)
	if status == 2 {
		filter = append(filter, "status", 3)
	} else {
		filter = append(filter, "status__lt", 3)
	}

	list, count := service.ResourceService.GetListByFilter(page, this.pageSize, filter...)
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("ResourceController.List"), true).ToString()
	this.Data["pageTitle"] = "资源列表"
	this.display()
}

// 添加资源
func (this *ResourceController) Add() {

	if this.isPost() {
		p := &entity.Resource{}
		p.Name = this.GetString("resource_name")
		p.Domain = this.GetString("resource_domain")
		p.MD5 = this.GetString("resource_md5")
		p.Version = this.GetString("resource_version")
		p.TaskReview, _ = this.GetInt("task_review")

		if err := this.validResource(p); err != nil {
			this.showMsg(err.Error(), MSG_ERR)
		}

		err := service.ResourceService.AddResource(p)
		this.checkError(err)

		service.ActionService.Add("add_resource", this.auth.GetUserName(), "resource", p.Id, "")

		this.redirect(beego.URLFor("ResourceController.List"))
	}

	serverList, _ := service.ServerService.GetServerList(1, -1)
	//this.checkError(err)
	this.Data["pageTitle"] = "添加资源"
	this.Data["serverList"] = serverList
	this.display()
}

// 编辑资源
func (this *ResourceController) Edit() {
	id, _ := this.GetInt("id")
	p, err := service.ResourceService.GetResource(id)
	this.checkError(err)

	if this.isPost() {
		p.Name = this.GetString("resource_name")
		p.Domain = this.GetString("resource_domain")
		p.MD5 = this.GetString("resource_md5")
		p.Version = this.GetString("resource_version")
		p.TaskReview, _ = this.GetInt("task_review")

		if err := this.validResource(p); err != nil {
			this.showMsg(err.Error(), MSG_ERR)
		}

		err := service.ResourceService.UpdateResource(p, "Name", "Version", "Domain", "MD5", "TaskReview")
		this.checkError(err)

		service.ActionService.Add("edit_resource", this.auth.GetUserName(), "resource", p.Id, "")

		this.redirect(beego.URLFor("ResourceController.List"))
	}

	serverList, err := service.ServerService.GetServerList(1, -1)
	this.checkError(err)

	this.Data["resource"] = p
	this.Data["serverList"] = serverList
	this.Data["pageTitle"] = "编辑资源"
	this.display()
}

// 删除资源
func (this *ResourceController) Del() {
	id, _ := this.GetInt("id")

	err := service.ResourceService.DeleteResource(id)
	this.checkError(err)

	service.ActionService.Add("del_resource", this.auth.GetUserName(), "resource", id, "")

	this.redirect(beego.URLFor("ResourceController.List"))
}

// 重新发布资源
func (this *ResourceController) RePublish() {
	id, _ := this.GetInt("id")

	m := make(map[string]string)  // package name -> url
	m2 := make(map[string]string) // package name -> md5
	res,_ := service.ResourceService.GetResource(id)
	uploadFileName := res.UploadFileName
	go func() {
		fmt.Println("重新同步失败文件 " + res.Domain)
		resp, err := http.Get("http://ams.lenovomm.com/ams/3.0/appdownaddress.do?dt=0&ty=2&pn=" + res.Domain + "&cid=12654&tcid=12654&ic=0")
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			service.ActionService.Add("publish", "admin", "publish", 1000, res.Domain+" 获取 apk 下载地址失败 ！")
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			service.ActionService.Add("publish", "admin", "publish", 1000, res.Domain+" 获取 apk 下载地址失败 ！")
			return
		}
		service.ActionService.Add("publish", "admin", "publish", 1000, res.Domain+" 获取 apk 地址成功 ！")
		var app App
		if err := json.Unmarshal(body, &app); err == nil {
			if app.ErrorMsg !=""{
				service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, uploadFileName+"  "+app.ErrorMsg)
				return
			}
			service.ActionService.Add("publish", "admin", "publish", 1000, res.Domain+" 获取 MD5 "+app.MD5+"成功 ！")
			m2[res.Domain] = app.MD5
			for _, v := range app.Urls {
				m[res.Domain] = v.DownUrl
				break
			}
		}
		for k, v := range m {
			name := strings.Split(filepath.Base(v), "?")[0]
			pub(name, v, k, uploadFileName, m2[k])
		}
	}()

	this.redirect(beego.URLFor("ResourceController.List"))
}

// 验证提交
func (this *ResourceController) validResource(p *entity.Resource) error {
	errorMsg := ""
	if p.Name == "" {
		errorMsg = "请输入资源名称"
	} else if p.Domain == "" {
		errorMsg = "请输入资源标识"
	} else if p.Version == "" {
		errorMsg = "请输入版本"
	} else if p.MD5 == "" {
		errorMsg = "请输入MD5"
	} else {

	}

	if errorMsg != "" {
		return fmt.Errorf(errorMsg)
	}
	return nil
}

// 验证提交
func (this *ResourceController) Download() {
	uploadFileName := this.GetString("fileName")
	uploadFileNames, _ := service.ResourceService.GetAllResourceByName(uploadFileName)
	var buffer bytes.Buffer
	buffer.WriteString("update ams_ipfs_conf set AVAILABLE='0';\r\n")
	buffer.WriteString("commit;\r\n")
	for _, v := range uploadFileNames {
		buffer.WriteString("insert into ams_ipfs_conf(id,pn,url) values(s_ams_ipfs_conf.nextval,'")
		buffer.WriteString(v.Domain)
		buffer.WriteString("','")
		buffer.WriteString("http://127.0.0.1:8080/ipfs/" + v.Hash + "?channel=lestore&ftype=apk")
		buffer.WriteString("'||'&'||'ftype=apk');\r\n")
	}
	buffer.WriteString("commit;\r\n")
	f, _ := os.Create(beego.AppConfig.String("pub_dir") + uploadFileName)
	w := bufio.NewWriter(f)
	w.WriteString(buffer.String())
	w.Flush()
	f.Close()
	this.Ctx.Output.Download(beego.AppConfig.String("pub_dir")+uploadFileName, uploadFileName+".txt")
}

// 查询状态
func (this *ResourceController) Query() {
	uploadFileName := this.GetString("fileName")
	uploadFileNames, _ := service.ResourceService.GetAllResourceByName(uploadFileName)
	length := 0
	flag := true
	var datas [1024]Data
	m := make(map[string]string)  // package name -> url
	m2 := make(map[string]string) // package name -> md5
	for _, vv := range uploadFileNames {
		if vv.Status == 3 {

			datas[length].Pn = vv.Domain
			datas[length].Url = "http://127.0.0.1:8080/ipfs/" + vv.Hash + "?channel=lestore&ftype=apk'||'&'||'ftype=apk"
			length++
		} else if vv.Status == -1 {
			flag = false
			go func() {
				fmt.Println("重新同步失败文件 " + vv.Domain)
				resp, err := http.Get("http://ams.lenovomm.com/ams/3.0/appdownaddress.do?dt=0&ty=2&pn=" + vv.Domain + "&cid=12654&tcid=12654&ic=0")
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					service.ActionService.Add("publish", "admin", "publish", 1000, vv.Domain+" 获取 apk 下载地址失败 ！")
					return
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					service.ActionService.Add("publish", "admin", "publish", 1000, vv.Domain+" 获取 apk 下载地址失败 ！")
					return
				}
				service.ActionService.Add("publish", "admin", "publish", 1000, vv.Domain+" 获取 apk 地址成功 ！")
				var app App
				if err := json.Unmarshal(body, &app); err == nil {
					if app.ErrorMsg !=""{
						service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, uploadFileName+"  "+app.ErrorMsg)
						return
					}
					service.ActionService.Add("publish", "admin", "publish", 1000, vv.Domain+" 获取 MD5 "+app.MD5+"成功 ！")
					m2[vv.Domain] = app.MD5
					for _, v := range app.Urls {
						m[vv.Domain] = v.DownUrl
						break
					}
				}
				for k, v := range m {
					name := strings.Split(filepath.Base(v), "?")[0]
					pub(name, v, k, uploadFileName, m2[k])
				}
			}()
		} else {
			flag = false
		}

	}
	var resp Resp

	if flag {
		resp.Status = 1
	} else {
		resp.Status = 1
	}
	resp.Length = length
	resp.Datas = datas[0:length]
	this.jsonResult(resp)
	return
}
