package entity

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"gorm.io/gorm"
	"time"
)

type SchemaOrganization schema.Organization

func (a SchemaOrganization) ToOrganization() *Organization {
	item := new(Organization)
	_ = common.Copy(a, item)
	return item
}

type SchemaOrganizations schema.Organizations

func (a SchemaOrganizations) ToOrganization() *Organizations {
	items := new(Organizations)
	_ = common.Copy(a, items)
	return items
}

type Organization struct {
	ID               uint64         `gorm:"primaryKey,autoIncrement;"`
	Name             string         `gorm:"index;default:'';not null;unique;"` // 组织名称
	RoleID           uint64         `gorm:"index;"`                            // 组织的基础角色ID
	Sequence         int            `gorm:"index;default:0;not null;"`         // 排序值
	ParentID         *uint64        `gorm:"index;"`                            // 父级组织ID
	Status           int            `gorm:"index;default:2;not null;"`         // 状态(1:启用 2:禁用)
	Memo             string         `gorm:"size:1024;"`                        // 备注
	Creator          uint64         `gorm:"index"`
	CreatedAt        time.Time      `gorm:"index"`
	UpdatedAt        time.Time      `gorm:"index"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	Positions        []Position     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`                      // 组织的职位列表
	SonOrganizations []Organization `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // 下属组织列表
}

func (a Organization) ToSchemaOrganization() *schema.Organization {
	item := new(schema.Organization)
	_ = common.Copy(a, item)
	return item.Init()
}

type Organizations []*Organization

func (a Organizations) ToSchemaOrganizations() schema.Organizations {
	list := make(schema.Organizations, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaOrganization()
	}
	return list
}
