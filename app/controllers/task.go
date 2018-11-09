package controllers

/*
 ============================================================================
 Name        : task.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 任务
 ============================================================================
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/astaxie/beego"
	"github.com/xipfs/ipfsadmin/app/entity"
	"github.com/xipfs/ipfsadmin/app/libs"
	"github.com/xipfs/ipfsadmin/app/service"
)

type TaskController struct {
	BaseController
}

type App struct {
	Version     string `json:"app_version"`
	VersionName string `json:"app_versionname"`
	MD5         string `json:"MD5"`
	downurls    string `json:"MD5"`
	Urls        []Url  `json:"downurls"`
}

type Url struct {
	DownUrl string `json:"downurl"`
	AppSize int32  `json:"appSize"`
}

// 列表
func (this *TaskController) List() {
	status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	startDate := this.GetString("start_date")
	endDate := this.GetString("end_date")
	resourceId, _ := this.GetInt("resource_id")
	if page < 1 {
		page = 1
	}
	filter := make([]interface{}, 0, 6)
	if resourceId > 0 {
		filter = append(filter, "resource_id", resourceId)
	}
	if startDate != "" {
		filter = append(filter, "start_date", startDate)
	}
	if endDate != "" {
		filter = append(filter, "end_date", endDate)
	}
	if status == 1 {
		filter = append(filter, "pub_status", 3)
	} else {
		filter = append(filter, "pub_status__lt", 3)
	}

	list, count := service.TaskService.GetList(page, this.pageSize, filter...)
	resourceList, _ := service.ResourceService.GetAllResource()

	this.Data["pageTitle"] = "发布单列表"
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["resourceList"] = resourceList
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("TaskController.List", "status", status, "resource_id", resourceId, "start_date", startDate, "end_date", endDate), true).ToString()
	this.Data["resourceId"] = resourceId
	this.Data["startDate"] = startDate
	this.Data["endDate"] = endDate
	this.display()
}

// 新建发布单
func (this *TaskController) Create() {

	if this.isPost() {
		resourceId, _ := this.GetInt("resource_id")
		envId, _ := this.GetInt("envId")
		message := this.GetString("editor_content")
		if envId < 1 {
			this.showMsg("请选择发布环境", MSG_ERR)
		}

		resource, err := service.ResourceService.GetResource(resourceId)
		this.checkError(err)
		task := new(entity.Task)
		task.ResourceId = resource.Id
		task.Message = message
		task.UserId = this.userId
		task.UserName = this.auth.GetUser().UserName
		task.FileName = resource.Name
		task.PubEnvId = envId

		err = service.TaskService.AddTask(task)
		this.checkError(err)

		// 构建任务
		go service.TaskService.BuildTask(task)

		service.ActionService.Add("create_task", this.auth.GetUserName(), "task", task.Id, "")

		this.redirect(beego.URLFor("TaskController.List"))
	}

	resourceId, _ := this.GetInt("resource_id")
	this.Data["pageTitle"] = "新建发布单"

	if resourceId < 1 {
		resourceList, _ := service.ResourceService.GetAllResource()
		this.Data["list"] = resourceList
		this.display("task/create_step1")
	} else {
		envList, _ := service.EnvService.GetEnvList()
		this.Data["resourceId"] = resourceId
		this.Data["envList"] = envList
		this.display()
	}
}

// 标签列表
func (this *TaskController) GetTags() {
	resourceId, _ := this.GetInt("resource_id")

	list, err := service.RepositoryService.GetTags(resourceId, 10)
	this.checkError(err)

	out := make(map[string]interface{})
	out["list"] = list
	this.jsonResult(out)
}

// 任务详情
func (this *TaskController) Detail() {
	taskId, _ := this.GetInt("id")
	task, err := service.TaskService.GetTask(taskId)
	this.checkError(err)
	env, err := service.EnvService.GetEnv(task.PubEnvId)
	this.checkError(err)
	this.Data["env"] = env
	this.Data["task"] = task
	this.Data["pageTitle"] = "发布单详情"
	this.display()
}

// 获取状态
func (this *TaskController) GetStatus() {
	taskId, _ := this.GetInt("id")
	tp := this.GetString("type")

	task, err := service.TaskService.GetTask(taskId)
	this.checkError(err)

	out := make(map[string]interface{})
	switch tp {
	case "pub":
		out["status"] = task.PubStatus
		if task.PubStatus < 0 {
			out["msg"] = task.ErrorMsg
		} else {
			out["msg"] = task.PubLog
		}

	default:
		out["status"] = task.BuildStatus
		out["msg"] = task.ErrorMsg
	}

	this.jsonResult(out)
}

// 开始发布
func (this *TaskController) StartPub() {
	taskId, _ := this.GetInt("id")

	if !this.auth.HasAccessPerm(this.controllerName, "publish") {
		this.showMsg("您没有执行该操作的权限", MSG_ERR)
	}
	err := service.DeployService.DeployTask(taskId)
	this.checkError(err)
	service.ActionService.Add("pub_task", this.auth.GetUserName(), "task", taskId, "")

	this.showMsg("", MSG_OK)
}

// 删除发布单
func (this *TaskController) Del() {
	taskId, _ := this.GetInt("id")
	refer := this.Ctx.Request.Referer()

	err := service.TaskService.DeleteTask(taskId)
	this.checkError(err)

	service.ActionService.Add("del_task", this.auth.GetUserName(), "task", taskId, "")

	if refer != "" {
		this.redirect(refer)
	} else {
		this.redirect(beego.URLFor("TaskController.List"))
	}
}

// 新建发布任务
func (this *TaskController) Publish() {
	if this.isPost() {
		//image，这是一个key值，对应的是html中input type-‘file’的name属性值
		f, h, _ := this.GetFile("file")
		//得到文件的名称
		fileName := h.Filename
		arr := strings.Split(fileName, ":")
		if len(arr) > 1 {
			index := len(arr) - 1
			fileName = arr[index]
		}
		fmt.Println("文件名称:")
		fmt.Println(fileName)
		//关闭上传的文件，不然的话会出现临时文件不能清除的情况
		f.Close()
		//保存文件到指定的位置
		//static/uploadfile,这个是文件的地址，第一个static前面不要有/
		this.SaveToFile("file", path.Join(beego.AppConfig.String("pub_dir"), fileName))
		go func() {
			fi, err := os.Open(path.Join(beego.AppConfig.String("pub_dir"), fileName))
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}
			defer fi.Close()
			br := bufio.NewReader(fi)
			m := make(map[string]string)
			for {
				a, _, c := br.ReadLine()
				if c == io.EOF {
					break
				}
				fmt.Println(string(a))
				resp, err := http.Get("http://ams.lenovomm.com/ams/3.0/appdownaddress.do?dt=0&ty=2&pn=" + string(a) + "&cid=12654&tcid=12654&ic=0")
				if err != nil {
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
				}
				fmt.Println(string(body))
				//json str 转struct
				var app App
				if err := json.Unmarshal(body, &app); err == nil {
					fmt.Println("================json str 转struct==")
					fmt.Println(app)
					fmt.Println(app.MD5)
					for _, v := range app.Urls {
						m[string(a)] = v.DownUrl
						break
					}

				}
			}
		}()
		this.redirect(beego.URLFor("TaskController.List"))
	} else {
		this.Data["pageTitle"] = "新建发布任务"
		this.display("task/publish")
	}

}
