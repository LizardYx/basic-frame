package bll

import (
	"basic-frame/middleware"
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"basic-frame/util/ginx/errors"
	"basic-frame/util/mysql"
	"basic-frame/util/secret"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

var UserBll = &User{
	RoleModel:           model.RoleModel,
	OrganizationModel:   model.OrganizationModel,
	PositionModel:       model.PositionModel,
	UserModel:           model.UserModel,
	UserExtendInfoModel: model.UserExtendInfoModel,
}

type User struct {
	RoleModel           *model.Role
	OrganizationModel   *model.Organization
	PositionModel       *model.Position
	UserModel           *model.User
	UserExtendInfoModel *model.UserExtendInfo
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

func (a *User) GetSubUsers(c *gin.Context, userID uint64, params schema.SubUserQueryParam) (*schema.UserQueryResult, error) {
	users := make(schema.Users, 0)

	// 获取当前用户的信息
	userInfo, err := a.Get(c, userID)
	if err != nil {
		return nil, errors.WithMessage(err, "获取当前用户信息失败")
	}

	// 获取 拥有用户同组织子职位的用户
	if params.CurrentOrg && len(userInfo.Positions) != 0 {
		var sonPositionIDs []uint64
		// 获取所有的职位
		PositionQueryResult, err := a.PositionModel.Query(schema.PositionQueryParam{
			Status:  consts.BaseStatusEnable,
			FindAll: true,
		})
		if err != nil {
			return nil, errors.WithMessage(err, "获取职位列表失败")
		}

		// 获取用户同组织的子职位ID
		for _, item := range userInfo.Positions {
			for _, position := range PositionQueryResult.Data {
				// 如果职位属于同一组织且权重小于用户的职位
				if item.OrganizationID == position.OrganizationID && position.Sequence < item.Sequence &&
					!common.ContainsUint64(sonPositionIDs, position.ID) {
					sonPositionIDs = append(sonPositionIDs, position.ID)
				}
			}
		}

		// 获取拥有该子职位的用户
		if len(sonPositionIDs) != 0 {
			if UserQueryResult, err := a.UserModel.Query(schema.UserQueryParam{
				Status:       consts.BaseStatusEnable,
				PositionIDs:  common.UintSliceToString(sonPositionIDs, ","),
				FindAll:      true,
				OmitPassword: true,
			}); err != nil {
				return nil, errors.WithMessage(err, "获取用户信息失败")
			} else {
				for _, user := range UserQueryResult.Data {
					if len(users) == 0 {
						users = append(users, user)
						continue
					}
					if !common.ContainsUint64(users.GetIDs(), user.ID) {
						users = append(users, user)
					}
				}
			}
		}
	}

	return &schema.UserQueryResult{
		PageResult: nil,
		Data:       *users.SortUser(),
	}, nil
}

func (a *User) Init(c *gin.Context) error {
	// 检查超管用户是否创建
	if UserQueryResult, err := a.Query(c, schema.UserQueryParam{
		UserName:     consts.AdminName,
		OmitPassword: true,
	}); err != nil {
		return errors.WithMessage(err, "检查超管用户是否创建失败")
	} else if len(UserQueryResult.Data) != 0 {
		// 超管已创建
		return nil
	} else {
		// 超管未创建
		if _, err = a.Create(c, schema.User{
			UserName: consts.AdminName,
			Password: consts.AdminPassword,
			Status:   consts.BaseStatusEnable,
		}); err != nil {
			return errors.WithMessage(err, "创建超管用户失败")
		}
	}
	return nil
}

func (a *User) Get(c *gin.Context, id uint64) (*schema.User, error) {
	// 获取用户信息
	item, err := a.UserModel.Get(id)
	if err != nil {
		return nil, errors.WithMessage(err, "获取用户信息失败")
	} else if item == nil {
		return nil, errors.New("用户不存在")
	}

	item.Password = ""
	return item, nil
}

func (a *User) Create(c *gin.Context, item schema.User) (*common.IDResult, error) {
	// 检查组织、职位、角色、用户组是否被禁用
	if err := a.UserParamsCheck(c, &item); err != nil {
		return nil, err
	}

	item.Password = secret.SHA1String(item.Password)
	var IDResult *common.IDResult
	var err error
	_ = mysql.DB.Transaction(func(tx *gorm.DB) error {
		if IDResult, err = a.UserModel.Create(item); err != nil {
			return err
		}
		if _, err = a.UserExtendInfoModel.Create(schema.UserExtendInfo{UserID: IDResult.ID, RealName: item.UserName}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return IDResult, nil
}

// ---------------------------------------- Params Validate --------------------------------------

func (a *User) UserParamsCheck(c *gin.Context, item *schema.User) error {
	// 检查用户名或密码是否为空
	item.UserName = strings.TrimSpace(item.UserName)
	item.Password = strings.TrimSpace(item.Password)
	if item.UserName == "" || item.Password == "" {
		return errors.New("用户名或密码不能为空")
	}

	// 检查用户名是否存在
	UserQueryResult, err := a.Query(c, schema.UserQueryParam{
		UserName:     item.UserName,
		OmitPassword: true,
	})
	if err != nil {
		return errors.WithMessage(err, "检查用户名是否存在失败")
	} else if len(UserQueryResult.Data) != 0 {
		if UserQueryResult.Data[0].ID != item.ID {
			return errors.New("该用户已存在")
		}
	}

	// 检查组织是否被禁用
	if len(item.Organizations) != 0 {
		if OrgQueryResult, err := a.OrganizationModel.Query(schema.OrganizationQueryParam{
			IDs:    common.UintSliceToString(item.Organizations.GetIDs(), ","),
			Status: consts.BaseStatusDisabled,
		}); err != nil {
			return err
		} else if len(OrgQueryResult.Data) != 0 {
			var orgName []string
			for _, orgInfo := range OrgQueryResult.Data {
				orgName = append(orgName, orgInfo.Name)
			}
			return errors.New(fmt.Sprintf("组织: %s 已被禁用", common.StringSliceToString(orgName, ",")))
		}
	}

	// 检查职位是否被禁用
	if len(item.Positions) != 0 {
		if PositionQueryResult, err := a.PositionModel.Query(schema.PositionQueryParam{
			IDs:    common.UintSliceToString(item.Positions.GetIDs(), ","),
			Status: consts.BaseStatusDisabled,
		}); err != nil {
			return err
		} else if len(PositionQueryResult.Data) != 0 {
			var positionNames []string
			for _, position := range PositionQueryResult.Data {
				positionNames = append(positionNames, position.Name)
			}
			return errors.New(fmt.Sprintf("职位: %s 已被禁用", common.StringSliceToString(positionNames, ",")))
		}
	}

	// 检查角色是否被禁用、类型是否正确
	if len(item.Roles) != 0 {
		if RoleQueryResult, err := a.RoleModel.Query(schema.RoleQueryParam{
			IDs:    common.UintSliceToString(item.Roles.GetIDs(), ","),
			Status: consts.BaseStatusDisabled,
		}); err != nil {
			return err
		} else if len(RoleQueryResult.Data) != 0 {
			var disabledRoleNames []string
			var ErrTypeRoleNames []string
			for _, roleInfo := range RoleQueryResult.Data {
				if roleInfo.Status == consts.BaseStatusDisabled {
					disabledRoleNames = append(disabledRoleNames, roleInfo.Name)
				}
				if roleInfo.Type != consts.RoleTypeForUser {
					ErrTypeRoleNames = append(ErrTypeRoleNames, roleInfo.Name)
				}
			}
			if len(disabledRoleNames) != 0 {
				return errors.New(fmt.Sprintf("角色: %s 已被禁用", common.StringSliceToString(disabledRoleNames, ",")))
			}
			if len(ErrTypeRoleNames) != 0 {
				return errors.New(fmt.Sprintf("角色: %s 的角色类型错误", common.StringSliceToString(ErrTypeRoleNames, ",")))
			}
		}
	}

	// 检查用户组是否被禁用
	// TODO: 等待用户组表完成
	//if len(item.UserGroups) != 0 {
	//	if UserGroupQueryResult, err := a.UserGroupModel.Query(schema.UserGroupQueryParam{
	//		IDs:             common.UintSliceToString(item.UserGroups.GetIDs(), ","),
	//		Status:          consts.BaseStatusDisabled,
	//	}); err != nil {
	//		return err
	//	} else if len(UserGroupQueryResult.Data) != 0 {
	//		var userGroupNames []string
	//		for index, userGroup := range UserGroupQueryResult.Data {
	//			userGroupNames = append(userGroupNames, userGroup.Name)
	//		}
	//		return errors.New(fmt.Sprintf("角色组: %s 已被禁用", common.StringSliceToString(userGroupNames, ",")))
	//	}
	//}
	return nil
}

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
