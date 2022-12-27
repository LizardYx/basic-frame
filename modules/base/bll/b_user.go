package bll

import (
	"basic-frame/middleware"
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"basic-frame/util/ginx/errors"
	"basic-frame/util/secret"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

var UserBll = &User{
	UserModel: model.UserModel,
	RoleModel: model.RoleModel,
}

type User struct {
	UserModel *model.User
	RoleModel *model.Role
}

func (a *User) Login(c *gin.Context, item schema.LoginParam) (*schema.LoginTokenInfo, *schema.User, error) {
	// 验证用户名密码
	user, err := a.Verify(c, item.UserName, item.Password)
	if err != nil {
		return nil, nil, err
	}
	user.Password = ""

	// 生成令牌
	newTokenInfo := schema.LoginTokenInfo{}
	if token, err := middleware.GenerateToken(user.ID, user.UserName); err != nil {
		return nil, nil, errors.WithMessage(err, "生成令牌失败")
	} else {
		expiredTime, _ := strconv.Atoi(consts.JWTAuthExpired)
		newTokenInfo = schema.LoginTokenInfo{
			AccessToken: token,
			TokenType:   consts.JWTAuthType,
			ExpiresAt:   time.Now().Add(time.Duration(expiredTime) * time.Second).Unix(),
		}
	}
	return &newTokenInfo, user, nil
}

func (a *User) Query(c *gin.Context, params schema.UserQueryParam) (*schema.UserQueryResult, error) {
	// 如果要查询某种审核类型的用户，先获取审核角色ID
	if params.AuditorTypes != "" {
		if roleQueryResult, err := a.RoleModel.Query(schema.RoleQueryParam{
			AuditorTypes: params.AuditorTypes,
			Status:       consts.BaseStatusEnable,
			FindAll:      true,
		}); err != nil {
			return nil, errors.WithMessage(err, "获取用户列表失败")
		} else if len(roleQueryResult.Data) != 0 {
			var roleIDs []uint64
			if params.RoleIDs != "" {
				roleIDs = common.SplitStringToUint64(params.RoleIDs, ",")
			}
			roleIDs = append(roleIDs, roleQueryResult.Data.GetIDs()...)
			params.RoleIDs = common.UintSliceToString(roleIDs, ",")
		}
	}

	return a.UserModel.Query(params)
}

// ---------------------------------------- Params Validate --------------------------------------

// Verify 检查用户是否存在
func (a *User) Verify(c *gin.Context, userName, password string) (*schema.User, error) {
	UserQueryResult, err := a.Query(c, schema.UserQueryParam{
		UserName: userName,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "检查用户是否存在失败")
	} else if len(UserQueryResult.Data) == 0 {
		return nil, errors.New("未找到该用户")
	}

	item := UserQueryResult.Data[0]
	if item.Password != secret.SHA1String(password) {
		return nil, errors.New("用户密码不正确，请使用正确的密码")
	} else if item.Status != consts.BaseStatusDisabled {
		return nil, errors.New("用户已被禁用")
	}
	return item, nil
}
