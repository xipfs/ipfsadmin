package service

import (
	"github.com/xipfs/ipfsadmin/app/entity"
)

type apiService struct{}

// 表名
func (this *apiService) table() string {
	return tableName("api")
}

// 获取所有 api
func (this *apiService) GetAllApi() ([]entity.Api, error) {
	return this.GetList(1, -1)
}

// 获取 api 列表
func (this *apiService) GetList(page, pageSize int) ([]entity.Api, error) {
	var list []entity.Api
	offset := 0
	if pageSize == -1 {
		pageSize = 100000
	} else {
		offset = (page - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
	}

	_, err := o.QueryTable(this.table()).Offset(offset).Limit(pageSize).All(&list)
	return list, err
}

// 获取 api 信息
func (this *apiService) GetApi(id int) (*entity.Api, error) {
	var err error
	api := &entity.Api{}
	api.Id = id
	err = o.Read(api)
	return api, err
}

// 获取API总数
func (this *apiService) GetTotal() (int64, error) {
	return o.QueryTable(this.table()).Count()
}
