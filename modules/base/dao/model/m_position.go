package model

import (
	"basic-frame/modules/base/dao/entity"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/mysql"
	"github.com/pkg/errors"
	"strings"
)

var PositionModel = &Position{}

type Position struct {
}

func (a *Position) Query(params schema.PositionQueryParam) (*schema.PositionQueryResult, error) {
	db := mysql.DB.Model(entity.Position{})
	if v := params.IDs; v != "" {
		db = db.Where("id IN (?)", strings.Split(v, ","))
	}
	if v := params.RoleID; v != 0 {
		db = db.Where("role_id=?", v)
	}
	if v := params.OrganizationID; v != 0 {
		db = db.Where("organization_id=?", v)
	}
	if v := params.Status; v != 0 {
		db = db.Where("status=?", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("name = ?", v)
	}
	if v := params.Memo; v != "" {
		v = "%" + strings.ToLower(v) + "%"
		db = db.Where("lower(memo) LIKE ?", v)
	}
	if v := strings.TrimSpace(params.QueryValue); v != "" {
		v = "%" + strings.ToLower(v) + "%"
		db = db.Where("lower(name) LIKE ? OR lower(memo) LIKE ?", v, v)
	}
	if v := params.SequenceSort; v == 1 || v == 2 {
		if v == 1 {
			db.Order("sequence ASC")
		} else {
			db.Order("sequence DESC")
		}
	} else {
		db.Order("id DESC")
	}

	var list entity.Positions
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.PositionQueryResult{
		Data:       list.ToSchemaPositions(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *Position) Get(id uint64) (*schema.Position, error) {
	db := mysql.DB.Model(entity.Position{ID: id})

	var item entity.Position
	if ok, err := mysql.FindOne(db, &item); err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaPosition(), nil
}

func (a *Position) Create(item schema.Position) (*common.IDResult, error) {
	eitem := entity.SchemaPosition(item).ToPosition()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

func (a *Position) UpdateByID(id uint64, item map[string]interface{}) error {
	result := mysql.DB.Model(entity.Position{ID: id}).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *Position) Delete(id uint64) error {
	result := mysql.DB.Model(entity.Position{ID: id}).Delete(&entity.Position{})
	return errors.WithStack(result.Error)
}
