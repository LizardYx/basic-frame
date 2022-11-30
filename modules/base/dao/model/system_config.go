package model

import (
	"basic-frame/modules/base/dao/entity"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/ginx"
	"basic-frame/util/ginx/errors"
	"basic-frame/util/mysql"
)

var SystemConfigModel = &SystemConfig{}

type SystemConfig struct {
}

func (a *SystemConfig) Query(params schema.SystemConfigQueryParam) (*schema.SystemConfigQueryResult, error) {
	db := mysql.DB.Model(entity.SystemConfig{})

	var list entity.SystemConfigs
	paginationResult, err := ginx.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, err
	}
	qr := &schema.SystemConfigQueryResult{
		Data:       list.ToSchemaSystemConfigPres(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *SystemConfig) Create(item schema.SystemConfig) (*common.IDResult, error) {
	eitem := entity.SchemaSystemConfig(item).ToSystemConfig()
	result := mysql.DB.Create(eitem)
	return &common.IDResult{ID: eitem.ID}, errors.New(result.Error.Error())
}
