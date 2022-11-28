package ginx

import (
	"basic-frame/util/mysql"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int
	Data    ListResult
	Message string
}

// ListResult 响应分页查询数据
type ListResult struct {
	List       interface{}             `json:"list"`
	Pagination *mysql.PaginationResult `json:"pagination,omitempty"`
}

// ResPagination 分页查询返回
func ResPagination(ctx *gin.Context, HttpStatus int, item Response) {

	ResJSON(ctx, HttpStatus, item)
}

// ResSuccess 返回列表数据
func ResSuccess(ctx *gin.Context, HttpStatus int, v interface{}) {
	ResJSON(ctx, HttpStatus, v)
}

// ResJSON 返回JSON数据
func ResJSON(ctx *gin.Context, HttpStatus int, v interface{}) {
	ctx.JSON(HttpStatus, v)
}
