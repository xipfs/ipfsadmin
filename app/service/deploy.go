package service

/*
 ============================================================================
 Name        : deploy.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 部署
 ============================================================================
*/
import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/xipfs/ipfsadmin/app/entity"
	"github.com/xipfs/ipfsadmin/app/libs"
	"github.com/xipfs/ipfsadmin/app/mail"
)

type deployService struct{}

var cnum chan int

// 执行部署任务
func (this *deployService) DeployTask(taskId int) error {
	task, err := TaskService.GetTask(taskId)
	if err != nil {
		return err
	}
	fmt.Println("deploy")
	task.PubStatus = 1
	task.ErrorMsg = ""
	TaskService.UpdateTask(task, "PubStatus", "ErrorMsg")

	go this.DoDeploy(task)

	return nil
}

func (this *deployService) DoDeploy(task *entity.Task) {
	// 1. 添加到本地库
	resource, err := ResourceService.GetResource(task.ResourceId)
	out, stderr, err := libs.ExecCmdDir(beego.AppConfig.String("pub_dir"), "baize", "add", "-w", task.FileName)
	fmt.Println("out", out)
	fmt.Println("stderr", stderr)
	fmt.Println("err", err)
	if err != nil {
		task.ErrorMsg = fmt.Sprintf(task.UploadFileName+" 添加到本地节点失败：%v", err)
		task.PubStatus = -2
		TaskService.UpdateTask(task, "PubStatus", "ErrorMsg")
		resource.Status = -1
		ResourceService.UpdateResource(resource, "Status")
		return
	}
	ActionService.Add("publish", "admin", "publish", 1000, task.UploadFileName+" 添加"+task.FileName+" 到本地成功 ！")
	words := strings.Split(string(out[:]), " ")
	hash := words[1]

	if err != nil {
		task.ErrorMsg = fmt.Sprintf("获取资源失败：%v", err)
		task.PubStatus = -2
		TaskService.UpdateTask(task, "PubStatus", "ErrorMsg")
		return
	}
	fmt.Println(hash)
	resource.Hash = hash
	resource.Status = 2
	ResourceService.UpdateResource(resource, "Hash", "Status")
	// 2. 发布到服务器
	env, _ := EnvService.GetEnv(task.PubEnvId)
	num := len(env.ServerList)
	if num < 1 {
		fmt.Errorf("服务器列表为空")
		return
	}
	cnum = make(chan int, num) //make一个chan,缓存为num
	for _, v := range env.ServerList {
		go this.PubToServer(task, v.Ip, v.SshPort, v.SshUser, v.SshPwd, hash)
	}
	flag := true
	for i := 0; i < num; i++ {
		value := <-cnum
		if value == -1 {
			flag = false
		}
	}
	if flag {
		task.PubTime = time.Now()
		task.PubStatus = 3
		task.ErrorMsg = ""
		TaskService.UpdateTask(task, "PubTime", "PubLog", "PubStatus", "ErrorMsg")
		resource.Status = 3
		ResourceService.UpdateResource(resource, "Status")
		ActionService.Add("publish", "admin", "publish", 1000, task.UploadFileName+"添加 "+task.FileName+" 到服务器成功！")
	} else {
		task.PubStatus = -3
		TaskService.UpdateTask(task, "PubTime", "PubLog", "PubStatus", "ErrorMsg")
		return
	}
	go func() {
		//3. 发送邮件
		fmt.Println("Send Mail :", env.SendMail)
		if env.SendMail > 0 {
			mailTpl, err := MailService.GetMailTpl(env.MailTplId)
			fmt.Println(err)
			if err == nil {
				replace := make(map[string]string)
				replace["{project}"] = resource.Name
				replace["{domain}"] = resource.Domain
				replace["{version}"] = resource.Version
				replace["{env}"] = env.Name
				replace["{description}"] = libs.Nl2br(html.EscapeString(task.Message))

				subject := mailTpl.Subject
				content := mailTpl.Content

				for k, v := range replace {
					subject = strings.Replace(subject, k, v, -1)
					content = strings.Replace(content, k, v, -1)
				}
				mailTo := strings.Split(mailTpl.MailTo+"\n"+env.MailTo, "\n")
				mailCc := strings.Split(mailTpl.MailCc+"\n"+env.MailCc, "\n")
				if err := mail.SendMail(subject, content, mailTo, mailCc); err != nil {
					beego.Error("邮件发送失败：", err)
				} else {
					beego.Info("发送邮件成功")
				}
			}
		}
	}()

}

func (this *deployService) Build(task *entity.Task) error {

	return nil
}

// 发布到服务器
func (this *deployService) PubToServer(task *entity.Task, ip string, port int, user string, pwd string, hash string) {
	// 连接到服务器
	addr := fmt.Sprintf("%s:%d", ip, port)
	server := libs.NewServerConn(addr, user, pwd)
	defer server.Close()
	beego.Debug("连接服务器: ", addr, ", 用户: ", user)
	ActionService.Add("publish", "admin", "publish", 1000, task.UploadFileName+"添加 "+task.FileName+" 连接到服务器 "+ip+" 成功！")
	// 执行命令
	result, err := server.RunCmd("/usr/local/sbin/baize pin add " + hash)
	beego.Debug("执行命令 : baize pin add ", hash, ", 结果: ", result)
	tmpErrorMsg := task.ErrorMsg
	tmpPubLog := task.PubLog
	if err != nil {
		cnum <- -1
		task.ErrorMsg = fmt.Sprintf("发布到 %s:%d ：%v\n", ip, port, err) + tmpErrorMsg
		task.PubStatus = -3
		TaskService.UpdateTask(task, "PubStatus", "ErrorMsg")
		ActionService.Add("publish", "admin", "publish", 1000, task.UploadFileName+"发布 "+task.FileName+" 服务器 "+ip+" 失败！")
		return
	}
	cnum <- 1
	task.PubLog = fmt.Sprintf("发布到 %s:%d ：%s\n", ip, port, result) + tmpPubLog
	ActionService.Add("publish", "admin", "publish", 1000, task.UploadFileName+"发布 "+task.FileName+" 服务器 "+ip+" 成功！")
	TaskService.UpdateTask(task, "PubLog")
	return
}
