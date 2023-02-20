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
	UserModel:         model.UserModel,
}

type Organization struct {
	OrganizationModel *model.Organization
	RoleModel         *model.Role
	PositionModel     *model.Position
	UserModel         *model.User
}

func (a *Organization) Query(c *gin.Context, params schema.OrganizationQueryParam) (*schema.OrganizationQueryResult, error) {
	return a.OrganizationModel.Query(params)
}

// Get 获取组织的基本信息
func (a *Organization) Get(c *gin.Context, id uint64) (*schema.Organization, error) {
	// 获取组织基本信息
	item, err := a.OrganizationModel.Get(id)
	if err != nil {
		return nil, errors.WithMessage(err, "获取组织基本信息失败")
	} else if item == nil {
		return nil, errors.New("组织不存在")
	}

	return item, nil
}

// GetPre 获取组织、职位信息
func (a *Organization) GetPre(c *gin.Context, id uint64, includeSonOrg bool) (*schema.Organization, error) {
	// 获取组织、职位信息
	item, err := a.OrganizationModel.GetPre(id, includeSonOrg)
	if err != nil {
		return nil, errors.WithMessage(err, "获取组织信息失败")
	} else if item == nil {
		return nil, errors.New("组织不存在")
	}

	return item, nil
}

// Create 创建组织信息(职位没有ID且职位的组织ID为0，则创建职位)
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

