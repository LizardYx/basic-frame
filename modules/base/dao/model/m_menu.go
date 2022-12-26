package model

import (
	"basic-frame/modules/base/dao/entity"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"basic-frame/util/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

var MenuModel = &Menu{}

type Menu struct {
}

func (a *Menu) Query(params schema.MenuQueryParam) (*schema.MenuQueryResult, error) {
	db := mysql.DB.Model(entity.Menu{})
	if v := params.ID; v != 0 {
		db.Where("id=?", v)
	}
	if v := params.IDs; v != "" {
		db.Where("id IN (?)", strings.Split(v, ","))
	}
	if v := params.UUID; v != "" {
		db.Where("uuid=?", v)
	}
	if v := params.ParentID; v != 0 {
		db.Where("parent_id=?", v)
	}
	if v := params.Status; v != 0 {
		db.Where("status=?", v)
	}
	if v := params.ShowStatus; v != 0 {
		db.Where("show_status=?", v)
	}
	db.Order("id DESC")

	var list entity.Menus
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.MenuQueryResult{
		Data:       list.ToSchemaMenus(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *Menu) Get(id uint64) (*schema.Menu, error) {
	db := mysql.DB.Model(entity.Menu{ID: id})

	var item entity.Menu
	if ok, err := mysql.FindOne(db, &item); err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaMenu(), nil
}

// GetRoleRestfulApis 获取角色的RestfulApis信息
func (a *Menu) GetRoleRestfulApis(ids []uint64) (*schema.MenuTrees, error) {
	db := mysql.DB.Model(entity.Menu{}).Where("id IN (?)", ids).Preload("RestfulApis")

	var items entity.Menus
	if err := db.Find(&items).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	menuTrees := items.ToSchemaMenuTrees()
	return &menuTrees, nil
}

func (a *Menu) Create(item schema.Menu) (*common.IDResult, error) {
	eitem := entity.SchemaMenu(item).ToMenu()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

func (a *Menu) UpdateByID(id uint64, item map[string]interface{}) error {
	result := mysql.DB.Model(entity.Menu{ID: id}).Updates(item)
	return errors.WithStack(result.Error)
}

// Delete 删除菜单、菜单调用的Api关联以及菜单的按钮
func (a *Menu) Delete(id uint64) error {
	result := mysql.DB.Model(entity.Menu{ID: id}).Unscoped().Delete(&entity.Menu{})
	return errors.WithStack(result.Error)
}

// ----------------------------------------MenuTrees--------------------------------------

// QueryMenuTree 获取所有的菜单树(包含禁用)
func (a *Menu) QueryMenuTree() (*schema.MenuTrees, error) {
	db := mysql.DB.Model(entity.Menu{}).Order("sequence DESC")
	db.Where("parent_id IS NULL").
		Preload(clause.Associations).
		Preload("SonMenus", PreloadMenuAll).
		Preload("Buttons", "parent_id IS NULL", PreloadAll)

	var list entity.Menus
	if err := db.Find(&list).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	MenuTrees := list.ToSchemaMenuTrees()
	return &MenuTrees, nil
}

// QueryMenuTreeForCreateRole 获取创建用户的菜单树(不包含禁用)
func (a *Menu) QueryMenuTreeForCreateRole() (*schema.MenuTrees, error) {
	db := mysql.DB.Model(entity.Menu{}).Order("sequence DESC")
	db.Where("parent_id IS NULL AND status = ?", consts.BaseStatusEnable).
		Preload(clause.Associations).
		Preload("SonMenus", PreloadMenuAllForCreateRole).
		Preload("Buttons", "parent_id IS NULL AND status = 1", PreloadAllForCreateRole)

	var list entity.Menus
	if err := db.Find(&list).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	MenuTrees := list.ToSchemaMenuTrees()
	return &MenuTrees, nil
}

func (a *Menu) CreateMenuTrees(items schema.MenuTrees) error {
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		db := mysql.DB.Model(entity.Menu{})
		for _, item := range items {
			eitem := entity.SchemaMenuTree(*item).ToMenu()
			if err := db.Create(eitem).Error; err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})
}

func (a *Menu) CreateMenuTree(item schema.MenuTree) (*common.IDResult, error) {
	eitem := entity.SchemaMenuTree(item).ToMenu()
	result := mysql.DB.Create(eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

// UpdateMenuRestfulApis 更新菜单关联的Api
func (a *Menu) UpdateMenuRestfulApis(id uint64, items schema.RestfulApis) error {
	eitem := entity.SchemaRestfulApis(items).ToRestfulApis()
	if err := mysql.DB.Model(&entity.Menu{ID: id}).
		Association("RestfulApis").
		Replace(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ----------------------------- Associations Action -----------------------------

func PreloadAll(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations, PreloadAll)
}

func PreloadAllForCreateRole(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations).
		Preload("SonButtons", "status = 1", PreloadAllForCreateRole)
}

func PreloadMenuAll(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations).
		Preload("SonMenus", PreloadMenuAll).
		Preload("Buttons", "parent_id IS NULL", PreloadAll)
}

func PreloadMenuAllForCreateRole(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", consts.BaseStatusEnable).
		Preload(clause.Associations).
		Preload("SonMenus", PreloadMenuAllForCreateRole).
		Preload("Buttons", "parent_id IS NULL AND status = 1", PreloadAll)
}
