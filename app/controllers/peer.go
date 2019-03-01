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
	"time"

	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xipfs/ipfsadmin/app/entity"
	"github.com/xipfs/ipfsadmin/app/libs"
	"github.com/xipfs/ipfsadmin/app/service"
)

type PeerController struct {
	BaseController
}

var (
	client sarama.SyncProducer
	flag   = false
)

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
	requestBody := this.GetRequestBody()
	msg := string(requestBody)
	logs.Info("kafka report ", msg)
	if strings.HasPrefix(msg, "3.0") {
		req := this.Ctx.Request
		addr := req.RemoteAddr
		time := time.Now().Format("2006-01-02 15:04:05")
		msg = time + "\u0003" + addr + "\u0003" + msg
		if flag {
			SendToKafka(msg, "test")
		} else {
			err := InitKafka()
			if err != nil {

			} else {
				SendToKafka(msg, "test")
			}

		}
	} else {

	}
	out := make(map[string]interface{})
	out["status"] = "1"
	out["msg"] = "ok"
	this.jsonResult(out)
}

/*初始化kafka*/
func InitKafka() (err error) {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err = sarama.NewSyncProducer([]string{"172.31.31.252:9092"}, config)
	if err != nil {
		logs.Error("init kafka producer failed, err:", err)
		return err
	}
	//记录步骤信息
	logs.Info("init kafka success")
	flag = true
	return nil
}

/*
   发送到kafak
*/
func SendToKafka(data, topic string) (err error) {

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		logs.Error("send message failed, err:%v data:%v topic:%v", err, data, topic)
		return
	}

	logs.Info("send succ, pid:%v offset:%v, topic:%v\n", pid, offset, topic)
	return
}

// 发布资源
func (this *PeerController) Pub() {
	//下载文件
	fmt.Println("pub~~~~~")
	fileName := this.GetString("fileName")
	//requestBody := this.GetRequestBody()

	out := make(map[string]interface{})
	f, err := os.Create(beego.AppConfig.String("pub_dir") + fileName)
	if err != nil {
		out["status"] = "-1"
		out["msg"] = "error"
		this.jsonResult(out)
		return
	}
	_, err2 := io.Copy(f, this.Controller.Ctx.Request.Body)
	if err2 != nil {
		out["status"] = "-1"
		out["msg"] = "error"
		this.jsonResult(out)
		return
	}
	fmt.Println("获取下载列表成功 ！")
	defer f.Close()

	go func() {
		fi, err := os.Open(path.Join(beego.AppConfig.String("pub_dir"), fileName))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, fileName+" 保存地址文件失败 ！")
			return
		}
		defer fi.Close()
		br := bufio.NewReader(fi)
		m := make(map[string]string)  // package name -> url
		m2 := make(map[string]string) // package name -> md5
		total := 0
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			packageName := string(a)
			total++
			fmt.Println(packageName)

			p := &entity.Resource{}
			p.Name = ""
			p.Domain = packageName
			p.MD5 = ""
			p.Version = ""
			p.RepoUrl = ""
			p.TaskReview = 0
			p.Status = 0
			p.UploadFileName = fileName

			add_err := service.ResourceService.AddResource(p)
			if add_err != nil {
				fmt.Printf("Error: %s\n", err)
				service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, fileName+"处理同步配置["+packageName+"]失败 ！")
				continue
			}

			resp, err := http.Get("http://ams.lenovomm.com/ams/3.0/appdownaddress.do?dt=0&ty=2&pn=" + string(a) + "&cid=12654&tcid=12654&ic=0")
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, fileName+" 获取 apk 下载地址失败 ！")
				continue
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, fileName+" 获取 apk 下载地址失败 ！")
				continue
			}
			service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, fileName+" 获取 apk 地址成功 ！")
			var app App
			if err := json.Unmarshal(body, &app); err == nil {
				service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, fileName+" 获取 MD5 "+app.MD5+"成功 ！")
				m2[packageName] = app.MD5
				for _, v := range app.Urls {
					m[packageName] = v.DownUrl
					break
				}
			}
		}
		service.ActionService.Add("publish", this.auth.GetUserName(), "publish", 1000, fileName+" 获取到 "+strconv.Itoa(total)+" 个待更新 APK 信息 !")
		fmt.Println(fileName + " 获取到 " + strconv.Itoa(total) + " 个待更新 APK 信息 !")
		for k, v := range m {
			name := strings.Split(filepath.Base(v), "?")[0]
			pub(name, v, k, fileName, m2[k])
		}
	}()
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
