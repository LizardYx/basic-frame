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

// Query 查询职位列表
//
//	@Summary		查询职位列表
//	@Description	查询职位列表
//	@Tags			Position
//	@Param			q	query		schema.PositionQueryParam	false	"查询职位列表参数"
//	@Success		200	{object}	ginx.ListResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/positions [get]
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

// Get 获取职位基础信息
//
//	@Summary		获取职位基础信息
//	@Description	获取职位基础信息
//	@Tags			Position
//	@Param			id	path		int	true	"职位ID"
//	@Success		200	{object}	schema.Position
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/positions/:id [get]
func (a *Position) Get(c *gin.Context) {
	item, err := a.PositionBll.Get(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

// Create 创建职位
//
//	@Summary		创建职位
//	@Description	创建职位
//	@Tags			Position
//	@Param			body	body		schema.Position	true	"创建职位参数"
//	@Success		200		{object}	common.IDResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/positions [post]
func (a *Position) Create(c *gin.Context) {
	var item schema.Position
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.PositionBll.Create(c, item)
	if err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResSuccess(c, item, result)
}

// Update 更新职位基础信息
//
//	@Summary		更新职位基础信息
//	@Description	更新职位基础信息
//	@Tags			Position
//	@Param			id		path		int				true	"职位ID"
//	@Param			body	body		schema.Position	true	"更新职位基础信息参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/positions/:id [put]
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

// PositionAddUser 用户新增职位
//
//	@Summary		用户新增职位
//	@Description	用户新增职位
//	@Tags			Position
//	@Param			id		path		int		true	"职位ID"
//	@Param			body	body		[]int	true	"用户ID集合"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/positions/user-join/:id [put]
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

// PositionRemoveUser 用户移除职位
//
//	@Summary		用户移除职位
//	@Description	用户移除职位
//	@Tags			Position
//	@Param			id		path		int		true	"职位ID"
//	@Param			body	body		[]int	true	"用户ID集合"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/positions/user-remove/:id [put]
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

// Delete 删除职位
//
//	@Summary		删除职位
//	@Description	删除职位
//	@Tags			Position
//	@Param			id	path		int	true	"职位ID"
//	@Success		200	{object}	ginx.OperateResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/positions/:id [delete]
func (a *Position) Delete(c *gin.Context) {
	if err := a.PositionBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}
