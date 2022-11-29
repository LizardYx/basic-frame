package ginx

import (
	"basic-frame/util/common"
	"basic-frame/util/ginx/errors"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// ParseJSON 解析请求JSON
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.Wrap400Response(err, fmt.Sprintf("解析参数发生错误 - %s", err.Error()))
	}
	return nil
}

// ParseQuery 解析Query参数
func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.Wrap400Response(err, fmt.Sprintf("解析参数发生错误 - %s", err.Error()))
	}
	return nil
}

// Paginate 分页查询
func Paginate(db *gorm.DB, params common.PaginationParam, out interface{}) (*common.PaginationResult, error) {
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, err
	}

	if params.OnlyCount {
		// 仅查询count
		return &common.PaginationResult{Total: count}, nil
	} else if !params.Pagination {
		// 查询所有数据
		err := db.Find(out).Error
		return nil, err
	} else {
		// 分页查询
		if params.Current == 0 {
			params.Current = 1
		} else if params.PageSize <= 0 {
			params.PageSize = 10
		}
		db = db.Offset((params.Current - 1) * params.PageSize).Limit(params.PageSize)
		err = db.Find(out).Error
		return &common.PaginationResult{
			Total:    count,
			Current:  params.Current,
			PageSize: params.PageSize,
		}, nil
	}
}

// Check 检查数据是否存在
func Check(db *gorm.DB) (bool, error) {
	var count int64
	result := db.Count(&count)
	if err := result.Error; err != nil {
		return false, err
	}
	return count > 0, nil
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
func ResPaginate(ctx *gin.Context, v interface{}, pr *common.PaginationResult) {
	list := ListResult{
		List:       v,
		Pagination: pr,
	}
	ResSuccess(ctx, list)
}

// ResList 响应列表数据
func ResList(ctx *gin.Context, v interface{}) {
	ResSuccess(ctx, ListResult{List: v})
}

// ResSuccess 响应成功
func ResSuccess(ctx *gin.Context, v interface{}) {
	ResJSON(ctx, http.StatusOK, v)
}

// ResJSON 响应JSON数据
func ResJSON(ctx *gin.Context, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	ctx.Set(fmt.Sprintf("%s/res-body", common.SysConfig.AppName), buf)
	ctx.Data(status, "application/json; charset=utf-8", buf)
	ctx.Abort()
}

// ResError 响应错误
func ResError(ctx *gin.Context, err error, status ...int) {
	var res *common.ResponseError

	if err != nil {
		if e, ok := err.(*common.ResponseError); ok {
			res = e
		} else {
			res = errors.UnWrapResponse(errors.New500Response("服务器发生错误"))
			res.ERR = err
		}
	} else {
		res = errors.UnWrapResponse(errors.New500Response("服务器发生错误"))
	}

	if len(status) > 0 {
		res.StatusCode = status[0]
	}

	if err = res.ERR; err != nil {
		if res.Message == "" {
			res.Message = err.Error()
		}
	}

	eitem := ErrorItem{
		Code:    res.Code,
		Message: res.Message,
	}
	ResJSON(ctx, res.StatusCode, ErrorResult{Error: eitem})
}
