package api

import (
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/ginx"
	"github.com/gin-gonic/gin"
	"strings"
)

var MenuApi = &Menu{
	MenuBll: bll.MenuBll,
}

type Menu struct {
	MenuBll *bll.Menu
}

// BatchUpdateMenus 批量更新菜单基本信息
//
//	@Summary		批量更新菜单基本信息
//	@Description	批量更新菜单基本信息
//	@Tags			Menu
//	@Param			body	body		schema.Menus	true	"菜单信息参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/menu//basic-info-update [put]
func (a *Menu) BatchUpdateMenus(c *gin.Context) {
	var items schema.Menus
	if err := ginx.ParseJSON(c, &items); err != nil {
		ginx.ResError(c, items, err)
		return
	}

	if err := a.MenuBll.BatchUpdateMenus(c, items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	ginx.ResOperateSuccess(c, items)
}

// Delete 删除菜单
//
//	@Summary		删除菜单
//	@Description	删除菜单、菜单的Api、菜单按钮、菜单按钮的Api
//	@Tags			Menu
//	@Param			id	path		int	true	"菜单ID"
//	@Success		200	{object}	ginx.OperateResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/menu/:id [delete]
func (a *Menu) Delete(c *gin.Context) {
	if err := a.MenuBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

// ----------------------------------------PermissionTree--------------------------------------

// GetPermissionTree 获取编辑菜单的菜单树
//
//	@Summary		获取编辑菜单的菜单树
//	@Description	获取所有的菜单树、特殊接口(包含禁用)
//	@Tags			Menu
//	@Success		200	{object}	schema.PermissionTree
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/menu/edit [get]
func (a *Menu) GetPermissionTree(c *gin.Context) {
	permissionTree, err := a.MenuBll.GetPermissionTree(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", permissionTree)
}

// UpdatePermissionTree 更新菜单树
//
//	@Summary		更新菜单树
//	@Description	菜单、按钮及可禁用字段的基础信息(创建、更新)。菜单和按钮调用的api(创建、更新、删除)
//	@Tags			Menu
//	@Param			body	body		schema.PermissionTree	true	"更新菜单树参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/menu [put]
func (a *Menu) UpdatePermissionTree(c *gin.Context) {
	var item schema.PermissionTree
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.MenuBll.UpdatePermissionTree(c, item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, item)
}

// DownloadPermissionTree 下载菜单树
//
//	@Summary		下载菜单树
//	@Description	下载菜单树
//	@Tags			Menu
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/menu/download [get]
func (a *Menu) DownloadPermissionTree(c *gin.Context) {
	// 更新配置文件内容
	if err := a.MenuBll.UpdatePermissionTreeFile(c); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	filePath := common.SysConfig.MenuFile
	idx := strings.LastIndex(filePath, "/")
	fileName := filePath[idx+1:]
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(filePath)
}

// GetPermissionTreeForCreateRole 获取创建角色的菜单树
//
//	@Summary		获取创建角色的菜单树
//	@Description	获取所有的菜单树、按钮树、特殊接口(不包含禁用的菜单和按钮)
//	@Tags			Menu
//	@Success		200	{object}	schema.PermissionTree
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/menu/create-role [get]
func (a *Menu) GetPermissionTreeForCreateRole(c *gin.Context) {
	permissionTree, err := a.MenuBll.GetPermissionTreeForCreateRole(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", permissionTree)
}
