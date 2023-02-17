package model

import (
	"basic-frame/modules/base/dao/entity"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

var ButtonModel = &Button{}

type Button struct {
}

func (a *Button) Query(params schema.ButtonQueryParam) (*schema.ButtonQueryResult, error) {
	db := mysql.DB.Model(&entity.Button{})
	if v := params.IDs; v != "" {
		db.Where("id IN (?)", strings.Split(v, ","))
	}
	if v := params.UUID; v != "" {
		db.Where("uuid = ?", v)
	}
	if v := params.ParentID; v != 0 {
		db.Where("parent_id = ?", v)
	}
	if v := params.ShowStatus; v != 0 {
		db.Where("show_status = ?", v)
	}
	if v := params.Status; v != 0 {
		db.Where("status = ?", v)
	}
	db.Order("id DESC")

	var list entity.Buttons
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.ButtonQueryResult{
		Data:       list.ToSchemaButtons(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *Button) Get(id uint64) (*schema.Button, error) {
	db := mysql.DB.Model(&entity.Button{}).Where("id = ?", id)

	var item entity.Button
	if ok, err := mysql.FindOne(db, &item); err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaButton(), nil
}

func (a *Button) GetButtonRestfulApis(id uint64) (*schema.ButtonPre, error) {
	db := mysql.DB.Model(&entity.Button{}).Where("id = ?", id).Preload("RestfulApis")

	var item entity.Button
	if ok, err := mysql.FindOne(db, &item); err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaButtonPre(), nil
}

// GetButtonsRestfulApis 获取多个按钮的所有RestfulApis
func (a *Button) GetButtonsRestfulApis(ids []uint64) (*schema.ButtonPres, error) {
	db := mysql.DB.Model(&entity.Button{}).Where("id IN (?)", ids).Preload("RestfulApis")

	var items entity.Buttons
	if err := db.Find(&items).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	buttonPres := items.ToSchemaButtonPres()

	return &buttonPres, nil
}

func (a *Button) Create(item schema.Button) (*common.IDResult, error) {
	eitem := entity.SchemaButton(item).ToButton()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

// CreateButtonPre 创建按钮和按钮关联的RestfulApi(不会创建子按钮)
func (a *Button) CreateButtonPre(item schema.ButtonPre) (*common.IDResult, error) {
	eitem := entity.SchemaButtonPre(item).ToButton()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

// UpdateButtonRestfulApis 更新指定按钮的RestfulApis
func (a *Button) UpdateButtonRestfulApis(id uint64, items schema.RestfulApis) error {
	eitem := entity.SchemaRestfulApis(items).ToRestfulApis()
	if err := mysql.DB.Model(&entity.Button{ID: id}).
		Association("RestfulApis").
		Replace(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateButtonPre 更新按钮和按钮关联的RestfulApis信息
func (a *Button) UpdateButtonPre(item schema.ButtonPre) error {
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 检查Button是否存在
		if oldItem, err := a.Get(item.ID); err != nil {
			return err
		} else if oldItem == nil {
			return errors.New("未找到该按钮")
		} else if item.ParentID != nil {
			parentID := uint64(0)
			if item.ParentID == &parentID {
				item.ParentID = nil
			}
		}

		// 更新Button
		if err := a.UpdateByID(item.ID, map[string]interface{}{
			"btn_id":    item.BtnID,
			"name":      item.Name,
			"icon":      item.Icon,
			"class":     item.Class,
			"menu_id":   item.MenuID,
			"sequence":  item.Sequence,
			"parent_id": item.ParentID,
			"status":    item.Status,
			"memo":      item.Memo,
		}); err != nil {
			return err
		}

		// 更新button关联的API
		if err := a.UpdateButtonRestfulApis(item.ID, item.RestfulApis); err != nil {
			return err
		}
		return nil
	})
}

func (a *Button) UpdateByID(id uint64, item map[string]interface{}) error {
	result := mysql.DB.Model(&entity.Button{}).Where("id = ?", id).Updates(item)
	return errors.WithStack(result.Error)
}

// Delete 删除按钮和按钮调用的Api关联
func (a *Button) Delete(id uint64) error {
	result := mysql.DB.Model(&entity.Button{}).
		Unscoped().
		Select(clause.Associations).
		Delete(&entity.Button{}, id)
	return errors.WithStack(result.Error)
}

// BatchDelete 批量删除按钮和按钮调用的Api关联
func (a *Button) BatchDelete(ids []uint64) error {
	result := mysql.DB.Model(&entity.Button{}).
		Unscoped().
		Select(clause.Associations).
		Delete(&entity.Button{}, "id IN (?)", ids)
	return errors.WithStack(result.Error)
}
