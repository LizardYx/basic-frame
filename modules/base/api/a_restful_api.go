package api

import "basic-frame/modules/base/bll"

var RestfulApiApi = &RestfulApi{
	RestfulApiBll: bll.RestfulApiBll,
}

type RestfulApi struct {
	RestfulApiBll *bll.RestfulApi
}
