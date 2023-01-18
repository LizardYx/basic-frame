package bll

import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/ginx/errors"
	"github.com/gin-gonic/gin"
	"strings"
)

var TagManageBll = &TagManage{
	TagManageModel: model.TagManageModel,
}

type TagManage struct {
	TagManageModel *model.TagManage
}

func (a *TagManage) Query(c *gin.Context, params schema.TagManageQueryParam) (*schema.TagManageQueryResult, error) {
	return a.TagManageModel.Query(params)
}

func (a *TagManage) Get(c *gin.Context, id uint64) (*schema.TagManage, error) {
	item, err := a.TagManageModel.Get(id)
	if err != nil {
		return nil, errors.WithMessage(err, "检查标签是否存在失败")
	} else if item == nil {
		return nil, errors.New("未找到该标签")
	}

	return item, nil
}

func (a *TagManage) Create(c *gin.Context, item schema.TagManage) (*common.IDResult, error) {
	// 标签参数验证
	if err := a.TagParamsValidate(c, &item); err != nil {
		return nil, err
	}

	// 检查标签是否存在
	if TagManageQueryResult, err := a.Query(c, schema.TagManageQueryParam{
		Type:  item.Type,
		Value: item.Value,
	}); err != nil {
		return nil, errors.WithMessage(err, "检查标签是否存在失败")
	} else if len(TagManageQueryResult.Data) != 0 {
		return nil, errors.New("标签已存在")
	}

	// 创建标签
	return a.TagManageModel.Create(item)
}

func (a *TagManage) Update(c *gin.Context, id uint64, item schema.TagManage) error {
	// 标签参数检查
	if err := a.TagParamsValidate(c, &item); err != nil {
		return err
	}

	// 检查标签是否存在
	oldItem, err := a.TagManageModel.Get(id)
	if err != nil {
		return errors.WithMessage(err, "检查标签是否存在失败")
	} else if oldItem == nil {
		return errors.New("未找到该标签")
	} else {
		item.Type = oldItem.Type
	}

	// 更新标签
	return a.TagManageModel.UpdateByID(id, map[string]interface{}{
		"key_name": item.KeyName,
		"value":    item.Value,
	})
}

func (a *TagManage) Delete(c *gin.Context, id uint64) error {
	// 检查标签是否存在
	oldItem, err := a.TagManageModel.Get(id)
	if err != nil {
		return errors.WithMessage(err, "检查标签是否存在失败")
	} else if oldItem == nil {
		return errors.New("未找到该标签")
	}

	// 删除标签
	return a.TagManageModel.Delete(id)
}

// ---------------------------------------- Params  Validate --------------------------------------

func (a *TagManage) TagParamsValidate(c *gin.Context, item *schema.TagManage) error {
	// 检查标签type是否为空
	if item.Type = strings.TrimSpace(item.Type); item.Type == "" {
		return errors.New("标签Type不能为空")
	}

	// 检查标签key是否为空
	if item.KeyName = strings.TrimSpace(item.KeyName); item.KeyName == "" {
		return errors.New("标签Key不能为空")
	}

	// 检查标签值是否为空
	if item.Value = strings.TrimSpace(item.Value); item.Value == "" {
		return errors.New("标签值不能为空")
	}
	return nil
}
