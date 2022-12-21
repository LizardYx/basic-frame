package entity

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"gorm.io/gorm"
	"time"
)

type SchemaPosition schema.Position

func (a SchemaPosition) ToPosition() *Position {
	item := new(Position)
	_ = common.Copy(a, item)
	return item
}

type SchemaPositions schema.Positions

func (a SchemaPositions) ToPosition() *Positions {
	item := new(Positions)
	_ = common.Copy(a, item)
	return item
}

type Position struct {
	ID             uint64         `gorm:"primaryKey,autoIncrement;"`
	Name           string         `gorm:"index;default:'';not null;unique;"` // 职位名称
	RoleID         uint64         `gorm:"index;"`                            // 职位的角色ID
	OrganizationID uint64         `gorm:"index;"`                            // 职位的组织ID
	Sequence       int            `gorm:"index;default:0;not null"`          // 排序值
	Status         int            `gorm:"index;default:2;not null"`          // 状态(1:启用 2:禁用)
	Memo           string         `gorm:"size:1024;"`                        // 备注
	Creator        uint64         `gorm:"index"`
	CreatedAt      time.Time      `gorm:"index"`
	UpdatedAt      time.Time      `gorm:"index"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (a Position) ToSchemaPosition() *schema.Position {
	item := new(schema.Position)
	_ = common.Copy(a, item)
	return item
}

type Positions []*Position

func (a Positions) ToSchemaPositions() schema.Positions {
	list := make(schema.Positions, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaPosition()
	}
	return list
}
