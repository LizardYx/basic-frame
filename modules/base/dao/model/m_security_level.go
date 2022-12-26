package model

import (
	"basic-frame/modules/base/dao/entity"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/ginx/errors"
	"basic-frame/util/mysql"
	"strings"
)

var SecurityLevelModel = &SecurityLevel{}

type SecurityLevel struct {
}

func (a *SecurityLevel) Query(params schema.SecurityLevelQueryParam) (*schema.SecurityLevelQueryResult, error) {
	db := mysql.DB.Model(entity.SecurityLevel{})
	if v := params.ID; v != 0 {
		db = db.Where("id=?", v)
	}
	if v := params.IDs; v != "" {
		db = db.Where("id IN (?)", strings.Split(v, ","))
	}
	if v := params.Status; v != 0 {
		db = db.Where("status=?", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("lower(name) LIKE ?", "%"+strings.ToLower(v)+"%")
	}
	if v := params.RoleIDs; v != "" {
		db.Where("id IN (?)", mysql.DB.Table("security_level_role").
			Select("security_level_id").
			Where("role_id IN (?)", strings.Split(v, ",")))
	}
	if v := params.QueryValue; v != "" {
		v = "%" + strings.ToLower(v) + "%"
		db.Where("lower(name) LIKE ? OR lower(memo) LIKE ?", v, v)
	}
	db.Order("id DESC")

	var list entity.SecurityLevels
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	qr := &schema.SecurityLevelQueryResult{
		Data:       list.ToSchemaSecurityLevels(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *SecurityLevel) Get(id uint64) (*schema.SecurityLevel, error) {
	db := mysql.DB.Model(entity.SecurityLevel{ID: id})

	var item entity.SecurityLevel
	if ok, err := mysql.FindOne(db, &item); err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaSecurityLevel(), nil
}

func (a *SecurityLevel) Create(item schema.SecurityLevel) (*common.IDResult, error) {
	eitem := entity.SchemaSecurityLevel(item).ToSecurityLevel()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

func (a *SecurityLevel) UpdateByID(id uint64, item map[string]interface{}) error {
	result := mysql.DB.Model(entity.SecurityLevel{ID: id}).Updates(item)
	return errors.WithStack(result.Error)
}

// ReplaceSecurityLevelRoles 更新安全等级关联的角色
func (a *SecurityLevel) ReplaceSecurityLevelRoles(id uint64, items schema.Roles) error {
	eitem := entity.SchemaRoles(items).ToRole()
	if err := mysql.DB.Model(entity.SecurityLevel{ID: id}).
		Association("Roles").
		Replace(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *SecurityLevel) Delete(id uint64) error {
	result := mysql.DB.Model(entity.SecurityLevel{ID: id}).Delete(&entity.SecurityLevel{})
	return errors.WithStack(result.Error)
}

// SecurityLevelUsed 检查安全等级是否被使用
func (a *SecurityLevel) SecurityLevelUsed(id uint64) error {
	var count int64
	// 检查安全等级是否被任务使用
	if err := mysql.DB.Table("tk_safe_level").Where("level = ?", id).Count(&count).Error; err != nil {
		return errors.WithStack(err)
	} else if count != 0 {
		return errors.New("安全等级已被任务使用")
	}

	// 检查安全等级是否被项目使用
	if err := mysql.DB.Table("project ").Where("sec_level = ?", id).Count(&count).Error; err != nil {
		return errors.WithStack(err)
	} else if count != 0 {
		return errors.New("安全等级已被项目使用")
	}

	// 检查安全等级是否被项目集使用
	if err := mysql.DB.Table("project_group").Where("sec_level = ?", id).Count(&count).Error; err != nil {
		return errors.WithStack(err)
	} else if count != 0 {
		return errors.New("安全等级已被项目集使用")
	}
	return nil
}
