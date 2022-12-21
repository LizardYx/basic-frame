package api

import (
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"basic-frame/util/ginx"
	"github.com/gin-gonic/gin"
)

var PositionApi = &Position{
	PositionBll: bll.PositionBll,
}

type Position struct {
	PositionBll *bll.Position
}

func (a *Position) Query(c *gin.Context) {
	var params schema.PositionQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, params, err)
		return
	}

	params.Pagination = true
	result, err := a.PositionBll.Query(c, params)
	if err != nil {
		ginx.ResError(c, params, err)
		return
	}
	ginx.ResList(c, params, result.Data)
}

func (a *Position) Get(c *gin.Context) {
	item, err := a.PositionBll.Get(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

func (a *Position) Create(c *gin.Context) {
	var item schema.Position
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, "", err)
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.PositionBll.Create(c, item)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, item, result)
}

func (a *Position) Update(c *gin.Context) {
	var item schema.Position
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.PositionBll.Update(c, ginx.ParseParamID(c, "id"), item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, item)
}

func (a *Position) PositionAddUser(c *gin.Context) {
	var userIDs []uint64
	if err := ginx.ParseJSON(c, &userIDs); err != nil {
		ginx.ResError(c, userIDs, err)
		return
	}
	if err := a.PositionBll.PositionAddUser(c, ginx.ParseParamID(c, "id"), userIDs); err != nil {
		ginx.ResError(c, userIDs, err)
		return
	}
	ginx.ResOperateSuccess(c, userIDs)
}

func (a *Position) PositionRemoveUser(c *gin.Context) {
	var userIDs []uint64
	if err := ginx.ParseJSON(c, &userIDs); err != nil {
		ginx.ResError(c, userIDs, err)
		return
	}
	if err := a.PositionBll.PositionRemoveUser(c, ginx.ParseParamID(c, "id"), userIDs); err != nil {
		ginx.ResError(c, userIDs, err)
		return
	}
	ginx.ResOperateSuccess(c, userIDs)
}

func (a *Position) Delete(c *gin.Context) {
	var params = ginx.ParamsID{ID: ginx.ParseParamID(c, "id")}

	if err := a.PositionBll.Delete(c, params.ID); err != nil {
		ginx.ResError(c, params, err)
		return
	}
	ginx.ResOperateSuccess(c, params)
}
