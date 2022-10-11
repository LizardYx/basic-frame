package api

import (
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ResApi = &RestfulApi{
	RestfulApiBll: bll.ResApiBll,
}

type RestfulApi struct {
	RestfulApiBll *bll.RestfulApi
}

func (a *RestfulApi) Create(ctx *gin.Context) {
	var item schema.RestfulApi
	if err := ctx.ShouldBindJSON(&item); err != nil {
		// 封装Gin response
		// TODO:
		ctx.JSON(http.StatusOK, err.Error())
		return
	}
	id, err := a.RestfulApiBll.Create(ctx, item)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}
	ctx.JSON(http.StatusOK, id)
}
