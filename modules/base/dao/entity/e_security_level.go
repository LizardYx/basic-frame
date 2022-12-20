package entity

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"gorm.io/gorm"
	"time"
)

type SchemaSecurityLevel schema.SecurityLevel

func (a SchemaSecurityLevel) ToSecurityLevel() *SecurityLevel {
	item := new(SecurityLevel)
	_ = common.Copy(a, item)
	return item
}

type SecurityLevel struct {
	ID        uint64         `gorm:"primaryKey,autoIncrement;"`
	Name      string         `gorm:"index;default:'';not null;unique;"` // 安全级别名称
	Sequence  int            `gorm:"index;default:0;not null"`          // 排序值
	Status    int            `gorm:"index;default:2;not null"`          // 状态(1:启用 2:禁用)
	Memo      string         `gorm:"size:1024;"`                        // 备注
	Creator   uint64         `gorm:"index"`
	CreatedAt time.Time      `gorm:"index"`
	UpdatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Roles     Roles          `gorm:"many2many:security_level_role;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (a SecurityLevel) ToSchemaSecurityLevel() *schema.SecurityLevel {
	item := new(schema.SecurityLevel)
	_ = common.Copy(a, item)
	return item.Init()
}

type SecurityLevels []*SecurityLevel

func (a SecurityLevels) ToSchemaSecurityLevels() schema.SecurityLevels {
	list := make(schema.SecurityLevels, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaSecurityLevel()
	}
	return *list.Sort()
}
