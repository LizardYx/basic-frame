package api

import (
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"basic-frame/util/ginx"
	"github.com/gin-gonic/gin"
)

var SystemConfigApi = &SystemConfig{
	SystemConfigBll: bll.SystemConfigBll,
}

type SystemConfig struct {
	SystemConfigBll *bll.SystemConfig
}

func (a *SystemConfig) Query(ctx *gin.Context) {
	var params schema.SystemConfigQueryParam
	if err := ginx.ParseQuery(ctx, &params); err != nil {
		ginx.ResError(ctx, err)
		return
	}

	params.Pagination = true
	result, err := a.SystemConfigBll.Query(params)
	if err != nil {
		ginx.ResError(ctx, err)
		return
	}
	ginx.ResPaginate(ctx, result.Data, result.PageResult)
}

func (a *SystemConfig) Create(ctx *gin.Context) {
	var item schema.SystemConfig
	if err := ginx.ParseJSON(ctx, &item); err != nil {
		ginx.ResError(ctx, err)
		return
	}
	result, err := a.SystemConfigBll.Create(ctx, item)
	if err != nil {
		ginx.ResError(ctx, err)
		return
	}
	ginx.ResSuccess(ctx, result)
}
