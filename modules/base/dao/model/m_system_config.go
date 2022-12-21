package model

import (
	"basic-frame/modules/base/dao/entity"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/ginx/errors"
	"basic-frame/util/mysql"
)

var SystemConfigModel = &SystemConfig{}

type SystemConfig struct {
}

func (a *SystemConfig) Query(params schema.SystemConfigQueryParam) (*schema.SystemConfigQueryResult, error) {
	db := mysql.DB.Model(entity.SystemConfig{})
	db.Order("id DESC")

	var list entity.SystemConfigs
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.SystemConfigQueryResult{
		Data:       list.ToSchemaSystemConfigPres(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *SystemConfig) First() (*schema.SystemConfig, error) {
	var item entity.SystemConfig
	if err := mysql.DB.First(&item).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return item.ToSchemaSystemConfig(), nil
}

func (a *SystemConfig) Create(item schema.SystemConfig) (*common.IDResult, error) {
	eitem := entity.SchemaSystemConfig(item).ToSystemConfig()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

func (a *SystemConfig) Get(id uint64) (*schema.SystemConfigPre, error) {
	db := mysql.DB.Model(entity.SystemConfig{})

	db.Where("id = ?", id)
	var item entity.SystemConfig
	ok, err := mysql.FindOne(db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, errors.New("未找到ID匹配的系统配置项")
	}
	return item.ToSchemaSystemConfigPre(), nil
}

func (a *SystemConfig) UpdateByID(id uint64, item map[string]interface{}) error {
	result := mysql.DB.Model(entity.SystemConfig{ID: id}).Updates(item)
	return errors.WithStack(result.Error)
}
