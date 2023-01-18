package entity

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"gorm.io/gorm"
	"time"
)

type SchemaUser schema.User

func (a SchemaUser) ToUser() *User {
	item := new(User)
	_ = common.Copy(a, item)
	return item
}

type SchemaUsers schema.Users

func (a SchemaUsers) ToUsers() *Users {
	items := new(Users)
	_ = common.Copy(a, items)
	return items
}

type User struct {
	ID            uint64         `gorm:"primaryKey,autoIncrement;"`
	UserName      string         `gorm:"index;default:'';not null;unique;"` // 用户名称
	Password      string         `gorm:"index;default:'';not null;"`        // 用户密码
	Status        int            `gorm:"index;default:2;not null"`          // 状态(1:启用 2:禁用)
	Sequence      int            `gorm:"index"`                             // 排序值
	Creator       uint64         `gorm:"index"`
	CreatedAt     time.Time      `gorm:"index"`
	UpdatedAt     time.Time      `gorm:"index"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	ExtendInfo    UserExtendInfo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Organizations Organizations  `gorm:"many2many:user_organization;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Positions     Positions      `gorm:"many2many:user_position;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Roles         Roles          `gorm:"many2many:user_role;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserGroups    UserGroups     `gorm:"many2many:user_user_group;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (a User) TableName() string {
	return "user"
}

func (a User) ToSchemaUser() *schema.User {
	item := new(schema.User)
	_ = common.Copy(a, item)
	return item.Init()
}

type Users []*User

func (a Users) ToSchemaUsers() schema.Users {
	list := make(schema.Users, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUser()
	}
	return list
}
