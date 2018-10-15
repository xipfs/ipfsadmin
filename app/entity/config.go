package entity

/*
 ============================================================================
 Name        : config.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 配置
 ============================================================================
*/

type Config struct {
	Id    int
	Key   string `orm:"size(255)"`
	Value string `orm:"size(255)"`
}
