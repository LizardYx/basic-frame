package model

import (
	"basic-frame/modules/base/dao/entity"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/mysql"
	"github.com/pkg/errors"
)

var DisabledFieldModel = &DisabledField{}

type DisabledField struct {
}

func (a *DisabledField) Query(params schema.DisabledFieldQueryParam) (*schema.DisabledFieldQueryResult, error) {
	db := mysql.DB.Model(entity.DisabledField{})
	if v := params.UUID; v != "" {
		db.Where("uuid = ?", v)
	}

	var list entity.DisabledFields
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	qr := &schema.DisabledFieldQueryResult{
		Data:       list.ToSchemaDisabledFields(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *DisabledField) Get(id uint64) (*schema.DisabledField, error) {
	db := mysql.DB.Model(entity.DisabledField{ID: id})

	var item entity.DisabledField
	if ok, err := mysql.FindOne(db, &item); err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaDisabledField(), nil
}

func (a *DisabledField) Create(item schema.DisabledField) (*common.IDResult, error) {
	eitem := entity.SchemaDisabledField(item).ToDisabledField()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

func (a *DisabledField) UpdateByID(id uint64, item map[string]interface{}) error {
	result := mysql.DB.Model(entity.DisabledField{ID: id}).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *DisabledField) Delete(id uint64) error {
	result := mysql.DB.Model(entity.DisabledField{ID: id}).Delete(&entity.DisabledField{})
	return errors.WithStack(result.Error)
}
