package api

import (
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"basic-frame/util/ginx"
	"github.com/gin-gonic/gin"
)

var OrganizationApi = &Organization{
	OrganizationBll: bll.OrganizationBll,
}

type Organization struct {
	OrganizationBll *bll.Organization
}

func (a *Organization) Query(c *gin.Context) {
	var params schema.OrganizationQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, params, err)
		return
	}

	if params.FindAll {
		params.Pagination = true
	}
	result, err := a.OrganizationBll.Query(c, params)
	if err != nil {
		ginx.ResError(c, params, err)
		return
	}

	ginx.ResList(c, params, result.Data)
}

func (a *Organization) Get(c *gin.Context) {
	item, err := a.OrganizationBll.Get(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

// Create 创建组织信息(职位没有ID且职位的组织ID为0，则创建职位)
func (a *Organization) Create(c *gin.Context) {
	var items schema.Organizations
	if err := ginx.ParseJSON(c, &items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	items.SetCreator(ginx.GetUserID(c))
	if err := a.OrganizationBll.Create(c, items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	ginx.ResOperateSuccess(c, items)
}

func (a *Organization) Update(c *gin.Context) {
	var item schema.Organization
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.OrganizationBll.Update(c, ginx.ParseParamID(c, "id"), item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, item)
}

func (a *Organization) UpdateOrganizations(c *gin.Context) {
	var items schema.Organizations
	if err := ginx.ParseJSON(c, &items); err != nil {
		ginx.ResError(c, items, err)
		return
	}

	if err := a.OrganizationBll.UpdateOrganizations(c, items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	ginx.ResOperateSuccess(c, items)
}

func (a *Organization) Delete(c *gin.Context) {
	if err := a.OrganizationBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

func (a *Organization) UserJoinOrganization(c *gin.Context) {
	var items []uint64
	if err := ginx.ParseJSON(c, &items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	if err := a.OrganizationBll.UserJoinOrganization(c, ginx.ParseParamID(c, "id"), items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	ginx.ResOperateSuccess(c, items)
}

func (a *Organization) UserRemoveOrganization(c *gin.Context) {
	var items []uint64
	if err := ginx.ParseJSON(c, &items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	if err := a.OrganizationBll.UserRemoveOrganization(c, ginx.ParseParamID(c, "id"), items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	ginx.ResOperateSuccess(c, items)
}

// ----------------------------------------OrganizationTree--------------------------------------

func (a *Organization) GetOrganizationTree(c *gin.Context) {
	Organizations, err := a.OrganizationBll.GetOrganizationTree(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", Organizations)
}

func (a *Organization) GetOrganizationTreeForCreateUser(c *gin.Context) {
	Organizations, err := a.OrganizationBll.GetOrganizationTreeForCreateUser(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", Organizations)
}

func (a *Organization) GetOrgTreeForCreateNotifications(c *gin.Context) {
	Organizations, err := a.OrganizationBll.GetOrgTreeForCreateNotifications(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", Organizations)
}

func (a *Organization) GetOrganizationTreeWithUser(c *gin.Context) {
	OrganizationTrees, err := a.OrganizationBll.GetOrganizationTreeWithUser(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", OrganizationTrees)
}
