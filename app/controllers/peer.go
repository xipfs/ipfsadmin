package controllers

/*
 ============================================================================
 Name        : peer.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 节点
 ============================================================================
*/
import (
	"fmt"
	"time"

	"context"
	"io"
	"net/http"
	"os"

	"github.com/astaxie/beego"
	"github.com/gogo/protobuf/proto"
	"github.com/xipfs/ipfsadmin/app/entity"
	"github.com/xipfs/ipfsadmin/app/libs"
	"github.com/xipfs/ipfsadmin/app/msg"
	"github.com/xipfs/ipfsadmin/app/service"
)

type PeerController struct {
	BaseController
}

// 首页
func (this *PeerController) Index() {
	this.Data["pageTitle"] = "节点监控"
	this.display()
}

// 节点列表
func (this *PeerController) List() {
	status, _ := this.GetInt("status")
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	filter := make([]interface{}, 0, 6)
	if status == 0 {
		filter = append(filter, "status", 1)
	} else {
		filter = append(filter, "status", status)
	}

	list, count := service.PeerService.GetList(page, this.pageSize, filter...)
	this.Data["pageTitle"] = "节点列表"
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("PeerController.List", "status", status), true).ToString()
	this.display()
}

// 上报节点信息
func (this *PeerController) Report() {
	reportRecordList := &msg.ReportRecordList{}
	requestBody := this.GetRequestBody()
	proto.Unmarshal(requestBody, reportRecordList)
	fmt.Println(reportRecordList)
	const base_format = "2006-01-02 15:04:05"
	peerId := ""
	for _, v := range reportRecordList.Records {
		p := &entity.PeerLog{}
		p.EventAction = v.EventAction
		p.Goarch = v.CommonData["goarch"]
		p.Goos = v.CommonData["goos"]
		p.Mac = v.CommonData["mac"]
		p.PeerId = v.CommonData["peer_id"]
		peerId = p.PeerId
		p.CreateTime, _ = time.Parse(base_format, v.CommonData["timestr"])
		//fmt.Println(p)
		err := service.PeerLogService.AddPeerLog(p)
		if err != nil {
			out := make(map[string]interface{})
			out["status"] = "-1"
			out["msg"] = "error"
			this.jsonResult(out)
			return
		}
	}
	if peerId != "" {
		peer := &entity.Peer{}
		peer.Status = 1
		peer.PeerId = peerId
		peer.UpdateTime = time.Now()
		peer.CreateTime = time.Now()
		err := service.PeerService.AddPeer(peer, "PeerId", "Status", "UpdateTime", "CreateTime")
		if err != nil {
			fmt.Println(err)
		}
	}
	out := make(map[string]interface{})
	out["status"] = "1"
	out["msg"] = "ok"
	this.jsonResult(out)
}

// 发布资源
func (this *PeerController) Pub() {
	//下载文件
	fileUrl := this.GetString("fileUrl")
	fileName := this.GetString("fileName")
	out := make(map[string]interface{})
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", fileUrl, nil)

	if err != nil {
		fmt.Println(err)
		out["status"] = "-1"
		out["msg"] = "error"
		this.jsonResult(out)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	reqest = reqest.WithContext(ctx)
	response, err := client.Do(reqest)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		out["status"] = "-1"
		out["msg"] = "error"
		this.jsonResult(out)
		return
	}

	f, err := os.Create(beego.AppConfig.String("pub_dir") + fileName)
	if err != nil {
		out["status"] = "-1"
		out["msg"] = "error"
		this.jsonResult(out)
	}
	io.Copy(f, response.Body)
	defer response.Body.Close()
	defer f.Close()
	defer cancel()

	// 发布资源
	p := &entity.Resource{}
	p.Name = fileName
	p.Domain = ""
	p.MD5 = ""
	p.Version = ""
	p.RepoUrl = fileUrl
	p.TaskReview = 0
	err = service.ResourceService.AddResource(p)

	//构建任务
	task := new(entity.Task)
	task.ResourceId = p.Id
	task.Message = ""
	task.UserId = this.userId
	task.UserName = this.auth.GetUser().UserName
	task.FileName = p.Name
	task.PubEnvId = 1
	task.BuildStatus = 1

	err = service.TaskService.AddTask(task)
	service.ActionService.Add("create_task", this.auth.GetUserName(), "task", task.Id, "")
	service.DeployService.DoDeploy(task)
	out["status"] = "1"
	out["msg"] = "ok"
	this.jsonResult(out)
}

func (this *PeerController) GetPeerStat() {
	rangeType := this.GetString("range")
	result := service.PeerService.GetPeerStat(rangeType)

	ticks := make([]interface{}, 0)
	chart := make([]interface{}, 0)
	json := make(map[string]interface{}, 0)
	switch rangeType {
	case "this_month":
		year, month, _ := time.Now().Date()
		maxDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0).AddDate(0, 0, -1).Day()

		for i := 1; i <= maxDay; i++ {
			var row [3]interface{}
			row[0] = i
			row[1] = fmt.Sprintf("%02d", i)
			row[2] = fmt.Sprintf("%d-%02d-%02d", year, month, i)
			ticks = append(ticks, row)
			if v, ok := result[i]; ok {
				chart = append(chart, []int{i, v})
			} else {
				chart = append(chart, []int{i, 0})
			}
		}
	case "last_month":
		year, month, _ := time.Now().AddDate(0, -1, 0).Date()
		maxDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0).AddDate(0, 0, -1).Day()

		for i := 1; i <= maxDay; i++ {
			var row [3]interface{}
			row[0] = i
			row[1] = fmt.Sprintf("%02d", i)
			row[2] = fmt.Sprintf("%d-%02d-%02d", year, month, i)
			ticks = append(ticks, row)
			if v, ok := result[i]; ok {
				chart = append(chart, []int{i, v})
			} else {
				chart = append(chart, []int{i, 0})
			}
		}
	case "this_year":
		year := time.Now().Year()
		for i := 1; i <= 12; i++ {
			var row [3]interface{}
			row[0] = i
			row[1] = fmt.Sprintf("%d月", i)
			row[2] = fmt.Sprintf("%d年%d月", year, i)
			ticks = append(ticks, row)
			if v, ok := result[i]; ok {
				chart = append(chart, []int{i, v})
			} else {
				chart = append(chart, []int{i, 0})
			}
		}
	case "last_year":
		year := time.Now().Year() - 1
		for i := 1; i <= 12; i++ {
			var row [3]interface{}
			row[0] = i
			row[1] = fmt.Sprintf("%d月", i)
			row[2] = fmt.Sprintf("%d年%d月", year, i)
			ticks = append(ticks, row)
			if v, ok := result[i]; ok {
				chart = append(chart, []int{i, v})
			} else {
				chart = append(chart, []int{i, 0})
			}
		}
	}

	json["ticks"] = ticks
	json["chart"] = chart
	this.Data["json"] = json
	this.ServeJSON()
}
