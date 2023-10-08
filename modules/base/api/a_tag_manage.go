package api

import (
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"basic-frame/util/ginx"
	"github.com/gin-gonic/gin"
)

var TagManageApi = &TagManage{
	TagManageBll: bll.TagManageBll,
}

type TagManage struct {
	TagManageBll *bll.TagManage
}

// Query 查询标签列表
//
//	@Summary		查询标签列表
//	@Description	查询标签列表
//	@Tags			TagManage
//	@Param			q	query		schema.TagManageQueryParam	false	"查询标签列表参数"
//	@Success		200	{object}	ginx.ListResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/tag-manages [get]
func (a *TagManage) Query(c *gin.Context) {
	var params schema.TagManageQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, params, err)
		return
	}

	params.Pagination = true
	result, err := a.TagManageBll.Query(c, params)
	if err != nil {
		ginx.ResError(c, params, err)
		return
	}
	ginx.ResList(c, params, result.Data)
}

// Get 查询指定标签详情
//
//	@Summary		查询指定标签详情
//	@Description	查询指定标签详情
//	@Tags			TagManage
//	@Param			id	path		int	true	"标签ID"
//	@Success		200	{object}	schema.TagManage
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/tag-manages/:id [get]
func (a *TagManage) Get(c *gin.Context) {
	item, err := a.TagManageBll.Get(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

// Create 创建标签
//
//	@Summary		创建标签
//	@Description	创建标签
//	@Tags			TagManage
//	@Param			body	body		schema.TagManage	true	"创建标签参数"
//	@Success		200		{object}	common.IDResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/tag-manages [post]
func (a *TagManage) Create(c *gin.Context) {
	var item schema.TagManage
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.TagManageBll.Create(c, item)
	if err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResSuccess(c, item, result)
}

// Update 更新标签信息
//
//	@Summary		更新标签信息
//	@Description	更新标签信息
//	@Tags			TagManage
//	@Param			id		path		int					true	"标签ID"
//	@Param			body	body		schema.TagManage	true	"更新标签信息参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/tag-manages/:id [put]
func (a *TagManage) Update(c *gin.Context) {
	var item schema.TagManage
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.TagManageBll.Update(c, ginx.ParseParamID(c, "id"), item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, item)
}

// Delete 删除标签信息
//
//	@Summary		删除标签信息
//	@Description	删除标签信息
//	@Tags			TagManage
//	@Param			id	path		int	true	"标签ID"
//	@Success		200	{object}	ginx.OperateResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/tag-manages/:id [delete]
func (a *TagManage) Delete(c *gin.Context) {
	if err := a.TagManageBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}
