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

func (a *Menu) Delete(c *gin.Context) {
	var params = ginx.ParamsID{ID: ginx.ParseParamID(c, "id")}

	if err := a.MenuBll.Delete(c, params.ID); err != nil {
		ginx.ResError(c, params, err)
		return
	}
	ginx.ResOperateSuccess(c, params)
}

// ----------------------------------------PermissionTree--------------------------------------

// GetPermissionTree 编辑权限树时调用的接口
func (a *Menu) GetPermissionTree(c *gin.Context) {
	permissionTree, err := a.MenuBll.GetPermissionTree(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", permissionTree)
}

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

// GetPermissionTreeForCreateRole 创建角色时调用的接口
func (a *Menu) GetPermissionTreeForCreateRole(c *gin.Context) {
	permissionTree, err := a.MenuBll.GetPermissionTreeForCreateRole(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", permissionTree)
}
