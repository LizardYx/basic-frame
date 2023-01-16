package entity

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"gorm.io/gorm"
	"time"
)

type SchemaUserExtendInfo schema.UserExtendInfo

func (a SchemaUserExtendInfo) ToUserExtendInfo() *UserExtendInfo {
	item := new(UserExtendInfo)
	_ = common.Copy(a, item)
	return item
}

type UserExtendInfo struct {
	ID          uint64         `gorm:"primaryKey,autoIncrement;"`
	UserID      uint64         `gorm:"index;default:0;not null;unique;"` // 用户ID
	RealName    string         `gorm:"index;default:'';not null;"`       // 用户昵称
	MobilePhone int            `gorm:"index;"`                           // 移动手机
	QQAccount   int            `gorm:"index;"`                           // QQ账号
	Email       string         `gorm:"index;"`                           // 邮箱账号
	Creator     uint64         `gorm:"index"`
	CreatedAt   time.Time      `gorm:"index"`
	UpdatedAt   time.Time      `gorm:"index"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (a UserExtendInfo) ToSchemaUserExtendInfo() *schema.UserExtendInfo {
	item := new(schema.UserExtendInfo)
	_ = common.Copy(a, item)
	return item
}
