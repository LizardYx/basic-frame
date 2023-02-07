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

type RestfulApiQueryParam struct {
	common.PaginationParam
	UUID string `form:"uuid"` // 接口UUID
}

type RestfulApis []*RestfulApi

type RestfulApiQueryResult struct {
	Data       RestfulApis
	PageResult *common.PaginationResult
}

func (a RestfulApis) Init() *RestfulApis {
	items := make(RestfulApis, 0)
	items = append(items, a...)
	return &items
}

// SetCreator 设置Restful请求的创建人
func (a RestfulApis) SetCreator(creator uint64) *RestfulApis {
	if creator != 0 {
		for _, restfulApi := range a {
			restfulApi.Creator = creator
		}
	}
	return &a
}

// InitUUID 初始化Restful请求的UUID
func (a RestfulApis) InitUUID() RestfulApis {
	for _, restfulApi := range a {
		restfulApi.UUID = common.GetUUID()
	}
	return a
}
