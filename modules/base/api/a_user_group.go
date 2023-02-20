package api

import (
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"basic-frame/util/ginx"
	"github.com/gin-gonic/gin"
)

var UserGroupApi = &UserGroup{
	UserGroupBll: bll.UserGroupBll,
}

type UserGroup struct {
	UserGroupBll *bll.UserGroup
}

func (a *UserGroup) Query(c *gin.Context) {
	var params schema.UserGroupQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, params, err)
		return
	}

	params.Pagination = true
	result, err := a.UserGroupBll.Query(c, params)
	if err != nil {
		ginx.ResError(c, params, err)
		return
	}
	ginx.ResList(c, params, result.Data)
}

// Get 获取用户组基本信息
func (a *UserGroup) Get(c *gin.Context) {
	item, err := a.UserGroupBll.Get(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

// GetUserGroupUsers 获取指定用户组的用户列表
func (a *UserGroup) GetUserGroupUsers(c *gin.Context) {
	item, err := a.UserGroupBll.GetUserGroupUsers(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

// Create 创建用户组
func (a *UserGroup) Create(c *gin.Context) {
	var item schema.UserGroup
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.UserGroupBll.Create(c, item)
	if err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResSuccess(c, item, result)
}

// Update 更新用户组基本信息
func (a *UserGroup) Update(c *gin.Context) {
	var item schema.UserGroup
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.UserGroupBll.Update(c, ginx.ParseParamID(c, "id"), item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, item)
}

// UserJoinUserGroup 多个用户加入指定用户组
func (a *UserGroup) UserJoinUserGroup(c *gin.Context) {
	var items []uint64
	if err := ginx.ParseJSON(c, &items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	if err := a.UserGroupBll.UserJoinUserGroup(c, ginx.ParseParamID(c, "id"), items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

// UserGroupRemoveUser 从用户组中移出多个用户
func (a *UserGroup) UserGroupRemoveUser(c *gin.Context) {
	var items []uint64
	if err := ginx.ParseJSON(c, &items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	if err := a.UserGroupBll.UserGroupRemoveUser(c, ginx.ParseParamID(c, "id"), items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

// Delete 删除用户组
func (a *UserGroup) Delete(c *gin.Context) {
	if err := a.UserGroupBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}
