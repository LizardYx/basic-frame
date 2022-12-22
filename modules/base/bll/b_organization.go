package bll

import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"basic-frame/util/ginx/errors"
	"basic-frame/util/mysql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

var OrganizationBll = &Organization{
	OrganizationModel: model.OrganizationModel,
	RoleModel:         model.RoleModel,
	PositionModel:     model.PositionModel,
}

type Organization struct {
	OrganizationModel *model.Organization
	RoleModel         *model.Role
	PositionModel     *model.Position
}

func (a *Organization) Query(c *gin.Context, params schema.OrganizationQueryParam) (*schema.OrganizationQueryResult, error) {
	return a.OrganizationModel.Query(params)
}

func (a *Organization) Get(c *gin.Context, id uint64) (*schema.Organization, error) {
	// 获取组织详情
	item, err := a.OrganizationModel.Get(id)
	if err != nil {
		return nil, errors.WithMessage(err, "获取组织信息失败")
	} else if item == nil {
		return nil, errors.New("组织不存在")
	}

	return item.SortOrganization(), nil
}

func (a *Organization) Create(c *gin.Context, items schema.Organizations) error {
	// 参数角色和职位是否被禁用
	for index, item := range items {
		if err := a.OrgParamsValidate(items[index]); err != nil {
			return err
		}
		if err := a.OrgParamsCheck(c, *item); err != nil {
			return err
		}
	}
	if err := a.OrganizationModel.CreateOrganizations(items); err != nil {
		return errors.WithMessage(err, "创建组织失败")
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

func (a *Organization) Update(c *gin.Context, id uint64, item schema.Organization) error {
	// 检查参数正确性
	if err := a.OrgParamsValidate(&item); err != nil {
		return err
	}

	// 参数角色和职位是否被禁用
	if err := a.OrgParamsCheck(c, item); err != nil {
		return err
	}

	// 检查组织是否存在
	oldItem, err := a.OrganizationModel.Get(id)
	if err != nil {
		return errors.WithMessage(err, "获取组织信息失败")
	} else if oldItem == nil {
		return errors.New("组织不存在")
	}

	// 更新组织信息
	if *item.ParentID == 0 {
		item.ParentID = nil
	}
	if err = a.OrganizationModel.UpdateByID(id, map[string]interface{}{
		"name":      item.Name,
		"role_id":   item.RoleID,
		"sequence":  item.Sequence,
		"parent_id": item.ParentID,
		"status":    item.Status,
		"memo":      item.Memo,
	}); err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

func (a *Organization) UpdateOrganizations(c *gin.Context, items schema.Organizations) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := a.Update(c, item.ID, *item); err != nil {
				return err
			}
		}
		return nil
	})
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return err
}

func (a *Organization) UserJoinOrganization(c *gin.Context, OrgID uint64, userIDs []uint64) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 获取用户组信息
		OrgInfo, err := a.Get(c, OrgID)
		if err != nil {
			return errors.WithMessage(err, "获取组织信息失败")
		} else if OrgInfo == nil {
			return errors.New("组织不存在")
		}

		// 更新用户和组织的关联关系
		// TODO: 等待用户表完成
		//for _, userID := range userIDs {
		//	// 用户加入组织
		//	if err = a.UserModel.AppendUserOrganizations(userID, schema.Organizations{OrgInfo}); err != nil {
		//		return errors.WithMessage(err, "用户加入组织失败")
		//	}
		//}
		return nil
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

