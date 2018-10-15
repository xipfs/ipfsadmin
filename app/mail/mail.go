package mail

/*
 ============================================================================
 Name        : mail.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 邮箱
 ============================================================================
*/

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/go-gomail/gomail"
	"github.com/xipfs/ipfsadmin/app/libs"
)

var (
	host     string
	port     int
	username string
	password string
	from     string
)

func init() {
	host = beego.AppConfig.String("mail.host")
	port, _ = beego.AppConfig.Int("mail.port")
	username = beego.AppConfig.String("mail.user")
	password = beego.AppConfig.String("mail.password")
	from = beego.AppConfig.String("mail.from")
	if port == 0 {
		port = 25
	}
}

// SendMail 发送邮件
func SendMail(subject, content string, to, cc []string) error {
	toList := make([]string, 0, len(to))
	ccList := make([]string, 0, len(cc))

	for _, v := range to {
		v = strings.TrimSpace(v)
		if libs.IsEmail([]byte(v)) {
			exists := false
			for _, vv := range toList {
				if v == vv {
					exists = true
					break
				}
			}
			if !exists {
				toList = append(toList, v)
			}
		}
	}
	for _, v := range cc {
		v = strings.TrimSpace(v)
		if libs.IsEmail([]byte(v)) {
			exists := false
			for _, vv := range ccList {
				if v == vv {
					exists = true
					break
				}
			}
			if !exists {
				ccList = append(ccList, v)
			}
		}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", toList...)
	if len(ccList) > 0 {
		m.SetHeader("Cc", ccList...)
	}
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	d := gomail.NewPlainDialer(host, port, username, password)

	return d.DialAndSend(m)
}
