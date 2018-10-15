package controllers

/*
 ============================================================================
 Name        : types.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : 菜单
 ============================================================================
*/

type Menu struct {
	Name    string
	Route   string
	Icon    string
	Submenu []SubMenu
}

type SubMenu struct {
	Name   string
	Route  string
	Action string
}
