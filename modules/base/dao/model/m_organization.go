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

var OrganizationModel = &Organization{}

type Organization struct {
}

func (a *Organization) Query(params schema.OrganizationQueryParam) (*schema.OrganizationQueryResult, error) {
	db := mysql.DB.Model(&entity.Organization{})
	if v := params.ID; v != 0 {
		db.Where("id=?", v)
	}
	if v := params.IDs; v != "" {
		db.Where("id IN (?)", strings.Split(v, ","))
	}
	if v := params.RoleID; v != 0 {
		db.Where("role_id=?", v)
	}
	if v := params.ParentID; v != 0 {
		db.Where("parent_id=?", v)
	}
	if v := params.Status; v != 0 {
		db.Where("status=?", v)
	}
	if params.ShowPositions {
		db = db.Preload("Positions")
	}
	if params.ShowDetails {
		db = db.Preload(clause.Associations).Preload("SonOrganizations", PreloadOrganizationAllForCreateUser)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + strings.ToLower(v) + "%"
		db.Where("lower(name) LIKE ? OR lower(memo) LIKE ?", v, v)
	}
	db.Order("id DESC")

	var list entity.Organizations
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.OrganizationQueryResult{
		Data:       list.ToSchemaOrganizations(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *Organization) Get(id uint64) (*schema.Organization, error) {
	db := mysql.DB.Model(&entity.Organization{}).Where("id = ?", id)

	var item entity.Organization
	if ok, err := mysql.FindOne(db, &item); err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaOrganization(), nil
}

func (a *Organization) Create(item schema.Organization) (*common.IDResult, error) {
	eitem := entity.SchemaOrganization(item).ToOrganization()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

func (a *Organization) CreateOrganizations(item schema.Organizations) error {
	eitem := entity.SchemaOrganizations(item).ToOrganization()
	result := mysql.DB.Create(&eitem)
	return errors.WithStack(result.Error)
}

func (a *Organization) UpdateByID(id uint64, item map[string]interface{}) error {
	result := mysql.DB.Model(&entity.Organization{}).Where("id = ?", id).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *Organization) Delete(id uint64) error {
	result := mysql.DB.Model(&entity.Organization{}).Delete(&entity.Organization{}, id)
	return errors.WithStack(result.Error)
}

// ----------------------------------------OrganizationTree-----------------------

// GetOrganizationTree 获取组织树(包含禁用的组织)
func (a *Organization) GetOrganizationTree() (*schema.Organizations, error) {
	db := mysql.DB.Model(&entity.Organization{}).Order("id DESC")
	db = db.
		Where("parent_id IS NULL").
		Preload(clause.Associations).
		Preload("SonOrganizations", PreloadOrganizationAll)

	var list entity.Organizations
	if err := db.Find(&list).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	organizationTrees := list.ToSchemaOrganizations()
	return &organizationTrees, nil
}

// GetOrganizationTreeForCreateUser 获取组织树(不包含禁用的组织)
func (a *Organization) GetOrganizationTreeForCreateUser() (*schema.Organizations, error) {
	db := mysql.DB.Model(&entity.Organization{}).Order("id DESC")
	db = db.
		Where("parent_id IS NULL AND status = 1").
		Preload(clause.Associations).
		Preload("SonOrganizations", PreloadOrganizationAllForCreateUser)

	var list entity.Organizations
	if err := db.Find(&list).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	organizationTrees := list.ToSchemaOrganizations()
	return &organizationTrees, nil
}

// GetOrgTreeForCreateNotifications 获取组织树(不包含职位和禁用的组织)
func (a *Organization) GetOrgTreeForCreateNotifications() (*schema.Organizations, error) {
	db := mysql.DB.Model(&entity.Organization{}).Order("id DESC")

	db = db.
		Where("parent_id IS NULL AND status = 1").
		Preload("SonOrganizations", PreloadOrganizationAllForCreateNotifications)
	var list entity.Organizations
	if err := db.Find(&list).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	organizationTrees := list.ToSchemaOrganizations()
	return &organizationTrees, nil
}

// ----------------------------- Associations Action -----------------------------

func PreloadOrganizationAll(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations).
		Preload("SonOrganizations", PreloadOrganizationAll)
}

func PreloadOrganizationAllForCreateUser(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations).
		Preload("SonOrganizations", "status = 1", PreloadOrganizationAllForCreateUser)
}

func PreloadOrganizationAllForCreateNotifications(db *gorm.DB) *gorm.DB {
	return db.Preload("SonOrganizations", "status = 1", PreloadOrganizationAllForCreateNotifications)
}
