package model

import (
	"basic-frame/modules/base/dao/entity"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/mysql"
	"github.com/pkg/errors"
	"strings"
)

var UserGroupModel = &UserGroup{}

type UserGroup struct {
}

func (a *UserGroup) Query(params schema.UserGroupQueryParam) (*schema.UserGroupQueryResult, error) {
	db := mysql.DB.Model(&entity.UserGroup{})
	if v := params.IDs; v != "" {
		db = db.Where("id IN (?)", strings.Split(v, ","))
	}
	if v := params.Name; v != "" {
		db = db.Where("name = ?", v)
	}
	if v := params.RoleID; v != 0 {
		db = db.Where("role_id=?", v)
	}
	if v := params.Status; v != 0 {
		db = db.Where("status=?", v)
	}
	if v := params.Memo; v != "" {
		v = "%" + strings.ToLower(v) + "%"
		db = db.Where("lower(memo) LIKE ?", v)
	}
	if v := strings.TrimSpace(params.QueryValue); v != "" {
		v = "%" + strings.ToLower(v) + "%"
		db = db.Where("lower(name) LIKE ? OR lower(memo) LIKE ?", v, v)
	}
	if params.FindAll {
		params.Pagination = false
	}
	db.Order("id DESC")

	var list entity.UserGroups
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.UserGroupQueryResult{
		Data:       list.ToSchemaUserGroups(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *UserGroup) Get(id uint64) (*schema.UserGroup, error) {
	db := mysql.DB.Model(&entity.UserGroup{}).Where("id = ?", id)

	var item entity.UserGroup
	if ok, err := mysql.FindOne(db, &item); err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaUserGroup(), nil
}

func (a *UserGroup) Create(item schema.UserGroup) (*common.IDResult, error) {
	eitem := entity.SchemaUserGroup(item).ToUserGroup()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

func (a *UserGroup) UpdateByID(id uint64, item map[string]interface{}) error {
	result := mysql.DB.Model(&entity.UserGroup{}).Where("id = ?", id).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *UserGroup) Delete(id uint64) error {
	result := mysql.DB.Model(&entity.UserGroup{}).Unscoped().Delete(&entity.UserGroup{}, id)
	return errors.WithStack(result.Error)
}
