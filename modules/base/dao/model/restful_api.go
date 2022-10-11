package model

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/mysql"
	"context"
)

var ResApiModel = &RestfulApi{}

type RestfulApi struct {
}

func (a *RestfulApi) Create(ctx context.Context, item schema.RestfulApi) (uint64, error) {
	result := mysql.DB.Create(&item)
	return item.ID, result.Error
}
