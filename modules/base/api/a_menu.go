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

// Delete 删除菜单、菜单的Api、菜单按钮、菜单按钮的Api
func (a *Menu) Delete(c *gin.Context) {
	if err := a.MenuBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

// ----------------------------------------PermissionTree--------------------------------------

// GetPermissionTree 获取所有的菜单树、特殊接口(包含禁用)
func (a *Menu) GetPermissionTree(c *gin.Context) {
	permissionTree, err := a.MenuBll.GetPermissionTree(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", permissionTree)
}

// UpdatePermissionTree 菜单、按钮及可禁用字段的基础信息(创建、更新)。菜单和按钮调用的api(创建、更新、删除)
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

// GetPermissionTreeForCreateRole 获取所有的菜单树、按钮树、特殊接口(不包含禁用的菜单和按钮)
func (a *Menu) GetPermissionTreeForCreateRole(c *gin.Context) {
	permissionTree, err := a.MenuBll.GetPermissionTreeForCreateRole(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", permissionTree)
}
