package api

import (
	"basic-frame/modules/base/bll"
	"github.com/gin-gonic/gin"
)

var ButtonAPI = &Button{}

type Button struct {
	ButtonBll *bll.Button
}

func (a *Button) Create(c *gin.Context) {

}
