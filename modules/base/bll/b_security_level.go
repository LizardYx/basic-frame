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

var SecurityLevelBll = &SecurityLevel{
	SecurityLevelModel: model.SecurityLevelModel,
	RoleModel:          model.RoleModel,
}

type SecurityLevel struct {
	SecurityLevelModel *model.SecurityLevel
	RoleModel          *model.Role
}

func (a *SecurityLevel) Query(c *gin.Context, params schema.SecurityLevelQueryParam) (*schema.SecurityLevelQueryResult, error) {
	return a.SecurityLevelModel.Query(params)
}

func (a *SecurityLevel) Get(c *gin.Context, id uint64) (*schema.SecurityLevel, error) {
	// 获取安全等级信息
	item, err := a.SecurityLevelModel.Get(id)
	if err != nil {
		return nil, errors.WithMessage(err, "获取安全等级信息失败")
	} else if item == nil {
		return nil, errors.New("未找到该安全等级")
	}

	return item, nil
}

func (a *SecurityLevel) GetUserSecurityLevels(c *gin.Context, userID uint64) (schema.SecurityLevels, error) {
	// 获取用户的角色ID集合
	// TODO: 等待用户表完成
	//userInfo, err := a.UserModel.Get(userID)
	//if err != nil {
	//	return nil, errors.WithMessage(err, "获取用户信息失败")
	//}

	var roleIDs []uint64
	//userInfo.Organizations.GetRoleIds(&roleIDs)
	//userInfo.Positions.GetRoleIds(&roleIDs)
	//userInfo.UserGroups.GetRoleIds(&roleIDs)
	//roleIDs = append(roleIDs, userInfo.Roles.GetIDs()...)

	// 获取安全级别集合
	securityLevels := make(schema.SecurityLevels, 0)
	if SecurityLevelQueryResult, err := a.Query(c, schema.SecurityLevelQueryParam{
		PaginationParam: common.PaginationParam{
			Pagination: false,
		},
		RoleIDs: common.UintSliceToString(roleIDs, ","),
		Status:  consts.BaseStatusEnable,
		FindAll: true,
	}); err != nil {
		return nil, errors.WithMessage(err, "获取用户的安全等级失败")
	} else if len(SecurityLevelQueryResult.Data) != 0 {
		securityLevels = append(securityLevels, SecurityLevelQueryResult.Data...)
	}
	return securityLevels, nil
}

func (a *SecurityLevel) Create(c *gin.Context, item schema.SecurityLevel) (*common.IDResult, error) {
	// 参数检查
	if err := a.SecurityLevelParamsCheck(c, &item); err != nil {
		return nil, err
	}
	return a.SecurityLevelModel.Create(item)
}

func (a *SecurityLevel) Update(c *gin.Context, id uint64, item schema.SecurityLevel) error {
	// 检查安全级别是否存在
	oldItem, err := a.SecurityLevelModel.Get(id)
	if err != nil {
		return errors.WithMessage(err, "获取安全等级信息失败")
	} else if oldItem == nil {
		return errors.New("未找到该安全等级")
	}

	// 参数检查
	if err = a.SecurityLevelParamsCheck(c, &item); err != nil {
		return err
	}

	// 更新安全级别和关联的角色信息
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 更新安全等级基本信息
		if err = a.SecurityLevelModel.UpdateByID(id, map[string]interface{}{
			"name":     item.Name,
			"sequence": item.Sequence,
			"status":   item.Status,
			"memo":     item.Memo,
		}); err != nil {
			return err
		}

		// 更新安全等级关联的角色
		if err = a.ReplaceSecurityLevelRoles(c, id, item.Roles); err != nil {
			return err
		}
		return nil
	})
}

// ReplaceSecurityLevelRoles 更新安全等级关联的角色
func (a *SecurityLevel) ReplaceSecurityLevelRoles(c *gin.Context, id uint64, items schema.Roles) error {
	// 检查安全级别是否存在
	if _, err := a.Get(c, id); err != nil {
		return err
	}
	return a.SecurityLevelModel.ReplaceSecurityLevelRoles(id, items)
}

func (a *SecurityLevel) Delete(c *gin.Context, id uint64) error {
	// 检查安全级别是否存在
	oldItem, err := a.SecurityLevelModel.Get(id)
	if err != nil {
		return errors.WithMessage(err, "获取安全等级信息失败")
	} else if oldItem == nil {
		return errors.New("未找到该安全等级")
	}

	// 判断安全级别是否被任务、项目、项目集等使用
	if err = a.SecurityLevelModel.SecurityLevelUsed(id); err != nil {
		return err
	}

	return a.SecurityLevelModel.Delete(id)
}

// ---------------------------------------- Params  Validate --------------------------------------

func (a *SecurityLevel) SecurityLevelParamsCheck(c *gin.Context, item *schema.SecurityLevel) error {
	// 检查名称是为空
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return errors.New("安全级别的名称不能为空")
	}
	// 检查角色是否被禁用
	RoleQueryResult, err := a.RoleModel.Query(schema.RoleQueryParam{
		PaginationParam: common.PaginationParam{
			Pagination: false,
		},
		IDs:     common.UintSliceToString(item.Roles.GetIDs(), ","),
		Status:  consts.BaseStatusDisabled,
		FindAll: true,
	})
	if err != nil {
		return errors.WithMessage(err, "检查角色是否被禁用失败")
	} else if len(RoleQueryResult.Data) != 0 {
		var roleNames string
		for index, roleInfo := range RoleQueryResult.Data {
			roleNames += roleInfo.Name
			if index != (len(RoleQueryResult.Data) - 1) {
				roleNames += ", "
			}
		}
		return errors.New(fmt.Sprintf("角色: %s 已被禁用", roleNames))
	}
	return nil
}
