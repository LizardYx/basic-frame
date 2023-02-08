package ginx

import (
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"basic-frame/util/ginx/errors"
	"basic-frame/util/i18n/Localizer"
	"basic-frame/util/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type ParamsID struct {
	ID uint64
}

// SetUserInfo 将用户ID、用户名称设置在context上下文中
func SetUserInfo(c *gin.Context, userID uint64, userName string) {
	c.Set("user_id", userID)
	c.Set("user_name", userName)
}

// GetUserID 从context上下文中获取用户ID
func GetUserID(c *gin.Context) uint64 {
	v, isExists := c.Get("user_id")
	if !isExists {
		return 0
	}
	if userID, ok := v.(uint64); ok {
		return userID
	}
	return 0
}

// GetUserName 从context上下文中获取用户名称
func GetUserName(c *gin.Context) string {
	v, isExists := c.Get("user_name")
	if !isExists {
		return ""
	}
	if userName, ok := v.(string); ok {
		return userName
	}
	return ""
}

// GetToken 获取用户令牌
func GetToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

// ParseParamID 获取请求Url中的ID参数
func ParseParamID(c *gin.Context, key string) uint64 {
	val := c.Param(key)
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

// ParseJSON 解析请求JSON
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.Wrap400Response(err, "", fmt.Sprintf("解析参数发生错误 - %s", err.Error()))
	}
	return nil
}

// ParseQuery 解析Query参数
func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.Wrap400Response(err, "", fmt.Sprintf("解析参数发生错误 - %s", err.Error()))
	}
	return nil
}

// ------------------------ Response ----------------------------

// ListResult 响应列表数据
type ListResult struct {
	List       interface{}              `json:"list"`
	Pagination *common.PaginationResult `json:"pagination,omitempty"`
}

// OperateResult 操作响应
type OperateResult struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
}

// ResPaginate 响应分页数据
func ResPaginate(c *gin.Context, params, v interface{}, pr *common.PaginationResult) {
	list := ListResult{
		List:       v,
		Pagination: pr,
	}
	ResSuccess(c, params, list)
}

// ResList 响应列表数据
func ResList(c *gin.Context, params, v interface{}) {
	ResSuccess(c, params, ListResult{List: v})
}

func ResOperateSuccess(c *gin.Context, params interface{}) {
	ResSuccessString(c, params, Localizer.I18n.Translate(consts.ApiOperateSuccess))
}

// ResSuccessString 响应操作成功
func ResSuccessString(c *gin.Context, params interface{}, msg string, args ...interface{}) {
	ResSuccess(c, params, OperateResult{Message: Localizer.I18n.Translate(msg, args...)})
}

// ResSuccess 响应成功
func ResSuccess(c *gin.Context, params interface{}, v interface{}) {
	ResJSON(c, http.StatusOK, params, v)
}

// ResError 响应错误
func ResError(c *gin.Context, params interface{}, err error) {
	var res *common.ResponseError

	if err != nil {
		if e, ok := err.(*common.ResponseError); ok {
			res = e
		} else {
			res = errors.UnWrapResponse(errors.NewIANAResponse(http.StatusInternalServerError, err.Error()))
			res.ERR = err
		}
	} else {
		res = errors.UnWrapResponse(errors.NewIANAResponse(http.StatusInternalServerError, ""))
	}

	// 日志记录
	if common.SysConfig.RunMode == consts.RunModeDebug {
		if res.ERR != nil {
			logger.Log.Warningf("%+v", res.ERR) // 详细的链式错误日志
			fmt.Printf("%+v\n", res.ERR)        // 详细的链式错误日志
		} else {
			logger.Log.Warningf("%+v", res.Message) // 详细的链式错误日志
			fmt.Printf("%+v\n", res.Message)        // 详细的链式错误日志
		}
	} else {
		if res.ERR != nil {
			logger.Log.Warningf("%v", res.ERR) // 基本错误日志
			fmt.Printf("%v\n", res.ERR)        // 基本错误日志
		} else {
			logger.Log.Warningf("%v", res.Message) // 基本错误日志
			fmt.Printf("%v\n", res.Message)        // 基本错误日志
		}
	}
	ResJSON(c, res.StatusCode, params, OperateResult{Code: res.Code, Message: res.Message})
}

// ResJSON 响应JSON数据
func ResJSON(c *gin.Context, status int, params, response interface{}) {
	if params != "" && params != nil {
		logger.Log.Infof("Params: %+v", params)
	}
	if response != "" && response != nil {
		logger.Log.Infof("Response: %+v", response)
	}
	c.JSON(status, response)
}
