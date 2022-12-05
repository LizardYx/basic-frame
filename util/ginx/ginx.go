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
)

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

// ErrorResult 响应错误
type ErrorResult struct {
	Error ErrorItem `json:"error"` // 错误项
}

// ErrorItem 响应错误项
type ErrorItem struct {
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

// ResSuccessString 响应操作成功
func ResSuccessString(c *gin.Context, params interface{}, msg string, args ...interface{}) {
	ResSuccess(c, params, Localizer.I18n.Translate(msg, args...))
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
			res = errors.UnWrapResponse(errors.New500Response(err.Error()))
			res.ERR = err
		}
	} else {
		res = errors.UnWrapResponse(errors.New500Response("服务器发生错误"))
	}

	// 日志记录
	if common.SysConfig.RunMode == consts.RunModeDebug {
		logger.Log.Warningf("%+v", res.ERR) // 详细的链式错误日志
		fmt.Printf("%+v\n", res.ERR)        // 详细的链式错误日志
	} else {
		logger.Log.Warningf("%v", res.ERR) // 基本错误日志
		fmt.Printf("%v\n", res.ERR)        // 基本错误日志
	}

	// 创建Response
	eitem := ErrorItem{
		Code:    res.Code,
		Message: res.Message,
	}
	ResJSON(c, res.StatusCode, params, ErrorResult{Error: eitem})
}

// ResJSON 响应JSON数据
func ResJSON(c *gin.Context, status int, params, response interface{}) {
	logger.Log.Infof("Params: %+v", params)
	logger.Log.Infof("Response: %+v", response)
	c.JSON(status, response)
}
