package service

import (
	"os"

	"github.com/xipfs/ipfsadmin/app/entity"
)

type resourceService struct{}

// 表名
func (this *resourceService) table() string {
	return tableName("resource")
}

// 获取一个资源信息
func (this *resourceService) GetResource(id int) (*entity.Resource, error) {
	resource := &entity.Resource{}
	resource.Id = id
	if err := o.Read(resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// 获取所有资源
func (this *resourceService) GetAllResource() ([]entity.Resource, error) {
	return this.GetList(1, -1)
}

// 根据标志获取资源
func (this *resourceService) GetAllResourceByName(uploadFileName string) ([]entity.Resource, error) {
	var list []entity.Resource
	offset := 0
	pageSize := 100000
	_, err := o.QueryTable(this.table()).OrderBy("-create_time").Filter("UploadFileName", uploadFileName).Offset(offset).Limit(pageSize).All(&list)
	return list, err
}

// 获取资源列表
func (this *resourceService) GetList(page, pageSize int) ([]entity.Resource, error) {
	var list []entity.Resource
	offset := 0
	if pageSize == -1 {
		pageSize = 100000
	} else {
		offset = (page - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
	}

	_, err := o.QueryTable(this.table()).OrderBy("-create_time").Offset(offset).Limit(pageSize).All(&list)
	return list, err
}

// 获取资源总数
func (this *resourceService) GetTotal() (int64, error) {
	return o.QueryTable(this.table()).Count()
}

// 添加资源
func (this *resourceService) AddResource(resource *entity.Resource) error {
	offset := 0
	pageSize := 100000
	var list []entity.Resource

	_, err := o.QueryTable(this.table()).Filter("domain", resource.Domain).Filter("upload_file_name", resource.UploadFileName).OrderBy("-create_time").Offset(offset).Limit(pageSize).All(&list)
	if len(list) > 0 {
		resource = &list[0]
		return err
	} else {
		_, err := o.Insert(resource)
		return err
	}

}

// 更新资源信息
func (this *resourceService) UpdateResource(resource *entity.Resource, fields ...string) error {
	_, err := o.Update(resource, fields...)
	return err
}

// 删除一个资源
func (this *resourceService) DeleteResource(resourceId int) error {
	resource, err := this.GetResource(resourceId)
	if err != nil {
		return err
	}
	// 删除目录
	path := GetResourcePath(resource.Domain)
	os.RemoveAll(path)

	// 删除任务
	TaskService.DeleteByResourceId(resource.Id)
	// 删除资源
	o.Delete(resource)
	return nil
}
