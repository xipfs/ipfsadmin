package entity

/*
 ============================================================================
 Name        : Perm.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description :
 ============================================================================
*/

type Perm struct {
	Module string `orm:"size(20)"`
	Action string `orm:"size(20)"`
	Key    string `orm:"-"` // Module.Action
}

func (p *Perm) TableUnique() [][]string {
	return [][]string{
		[]string{"Module", "Action"},
	}
}
