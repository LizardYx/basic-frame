package schema

import (
	"basic-frame/util/common"
	"time"
)

type TagManage struct {
	ID        uint64    `json:"id"`                       // 唯一标识
	Type      string    `json:"type" binding:"required"`  // 标签类型
	KeyName   string    `json:"key_name"`                 // 标签名称
	Value     string    `json:"value" binding:"required"` // 标签值
	Creator   uint64    `json:"creator"`                  // 创建者
	CreatedAt time.Time `json:"created_at"`               // 创建时间
	UpdatedAt time.Time `json:"updated_at"`               // 更新时间

}

type TagManageQueryParam struct {
	common.PaginationParam
	Type       string `form:"type"`        // 标签类型
	KeyName    string `form:"key_name"`    // 标签名称
	Value      string `form:"value"`       // 标签值
	QueryValue string `form:"query_value"` // 模糊搜索(模糊查询:标签名称、标签值)
}

type TagManages []*TagManage

type TagManageQueryResult struct {
	Data       TagManages
	PageResult *common.PaginationResult
}

func (a TagManages) GetIDs() []uint64 {
	l := make([]uint64, len(a))
	for i, j := range a {
		l[i] = j.ID
	}
	return l
}

func (a TagManages) Init() *TagManages {
	if len(a) == 0 {
		a = make(TagManages, 0)
	} else {
		for index, role := range a {
			a[index] = role
		}
	}
	return &a
}
