package main

/*
 ============================================================================
 Name        : main.go
 Author      : xiehui
 Version     : 1.0
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 入口程序
 ============================================================================
*/
import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/astaxie/beego/logs"

	"bufio"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"github.com/robfig/cron"
	"github.com/xipfs/ipfsadmin/app/controllers"
	_ "github.com/xipfs/ipfsadmin/app/mail"
	"github.com/xipfs/ipfsadmin/app/service"

	"github.com/xipfs/ipfsadmin/app/entity"
)

const VERSION = "1.0.0"

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

func main() {
	service.Init()
	fmt.Println("start app .....")
	i := 0
	c := cron.New()
	spec := "0 0 0,5,10,15,20 * * ?"
	c.AddFunc(spec, func() {
		i++
		fmt.Println("cron running:", i)
		logs.Info("开始定时同步最新的 APP !")
		fi, err := os.Open(path.Join(beego.AppConfig.String("pub_dir"), "pkgs-abtest.txt"))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			service.ActionService.Add("publish", "root", "publish", 1000, "pkgs-abtest.txt 保存地址文件失败 ！")
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
				service.ActionService.Add("publish", "root", "publish", 1000, "pkgs-abtest.txt 获取 apk 下载地址失败 ！")
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				service.ActionService.Add("publish", "root", "publish", 1000, "pkgs-abtest.txt 获取 apk 下载地址失败 ！")
				return
			}
			//json str 转struct
			service.ActionService.Add("publish", "root", "publish", 1000, "pkgs-abtest.txt 获取 apk 地址成功 ！")
			var app App
			if err := json.Unmarshal(body, &app); err == nil {
				service.ActionService.Add("publish", "root", "publish", 1000, "pkgs-abtest.txt 获取 MD5 "+app.MD5+"成功 ！")
				m2[packageName] = app.MD5
				for _, v := range app.Urls {
					m[packageName] = v.DownUrl
					break
				}
			}
		}
		for k, v := range m {
			name := strings.Split(filepath.Base(v), "?")[0]
			pub(name, v, k, "pkgs-abtest.txt", m2[k])
		}
	})
	c.Start()

	c2 := cron.New()
	spec2 := "0 0 1,6,11,16,21 * * ?"
	//spec2 := "*/5 * * * * ?"
	c2.AddFunc(spec2, func() {
		i++
		fmt.Println("cron running:", i)
		logs.Info("开始将生成的 IPFS 地址同步到服务器 ！")
		uploadFileNames, _ := service.ResourceService.GetAllResourceByName("pkgs-abtest.txt")
		fmt.Println("file length:", len(uploadFileNames))
		var buffer bytes.Buffer
		buffer.WriteString("{\"pkgs\":[")
		for _, v := range uploadFileNames {
			buffer.WriteString("{\"pn\":\"")
			buffer.WriteString(v.Domain)
			buffer.WriteString("\",")
			buffer.WriteString("\"url\":\"")
			buffer.WriteString("http://127.0.0.1:8080/ipfs/" + v.Hash + "?channel=lestore&ftype=apk")
			buffer.WriteString("'||'&'||'ftype=apk'")
			buffer.WriteString("\"},")
		}
		buffer.WriteString("]}")
		var jsonStr = []byte(buffer.String())
		fmt.Println("jsonStr", buffer.String())
		req, err := http.NewRequest("POST", "http://sp.lenovomm.com/fui/operator/dataAnalysis/service/baizeAppUpdate", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Printf("发送请求到 baizeAppUpdate 失败，原因: %s\n", err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("response Body: %s", string(body))
	})
	c2.Start()
	beego.AppConfig.Set("version", VERSION)
	if beego.AppConfig.String("runmode") == "dev" {
		beego.SetLevel(beego.LevelDebug)
	} else {
		beego.SetLevel(beego.LevelInformational)
		beego.SetLogger("file", `{"maxdays":90,"filename":"`+beego.AppConfig.String("log_file")+`"}`)
		beego.BeeLogger.DelLogger("console")
	}

	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.Router("/login", &controllers.MainController{}, "*:Login")
	beego.Router("/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/profile", &controllers.MainController{}, "*:Profile")

	beego.AutoRouter(&controllers.ResourceController{})
	beego.AutoRouter(&controllers.TaskController{})
	beego.AutoRouter(&controllers.ServerController{})
	beego.AutoRouter(&controllers.EnvController{})
	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.RoleController{})
	beego.AutoRouter(&controllers.MailTplController{})
	beego.AutoRouter(&controllers.PeerController{})
	beego.AutoRouter(&controllers.PeerLogController{})
	beego.AutoRouter(&controllers.DownloadLogController{})
	beego.AutoRouter(&controllers.ConfigController{})
	beego.AutoRouter(&controllers.ApiController{})
	beego.AutoRouter(&controllers.MainController{})

	beego.AppConfig.Set("up_time", fmt.Sprintf("%d", time.Now().Unix()))

	beego.AddFuncMap("i18n", i18n.Tr)

	beego.SetStaticPath("/assets", "assets")
	beego.Run()

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
