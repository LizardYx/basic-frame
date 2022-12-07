package entity

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"gorm.io/gorm"
	"time"
)

type SchemaDisabledField schema.DisabledField

func (a SchemaDisabledField) ToDisabledField() *DisabledField {
	item := new(DisabledField)
	_ = common.Copy(a, item)
	return item
}

type DisabledField struct {
	ID        uint64         `gorm:"primaryKey,autoIncrement;"`
	UUID      string         `gorm:"index;"`                     // 可禁用字段UUID
	KeyName   string         `gorm:"index;default:'';not null;"` // 字段名称
	KeyValue  string         `gorm:"index;default:'';not null;"` // 字段值
	Memo      string         `gorm:"size:1024;"`                 // 备注
	Creator   uint64         `gorm:"index"`
	CreatedAt time.Time      `gorm:"index"`
	UpdatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (a DisabledField) ToSchemaDisabledField() *schema.DisabledField {
	item := new(schema.DisabledField)
	_ = common.Copy(a, item)
	return item
}

type DisabledFields []*DisabledField

func (a DisabledFields) ToSchemaDisabledFields() schema.DisabledFields {
	item := new(schema.DisabledFields)
	_ = common.Copy(a, item)
	return *item.Init()
}
