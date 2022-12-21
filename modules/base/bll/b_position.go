package bll

import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"basic-frame/util/mysql"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

var PositionBll = &Position{
	PositionModel: model.PositionModel,
	RoleModel:     model.RoleModel,
}

type Position struct {
	PositionModel *model.Position
	RoleModel     *model.Role
}

func (a *Position) Query(c *gin.Context, params schema.PositionQueryParam) (*schema.PositionQueryResult, error) {
	return a.PositionModel.Query(params)
}

func (a *Position) Get(c *gin.Context, id uint64) (*schema.Position, error) {
	item, err := a.PositionModel.Get(id)
	if err != nil {
		return nil, errors.WithMessage(err, "检查职位是否存在失败")
	} else if item == nil {
		return nil, errors.New("未找到该职位")
	}

	return item, nil
}

func (a *Position) Create(c *gin.Context, item schema.Position) (*common.IDResult, error) {
	// 检查角色ID和组织ID是否被禁用
	if err := a.PositionParamsCheck(c, &item); err != nil {
		return nil, err
	}
	IDResult, err := a.PositionModel.Create(item)
	if err != nil {
		return nil, err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return IDResult, nil
}

func (a *Position) Update(c *gin.Context, id uint64, item schema.Position) error {
	// 检查职位是否存在
	oldItem, err := a.PositionModel.Get(id)
	if err != nil {
		return errors.WithMessage(err, "检查职位是否存在失败")
	} else if oldItem == nil {
		return errors.New("未找到该职位")
	} else {
		item.OrganizationID = oldItem.OrganizationID
	}

	// 检查角色ID和组织ID是否被禁用
	if err = a.PositionParamsCheck(c, &item); err != nil {
		return err
	}

	// 更新职位信息
	if err = a.PositionModel.UpdateByID(id, map[string]interface{}{
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

func (a *Position) PositionAddUser(c *gin.Context, positionID uint64, userIDs []uint64) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 获取职位信息
		//positionInfo, err := a.Get(c, positionID)
		//if err != nil {
		//	return err
		//}

		// 更新用户和职位的关联关系
		// TODO: 等待用户表完成
		//for _, userID := range userIDs {
		//	if err = a.UserModel.AppendUserPositions(userID, schema.Positions{positionInfo}); err != nil {
		//		return err
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

func (a *Position) PositionRemoveUser(c *gin.Context, positionID uint64, userIDs []uint64) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 获取职位信息
		//positionInfo, err := a.Get(c, positionID)
		//if err != nil {
		//	return err
		//}

		// 更新用户和职位的关联关系
		// TODO: 等待用户表完成
		//for _, userID := range userIDs {
		//	if err = a.UserModel.UserRemovePosition(userID, *positionInfo); err != nil {
		//		return err
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

func (a *Position) Delete(c *gin.Context, id uint64) error {
	// 检查职位是否存在
	oldItem, err := a.PositionModel.Get(id)
	if err != nil {
		return errors.WithMessage(err, "检查职位是否存在失败")
	} else if oldItem == nil {
		return errors.New("未找到该职位")
	}

	// 删除角色
	if err = a.PositionModel.Delete(id); err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// ---------------------------------------- Params  Validate --------------------------------------

func (a *Position) PositionParamsCheck(c *gin.Context, item *schema.Position) error {
	// 检查职位名称
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return errors.New("职位名称不能为空")
	}

	// 检查职位名称是否存在
	if PositionQueryResult, err := a.Query(c, schema.PositionQueryParam{
		Name: item.Name,
	}); err != nil {
		return errors.WithMessage(err, "检查职位名称是否存在失败")
	} else if len(PositionQueryResult.Data) != 0 {
		if PositionQueryResult.Data[0].ID != item.ID {
			return errors.New("该职位名称已经存在")
		}
	}

	// 检查角色是否可用
	if item.RoleID != 0 {
		if RoleQueryResult, err := a.RoleModel.Query(schema.RoleQueryParam{
			PaginationParam: common.PaginationParam{},
			ID:              item.RoleID,
			Status:          consts.BaseStatusEnable,
		}); err != nil {
			return errors.WithMessage(err, "检查角色是否可用失败")
		} else if len(RoleQueryResult.Data) == 0 {
			return errors.New("该角色已被禁用")
		} else if RoleQueryResult.Data[0].Type != consts.RoleTypeForPosition {
			return errors.New("角色类型错误")
		}
	}
	// 检查组织是否可用
	// TODO: 等待组织表完成
	//if item.OrganizationID != 0 {
	//	if OrganizationQueryResult, err := a.OrganizationModel.Query(schema.OrganizationQueryParam{
	//		PaginationParam: common.PaginationParam{
	//			OnlyCount: true,
	//		},
	//		ID:     item.OrganizationID,
	//		Status: consts.BaseStatusDisabled,
	//	}); err != nil {
	//		return err
	//	} else if OrganizationQueryResult.PageResult.Total != 0 {
	//		return errors.New("该组织已被禁用")
	//	}
	//}
	return nil
}
