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
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
	ErrorMsg    string `json:"errorMsg"`
}

type Url struct {
	DownUrl string `json:"downurl"`
	AppSize int32  `json:"appSize"`
}

// 列表
func (this *TaskController) List() {
	//status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	list, count := service.ActionService.GetList(page, 10)
	this.Data["pageTitle"] = "状态列表"
	this.Data["count"] = count
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("TaskController.List"), true).ToString()
	this.Data["list"] = list
	//	startDate := this.GetString("start_date")
	//	endDate := this.GetString("end_date")
	//	resourceId, _ := this.GetInt("resource_id")
	//	if page < 1 {
	//		page = 1
	//	}
	//	filter := make([]interface{}, 0, 6)
	//	if resourceId > 0 {
	//		filter = append(filter, "resource_id", resourceId)
	//	}
	//	if startDate != "" {
	//		filter = append(filter, "start_date", startDate)
	//	}
	//	if endDate != "" {
	//		filter = append(filter, "end_date", endDate)
	//	}
	//	if status == 1 {
	//		filter = append(filter, "pub_status", 3)
	//	} else {
	//		filter = append(filter, "pub_status__lt", 3)
	//	}

	//	list, count := service.TaskService.GetList(page, this.pageSize, filter...)
	//	resourceList, _ := service.ResourceService.GetAllResource()

	//	this.Data["pageTitle"] = "发布单列表"
	//	this.Data["status"] = status
	//	this.Data["count"] = count
	//	this.Data["list"] = list
	//	this.Data["resourceList"] = resourceList
	//	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("TaskController.List", "status", status, "resource_id", resourceId, "start_date", startDate, "end_date", endDate), true).ToString()
	//	this.Data["resourceId"] = resourceId
	//	this.Data["startDate"] = startDate
	//	this.Data["endDate"] = endDate
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
		//image，这是一个key值，对应的是html中input type=‘file’的name属性值
		f, h, _ := this.GetFile("file")
		uploadFileName := this.GetString("uploadFileName")
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
		service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, uploadFileName+" 上传下载列表成功 ！")
		go func() {
			fi, err := os.Open(path.Join(beego.AppConfig.String("pub_dir"), fileName))
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, uploadFileName+" 保存地址文件失败 ！")
				return
			}
			defer fi.Close()
			br := bufio.NewReader(fi)
			m := make(map[string]string)  // package name -> url
			m2 := make(map[string]string) // package name -> md5
			for {
				a, _, c := br.ReadLine()
				if c == io.EOF {
					break
				}
				packageName := string(a)
				resp, err := http.Get("http://ams.lenovomm.com/ams/3.0/appdownaddress.do?dt=0&ty=2&pn=" + string(a) + "&cid=12654&tcid=12654&ic=0")
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, uploadFileName+" 获取 apk 下载地址失败 ！")
					return
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, uploadFileName+" 获取 apk 下载地址失败 ！")
					return
				}
				fmt.Println(string(body))
				//json str 转struct
				service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, uploadFileName+" 获取 apk 地址成功 ！")
				var app App
				if err := json.Unmarshal(body, &app); err == nil {
					fmt.Println("================json str 转struct==")
					fmt.Println(app.ErrorMsg)
					if app.ErrorMsg !=""{
						service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, uploadFileName+"  "+app.ErrorMsg)
						return
					}
					fmt.Println(app)
					fmt.Println(app.MD5)
					service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, uploadFileName+" 获取 MD5 "+app.MD5+"成功 ！")
					m2[packageName] = app.MD5
					for _, v := range app.Urls {
						m[packageName] = v.DownUrl
						break
					}
				}
			}
			for k, v := range m {
				name := strings.Split(filepath.Base(v), "?")[0]
				pub(name, v, k, uploadFileName, m2[k])
			}
		}()
		this.redirect(beego.URLFor("TaskController.List"))
	} else {
		this.Data["pageTitle"] = "新建发布任务"
		this.display("task/publish")
	}

}

// 发布资源
func pub(fileName string, fileUrl string, domain string, uploadFileName string, md5Original string) {
	//下载文件
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		fmt.Println(err)
		service.ActionService.Add("publish", "admin", "publish", 1000, uploadFileName+" 下载 APK "+fileName+" 失败 ！")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	reqest = reqest.WithContext(ctx)
	response, err := client.Do(reqest)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		service.ActionService.Add("publish", "admin", "publish", 1000, uploadFileName+" 下载 APK "+fileName+" 失败 ！")
		return
	}
	f, err := os.Create(beego.AppConfig.String("pub_dir") + fileName)
	if err != nil {
		fmt.Println(err)
		service.ActionService.Add("publish", "admin", "publish", 1000, uploadFileName+" 下载 APK "+fileName+" 失败 ！")
		return
	}
	n, err := io.Copy(f, response.Body)
	fmt.Printf("\n write %d err %v \n", n, err)
	defer response.Body.Close()
	defer f.Close()
	defer cancel()
	if err != nil {
		service.ActionService.Add("publish", "admin", "publish", 1000, uploadFileName+" 下载 APK "+fileName+"失败 ！")
		return
	}
	service.ActionService.Add("publish", "admin", "publish", 1000, uploadFileName+" 下载 APK "+fileName+"成功 ！")
	fmt.Println("download ok~~~")
	file, inerr := os.Open(beego.AppConfig.String("pub_dir") + fileName)
	if inerr != nil {
		service.ActionService.Add("publish", "admin", "publish", 1000, uploadFileName+" 打开 APK "+fileName+"失败 ！")
		return
	}
	md5h := md5.New()
	io.Copy(md5h, file)
	fileMd5 := strings.ToUpper(hex.EncodeToString(md5h.Sum([]byte(""))))
	fmt.Println("MD5 : " + fileMd5)
	defer file.Close()
	if fileMd5 == md5Original {
		service.ActionService.Add("publish", "admin", "publish", 1000, uploadFileName+" 校验 MD5 "+fileName+" 成功 ！")
	} else {
		service.ActionService.Add("publish", "admin", "publish", 1000, uploadFileName+" 校验 MD5 "+fileName+" 失败 ！")
		return
	}

	// 发布资源
	p := &entity.Resource{}
	p.Name = fileName
	p.Domain = domain
	p.MD5 = ""
	p.Version = ""
	p.RepoUrl = ""
	p.TaskReview = 0
	p.Status = 1
	p.UploadFileName = uploadFileName
	err = service.ResourceService.AddResource(p)
	fmt.Println("resource id --->", p.Id)
	//构建任务
	task := new(entity.Task)
	task.ResourceId = p.Id
	task.Message = ""
	task.UserId = 1
	task.UserName = "admin"
	task.FileName = p.Name
	task.PubEnvId = 1
	task.BuildStatus = 1
	task.UploadFileName = uploadFileName

	err = service.TaskService.AddTask(task)
	service.ActionService.Add("create_task", "admin", "task", task.Id, uploadFileName+" 开始部署 apk"+fileName)
	service.DeployService.DoDeploy(task)
}
