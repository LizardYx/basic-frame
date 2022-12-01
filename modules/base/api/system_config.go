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
		ginx.ResError(ctx, params, err)
		return
	}

	params.Pagination = true
	result, err := a.SystemConfigBll.Query(params)
	if err != nil {
		ginx.ResError(ctx, params, err)
		return
	}
	ginx.ResList(ctx, params, result.Data)
}

func (a *SystemConfig) Update(ctx *gin.Context) {
	var item schema.SystemConfig
	if err := ginx.ParseQuery(ctx, &item); err != nil {
		ginx.ResError(ctx, item, err)
		return
	}

	err := a.SystemConfigBll.Update(ginx.ParseParamID(ctx, "id"), item)
	if err != nil {
		ginx.ResError(ctx, item, err)
		return
	}
	ginx.ResSuccessString(ctx, item, "更新成功")
}
