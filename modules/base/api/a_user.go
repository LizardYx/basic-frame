package api

import (
	"basic-frame/middleware"
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"basic-frame/util/consts"
	"basic-frame/util/ginx"
	"github.com/gin-gonic/gin"
)

var UserApi = &User{
	UserBll: bll.UserBll,
}

type User struct {
	UserBll *bll.User
}

// Login 用户登录
func (a *User) Login(c *gin.Context) {
	var item schema.LoginParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	tokenInfo, userInfo, err := a.UserBll.Login(c, item)
	if err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResSuccess(c, "", schema.LoginRes{
		UserInfo:  *userInfo,
		TokenInfo: *tokenInfo,
	})
}

// Logout 用户登出(前端需要自己删除Jwt Token字符串)
func (a *User) Logout(c *gin.Context) {
	userID := ginx.GetUserID(c)
	if userID != 0 {
		// 注销ws连接
		for client := range middleware.Manager.Clients {
			if client.ID == userID {
				middleware.Manager.Unregister <- client
				break
			}
		}
	}
	ginx.ResOperateSuccess(c, "")
}

// RefreshToken 刷新Jwt Token字符串
func (a *User) RefreshToken(c *gin.Context) {
	token, err := middleware.RefreshToken(ginx.GetToken(c))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", map[string]interface{}{
		"token_info": token,
	})
}

func (a *User) Query(c *gin.Context) {
	var params schema.UserQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, params, err)
		return
	}

	params.Pagination = true
	result, err := a.UserBll.Query(c, params)
	if err != nil {
		ginx.ResError(c, params, err)
		return
	}
	ginx.ResList(c, params, result.Data)
}

func (a *User) GetSubUsers(c *gin.Context) {
	var params schema.SubUserQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, params, err)
		return
	}

	result, err := a.UserBll.GetSubUsers(c, ginx.GetUserID(c), params)
	if err != nil {
		ginx.ResError(c, params, err)
		return
	}
	ginx.ResList(c, params, result.Data)
}

func (a *User) Get(c *gin.Context) {
	item, err := a.UserBll.Get(c, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

func (a *User) GetMenuTree(c *gin.Context) {
	item, err := a.UserBll.GetMenuTree(c)
	if err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResSuccess(c, "", item)
}

func (a *User) UpdateUserPermission(c *gin.Context) {
	var item schema.User
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.UserBll.UpdateUserPermission(c, ginx.ParseParamID(c, "id"), item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

func (a *User) BatchUpdateUserPermission(c *gin.Context) {
	var items schema.Users
	if err := ginx.ParseJSON(c, &items); err != nil {
		ginx.ResError(c, items, err)
		return
	}

	if err := a.UserBll.BatchUpdateUserPermission(c, items); err != nil {
		ginx.ResError(c, items, err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

func (a *User) Create(c *gin.Context) {
	var item schema.User
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.UserBll.Create(c, item)
	if err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResSuccess(c, item, result)
}

func (a *User) Update(c *gin.Context) {
	var item schema.User
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.UserBll.Update(c, ginx.ParseParamID(c, "id"), item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, item)
}

func (a *User) EnableUser(c *gin.Context) {
	if err := a.UserBll.UpdateStatus(c, ginx.ParseParamID(c, "id"), consts.BaseStatusEnable); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

func (a *User) DisabledUser(c *gin.Context) {
	if err := a.UserBll.UpdateStatus(c, ginx.ParseParamID(c, "id"), consts.BaseStatusDisabled); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

func (a *User) UpdatePassword(c *gin.Context) {
	var item schema.UpdatePasswordParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.UserBll.UpdatePassword(c, ginx.GetUserID(c), item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}

func (a *User) Delete(c *gin.Context) {
	if err := a.UserBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}
