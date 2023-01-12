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

func (a *TagManage) Get(c *gin.Context) {
	item, err := a.TagManageBll.Get(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

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

func (a *TagManage) Delete(c *gin.Context) {
	if err := a.TagManageBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}
