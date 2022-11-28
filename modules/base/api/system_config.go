package api

import (
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

var SystemConfigApi = &SystemConfig{
	SystemConfigBll: bll.SystemConfigBll,
}

type SystemConfig struct {
	SystemConfigBll *bll.SystemConfig
}

func (a *SystemConfig) Query(ctx *gin.Context) {
	var params schema.SystemConfigQueryParam
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusOK, err.Error())
		return
	}

	params.Pagination = true
	result, err := a.SystemConfigBll.Query(params)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (a *SystemConfig) Create(ctx *gin.Context) {
	var item schema.SystemConfig
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusOK, err.Error())
		return
	}
	id, err := a.SystemConfigBll.Create(ctx, item)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}
	ctx.JSON(http.StatusOK, id)
}
