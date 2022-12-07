package schema

import (
	"basic-frame/util/common"
	"time"
)

type DisabledField struct {
	ID        uint64    `json:"id"`                           // 唯一标识
	UUID      string    `json:"uuid"`                         // 可禁用字段UUID
	KeyName   string    `json:"key_name" binding:"required"`  // 字段名称
	KeyValue  string    `json:"key_value" binding:"required"` // 字段值
	Memo      string    `json:"memo"`                         // 备注
	Creator   uint64    `json:"creator"`                      // 创建者
	CreatedAt time.Time `json:"created_at"`                   // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                   // 更新时间

}

type DisabledFields []*DisabledField

type DisabledFieldQueryParam struct {
	common.PaginationParam
	UUID string `form:"uuid"` // 按钮UUID
}

type DisabledFieldQueryResult struct {
	Data       DisabledFields
	PageResult *common.PaginationResult
}

func (a DisabledFields) Init() *DisabledFields {
	if len(a) == 0 {
		a = make(DisabledFields, 0)
	}
	return &a
}
