package model

import (
	"basic-frame/modules/base/dao/entity"
	"basic-frame/modules/base/schema"
	"basic-frame/util/mysql"
)

var SystemConfigModel = &SystemConfig{}

type SystemConfig struct {
}

func (a *SystemConfig) Query(params schema.SystemConfigQueryParam) (*schema.SystemConfigQueryResult, error) {
	db := mysql.DB

	var list entity.SystemConfigs
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, err
	}
	qr := &schema.SystemConfigQueryResult{
		Data:       list.ToSchemaSystemConfigPres(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *SystemConfig) Create(item schema.SystemConfig) (uint64, error) {
	eitem := entity.SchemaSystemConfig(item).ToSystemConfig()
	result := mysql.DB.Create(eitem)
	return item.ID, result.Error
}
