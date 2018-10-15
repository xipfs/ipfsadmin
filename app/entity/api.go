package entity

/*
 ============================================================================
 Name        : api.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : api
 ============================================================================
*/

// API
type Api struct {
	Id   int
	Name string `orm:"size(32)"`  // ID
	Url  string `orm:"size(255)"` // 地址
	Desc string `orm:"size(255)"` // 描述
}
