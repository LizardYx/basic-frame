package bll

import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"basic-frame/util/mysql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

var UserGroupBll = &UserGroup{
	RoleModel:      model.RoleModel,
	UserGroupModel: model.UserGroupModel,
	UserModel:      model.UserModel,
}

type UserGroup struct {
	RoleModel      *model.Role
	UserGroupModel *model.UserGroup
	UserModel      *model.User
}

func (a *UserGroup) Query(c *gin.Context, params schema.UserGroupQueryParam) (*schema.UserGroupQueryResult, error) {
	return a.UserGroupModel.Query(params)
}

// Get 获取用户组基本信息
func (a *UserGroup) Get(c *gin.Context, id uint64) (*schema.UserGroup, error) {
	item, err := a.UserGroupModel.Get(id)
	if err != nil {
		return nil, errors.WithMessage(err, "获取用户组信息失败")
	} else if item == nil {
		return nil, errors.New("未找到该用户组")
	}

	return item, nil
}

// GetUserGroupUsers 获取指定用户组的用户列表
func (a *UserGroup) GetUserGroupUsers(c *gin.Context, id uint64) (*schema.Users, error) {
	UserQueryResult, err := a.UserModel.Query(schema.UserQueryParam{
		UserGroupID:  id,
		Status:       consts.BaseStatusEnable,
		OmitPassword: true,
		FindAll:      true,
	})
	if err != nil {
		return nil, err
	}
	return &UserQueryResult.Data, nil
}

// Create 创建用户组
func (a *UserGroup) Create(c *gin.Context, item schema.UserGroup) (*common.IDResult, error) {
	// 检查角色是否被禁用
	if err := a.UserGroupParamsCheck(c, &item); err != nil {
		return &common.IDResult{}, err
	}

	// 创建用户组
	IDResult, err := a.UserGroupModel.Create(item)
	if err != nil {
		return nil, err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return IDResult, nil
}

// Update 更新用户组基本信息
func (a *UserGroup) Update(c *gin.Context, id uint64, item schema.UserGroup) error {
	// 检查用户组是否存在
	if userGroupInfo, err := a.Get(c, id); err != nil {
		return err
	} else {
		item.ID = userGroupInfo.ID
	}

	// 检查角色是否被禁用
	if err := a.UserGroupParamsCheck(c, &item); err != nil {
		return err
	}

	// 更新用户信息
	if err := a.UserGroupModel.UpdateByID(id, map[string]interface{}{
		"name":     item.Name,
		"role_id":  item.RoleID,
		"sequence": item.Sequence,
		"status":   item.Status,
		"memo":     item.Memo,
	}); err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// UserJoinUserGroup 多个用户加入指定用户组
func (a *UserGroup) UserJoinUserGroup(c *gin.Context, userGroupID uint64, userIDs []uint64) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 获取用户组信息
		userGroupInfo, err := a.Get(c, userGroupID)
		if err != nil {
			return err
		}

		// 更新用户和用户组的关联关系
		for _, userID := range userIDs {
			var userInfo *schema.User
			userInfo, err = a.UserModel.Get(userID)
			if err != nil {
				return err
			}
			// 检查用户是否已加入该用户组
			if common.ContainsUint64(userInfo.UserGroups.GetIDs(), userGroupID) {
				return errors.New(fmt.Sprintf("用户: %s 已加入该用户组", userInfo.UserName))
			}

			// 用户加入用户组
			if err = a.UserModel.AppendUserUserGroup(userID, schema.UserGroups{userGroupInfo}); err != nil {
				return err
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

// UserGroupRemoveUser 从用户组中移出多个用户
func (a *UserGroup) UserGroupRemoveUser(c *gin.Context, userGroupID uint64, userIDs []uint64) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 获取用户组信息
		userGroupInfo, err := a.Get(c, userGroupID)
		if err != nil {
			return err
		}

		// 将用户从用户组中移除
		for _, userID := range userIDs {
			// 获取用户信息
			var userInfo *schema.User
			if userInfo, err = a.UserModel.Get(userID); err != nil {
				return err
			} else if len(userInfo.UserGroups) != 0 {
				if !common.ContainsUint64(userInfo.UserGroups.GetIDs(), userGroupID) {
					// 如果用户没有加入该用户组
					return errors.New(fmt.Sprintf("用户: %s 未加入到用户组 %s，无法移出", userInfo.UserName, userGroupInfo.Name))
				}
			}

			// 将用户从用户组中移除
			if err = a.UserModel.UserRemoveUserGroup(userID, *userGroupInfo); err != nil {
				return err
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

// Delete 删除用户组
func (a *UserGroup) Delete(c *gin.Context, id uint64) error {
	// 检查用户组是否存在
	if _, err := a.Get(c, id); err != nil {
		return err
	}

	// 检查用户组是否正在被用户使用
	if userQueryResult, err := a.UserModel.Query(schema.UserQueryParam{
		UserGroupID:  id,
		OmitPassword: true,
		FindAll:      true,
	}); err != nil {
		return errors.WithMessage(err, "检查用户组是否被使用失败")
	} else if len(userQueryResult.Data) != 0 {
		return errors.New(fmt.Sprintf("该用户组正在被用户: %s 使用", common.StringSliceToString(userQueryResult.Data.GetUserNames(), ",")))
	}

	// 删除用户组
	if err := a.UserGroupModel.Delete(id); err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// ---------------------------------------- Params  Validate --------------------------------------

func (a *UserGroup) UserGroupParamsCheck(c *gin.Context, item *schema.UserGroup) error {
	// 检查用户组名称是否为空
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return errors.New("用户组名称不能为空")
	}

	// 检查用户组名称是否已经存在
	if UserGroupQueryResult, err := a.Query(c, schema.UserGroupQueryParam{
		Name: item.Name,
	}); err != nil {
		return err
	} else if len(UserGroupQueryResult.Data) != 0 {
		if UserGroupQueryResult.Data[0].ID != item.ID {
			return errors.New("该用户组名称已被使用")
		}
	}

	// 检查角色是否可用
	if item.RoleID != 0 {
		if RoleQueryResult, err := a.RoleModel.Query(schema.RoleQueryParam{
			ID: item.RoleID,
		}); err != nil {
			return err
		} else if len(RoleQueryResult.Data) != 0 {
			if RoleQueryResult.Data[0].Status == consts.BaseStatusDisabled {
				return errors.New("该角色已被禁用")
			}
			if RoleQueryResult.Data[0].Type != consts.RoleTypeForUserGroup {
				return errors.New("角色类型错误")
			}
		} else if len(RoleQueryResult.Data) == 0 {
			return errors.New("未找到该角色")
		}
	}
	return nil
}
