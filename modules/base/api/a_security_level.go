package api

import (
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"basic-frame/util/ginx"
	"github.com/gin-gonic/gin"
)

var SecurityLevelApi = &SecurityLevel{
	SecurityLevelBll: bll.SecurityLevelBll,
}

type SecurityLevel struct {
	SecurityLevelBll *bll.SecurityLevel
}

func (a *SecurityLevel) Query(c *gin.Context) {
	var params schema.SecurityLevelQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, params, err)
		return
	}

	params.Pagination = true
	result, err := a.SecurityLevelBll.Query(c, params)
	if err != nil {
		ginx.ResError(c, params, err)
		return
	}
	ginx.ResList(c, params, result.Data)
}

// Get 获取安全级别和绑定的角色信息
func (a *SecurityLevel) Get(c *gin.Context) {
	item, err := a.SecurityLevelBll.Get(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

func (a *SecurityLevel) GetUserSecurityLevels(c *gin.Context) {
	item, err := a.SecurityLevelBll.GetUserSecurityLevels(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

// Create 创建安全级别(如果绑定的角色没有ID，则创建角色)
func (a *SecurityLevel) Create(c *gin.Context) {
	var item schema.SecurityLevel
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.SecurityLevelBll.Create(c, item)
	if err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResSuccess(c, item, result)
}

// Update 更新安全级别基础信息和安全级别与角色的关联关系
func (a *SecurityLevel) Update(c *gin.Context) {
	var item schema.SecurityLevel
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.SecurityLevelBll.Update(c, ginx.ParseParamID(c, "id"), item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, item)
}

// Delete 删除安全级别
func (a *SecurityLevel) Delete(c *gin.Context) {
	if err := a.SecurityLevelBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}
