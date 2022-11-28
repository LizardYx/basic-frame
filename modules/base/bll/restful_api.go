package bll

import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"context"
)

var ResApiBll = &RestfulApi{
	ResApiModel: model.ResApiModel,
}

type RestfulApi struct {
	ResApiModel *model.RestfulApi
}

func (a *RestfulApi) Create(ctx context.Context, item schema.RestfulApi) (uint64, error) {
	// 参数检查
	// 其它的业务逻辑
	item.UUID = common.GetUUID()
	return a.ResApiModel.Create(item)
}
