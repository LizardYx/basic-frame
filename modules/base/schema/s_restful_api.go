package schema

import (
	"basic-frame/util/common"
	"time"
)

type RestfulApi struct {
	ID        uint64    `json:"id"`                        // 唯一标识
	UUID      string    `json:"UUID"`                      // 接口UUID
	Method    string    `json:"method" binding:"required"` // 资源请求方式(支持正则)
	Path      string    `json:"path" binding:"required"`   // 资源请求路径（支持/:id匹配）
	Memo      string    `json:"memo"`                      // 备注
	Creator   uint64    `json:"creator"`                   // 创建者
	CreatedAt time.Time `json:"created_at"`                // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                // 更新时间
}

type RestfulApis []*RestfulApi

type RestfulApiQueryParam struct {
	common.PaginationParam
}

type RestfulApiQueryResult struct {
	Data       RestfulApis
	PageResult *common.PaginationResult
}

func (a RestfulApis) Init() *RestfulApis {
	if len(a) == 0 {
		a = make(RestfulApis, 0)
	}
	return &a
}