// Update 更新组织的基础信息
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
	if _, err := a.Get(c, id); err != nil {
		return err
	}

	// 更新组织信息
	if *item.ParentID == 0 {
		item.ParentID = nil
	}
	if err := a.OrganizationModel.UpdateByID(id, map[string]interface{}{
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

// UserJoinOrganization 用户加入指定组织
func (a *Organization) UserJoinOrganization(c *gin.Context, OrgID uint64, userIDs []uint64) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 获取组织信息
		OrgInfo, err := a.Get(c, OrgID)
		if err != nil {
			return err
		}

		// 更新用户和组织的关联关系
		for _, userID := range userIDs {
			// 用户加入组织
			if err = a.UserModel.AppendUserOrganizations(userID, schema.Organizations{OrgInfo}); err != nil {
				return errors.WithMessage(err, "用户加入组织失败")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// UserRemoveOrganization 用户移出指定组织
func (a *Organization) UserRemoveOrganization(c *gin.Context, OrgID uint64, userIDs []uint64) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 获取组织信息
		OrgInfo, err := a.Get(c, OrgID)
		if err != nil {
			return err
		}

		// 将用户从组织中移出
		for _, userID := range userIDs {
			if err = a.UserModel.UserRemoveOrganization(userID, *OrgInfo); err != nil {
				return errors.WithMessage(err, "将用户从组织中移除失败")
			}
		}
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
	organizationInfo, err := a.GetPre(c, id, false)
	if err != nil {
		return err
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

	// 删除组织和职位
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 删除组织的职位
		if err = a.PositionModel.BatchDelete(organizationInfo.Positions.GetIDs()); err != nil {
			return errors.WithMessage(err, "删除组织的职位失败")
		}

		// 删除组织
		if err = a.OrganizationModel.Delete(id); err != nil {
			return errors.WithMessage(err, "删除组织失败")
		}
		return nil
	})
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// ----------------------------------------OrganizationTree--------------------------------------

// GetOrganizationTree 获取组织树、职位列表(包含禁用的组织、职位)
func (a *Organization) GetOrganizationTree(c *gin.Context) (*schema.Organizations, error) {
	Organizations, err := a.OrganizationModel.GetOrganizationTree()
	if err != nil {
		return &schema.Organizations{}, errors.WithMessage(err, "获取组织结构失败")
	}
	return Organizations.SortOrganizations(), nil
}

// GetOrganizationTreeForCreateUser 获取组织树、职位列表(不包含禁用的组织和职位)
func (a *Organization) GetOrganizationTreeForCreateUser(c *gin.Context) (*schema.Organizations, error) {
	Organizations, err := a.OrganizationModel.GetOrganizationTreeForCreateUser()
	if err != nil {
		return &schema.Organizations{}, errors.WithMessage(err, "获取组织结构失败")
	}
	return Organizations.SortOrganizations(), nil
}

// GetOrgTreeForCreateNotifications 获取组织树(不包含职位和禁用的组织)
func (a *Organization) GetOrgTreeForCreateNotifications(c *gin.Context) (*schema.Organizations, error) {
	Organizations, err := a.OrganizationModel.GetOrgTreeForCreateNotifications()
	if err != nil {
		return &schema.Organizations{}, errors.WithMessage(err, "获取组织结构失败")
	}
	return Organizations.SortOrganizations(), nil
}

// GetOrganizationTreeWithUser 获取包含用户的组织树、职位列表(不包含禁用的组织和职位)
func (a *Organization) GetOrganizationTreeWithUser(c *gin.Context) (*schema.OrganizationTrees, error) {
	// 获取组织树、职位列表(不包含禁用的组织和职位)
	Organizations, err := a.OrganizationModel.GetOrganizationTreeForCreateUser()
	if err != nil {
		return nil, errors.WithMessage(err, "获取组织结构失败")
	}

	// 获取所有用户
	UserQueryResult, err := a.UserModel.Query(schema.UserQueryParam{
		Status:       consts.BaseStatusEnable,
		SequenceSort: consts.BaseDescSort,
		ShowDetails:  true,
		OmitPassword: true,
		FindAll:      true,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "获取所有用户失败")
	}
	OrgTrees := Organizations.SortOrganizations().ToSchemaOrgTrees()
	OrgTrees.AddUserToTree(UserQueryResult.Data)
	return OrgTrees, nil
}

// ---------------------------------------- Params  Validate --------------------------------------

func (a *Organization) OrgParamsCheck(c *gin.Context, item schema.Organization) error {
	// 检查组织角色是否被禁用、是否为非组织类型角色
	var roleIds []uint64
	item.GetOrgRoleIDs(&roleIds)
	if len(roleIds) != 0 {
		if RoleQueryResult, err := a.RoleModel.Query(schema.RoleQueryParam{
			IDs:     common.UintSliceToString(roleIds, ","),
			FindAll: true,
		}); err != nil {
			return errors.WithMessage(err, "获取组织的角色信息失败")
		} else if len(RoleQueryResult.Data) != len(roleIds) {
			// 如果有未找到的角色ID
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
				return errors.New(fmt.Sprintf("角色ID: %s 的角色未找到", common.UintSliceToString(notFoundRoleId, ",")))
			}
		} else if len(RoleQueryResult.Data) != 0 {
			// 检查角色是否被禁用、是否为非组织类型
			var disabledRoleName []string
			var ErrTypeRoleName []string
			for _, role := range RoleQueryResult.Data {
				if role.Status == consts.BaseStatusDisabled {
					disabledRoleName = append(disabledRoleName, role.Name)
				}
				if role.Type != consts.RoleTypeForOrg {
					ErrTypeRoleName = append(ErrTypeRoleName, role.Name)
				}
			}
			if len(disabledRoleName) != 0 {
				return errors.New(fmt.Sprintf("角色: %s 已被禁用", common.StringSliceToString(disabledRoleName, ",")))
			}
			if len(ErrTypeRoleName) != 0 {
				return errors.New(fmt.Sprintf("角色: %s 的角色类型不是组织类型", common.StringSliceToString(ErrTypeRoleName, ",")))
			}
		}
	}

	// 检查职位的角色是否被禁用、是否为非职位类型的角色
	var positionRoleIDs []uint64
	item.GetPositionRoleIDs(&positionRoleIDs)
	if len(positionRoleIDs) != 0 {
		if roleQueryResult, err := a.RoleModel.Query(schema.RoleQueryParam{
			IDs:     common.UintSliceToString(positionRoleIDs, ","),
			FindAll: true,
		}); err != nil {
			return errors.WithMessage(err, "获取职位的角色信息失败")
		} else if len(roleQueryResult.Data) != len(positionRoleIDs) {
			// 如果有未找到的角色ID
			var notFoundRoleId []uint64
			for _, positionRoleID := range positionRoleIDs {
				for index, roleInfo := range roleQueryResult.Data {
					if roleInfo.ID == positionRoleID {
						continue
					}
					if index == (len(roleQueryResult.Data) - 1) {
						notFoundRoleId = append(notFoundRoleId, positionRoleID)
					}
				}
			}
			if len(notFoundRoleId) != 0 {
				return errors.New(fmt.Sprintf("角色ID: %s 的角色未找到", common.UintSliceToString(notFoundRoleId, ",")))
			}
		} else if len(roleQueryResult.Data) != 0 {
			// 检查角色是否被禁用、是否为非职位类型
			var disabledRoleName []string
			var ErrTypeRoleName []string
			for _, roleInfo := range roleQueryResult.Data {
				if roleInfo.Status == consts.BaseStatusDisabled {
					disabledRoleName = append(disabledRoleName, roleInfo.Name)
				}
				if roleInfo.Type != consts.RoleTypeForPosition {
					ErrTypeRoleName = append(ErrTypeRoleName, roleInfo.Name)
				}
			}
			if len(disabledRoleName) != 0 {
				return errors.New(fmt.Sprintf("角色: %s 已被禁用", common.StringSliceToString(disabledRoleName, ",")))
			}
			if len(ErrTypeRoleName) != 0 {
				return errors.New(fmt.Sprintf("角色: %s 的角色类型不是职位类型", common.StringSliceToString(ErrTypeRoleName, ",")))
			}
		}
	}

	// 检查职位是否被禁用
	var positionIds []uint64
	item.GetPositionIds(&positionIds)
	if len(positionIds) != 0 {
		if PositionQueryResult, err := a.PositionModel.Query(schema.PositionQueryParam{
			IDs:     common.UintSliceToString(positionIds, ","),
			Status:  consts.BaseStatusDisabled,
			FindAll: true,
		}); err != nil {
			return errors.WithMessage(err, "检查职位是否被禁用失败")
		} else if len(PositionQueryResult.Data) != 0 {
			var positionName []string
			for _, positionInfo := range PositionQueryResult.Data {
				positionName = append(positionName, positionInfo.Name)
			}
			return errors.New(fmt.Sprintf("职位: %s 已被禁用", common.StringSliceToString(positionName, ",")))
		}
	}
	return nil
}

func (a *Organization) OrgParamsValidate(item *schema.Organization) error {
	// 组织名称检查
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return errors.New("组织名称不能为空")
	}

	// 组织职位名称检查
	for index, position := range item.Positions {
		positionName := strings.TrimSpace(position.Name)
		if positionName == "" {
			return errors.New("职位名称不能为空")
		}
		item.Positions[index].Name = positionName
	}

	// 检查子组织
	if len(item.SonOrganizations) != 0 {
		for _, SonOrg := range item.SonOrganizations {
			if err := a.OrgParamsValidate(SonOrg); err != nil {
				return err
			}
		}
	}
	return nil
}
