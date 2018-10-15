package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/xipfs/ipfsadmin/app/entity"
)

type peerLogService struct{}

// 表名
func (this *peerLogService) table() string {
	return tableName("peer_log")
}

// 添加节点日志
func (this *peerLogService) AddPeerLog(peerLog *entity.PeerLog) error {
	_, err := o.Insert(peerLog)
	return err
}

// 节点统计
func (this *peerLogService) GetPeerStat(rangeType string) map[int]int {
	var sql string
	var maps []orm.Params

	switch rangeType {
	case "this_month":
		year, month, _ := time.Now().Date()
		startTime := fmt.Sprintf("%d-%02d-01 00:00:00", year, month)
		endTime := fmt.Sprintf("%d-%02d-31 23:59:59", year, month)
		sql = fmt.Sprintf("SELECT DAY(create_time) AS date, sum(peer_count) AS count FROM %s WHERE create_time BETWEEN '%s' AND '%s' GROUP BY DAY(create_time) ORDER BY `date` ASC", this.table(), startTime, endTime)
	case "last_month":
		year, month, _ := time.Now().AddDate(0, -1, 0).Date()
		startTime := fmt.Sprintf("%d-%02d-01 00:00:00", year, month)
		endTime := fmt.Sprintf("%d-%02d-31 23:59:59", year, month)
		sql = fmt.Sprintf("SELECT DAY(create_time) AS date, sum(peer_count) AS count FROM %s WHERE create_time BETWEEN '%s' AND '%s' GROUP BY DAY(create_time) ORDER BY `date` ASC", this.table(), startTime, endTime)
	case "this_year":
		year := time.Now().Year()
		startTime := fmt.Sprintf("%d-01-01 00:00:00", year)
		endTime := fmt.Sprintf("%d-12-31 23:59:59", year)
		sql = fmt.Sprintf("SELECT MONTH(create_time) AS date, sum(peer_count) AS count FROM %s WHERE create_time BETWEEN '%s' AND '%s' GROUP BY MONTH(create_time) ORDER BY `date` ASC", this.table(), startTime, endTime)
	case "last_year":
		year := time.Now().Year() - 1
		startTime := fmt.Sprintf("%d-01-01 00:00:00", year)
		endTime := fmt.Sprintf("%d-12-31 23:59:59", year)
		sql = fmt.Sprintf("SELECT MONTH(create_time) AS date, sum(peer_count) AS count FROM %s WHERE create_time BETWEEN '%s' AND '%s' GROUP BY MONTH(create_time) ORDER BY `date` ASC", this.table(), startTime, endTime)
	}

	num, err := o.Raw(sql).Values(&maps)

	result := make(map[int]int)
	if err == nil && num > 0 {
		for _, v := range maps {
			date, _ := strconv.Atoi(v["date"].(string))
			count, _ := strconv.Atoi(v["count"].(string))
			result[date] = count
		}
	}
	return result
}
