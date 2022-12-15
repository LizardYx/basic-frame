package model

import (
	"basic-frame/modules/base/dao/entity"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
	"strings"
)

var RoleModel = &Role{}

type Role struct {
}

func (a *Role) Query(params schema.RoleQueryParam) (*schema.RoleQueryResult, error) {
	db := mysql.DB.Model(entity.Role{})
	if v := params.ID; v != 0 {
		db = db.Where("id=?", v)
	}
	if v := params.IDs; len(v) != 0 {
		db = db.Where("id IN (?)", v)
	}
	if v := params.Type; v != 0 {
		db = db.Where("type=?", v)
	}
	if v := params.Types; len(v) != 0 {
		db = db.Where("type IN (?)", v)
	}
	if v := params.AuditorTypes; v != "" {
		auditorTypeList := strings.Split(v, ",")
		for _, auditorType := range auditorTypeList {
			db = db.Where("lower(auditor_types) LIKE ?", "%"+strings.ToLower(auditorType)+"%")
		}
	}
	if v := params.Status; v != 0 {
		db = db.Where("status=?", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("name = ?", v)
	}
	if v := params.Memo; v != "" {
		db = db.Where("lower(memo) LIKE ?", "%"+strings.ToLower(v)+"%")
	}
	if v := params.ShowDetails; v == true {
		db = db.Preload(clause.Associations)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + strings.ToLower(v) + "%"
		db.Where("lower(name) LIKE ? OR lower(name) LIKE ?", v, v)
	}
	db.Order("id DESC")

	var list entity.Roles
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	qr := &schema.RoleQueryResult{
		Data:       list.ToSchemaRoles(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *Role) Get(id uint64) (*schema.Role, error) {
	db := mysql.DB.Model(entity.Role{ID: id}).Preload(clause.Associations)

	var item entity.Role
	if ok, err := mysql.FindOne(db, &item); err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaRole(), nil
}

func (a *Role) Create(item schema.Role) (*common.IDResult, error) {
	eitem := entity.SchemaRole(item).ToRole()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

func (a *Role) UpdateByID(id uint64, item map[string]interface{}) error {
	result := mysql.DB.Model(entity.Role{ID: id}).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *Role) UpdateRoleMenu(id uint64, items schema.Menus) error {
	eitem := entity.SchemaMenus(items).ToMenu()
	if err := mysql.DB.Model(entity.Role{ID: id}).
		Association("Menus").
		Replace(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *Role) UpdateRoleButton(id uint64, items schema.Buttons) error {
	eitem := entity.SchemaButtons(items).ToButton()
	if err := mysql.DB.Model(entity.Role{ID: id}).
		Association("Buttons").
		Replace(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *Role) UpdateRoleDisabledFields(id uint64, items schema.DisabledFields) error {
	eitem := entity.SchemaDisabledFields(items).ToDisabledField()
	if err := mysql.DB.Model(entity.Role{ID: id}).
		Association("DisabledFields").
		Replace(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *Role) Delete(id uint64) error {
	result := mysql.DB.Model(entity.Role{ID: id}).Unscoped().Delete(&entity.Role{})
	return errors.WithStack(result.Error)
}
