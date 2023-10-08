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

// Query 查询用户组列表
//
//	@Summary		查询用户组列表
//	@Description	查询用户组列表
//	@Tags			UserGroup
//	@Param			q	query		schema.UserGroupQueryParam	false	"查询用户组列表参数"
//	@Success		200	{object}	ginx.ListResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/user-groups [get]
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
//
//	@Summary		获取用户组基本信息
//	@Description	获取用户组基本信息
//	@Tags			UserGroup
//	@Param			id	path		int	true	"用户组ID"
//	@Success		200	{object}	schema.UserGroup
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/user-groups/:id [get]
func (a *UserGroup) Get(c *gin.Context) {
	item, err := a.UserGroupBll.Get(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

// GetUserGroupUsers 获取指定用户组的用户列表
//
//	@Summary		获取指定用户组的用户列表
//	@Description	获取指定用户组的用户列表
//	@Tags			UserGroup
//	@Param			id	path		int	true	"用户组ID"
//	@Success		200	{object}	schema.Users
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/user-groups/:id/users [get]
func (a *UserGroup) GetUserGroupUsers(c *gin.Context) {
	item, err := a.UserGroupBll.GetUserGroupUsers(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

// Create 创建用户组
//
//	@Summary		创建用户组
//	@Description	创建用户组
//	@Tags			UserGroup
//	@Param			body	body		schema.UserGroup	true	"创建用户组参数"
//	@Success		200		{object}	common.IDResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/user-groups [post]
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
//
//	@Summary		更新用户组基本信息
//	@Description	更新用户组基本信息
//	@Tags			UserGroup
//	@Param			body	body		schema.UserGroup	true	"用户组基本信息参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/user-groups/:id [put]
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

// UserJoinUserGroup 用户加入指定用户组
//
//	@Summary		用户加入指定用户组
//	@Description	用户加入指定用户组
//	@Tags			UserGroup
//	@Param			id		path		int		true	"用户组ID"
//	@Param			body	body		[]int	true	"用户ID集合"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/user-groups/:id/user-join [put]
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

// UserGroupRemoveUser 用户组中移出多个用户
//
//	@Summary		用户组中移出多个用户
//	@Description	用户组中移出多个用户
//	@Tags			UserGroup
//	@Param			id		path		int		true	"用户组ID"
//	@Param			body	body		[]int	true	"用户ID集合"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/user-groups/:id/user-remove [put]
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
//
//	@Summary		删除用户组
//	@Description	删除用户组
//	@Tags			UserGroup
//	@Param			id	path		int	true	"用户组ID"
//	@Success		200	{object}	ginx.OperateResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/user-groups/:id [delete]
func (a *UserGroup) Delete(c *gin.Context) {
	if err := a.UserGroupBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}
