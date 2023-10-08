package api

import (
	"basic-frame/middleware"
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"basic-frame/util/consts"
	"basic-frame/util/ginx"
	"github.com/gin-gonic/gin"
)

var UserApi = &User{
	UserBll: bll.UserBll,
}

type User struct {
	UserBll *bll.User
}

// Login 用户登录
//
//	@Summary		用户登录
//	@Description	用户登录
//	@Tags			User
//	@Param			body	body		schema.LoginParam	true	"用户登录参数"
//	@Success		200		{object}	schema.LoginRes
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/login [post]
func (a *User) Login(c *gin.Context) {
	var item schema.LoginParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	tokenInfo, userInfo, err := a.UserBll.Login(c, item)
	if err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResSuccess(c, "", schema.LoginRes{
		UserInfo:  *userInfo,
		TokenInfo: *tokenInfo,
	})
}

// Logout 用户登出
//
//	@Summary		用户登出
//	@Description	用户登出(前端需要自己删除Jwt Token字符串)
//	@Tags			User
//	@Success		200	{object}	ginx.OperateResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/logout [post]
func (a *User) Logout(c *gin.Context) {
	userID := ginx.GetUserID(c)
	if userID != 0 {
		// 注销ws连接
		for client := range middleware.Manager.Clients {
			if client.ID == userID {
				middleware.Manager.Unregister <- client
				break
			}
		}
	}
	ginx.ResOperateSuccess(c, "")
}

// RefreshToken 更新用户登录缓存Token
//
//	@Summary		更新用户登录缓存Token
//	@Description	更新用户登录缓存Token
//	@Tags			User
//	@Success		200	{string}	json	"{"token_info": ""}"
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/refresh-token [get]
func (a *User) RefreshToken(c *gin.Context) {
	token, err := middleware.RefreshToken(ginx.GetToken(c))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", map[string]interface{}{
		"token_info": token,
	})
}

// Query 查询用户列表
//
//	@Summary		查询用户列表
//	@Description	查询用户列表
//	@Tags			User
//	@Param			q	query		schema.UserQueryParam	false	"查询用户列表参数"
//	@Success		200	{object}	ginx.ListResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/users [get]
func (a *User) Query(c *gin.Context) {
	var params schema.UserQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, params, err)
		return
	}

	params.Pagination = true
	params.OmitPassword = true
	result, err := a.UserBll.Query(c, params)
	if err != nil {
		ginx.ResError(c, params, err)
		return
	}
	ginx.ResList(c, params, result.Data)
}

// GetSubUsers 获取子用户列表
//
//	@Summary		获取子用户列表
//	@Description	获取子用户列表
//	@Tags			User
//	@Param			id	path		int							true	"用户ID"
//	@Param			q	query		schema.SubUserQueryParam	false	"获取子用户列表参数"
//	@Success		200	{object}	ginx.ListResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/users/:id/sub-users [get]
func (a *User) GetSubUsers(c *gin.Context) {
	var params schema.SubUserQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, params, err)
		return
	}

	result, err := a.UserBll.GetSubUsers(c, ginx.GetUserID(c), params)
	if err != nil {
		ginx.ResError(c, params, err)
		return
	}
	ginx.ResList(c, params, result.Data)
}

// Get 获取用户信息
//
//	@Summary		获取用户信息
//	@Description	获取用户和用户与组织、职位、角色的关联信息
//	@Tags			User
//	@Param			id	path		int	true	"用户ID"
//	@Success		200	{object}	schema.User
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/users/:id [get]
func (a *User) Get(c *gin.Context) {
	item, err := a.UserBll.Get(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

// GetMenuTree 获取用户菜单树
//
//	@Summary		获取用户菜单树
//	@Description	获取用户菜单树
//	@Tags			User
//	@Success		200	{object}	schema.MenuPres
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/menu/menu-tree [get]
func (a *User) GetMenuTree(c *gin.Context) {
	item, err := a.UserBll.GetMenuTree(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

// UpdateUserPermission 更新用户的信息和权限
//
//	@Summary		更新用户的信息和权限
//	@Description	更新用户基础信息及用户与组织、职位、角色、用户组的关联关系
//	@Tags			User
//	@Param			id		path		int			true	"用户ID"
//	@Param			body	body		schema.User	true	"更新用户信息和权限的参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/users/:id/permission [put]
func (a *User) UpdateUserPermission(c *gin.Context) {
	var item schema.User
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.UserBll.UpdateUserPermission(c, ginx.ParseParamID(c, "id"), item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

// BatchUpdateUserPermission 批量更新用户的信息和权限
//
//	@Summary		批量更新用户的信息和权限
//	@Description	更新多个的用户基础信息及用户与组织、职位、角色、用户组的关联关系
//	@Tags			User
//	@Param			body	body		schema.Users	true	"更新用户信息和权限的参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/users [put]
func (a *User) BatchUpdateUserPermission(c *gin.Context) {
	var items schema.Users
	if err := ginx.ParseJSON(c, &items); err != nil {
		ginx.ResError(c, items, err)
		return
	}

	if err := a.UserBll.BatchUpdateUserPermission(c, items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

// Create 创建用户
//
//	@Summary		创建用户
//	@Description	创建用户
//	@Tags			User
//	@Param			body	body		schema.User	true	"创建用户参数"
//	@Success		200		{object}	common.IDResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/users [post]
func (a *User) Create(c *gin.Context) {
	var item schema.User
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.UserBll.Create(c, item)
	if err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResSuccess(c, item, result)
}

// Update 更新用户基础信息
//
//	@Summary		更新用户基础信息
//	@Description	更新用户基础信息
//	@Tags			User
//	@Param			id		path		int			true	"用户ID"
//	@Param			body	body		schema.User	true	"更新用户基础信息参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/users/:id [put]
func (a *User) Update(c *gin.Context) {
	var item schema.User
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.UserBll.Update(c, ginx.ParseParamID(c, "id"), item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, item)
}

// EnableUser 启用指定用户
//
//	@Summary		启用指定用户
//	@Description	启用指定用户
//	@Tags			User
//	@Param			id	path		int	true	"用户ID"
//	@Success		200	{object}	ginx.OperateResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/users/:id/enable [put]
func (a *User) EnableUser(c *gin.Context) {
	if err := a.UserBll.UpdateStatus(c, ginx.ParseParamID(c, "id"), consts.BaseStatusEnable); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

// DisabledUser 禁用指定用户
//
//	@Summary		禁用指定用户
//	@Description	禁用指定用户
//	@Tags			User
//	@Param			id	path		int	true	"用户ID"
//	@Success		200	{object}	ginx.OperateResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/users/:id/disabled [put]
func (a *User) DisabledUser(c *gin.Context) {
	if err := a.UserBll.UpdateStatus(c, ginx.ParseParamID(c, "id"), consts.BaseStatusDisabled); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

// UpdatePassword 更新用户密码
//
//	@Summary		更新用户密码
//	@Description	更新用户密码
//	@Tags			User
//	@Param			body	body		schema.UpdatePasswordParam	true	"更新用户密码参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/update-password [put]
func (a *User) UpdatePassword(c *gin.Context) {
	var item schema.UpdatePasswordParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.UserBll.UpdatePassword(c, ginx.GetUserID(c), item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

// Delete 删除用户
//
//	@Summary		删除用户
//	@Description	删除用户、用户扩展信息
//	@Tags			User
//	@Param			id	path		int	true	"用户ID"
//	@Success		200	{object}	ginx.OperateResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/users/:id [delete]
func (a *User) Delete(c *gin.Context) {
	if err := a.UserBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}
