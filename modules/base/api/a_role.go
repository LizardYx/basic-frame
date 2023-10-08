package api

import (
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"basic-frame/util/ginx"
	"github.com/gin-gonic/gin"
)

var RoleApi = &Role{
	RoleBll: bll.RoleBll,
	UserBll: bll.UserBll,
}

type Role struct {
	RoleBll *bll.Role
	UserBll *bll.User
}

// Query 查询角色列表
//
//	@Summary		查询角色列表
//	@Description	查询角色列表
//	@Tags			Role
//	@Param			q	query		schema.RoleQueryParam	false	"查询角色列表参数"
//	@Success		200	{object}	ginx.ListResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/roles [get]
func (a *Role) Query(c *gin.Context) {
	var params schema.RoleQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, params, err)
		return
	}

	params.Pagination = true
	result, err := a.RoleBll.Query(c, params)
	if err != nil {
		ginx.ResError(c, params, err)
		return
	}
	ginx.ResList(c, params, result.Data)
}

// Get 获取角色所有信息
//
//	@Summary		获取角色所有信息
//	@Description	获取角色所有信息(包含菜单、按钮、禁用字段。以及菜单、按钮是否为全选)
//	@Tags			Role
//	@Param			id	path		int	true	"角色ID"
//	@Success		200	{object}	schema.RolePre
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/roles/:id [get]
func (a *Role) Get(c *gin.Context) {
	item, err := a.RoleBll.GetPreWithSelect(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

// Create 创建角色
//
//	@Summary		创建角色
//	@Description	创建角色
//	@Tags			Role
//	@Param			body	body		schema.RolePre	true	"创建角色参数"
//	@Success		200		{object}	common.IDResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/roles [post]
func (a *Role) Create(c *gin.Context) {
	var item schema.RolePre
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.RoleBll.Create(c, item)
	if err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResSuccess(c, item, result)
}

// Update 更新角色信息
//
//	@Summary		更新角色信息
//	@Description	更新角色基本信息、角色和菜单的关联、角色和按钮的关联、角色和可禁用字段的关联
//	@Tags			Role
//	@Param			body	body		schema.RolePre	true	"更新角色信息参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/roles/:id [put]
func (a *Role) Update(c *gin.Context) {
	var item schema.RolePre
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.RoleBll.UpdateDetails(c, item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, item)
}

// UpdateAuditorType 角色新增审核类型
//
//	@Summary		角色新增审核类型
//	@Description	角色新增审核类型
//	@Tags			Role
//	@Param			id		path		int								true	"角色ID"
//	@Param			body	body		schema.UpdateAuditorTypeParam	true	"角色新增审核类型参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/roles/:id/auditor_type [put]
func (a *Role) UpdateAuditorType(c *gin.Context) {
	var item schema.UpdateAuditorTypeParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.RoleBll.UpdateAuditorType(c, ginx.ParseParamID(c, "id"), item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, item)
}

// Delete 删除角色
//
//	@Summary		删除角色
//	@Description	删除角色
//	@Tags			Role
//	@Param			id	path		int	true	"角色ID"
//	@Success		200	{object}	ginx.OperateResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/roles/:id [delete]
func (a *Role) Delete(c *gin.Context) {
	if err := a.RoleBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

// UserAddRole 用户新增角色
//
//	@Summary		用户新增角色
//	@Description	用户新增角色
//	@Tags			Role
//	@Param			id		path		int		true	"角色ID"
//	@Param			body	body		[]int	true	"用户ID集合"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/roles/:id/user-join [put]
func (a *Role) UserAddRole(c *gin.Context) {
	var userIDs []uint64
	if err := ginx.ParseJSON(c, &userIDs); err != nil {
		ginx.ResError(c, userIDs, err)
		return
	}

	if err := a.RoleBll.UserAddRole(c, ginx.ParseParamID(c, "id"), userIDs); err != nil {
		ginx.ResError(c, userIDs, err)
		return
	}
	ginx.ResOperateSuccess(c, userIDs)
}

// UserRemoveRole 移除用户的指定角色
//
//	@Summary		移除用户的指定角色
//	@Description	移除用户的指定角色
//	@Tags			Role
//	@Param			id		path		int		true	"角色ID"
//	@Param			body	body		[]int	true	"用户ID集合"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/roles/:id/user-remove [put]
func (a *Role) UserRemoveRole(c *gin.Context) {
	var userIDs []uint64
	if err := ginx.ParseJSON(c, &userIDs); err != nil {
		ginx.ResError(c, userIDs, err)
		return
	}

	if err := a.RoleBll.UserRemoveRole(c, ginx.ParseParamID(c, "id"), userIDs); err != nil {
		ginx.ResError(c, userIDs, err)
		return
	}
	ginx.ResOperateSuccess(c, userIDs)
}

// GetRolePermissionTree 获取角色的菜单树
//
//	@Summary		获取角色的菜单树
//	@Description	获取角色的菜单树
//	@Tags			Role
//	@Param			id	path		int	true	"角色ID"
//	@Success		200	{object}	schema.MenuPres
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/roles/:id/permission-tree [get]
func (a *Role) GetRolePermissionTree(c *gin.Context) {
	// 获取角色所有信息(包含菜单、按钮、禁用字段)
	item, err := a.RoleBll.GetPre(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}

	// 将菜单和按钮列表，组装成树
	menuTrees := make(schema.MenuPres, 0)
	parentId := uint64(0)
	if err = a.UserBll.GetMenuTreesInfo(&menuTrees, &parentId, item.Menus, item.Buttons); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", menuTrees.Init().SortMenuPres())
}
