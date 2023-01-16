package bll

import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

var UserExtendInfoBll = &UserExtendInfo{
	UserModel:           model.UserModel,
	UserExtendInfoModel: model.UserExtendInfoModel,
}

type UserExtendInfo struct {
	UserModel           *model.User
	UserExtendInfoModel *model.UserExtendInfo
}

func (a *UserExtendInfo) Get(c *gin.Context, id uint64) (*schema.UserExtendInfo, error) {
	item, err := a.UserExtendInfoModel.Get(id)
	if err != nil {
		return nil, errors.WithMessage(err, "获取用户扩展信息失败")
	} else if item == nil {
		return nil, errors.New("用户扩展信息不存在")
	}

	return item, nil
}

// Update 更新用户扩展信息
func (a *UserExtendInfo) Update(c *gin.Context, id uint64, item schema.UserExtendInfo) error {
	// 检查用户是否存在
	if err := a.ExtendInfoParamsCheck(c, &item); err != nil {
		return err
	}

	// 检查扩展信息是否存在
	if _, err := a.Get(c, id); err != nil {
		return err
	}

	// 更新用户扩展信息
	return a.UserExtendInfoModel.UpdateByID(id, map[string]interface{}{
		"real_name":    item.RealName,
		"mobile_phone": item.MobilePhone,
		"qq_account":   item.QQAccount,
		"email":        item.Email,
	})
}

// ---------------------------------------- Params  Validate --------------------------------------

func (a *UserExtendInfo) ExtendInfoParamsCheck(c *gin.Context, item *schema.UserExtendInfo) error {
	// 格式化信息
	item.RealName = strings.TrimSpace(item.RealName)
	item.Email = strings.TrimSpace(item.Email)

	// 检查用户是否存在
	if item.UserID == 0 {
		return errors.New("用户ID不能为0")
	}
	UserQueryResult, err := a.UserModel.Query(schema.UserQueryParam{
		ID:           item.UserID,
		OmitPassword: true,
	})
	if err != nil {
		return errors.WithMessage(err, "检查用户是否存在失败")
	} else if len(UserQueryResult.Data) == 0 {
		return errors.New("未找到该用户")
	}

	// 检查新的用户昵称是否已经存在
	if userExtendInfo, err := a.UserExtendInfoModel.GetByRealName(item.RealName); err != nil {
		return errors.New("检查用户昵称是否被使用失败")
	} else if userExtendInfo != nil && userExtendInfo.ID != item.ID {
		return errors.New("用户昵称已被使用")
	}

	return nil
}
