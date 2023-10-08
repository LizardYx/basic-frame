package api

import (
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"basic-frame/util/ginx"
	"github.com/gin-gonic/gin"
)

var ButtonApi = &Button{
	ButtonBll: bll.ButtonBll,
}

type Button struct {
	ButtonBll *bll.Button
}

// Create 创建按钮
//
//	@Summary		创建按钮
//	@Description	创建按钮和按钮关联的RestfulApi(不会创建子按钮)
//	@Tags			Button
//	@Param			body	body		schema.ButtonPre	true	"创建按钮参数"
//	@Success		200		{object}	common.IDResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/button [post]
func (a *Button) Create(c *gin.Context) {
	var item schema.ButtonPre

	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	item.SetCreator(ginx.GetUserID(c), &item)
	result, err := a.ButtonBll.Create(c, item)
	if err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResSuccess(c, item, result)
}

// Update 更新按钮信息
//
//	@Summary		更新按钮信息
//	@Description	更新按钮基础信息、按钮关联的Restful接口信息、指定Restful接口基础信息
//	@Tags			Button
//	@Param			id		path		int					true	"按钮ID"
//	@Param			body	body		schema.ButtonPre	true	"更新按钮信息参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/button/:id [put]
func (a *Button) Update(c *gin.Context) {
	var item schema.ButtonPre
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	item.ID = ginx.ParseParamID(c, "id")
	item.SetCreator(ginx.GetUserID(c), &item)
	if err := a.ButtonBll.UpdateButtonPre(c, item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, item)
}

// BatchUpdate 批量更新按钮基础信息
//
//	@Summary		批量更新按钮基础信息
//	@Description	批量更新按钮基础信息
//	@Tags			Button
//	@Param			body	body		schema.Buttons	true	"更新按钮基础信息参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/button [put]
func (a *Button) BatchUpdate(c *gin.Context) {
	var items schema.Buttons
	if err := ginx.ParseJSON(c, &items); err != nil {
		ginx.ResError(c, items, err)
		return
	}

	if err := a.ButtonBll.BatchUpdate(c, items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	ginx.ResOperateSuccess(c, items)
}

// Delete 删除按钮
//
//	@Summary		删除按钮
//	@Description	删除按钮和按钮调用的Api
//	@Tags			Button
//	@Param			id	path		int	true	"按钮ID"
//	@Success		200	{object}	ginx.OperateResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/button/:id [delete]
func (a *Button) Delete(c *gin.Context) {
	if err := a.ButtonBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}
