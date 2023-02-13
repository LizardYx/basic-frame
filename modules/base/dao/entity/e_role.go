package entity

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"gorm.io/gorm"
	"time"
)

type SchemaRole schema.Role

func (a SchemaRole) ToRole() *Role {
	item := new(Role)
	_ = common.Copy(a, item)
	return item
}

type SchemaRoles schema.Roles

func (a SchemaRoles) ToRole() *Roles {
	item := new(Roles)
	_ = common.Copy(a, item)
	return item
}

type SchemaRolePre schema.RolePre

func (a SchemaRolePre) ToRole() *Role {
	item := new(Role)
	_ = common.Copy(a, item)
	return item
}

type SchemaRolePres schema.RolePre

func (a SchemaRolePres) ToRole() *Roles {
	item := new(Roles)
	_ = common.Copy(a, item)
	return item
}

type Role struct {
	ID             uint64          `gorm:"primaryKey,autoIncrement;"`
	Name           string          `gorm:"index;default:'';not null;unique;"` // 角色名称
	Sequence       int             `gorm:"index;default:0;not null"`          // 排序值
	Type           int             `gorm:"index;"`                            // 角色类型(1.用户角色 2.组织角色 3.职位角色 4.用户组角色)
	AuditorTypes   string          `gorm:"index;"`                            // 审核类型(逗号分隔)
	Status         int             `gorm:"index;default:2;not null"`          // 状态(1:启用 2:禁用)
	Memo           string          `gorm:"size:1024;"`                        // 备注
	Creator        uint64          `gorm:"index"`
	CreatedAt      time.Time       `gorm:"index"`
	UpdatedAt      time.Time       `gorm:"index"`
	DeletedAt      gorm.DeletedAt  `gorm:"index"`
	Menus          []Menu          `gorm:"many2many:role_menu;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`            // 角色关联的菜单
	Buttons        []Button        `gorm:"many2many:role_button;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`          // 角色关联的按钮
	DisabledFields []DisabledField `gorm:"many2many:role_disabled_fields;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 角色关联的可禁用字段
}

func (a Role) ToSchemaRole() *schema.Role {
	item := new(schema.Role)
	_ = common.Copy(a, item)
	return item
}

func (a Role) ToSchemaRolePre() *schema.RolePre {
	item := new(schema.RolePre)
	_ = common.Copy(a, item)
	item = item.Init()
	return item
}

type Roles []*Role

func (a Roles) ToSchemaRoles() schema.Roles {
	list := make(schema.Roles, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRole()
	}
	return list
}

func (a Roles) ToSchemaRolePres() schema.RolePres {
	list := make(schema.RolePres, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRolePre()
	}
	return list
}
