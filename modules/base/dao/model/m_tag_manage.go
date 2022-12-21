package model

import (
	"basic-frame/modules/base/dao/entity"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/ginx/errors"
	"basic-frame/util/mysql"
	"strings"
)

var TagManageModel = &TagManage{}

type TagManage struct {
}

func (a *TagManage) Query(params schema.TagManageQueryParam) (*schema.TagManageQueryResult, error) {
	db := mysql.DB.Model(entity.TagManage{})
	if v := params.Type; v != "" {
		db = db.Where("type=?", v)
	}
	if v := params.KeyName; v != "" {
		db = db.Where("key_name=?", v)
	}
	if v := params.Value; v != "" {
		db = db.Where("value=?", v)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + strings.ToLower(v) + "%"
		db = db.Where("lower(type) LIKE ? OR lower(key_name) LIKE ? OR lower(value) LIKE ?", v, v, v)
	}
	db.Order("id DESC")

	var list entity.TagManages
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.TagManageQueryResult{
		Data:       list.ToSchemaTagManages(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *TagManage) Get(id uint64) (*schema.TagManage, error) {
	db := mysql.DB.Model(entity.TagManage{ID: id})

	var item entity.TagManage
	if ok, err := mysql.FindOne(db, &item); err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaTagManage(), nil
}

func (a *TagManage) Create(item schema.TagManage) (*common.IDResult, error) {
	eitem := entity.SchemaTagManage(item).ToTagManage()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

func (a *TagManage) UpdateByID(id uint64, item map[string]interface{}) error {
	result := mysql.DB.Model(entity.TagManage{ID: id}).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *TagManage) Delete(id uint64) error {
	result := mysql.DB.Model(entity.TagManage{ID: id}).Delete(&entity.TagManage{})
	return errors.WithStack(result.Error)
}
