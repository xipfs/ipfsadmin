package service

import (
	"github.com/xipfs/ipfsadmin/app/entity"
)

type configService struct{}

// 表名
func (this *configService) table() string {
	return tableName("config")
}

// 添加配置
func (this *configService) AddConfig(config *entity.Config) error {
	_, err := o.Insert(config)
	return err
}

// 添加配置
func (this *configService) SaveConfig(config *entity.Config) error {
	_, err := o.Update(config)
	return err
}

// 删除配置
func (this *configService) DelConfig(id int) error {
	_, err := o.QueryTable(this.table()).Filter("id", id).Delete()
	return err
}

// 获取所有配置
func (this *configService) GetAllConfig() ([]entity.Config, error) {
	return this.GetList(1, -1)
}

// 获取资源列表
func (this *configService) GetList(page, pageSize int) ([]entity.Config, error) {
	var list []entity.Config
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

// 获取配置信息
func (this *configService) GetConfig(id int) (*entity.Config, error) {
	var err error
	config := &entity.Config{}
	config.Id = id
	err = o.Read(config)
	return config, err
}
