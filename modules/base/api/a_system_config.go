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

// Query 查询系统配置信息列表
//
//	@Summary		查询系统配置信息列表
//	@Description	查询系统配置信息列表
//	@Tags			SystemConfig
//	@Param			q	query		schema.SystemConfigQueryParam	false	"查询系统配置信息列表参数"
//	@Success		200	{object}	ginx.ListResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/system-config [get]
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

// Update 更新系统基础配置信息
//
//	@Summary		更新系统基础配置信息
//	@Description	更新系统基础配置信息
//	@Tags			SystemConfig
//	@Param			id		path		int					true	"系统基础配置ID"
//	@Param			body	body		schema.SystemConfig	true	"更新系统基础配置信息"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/system-config/:id [put]
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
