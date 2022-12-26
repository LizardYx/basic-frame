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

var UserModel = &User{}

type User struct {
}

func (a *User) Query(params schema.UserQueryParam) (*schema.UserQueryResult, error) {
	db := mysql.DB.Model(entity.User{})
	if v := params.ID; v != 0 {
		db = db.Where("id=?", v)
	}
	if v := params.IDs; v != "" {
		db = db.Where("id IN (?)", strings.Split(v, ","))
	}
	if v := params.UserName; v != "" {
		db = db.Where("user_name = ?", v)
	}
	if v := params.Status; v != 0 {
		db = db.Where("status=?", v)
	}
	if v := params.OrgID; v != 0 {
		db.Where("id IN (?)", mysql.DB.Table("user_organization").
			Select("user_id").
			Where("organization_id = ?", v))
	}
	if v := params.OrgIDs; v != "" {
		db.Where("id IN (?)", mysql.DB.Table("user_organization").
			Select("user_id").
			Where("organization_id IN (?)", strings.Split(v, ",")))
	}
	if v := params.PositionID; v != 0 {
		db.Where("id IN (?)", mysql.DB.Table("user_position").
			Select("user_id").
			Where("position_id = ?", v))
	}
	if v := params.PositionIDs; v != "" {
		db.Where("id IN (?)", mysql.DB.Table("user_position").
			Select("user_id").
			Where("position_id IN (?)", strings.Split(v, ",")))
	}
	if v := params.RoleID; v != 0 {
		db.Where("id IN (?)", mysql.DB.Table("user_role").
			Select("user_id").
			Where("role_id = ?", v))
	}
	if v := params.RoleIDs; v != "" {
		db.Where("id IN (?)", mysql.DB.Table("user_role").
			Select("user_id").
			Where("role_id IN (?)", strings.Split(v, ",")))
	}
	if v := params.UserGroupID; v != 0 {
		db.Where("id IN (?)", mysql.DB.Table("user_user_group").
			Select("user_id").
			Where("user_group_id = ?", v))
	}
	if v := params.UserGroupIDs; v != "" {
		db.Where("id IN (?)", mysql.DB.Table("user_user_group").
			Select("user_id").
			Where("user_group_id IN (?)", strings.Split(v, ",")))
	}
	if v := params.ShowExtendInfo; v {
		db.Preload("ExtendInfo")
	}
	if v := params.ShowDetails; v {
		db.Preload(clause.Associations)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + strings.ToLower(v) + "%"
		db.Where("lower(user_name) LIKE ? OR id IN (?)", v, mysql.DB.Table("user_extend_info").
			Select("user_id").
			Where("lower(real_name) LIKE ? OR lower(mobile_phone) LIKE ? OR lower(qq_account) LIKE ? OR lower(email) LIKE ?", v, v, v, v))
	}
	if v := params.OmitPassword; v {
		db.Omit("password")
	}
	if v := params.FindAll; v {
		params.PaginationParam.Pagination = false
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

	var list entity.Users
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.UserQueryResult{
		Data:       list.ToSchemaUsers(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *User) QueryIncludeDelete(params schema.UserQueryParam) (*schema.UserQueryResult, error) {
	db := mysql.DB.Model(entity.User{}).Unscoped()
	if v := params.ID; v != 0 {
		db = db.Where("id=?", v)
	}
	if v := params.IDs; v != "" {
		db = db.Where("id IN (?)", strings.Split(v, ","))
	}
	if v := params.UserName; v != "" {
		db = db.Where("user_name = ?", v)
	}
	if v := params.Status; v != 0 {
		db = db.Where("status=?", v)
	}
	if v := params.OrgID; v != 0 {
		db.Where("id IN (?)", mysql.DB.Table("user_organization").
			Select("user_id").
			Where("organization_id = ?", v))
	}
	if v := params.OrgIDs; v != "" {
		db.Where("id IN (?)", mysql.DB.Table("user_organization").
			Select("user_id").
			Where("organization_id IN (?)", strings.Split(v, ",")))
	}
	if v := params.PositionID; v != 0 {
		db.Where("id IN (?)", mysql.DB.Table("user_position").
			Select("user_id").
			Where("position_id = ?", v))
	}
	if v := params.PositionIDs; v != "" {
		db.Where("id IN (?)", mysql.DB.Table("user_position").
			Select("user_id").
			Where("position_id IN (?)", strings.Split(v, ",")))
	}
	if v := params.RoleID; v != 0 {
		db.Where("id IN (?)", mysql.DB.Table("user_role").
			Select("user_id").
			Where("role_id = ?", v))
	}
	if v := params.RoleIDs; v != "" {
		db.Where("id IN (?)", mysql.DB.Table("user_role").
			Select("user_id").
			Where("role_id IN (?)", strings.Split(v, ",")))
	}
	if v := params.UserGroupID; v != 0 {
		db.Where("id IN (?)", mysql.DB.Table("user_user_group").
			Select("user_id").
			Where("user_group_id = ?", v))
	}
	if v := params.UserGroupIDs; v != "" {
		db.Where("id IN (?)", mysql.DB.Table("user_user_group").
			Select("user_id").
			Where("user_group_id IN (?)", strings.Split(v, ",")))
	}
	if v := params.ShowExtendInfo; v {
		db.Preload("ExtendInfo")
	}
	if v := params.ShowDetails; v {
		db.Preload(clause.Associations)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + strings.ToLower(v) + "%"
		db.Where("lower(user_name) LIKE ? OR id IN (?)", v, mysql.DB.Table("user_extend_info").
			Select("user_id").
			Where("lower(real_name) LIKE ? OR lower(mobile_phone) LIKE ? OR lower(qq_account) LIKE ? OR lower(email) LIKE ?", v, v, v, v))
	}
	if v := params.OmitPassword; v {
		db.Omit("password")
	}
	if v := params.FindAll; v {
		params.PaginationParam.Pagination = false
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

	var list entity.Users
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.UserQueryResult{
		Data:       list.ToSchemaUsers(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *User) GetByFilterOr(params schema.UserQueryParam) (*schema.UserQueryResult, error) {
	db := mysql.DB.Model(entity.User{}).Omit("password").Preload("ExtendInfo")
	if v := params.ID; v != 0 {
		db = db.Where("id=?", v)
	}
	if v := params.IDs; v != "" {
		db = db.Where("id IN (?)", strings.Split(v, ","))
	}
	if v := params.UserName; v != "" {
		db = db.Where("user_name = ?", v)
	}
	if v := params.Status; v != 0 {
		db = db.Where("status=?", v)
	}
	if v := params.OrgID; v != 0 {
		db.Where("id IN (?)", mysql.DB.Table("user_organization").
			Select("user_id").
			Where("organization_id = ?", v))
	}
	if v := params.OrgIDs; v != "" {
		db.Where("id IN (?)", mysql.DB.Table("user_organization").
			Select("user_id").
			Where("organization_id IN (?)", strings.Split(v, ",")))
	}
	if v := params.PositionID; v != 0 {
		db.Where("id IN (?)", mysql.DB.Table("user_position").
			Select("user_id").
			Where("position_id = ?", v))
	}
	if v := params.PositionIDs; v != "" {
		db.Where("id IN (?)", mysql.DB.Table("user_position").
			Select("user_id").
			Where("position_id IN (?)", strings.Split(v, ",")))
	}
	if v := params.RoleID; v != 0 {
		db.Where("id IN (?)", mysql.DB.Table("user_role").
			Select("user_id").
			Where("role_id = ?", v))
	}
	if v := params.RoleIDs; v != "" {
		db.Where("id IN (?)", mysql.DB.Table("user_role").
			Select("user_id").
			Where("role_id IN (?)", strings.Split(v, ",")))
	}
	if v := params.UserGroupID; v != 0 {
		db.Where("id IN (?)", mysql.DB.Table("user_user_group").
			Select("user_id").
			Where("user_group_id = ?", v))
	}
	if v := params.UserGroupIDs; v != "" {
		db.Where("id IN (?)", mysql.DB.Table("user_user_group").
			Select("user_id").
			Where("user_group_id IN (?)", strings.Split(v, ",")))
	}
	db.Order("id DESC")

	var list entity.Users
	paginationResult, err := mysql.Paginate(db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.UserQueryResult{
		Data:       list.ToSchemaUsers(),
		PageResult: paginationResult,
	}
	return qr, nil
}

func (a *User) Get(id uint64) (*schema.User, error) {
	db := mysql.DB.Model(entity.User{ID: id}).Preload(clause.Associations)

	var item entity.User
	if ok, err := mysql.FindOne(db, &item); err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaUser(), nil
}

func (a *User) Create(item schema.User) (*common.IDResult, error) {
	eitem := entity.SchemaUser(item).ToUser()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

func (a *User) UpdateByID(id uint64, item map[string]interface{}) error {
	result := mysql.DB.Model(entity.User{ID: id}).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *User) Delete(id uint64) error {
	result := mysql.DB.Model(entity.User{ID: id}).Select(clause.Associations).Delete(&entity.User{})
	return errors.WithStack(result.Error)
}

// ---------------------------------------- User Permission --------------------------------------

func (a *User) ReplaceUserOrganizations(id uint64, items schema.Organizations) error {
	eitem := entity.SchemaOrganizations(items).ToOrganization()
	if err := mysql.DB.Model(&entity.User{ID: id}).
		Association("Organizations").
		Replace(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *User) ReplaceUserPositions(id uint64, items schema.Positions) error {
	eitem := entity.SchemaPositions(items).ToPosition()
	if err := mysql.DB.Model(&entity.User{ID: id}).
		Association("Positions").
		Replace(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *User) ReplaceUserRoles(id uint64, items schema.Roles) error {
	eitem := entity.SchemaRoles(items).ToRole()
	if err := mysql.DB.Model(&entity.User{ID: id}).
		Association("Roles").
		Replace(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *User) AppendUserOrganizations(id uint64, items schema.Organizations) error {
	eitem := entity.SchemaOrganizations(items).ToOrganization()
	if err := mysql.DB.Model(&entity.User{ID: id}).
		Association("Organizations").
		Append(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *User) AppendUserPositions(id uint64, items schema.Positions) error {
	eitem := entity.SchemaPositions(items).ToPosition()
	if err := mysql.DB.Model(&entity.User{ID: id}).
		Association("Positions").
		Append(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *User) AppendUserRoles(id uint64, items schema.Roles) error {
	eitem := entity.SchemaRoles(items).ToRole()
	if err := mysql.DB.Model(&entity.User{ID: id}).
		Association("Roles").
		Append(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// TODO: 等待用户组表完成
//func (a *User) AppendUserUserGroup(id uint64, items schema.UserGroups) error {
//	eitem := entity.SchemaUserGroups(items).ToUserGroup()
//	if err := mysql.DB.Model(&entity.User{ID: id}).
//		Association("UserGroups").
//		Append(eitem); err != nil {
//		return errors.WithStack(err)
//	}
//	return nil
//}

func (a *User) UserRemoveOrganization(userID uint64, organization schema.Organization) error {
	eitem := entity.SchemaOrganization(organization).ToOrganization()
	if err := mysql.DB.Model(entity.User{ID: userID}).
		Association("Organizations").
		Delete(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *User) UserRemovePosition(userID uint64, position schema.Position) error {
	eitem := entity.SchemaPosition(position).ToPosition()
	if err := mysql.DB.Model(entity.User{ID: userID}).
		Association("Positions").
		Delete(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *User) UserRemoveRole(userID uint64, role schema.Role) error {
	eitem := entity.SchemaRole(role).ToRole()
	if err := mysql.DB.Model(entity.User{ID: userID}).
		Association("Roles").
		Delete(eitem); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// TODO: 等待用户组表完成
//func (a *User) UserRemoveUserGroup(userID uint64, userGroup schema.UserGroup) error {
//	eitem := entity.SchemaUserGroup(userGroup).ToUserGroup()
//	if err := mysql.DB.Model(entity.User{ID: userID}).
//		Association("UserGroups").
//		Delete(eitem); err != nil {
//		return errors.WithStack(err)
//	}
//	return nil
//}
