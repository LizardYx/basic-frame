package model

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/mysql"
)

var ResApiModel = &RestfulApi{}

type RestfulApi struct {
}

func (a *RestfulApi) Create(item schema.RestfulApi) (uint64, error) {
	result := mysql.DB.Create(&item)
	return item.ID, result.Error
}
