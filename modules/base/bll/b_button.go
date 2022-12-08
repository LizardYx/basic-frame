package bll

import "basic-frame/modules/base/dao/model"

var ButtonBll = &Button{
	ButtonModel: model.ButtonModel,
}

type Button struct {
	ButtonModel *model.Button
}
