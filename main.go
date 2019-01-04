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
	"time"

	"fmt"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"github.com/xipfs/ipfsadmin/app/controllers"
	_ "github.com/xipfs/ipfsadmin/app/mail"
	"github.com/xipfs/ipfsadmin/app/service"
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
