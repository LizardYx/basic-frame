package schema

import "C"
import (
	"basic-frame/util/common"
	"sort"
	"time"
)

type Button struct {
	ID         uint64    `json:"id"`                         // 唯一标识
	UUID       string    `json:"UUID"`                       // 前端组装菜单需要的
	BtnID      int       `json:"btn_id" binding:"required"`  // 前端识别按钮用
	Select     bool      `json:"select"`                     // 是否被选中(前端用于判断菜单是否被选中)
	Name       string    `json:"name" binding:"required"`    // 按钮名称
	Icon       string    `json:"icon"`                       // 按钮图标
	Class      string    `json:"class"`                      // 按钮样式
	MenuID     uint64    `json:"menu_id" binding:"required"` // 菜单ID
	Sequence   int       `json:"sequence"`                   // 排序值
	ParentID   *uint64   `json:"parent_id"`                  // 父级按钮ID
	ShowStatus int       `json:"show_status"`                // 显示状态(1:显示 2:隐藏)
	Status     int       `json:"status"`                     // 状态(1:启用 2:禁用)
	Memo       string    `json:"memo"`                       // 备注
	Creator    uint64    `json:"creator"`                    // 创建者
	CreatedAt  time.Time `json:"created_at"`                 // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                 // 更新时间
}

func (a Button) Init() *Button {
	if a.ParentID == nil {
		parentID := uint64(0)
		a.ParentID = &parentID
	}
	return &a
}

type ButtonQueryParam struct {
	common.PaginationParam
	IDs        string `form:"ids"`         // ID集合(逗号分隔)
	UUID       string `form:"uuid"`        // 按钮UUID
	ParentID   uint64 `form:"parent_id"`   // 父级菜单ID
	ShowStatus int    `form:"show_status"` // 显示状态(1:显示 2:隐藏)
	Status     int    `form:"status"`      // 状态(1:启用 2:禁用)
}

type Buttons []*Button

type ButtonQueryResult struct {
	Data       Buttons
	PageResult *common.PaginationResult
}

func (a Buttons) GetIDs() []uint64 {
	l := make([]uint64, len(a))
	for i, j := range a {
		l[i] = j.ID
	}
	return l
}

func (a Buttons) Init() *Buttons {
	items := make(Buttons, 0)

	for _, button := range a {
		items = append(items, button.Init())
	}
	return &items
}

// ---------------------------------------- Response Struct --------------------------------------

type ButtonPre struct {
	Button
	RestfulApis RestfulApis `json:"restful_apis"` // 按钮调用的Api
	SonButtons  ButtonPres  `json:"son_buttons"`  // 子按钮
}

// Init 按钮树数据初始化
func (a ButtonPre) Init() *ButtonPre {
	if a.ParentID == nil {
		parentID := uint64(0)
		a.ParentID = &parentID
	}
	a.RestfulApis = *a.RestfulApis.Init()
	a.SonButtons = *a.SonButtons.Init()
	return &a
}

// SetCreator 设置按钮、子按钮的创建者
func (a ButtonPre) SetCreator(creator uint64, item *ButtonPre) {
	if creator != 0 {
		item.Creator = creator
		item.RestfulApis = *a.RestfulApis.SetCreator(creator)
		item.SonButtons = *a.SonButtons.SetCreator(creator)
	}
	return
}

// GetIDs 获取按钮、子按钮的按钮ID集合
func (a ButtonPre) GetIDs(buttonIDs *[]uint64) {
	*buttonIDs = append(*buttonIDs, a.ID)
	if len(a.SonButtons) != 0 {
		a.SonButtons.GetIDs(buttonIDs)
	}
}

func (a ButtonPre) ToSchemaButton() *Button {
	item := new(Button)
	_ = common.Copy(a, item)
	return item
}

// GetRestfulApis 获取按钮关联的restfulApi集合(去重)
func (a ButtonPre) GetRestfulApis(items *RestfulApis) {
	for _, restfulApi := range a.RestfulApis {
		if len(*items) == 0 {
			*items = append(*items, restfulApi)
		}
		for index, item := range *items {
			if item.UUID == restfulApi.UUID {
				break
			}
			if index == (len(*items) - 1) {
				*items = append(*items, restfulApi)
			}
		}
	}
}

type ButtonPres []*ButtonPre

func (a ButtonPres) Init() *ButtonPres {
	items := make(ButtonPres, 0)
	for _, item := range a {
		items = append(items, item.Init())
	}
	return &items
}

// GetIDs 获取按钮、子按钮的按钮ID集合
func (a ButtonPres) GetIDs(buttonIDs *[]uint64) {
	for _, buttonPreInfo := range a {
		*buttonIDs = append(*buttonIDs, buttonPreInfo.ID)
		if len(buttonPreInfo.SonButtons) != 0 {
			buttonPreInfo.SonButtons.GetIDs(buttonIDs)
		}
	}
}

// SetCreator 设置按钮、子按钮的创建者
func (a ButtonPres) SetCreator(creator uint64) *ButtonPres {
	if creator != 0 {
		for _, buttonPreInfo := range a {
			buttonPreInfo.SetCreator(creator, buttonPreInfo)
		}
	}
	return &a
}

// GetButtonIDsByButtonID 获取指定按钮的子按钮ID集合
func (a ButtonPres) GetButtonIDsByButtonID(buttonID uint64, sonButtonIDs *[]uint64) {
	for _, buttonPreInfo := range a {
		if buttonPreInfo.ID == buttonID {
			buttonPreInfo.GetIDs(sonButtonIDs)
			return
		} else if len(buttonPreInfo.SonButtons) != 0 {
			buttonPreInfo.SonButtons.GetButtonIDsByButtonID(buttonID, sonButtonIDs)
		}
	}
}

// SortButtonTrees 按钮树通过权重排序
func (a ButtonPres) SortButtonTrees() {
	if len(a) != 0 {
		for _, buttonPreInfo := range a {
			if len(buttonPreInfo.SonButtons) != 0 {
				buttonPreInfo.SonButtons.SortButtonTrees()
			}
		}
		sort.SliceStable(a, func(i, j int) bool {
			return a[i].Sequence > a[j].Sequence
		})
	}
}

// GetRestfulApis 获取按钮关联的restfulApi集合(去重)
func (a ButtonPres) GetRestfulApis() RestfulApis {
	items := make(RestfulApis, 0)

	for _, ButtonPreInfo := range a {
		ButtonPreInfo.GetRestfulApis(&items)
	}
	return items
}
