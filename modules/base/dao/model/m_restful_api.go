package model

import (
	"basic-frame/modules/base/dao/entity"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/mysql"
	"github.com/pkg/errors"
)

var RestfulApiModel = &RestfulApi{}

type RestfulApi struct {
}

func (a *RestfulApi) Query(params schema.RestfulApiQueryParam) (*schema.RestfulApiQueryResult, error) {
	db := mysql.DB.Model(entity.RestfulApi{})
	if v := params.UUID; v != "" {
		db.Where("uuid=?", v)
	}
	db.Order("id DESC")

	var list entity.RestfulApis
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RestfulApiQueryResult{
		Data:       list.ToSchemaRestfulApis(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *RestfulApi) GetByUUID(uuid string) (*schema.RestfulApi, error) {
	db := mysql.DB.Model(entity.RestfulApi{UUID: uuid})

	var item entity.RestfulApi
	ok, err := mysql.FindOne(db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, errors.New("未找到UUID匹配的RestfulApi")
	}
	return item.ToSchemaRestfulApi(), nil
}

func (a *RestfulApi) Create(item schema.RestfulApi) (*common.IDResult, error) {
	eitem := entity.SchemaRestfulApi(item).ToRestfulApi()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

func (a *RestfulApi) UpdateByID(id uint64, item map[string]interface{}) error {
	result := mysql.DB.Model(entity.RestfulApi{ID: id}).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *RestfulApi) UpdateByUUID(uuid string, item map[string]interface{}) error {
	result := mysql.DB.Model(entity.RestfulApi{UUID: uuid}).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *RestfulApi) Delete(id uint64) error {
	result := mysql.DB.Model(entity.RestfulApi{ID: id}).Delete(&entity.RestfulApi{})
	return errors.WithStack(result.Error)
}
