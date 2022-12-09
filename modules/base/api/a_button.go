package api

import "basic-frame/modules/base/bll"

var ButtonApi = &Button{
	ButtonBll: bll.ButtonBll,
}

type Button struct {
	ButtonBll *bll.Button
}
