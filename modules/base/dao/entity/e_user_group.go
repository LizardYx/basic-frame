package entity

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"gorm.io/gorm"
	"time"
)

type SchemaUserGroup schema.UserGroup

func (a SchemaUserGroup) ToUserGroup() *UserGroup {
	item := new(UserGroup)
	_ = common.Copy(a, item)
	return item
}

type SchemaUserGroups schema.UserGroups

func (a SchemaUserGroups) ToUserGroup() *UserGroups {
	item := new(UserGroups)
	_ = common.Copy(a, item)
	return item
}

type UserGroup struct {
	ID        uint64         `gorm:"primaryKey,autoIncrement;"`
	Name      string         `gorm:"index;default:'';not null;unique;"` // 用户组名称
	RoleID    uint64         `gorm:"index;"`                            // 用户组的基础角色
	Sequence  int            `gorm:"index;default:0;not null"`          // 排序值
	Status    int            `gorm:"index;default:2;not null"`          // 状态(1:启用 2:禁用)
	Memo      string         `gorm:"size:1024;"`                        // 备注
	Creator   uint64         `gorm:"index"`
	CreatedAt time.Time      `gorm:"index"`
	UpdatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (a UserGroup) ToSchemaUserGroup() *schema.UserGroup {
	item := new(schema.UserGroup)
	_ = common.Copy(a, item)
	return item
}

type UserGroups []*UserGroup

func (a UserGroups) ToSchemaUserGroups() schema.UserGroups {
	list := make(schema.UserGroups, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUserGroup()
	}
	return list
}
