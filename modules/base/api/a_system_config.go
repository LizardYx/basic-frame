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

func (a *SystemConfig) Query(c *gin.Context) {
	var params schema.SystemConfigQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, params, err)
		return
	}

	params.Pagination = true
	result, err := a.SystemConfigBll.Query(c, params)
	if err != nil {
		ginx.ResError(c, params, err)
		return
	}
	ginx.ResList(c, params, result.Data)
}

func (a *SystemConfig) Update(c *gin.Context) {
	var item schema.SystemConfig
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	err := a.SystemConfigBll.Update(c, ginx.ParseParamID(c, "id"), item)
	if err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResSuccessString(c, item, "更新成功")
}
