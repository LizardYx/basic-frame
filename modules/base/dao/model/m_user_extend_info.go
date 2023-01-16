package model

import (
	"basic-frame/modules/base/dao/entity"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/mysql"
	"github.com/pkg/errors"
)

var UserExtendInfoModel = &UserExtendInfo{}

type UserExtendInfo struct {
}

func (a *UserExtendInfo) Get(id uint64) (*schema.UserExtendInfo, error) {
	db := mysql.DB.Model(entity.UserExtendInfo{ID: id})

	var item entity.UserExtendInfo
	if ok, err := mysql.FindOne(db, &item); err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaUserExtendInfo(), nil
}

func (a *UserExtendInfo) GetByRealName(realName string) (*schema.UserExtendInfo, error) {
	db := mysql.DB.Model(entity.UserExtendInfo{RealName: realName})

	var item entity.UserExtendInfo
	if ok, err := mysql.FindOne(db, &item); err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaUserExtendInfo(), nil
}

func (a *UserExtendInfo) Create(item schema.UserExtendInfo) (*common.IDResult, error) {
	eitem := entity.SchemaUserExtendInfo(item).ToUserExtendInfo()
	result := mysql.DB.Create(&eitem)
	return &common.IDResult{ID: eitem.ID}, errors.WithStack(result.Error)
}

func (a *UserExtendInfo) UpdateByID(id uint64, item map[string]interface{}) error {
	result := mysql.DB.Model(entity.UserExtendInfo{ID: id}).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *UserExtendInfo) Delete(id uint64) error {
	result := mysql.DB.Model(entity.UserExtendInfo{ID: id}).Delete(&entity.UserExtendInfo{})
	return errors.WithStack(result.Error)
}
