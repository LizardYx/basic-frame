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

func (a *Role) Get(c *gin.Context) {
	item, err := a.RoleBll.GetPreWithSelect(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

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

func (a *Role) Delete(c *gin.Context) {
	if err := a.RoleBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

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

func (a *Role) GetRolePermissionTree(c *gin.Context) {
	// 获取角色所有信息
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