func (a *Organization) UserRemoveOrganization(c *gin.Context, OrgID uint64, userIDs []uint64) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 获取组织信息
		OrgInfo, err := a.Get(c, OrgID)
		if err != nil {
			return errors.WithMessage(err, "获取组织信息失败")
		} else if OrgInfo == nil {
			return errors.New("组织不存在")
		}

		// 将用户从组织中移除
		// TODO: 等待用户表完成
		//for _, userID := range userIDs {
		//	if err = a.UserModel.UserRemoveOrganization(userID, *OrgInfo); err != nil {
		//		return errors.WithMessage(err, "将用户从组织中移除失败")
		//	}
		//}
		return nil
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// Delete 删除组织和组织的职位
func (a *Organization) Delete(c *gin.Context, id uint64) error {
	// 检查组织是否存在
	oldItem, err := a.OrganizationModel.Get(id)
	if err != nil {
		return errors.WithMessage(err, "获取组织信息失败")
	} else if oldItem == nil {
		return errors.New("组织不存在")
	}

	// 检查组织是否有子项
	if OrganizationQueryResult, err := a.OrganizationModel.Query(schema.OrganizationQueryParam{
		PaginationParam: common.PaginationParam{
			OnlyCount: true,
		},
		ParentID: id,
	}); err != nil {
		return errors.WithMessage(err, "检查组织是否有子项失败")
	} else if OrganizationQueryResult.PageResult.Total != 0 {
		return errors.New("该组织有子级，请勿删除")
	}

	// 删除组织
	if err = a.OrganizationModel.Delete(id); err != nil {
		return errors.WithMessage(err, "删除组织失败")
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// ----------------------------------------OrganizationTree--------------------------------------

func (a *Organization) GetOrganizationTree(c *gin.Context) (*schema.Organizations, error) {
	Organizations, err := a.OrganizationModel.GetOrganizationTree()
	if err != nil {
		return &schema.Organizations{}, errors.WithMessage(err, "获取组织结构失败")
	}
	return Organizations.SortOrganizations(), nil
}

func (a *Organization) GetOrganizationTreeForCreateUser(c *gin.Context) (*schema.Organizations, error) {
	Organizations, err := a.OrganizationModel.GetOrganizationTreeForCreateUser()
	if err != nil {
		return &schema.Organizations{}, errors.WithMessage(err, "获取组织结构失败")
	}
	return Organizations.SortOrganizations(), nil
}

func (a *Organization) GetOrgTreeForCreateNotifications(c *gin.Context) (*schema.Organizations, error) {
	Organizations, err := a.OrganizationModel.GetOrgTreeForCreateNotifications()
	if err != nil {
		return &schema.Organizations{}, errors.WithMessage(err, "获取组织结构失败")
	}
	return Organizations.SortOrganizations(), nil
}

func (a *Organization) GetOrganizationTreeWithUser(c *gin.Context) (*schema.OrganizationTrees, error) {
	// 获取组织、职位树
	Organizations, err := a.OrganizationModel.GetOrganizationTreeForCreateUser()
	if err != nil {
		return nil, errors.WithMessage(err, "获取组织结构失败")
	}
	// 获取所有用户
	// TODO: 等待用户表完成
	//UserQueryResult, err := a.UserModel.Query(schema.UserQueryParam{
	//	PaginationParam: common.PaginationParam{
	//		Pagination: false,
	//	},
	//	Status:       consts.BaseStatusEnable,
	//	ShowDetails:  true,
	//	SequenceSort: 2,
	//	FindAll:      true,
	//}, true)
	//if err != nil {
	//	return nil, errors.WithMessage(err, "获取所有用户失败")
	//}
	OrgTrees := Organizations.SortOrganizations().ToSchemaOrgTrees()
	//OrgTrees.AddUserToTree(UserQueryResult.Data)
	return OrgTrees, nil
}

// ---------------------------------------- Params  Validate --------------------------------------

func (a *Organization) OrgParamsCheck(c *gin.Context, item schema.Organization) error {
	// 检查角色是否被禁用
	var roleIds []uint64
	item.GetRoleIds(&roleIds, true)
	if len(roleIds) != 0 {
		if RoleQueryResult, err := a.RoleModel.Query(schema.RoleQueryParam{
			PaginationParam: common.PaginationParam{},
			IDs:             common.UintSliceToString(roleIds, ","),
		}); err != nil {
			return errors.WithMessage(err, "检查角色是否被禁用失败")
		} else if len(RoleQueryResult.Data) != 0 {
			var disabledRoleName string
			var ErrTypeRoleName string
			for index, role := range RoleQueryResult.Data {
				if role.Status == consts.BaseStatusDisabled {
					disabledRoleName += role.Name
					if index != (len(RoleQueryResult.Data) - 1) {
						disabledRoleName += ", "
					}
				}
				if role.Type != consts.RoleTypeForOrg {
					ErrTypeRoleName += role.Name
					if index != (len(RoleQueryResult.Data) - 1) {
						disabledRoleName += ", "
					}
				}
			}
			if disabledRoleName != "" {
				return errors.New(fmt.Sprintf("角色: %s 已被禁用", disabledRoleName))
			}
			if ErrTypeRoleName != "" {
				return errors.New(fmt.Sprintf("角色: %s 的角色类型错误", disabledRoleName))
			}
		} else if len(RoleQueryResult.Data) != len(roleIds) {
			var notFoundRoleId []uint64
			for _, roleId := range roleIds {
				for index, roleInfo := range RoleQueryResult.Data {
					if roleInfo.ID == roleId {
						continue
					}
					if index == (len(RoleQueryResult.Data) - 1) {
						notFoundRoleId = append(notFoundRoleId, roleId)
					}
				}
			}
			if len(notFoundRoleId) != 0 {
				return errors.New(fmt.Sprintf("角色ID: %d 的角色未找到", notFoundRoleId))
			}
		}
	}
	// 检查职位是否被禁用
	var positionIds []uint64
	item.GetPositionIds(&positionIds)
	if len(positionIds) != 0 {
		if PositionQueryResult, err := a.PositionModel.Query(schema.PositionQueryParam{
			PaginationParam: common.PaginationParam{},
			IDs:             common.UintSliceToString(positionIds, ","),
			Status:          consts.BaseStatusDisabled,
		}); err != nil {
			return errors.WithMessage(err, "检查职位是否被禁用失败")
		} else if len(PositionQueryResult.Data) != 0 {
			var positionName string
			for index, positionInfo := range PositionQueryResult.Data {
				positionName += positionInfo.Name
				if index != (len(PositionQueryResult.Data) - 1) {
					positionName += ", "
				}
			}
			return errors.New(fmt.Sprintf("职位: %s 已被禁用", positionName))
		}
	}
	return nil
}

func (a *Organization) OrgParamsValidate(item *schema.Organization) error {
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return errors.New("组织名称不能为空")
	}
	for index, position := range item.Positions {
		positionName := strings.TrimSpace(position.Name)
		if positionName == "" {
			return errors.New("职位名称不能为空")
		}
		item.Positions[index].Name = positionName
	}
	if len(item.SonOrganizations) != 0 {
		for _, SonOrg := range item.SonOrganizations {
			if err := a.OrgParamsValidate(SonOrg); err != nil {
				return err
			}
		}
	}
	return nil
}
