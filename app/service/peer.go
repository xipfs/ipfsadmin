package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/xipfs/ipfsadmin/app/entity"
)

type peerService struct{}

// 表名
func (this *peerService) table() string {
	return tableName("peer")
}

// 获取节点信息
func (this *peerService) GetPeer(peerId string) (*entity.Peer, error) {
	peer := &entity.Peer{}
	query := o.QueryTable(this.table())
	err := query.Filter("PeerId", peerId).One(peer)
	return peer, err
}

// 获取所有节点
func (this *peerService) GetAllPeer() ([]entity.Peer, int64) {
	return this.GetList(1, -1)
}

// 获取节点列表
func (this *peerService) GetList(page, pageSize int, filters ...interface{}) ([]entity.Peer, int64) {
	var (
		list  []entity.Peer
		count int64
	)

	offset := 0
	if pageSize == -1 {
		pageSize = 100000
	} else {
		offset = (page - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
	}
	query := o.QueryTable(this.table())
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			_, ok := filters[k].(string)
			if !ok {
				continue
			}
			v := filters[k+1]
			query = query.Filter(filters[k].(string), v)
		}
	}
	count, _ = query.Count()
	query.Offset(offset).Limit(pageSize).All(&list)
	return list, count
}

// 获取节点总数
func (this *peerService) GetTotal() (int64, error) {
	return o.QueryTable(this.table()).Count()
}

// 添加节点
func (this *peerService) AddPeer(peer *entity.Peer, fields ...string) error {
	_, err := o.InsertOrUpdate(peer, fields...)
	return err
}

// 更新节点信息
func (this *peerService) UpdatePeer(peer *entity.Peer, fields ...string) error {
	_, err := o.Update(peer, fields...)
	return err
}

// 删除节点
func (this *peerService) DeletePeer(peerId int) error {

	return nil
}

// 节点统计
func (this *peerService) GetPeerStat(rangeType string) map[int]int {
	var sql string
	var maps []orm.Params

	switch rangeType {
	case "this_month":
		year, month, _ := time.Now().Date()
		//startTime := fmt.Sprintf("%d-%02d-01 00:00:00", year, month)
		endTime := fmt.Sprintf("%d-%02d-31 23:59:59", year, month)
		sql = fmt.Sprintf("SELECT DAY(create_time) AS date, sum(peer_id) AS count FROM %s WHERE create_time < '%s' GROUP BY DAY(create_time) ORDER BY `date` ASC", this.table(), endTime)
	case "last_month":
		year, month, _ := time.Now().AddDate(0, -1, 0).Date()
		//startTime := fmt.Sprintf("%d-%02d-01 00:00:00", year, month)
		endTime := fmt.Sprintf("%d-%02d-31 23:59:59", year, month)
		sql = fmt.Sprintf("SELECT DAY(create_time) AS date, sum(peer_id) AS count FROM %s WHERE create_time < '%s' GROUP BY DAY(create_time) ORDER BY `date` ASC", this.table(), endTime)
	case "this_year":
		year := time.Now().Year()
		//startTime := fmt.Sprintf("%d-01-01 00:00:00", year)
		endTime := fmt.Sprintf("%d-12-31 23:59:59", year)
		sql = fmt.Sprintf("SELECT MONTH(create_time) AS date, sum(peer_id) AS count FROM %s WHERE create_time < '%s' GROUP BY MONTH(create_time) ORDER BY `date` ASC", this.table(), endTime)
	case "last_year":
		year := time.Now().Year() - 1
		//startTime := fmt.Sprintf("%d-01-01 00:00:00", year)
		endTime := fmt.Sprintf("%d-12-31 23:59:59", year)
		sql = fmt.Sprintf("SELECT MONTH(create_time) AS date, sum(peer_id) AS count FROM %s WHERE create_time < '%s' GROUP BY MONTH(create_time) ORDER BY `date` ASC", this.table(), endTime)
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
