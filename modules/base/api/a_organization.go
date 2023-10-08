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

// Query 查询组织列表
//
//	@Summary		查询组织列表
//	@Description	查询组织列表
//	@Tags			Organization
//	@Param			q	query		schema.OrganizationQueryParam	false	"查询组织列表参数"
//	@Success		200	{object}	ginx.ListResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/organizations [get]
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

// Create 创建组织
//
//	@Summary		创建组织
//	@Description	创建组织信息(职位没有ID且职位的组织ID为0，则创建职位)
//	@Tags			Organization
//	@Param			body	body		schema.Organizations	true	"创建组织参数"
//	@Success		200		{object}	common.IDResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/organizations [post]
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

// Update 更新组织的基础信息
//
//	@Summary		更新组织的基础信息
//	@Description	更新组织的基础信息
//	@Tags			Organization
//	@Param			id		path		int					true	"组织ID"
//	@Param			body	body		schema.Organization	true	"组织基础信息参数"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/organizations/update-basic-info/:id [put]
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

// Delete 删除组织
//
//	@Summary		删除组织
//	@Description	删除组织和组织的职位
//	@Tags			Organization
//	@Param			id	path		int	true	"组织ID"
//	@Success		200	{object}	ginx.OperateResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/organizations/:id [delete]
func (a *Organization) Delete(c *gin.Context) {
	if err := a.OrganizationBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

// UserJoinOrganization 用户加入指定组织
//
//	@Summary		用户加入指定组织
//	@Description	用户加入指定组织
//	@Tags			Organization
//	@Param			id		path		int		true	"组织ID"
//	@Param			body	body		[]int	true	"用户ID参数集合"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/organizations/user-join/:id [put]
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

// UserRemoveOrganization 用户移出指定组织
//
//	@Summary		用户移出指定组织
//	@Description	用户移出指定组织
//	@Tags			Organization
//	@Param			id		path		int		true	"组织ID"
//	@Param			body	body		[]int	true	"用户D参数集合"
//	@Success		200		{object}	ginx.OperateResult
//	@Failure		500		{object}	ginx.OperateResult
//	@Router			/api/v1/base/organizations/user-remove/:id [put]
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

// GetOrganizationTree 获取可编辑的组织树
//
//	@Summary		获取可编辑的组织树
//	@Description	获取组织树、职位列表(包含禁用的组织、职位)
//	@Tags			Organization
//	@Success		200	{object}	schema.Organizations
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/organizations/edit [get]
func (a *Organization) GetOrganizationTree(c *gin.Context) {
	Organizations, err := a.OrganizationBll.GetOrganizationTree(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", Organizations)
}

// GetOrganizationTreeForCreateUser 获取创建用户的组织树
//
//	@Summary		获取创建用户的组织树
//	@Description	获取组织树、职位列表(不包含禁用的组织和职位)
//	@Tags			Organization
//	@Success		200	{object}	schema.Organizations
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/organizations/create-user [get]
func (a *Organization) GetOrganizationTreeForCreateUser(c *gin.Context) {
	Organizations, err := a.OrganizationBll.GetOrganizationTreeForCreateUser(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", Organizations)
}

// GetOrgTreeForCreateNotifications 获取创建通知公告的组织树
//
//	@Summary		获取创建通知公告的组织树
//	@Description	获取组织树(不包含职位和禁用的组织)
//	@Tags			Organization
//	@Success		200	{object}	schema.Organizations
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/organizations/create-notifications [get]
func (a *Organization) GetOrgTreeForCreateNotifications(c *gin.Context) {
	Organizations, err := a.OrganizationBll.GetOrgTreeForCreateNotifications(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", Organizations)
}

// GetOrganizationTreeWithUser 获取人员组织架构树
//
//	@Summary		获取人员组织架构树
//	@Description	获取包含用户的组织树、职位列表(不包含禁用的组织和职位)
//	@Tags			Organization
//	@Success		200	{object}	schema.OrganizationTrees
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/organizations/include-user [get]
func (a *Organization) GetOrganizationTreeWithUser(c *gin.Context) {
	OrganizationTrees, err := a.OrganizationBll.GetOrganizationTreeWithUser(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", OrganizationTrees)
}
