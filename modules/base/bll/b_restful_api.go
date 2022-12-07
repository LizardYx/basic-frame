package bll

import "basic-frame/modules/base/dao/model"

var RestfulApiBll = &RestfulApi{
	RestfulApiModel: model.RestfulApiModel,
}

type RestfulApi struct {
	RestfulApiModel *model.RestfulApi
}
