package service

import (
	"fmt"
	"net/url"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xipfs/ipfsadmin/app/entity"
)

var (
	o                  orm.Ormer
	tablePrefix        string              // 表前缀
	UserService        *userService        // 用户服务
	RoleService        *roleService        // 角色服务
	EnvService         *envService         // 发布环境服务
	ServerService      *serverService      // 服务器服务
	ResourceService    *resourceService    // 资源服务
	MailService        *mailService        // 邮件服务
	TaskService        *taskService        // 任务服务
	DeployService      *deployService      // 部署服务
	RepositoryService  *repositoryService  // 仓库服务
	PeerService        *peerService        // 节点
	PeerLogService     *peerLogService     // 节点日志
	DownloadLogService *downloadLogService // 下载日志
	ConfigService      *configService      //  配置
	SystemService      *systemService
	ActionService      *actionService // 系统动态
	ApiService         *apiService    // api
)

func Init() {
	dbHost := beego.AppConfig.String("db.host")
	dbPort := beego.AppConfig.String("db.port")
	dbUser := beego.AppConfig.String("db.user")
	dbPassword := beego.AppConfig.String("db.password")
	dbName := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	tablePrefix = beego.AppConfig.String("db.prefix")

	if dbPort == "" {
		dbPort = "3306"
	}
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8"
	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	orm.RegisterDataBase("default", "mysql", dsn)

	orm.RegisterModelWithPrefix(tablePrefix,
		new(entity.Action),
		new(entity.Env),
		new(entity.EnvServer),
		new(entity.MailTpl),
		new(entity.Resource),
		new(entity.Role),
		new(entity.Server),
		new(entity.Task),
		new(entity.User),
		new(entity.Peer),
		new(entity.PeerLog),
		new(entity.DownloadLog),
		new(entity.Config),
		new(entity.Api),
	)

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	o = orm.NewOrm()
	orm.RunCommand()

	// 创建代码
	os.MkdirAll(GetResourcesBasePath(), 0755)
	os.MkdirAll(GetTasksBasePath(), 0755)

	// 初始化服务对象
	initService()
}

func initService() {
	UserService = &userService{}
	RoleService = &roleService{}
	EnvService = &envService{}
	ServerService = &serverService{}
	ResourceService = &resourceService{}
	MailService = &mailService{}
	TaskService = &taskService{}
	DeployService = &deployService{}
	RepositoryService = &repositoryService{}
	SystemService = &systemService{}
	ActionService = &actionService{}
	PeerService = &peerService{}
	PeerLogService = &peerLogService{}
	DownloadLogService = &downloadLogService{}
	ConfigService = &configService{}
	ApiService = &apiService{}
}

// 返回真实表名
func tableName(name string) string {
	return tablePrefix + name
}

func debug(v ...interface{}) {
	beego.Debug(v...)
}

// 任务单根目录
func GetTasksBasePath() string {
	return fmt.Sprintf(beego.AppConfig.String("data_dir") + "/tasks")
}

// 所有资源根目录
func GetResourcesBasePath() string {
	return fmt.Sprintf(beego.AppConfig.String("data_dir") + "/resources")
}

// 任务单目录
func GetTaskPath(id int) string {
	return fmt.Sprintf(GetTasksBasePath()+"/task-%d", id)
}

// 某个资源的代码目录
func GetResourcePath(name string) string {
	return GetResourcesBasePath() + "/" + name
}

func concatenateError(err error, stderr string) error {
	if len(stderr) == 0 {
		return err
	}
	return fmt.Errorf("%v: %s", err, stderr)
}

func DBVersion() string {
	var lists []orm.ParamsList
	o.Raw("SELECT VERSION()").ValuesList(&lists)
	return lists[0][0].(string)
}
