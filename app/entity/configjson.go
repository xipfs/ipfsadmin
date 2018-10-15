package entity

/*
 ============================================================================
 Name        : configjson.go
 Author      : xiehui
 Version     : 1.0
 Email 	     : xiehui6@lenovo.com
 Copyright   : Copyright © 2018 Lenovo Services. All rights reserved.
 Description : JSON配置
 ============================================================================
*/

type ConfigJson struct {
	Version    string
	DataReport *DataReport
}

type DataReport struct {
	ReportUrls          []string
	BatchReportNum      int
	MaxCacheFileSize    int
	DisableEventActions []string
	DataGathers         []DataGathers
}

type DataGathers struct {
	EventAction string
	Param       *Param
}

type Param struct {
	period string
}
