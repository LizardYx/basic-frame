package schema

import (
	"basic-frame/util/common"
	"sort"
	"time"
)

type Menu struct {
	ID         uint64    `json:"id"`                        // 唯一标识
	UUID       string    `json:"UUID"`                      // 前端组装菜单需要的
	Select     bool      `json:"select"`                    // 是否被选中(前端用于判断菜单是否被选中)
	Name       string    `json:"name" binding:"required"`   // 菜单名称
	Icon       string    `json:"icon"`                      // 菜单图标
	Class      string    `json:"class"`                     // 菜单样式
	Router     string    `json:"router" binding:"required"` // 访问路由
	Sequence   int       `json:"sequence"`                  // 排序值
	ParentID   *uint64   `json:"parent_id"`                 // 父级菜单ID
	ShowStatus int       `json:"show_status"`               // 显示状态(1:显示 2:隐藏)
	Status     int       `json:"status"`                    // 状态(1:启用 2:禁用)
	Memo       string    `json:"memo"`                      // 备注
	Creator    uint64    `json:"creator"`                   // 创建者
	CreatedAt  time.Time `json:"created_at"`                // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                // 更新时间
}

func (a Menu) Init() *Menu {
	if a.ParentID == nil {
		parentID := uint64(0)
		a.ParentID = &parentID
	}
	return &a
}

type MenuQueryParam struct {
	common.PaginationParam
	ID         uint64   `form:"id"`          // 菜单ID
	IDs        []uint64 `form:"ids"`         // 菜单ID集合
	UUID       string   `form:"uuid"`        // 菜单UUID
	ParentID   uint64   `form:"parent_id"`   // 父级菜单ID
	Router     string   `form:"router"`      // 访问路由
	Status     int      `form:"status"`      // 状态(1:启用 2:禁用)
	ShowStatus int      `form:"show_status"` // 显示状态(1:显示 2:隐藏)
	FindAll    bool     `form:"find_all"`    // 是否查询所有数据
}

type Menus []*Menu

type MenuQueryResult struct {
	Data       Menus
	PageResult *common.PaginationResult
}

func (a Menus) GetIDs() []uint64 {
	l := make([]uint64, len(a))
	for i, j := range a {
		l[i] = j.ID
	}
	return l
}

func (a Menus) Init() *Menus {
	if len(a) == 0 {
		a = make(Menus, 0)
	} else {
		for index, menu := range a {
			a[index] = menu.Init()
		}
	}
	return &a
}

// ----------------------------------------MenuTree--------------------------------------

// MenuTree 菜单树
type MenuTree struct {
	Menu
	RestfulApis RestfulApis `json:"restful_apis"` // 页面调用的Api
	Buttons     ButtonPres  `json:"buttons"`      // 页面的按钮
	SonMenus    MenuTrees   `json:"son_menus"`    // 子菜单
}

type MenuTrees []*MenuTree

// MenuTreeJson 前端菜单Json文件解析struct
type MenuTreeJson struct {
	MenuVersion    float64        `json:"menu_version"`    // 菜单版本号
	MenuTrees      MenuTrees      `json:"menu_trees"`      // 菜单树
	DisabledFields DisabledFields `json:"disabled_fields"` // 可禁用的字段
}

// PermissionTree 创建角色的权限树
type PermissionTree struct {
	MenuTrees      MenuTrees      `json:"menu_trees"`      // 菜单树
	DisabledFields DisabledFields `json:"disabled_fields"` // 可禁用的字段
}

// SortMenuTrees 菜单树排序
func (a MenuTrees) SortMenuTrees() *MenuTrees {
	if len(a) != 0 {
		for _, menuTreeInfo := range a {
			menuTreeInfo.Buttons.SortButtonTrees()
			if len(menuTreeInfo.SonMenus) != 0 {
				menuTreeInfo.SonMenus.SortMenuTrees()
			}
		}
		sort.SliceStable(a, func(i, j int) bool {
			return a[i].Sequence > a[j].Sequence
		})
	}
	return &a
}

// Init 菜单树数据初始化
func (a MenuTree) Init() *MenuTree {
	if a.ParentID == nil {
		parentID := uint64(0)
		a.ParentID = &parentID
	}
	if len(a.RestfulApis) == 0 {
		a.RestfulApis = make(RestfulApis, 0)
	}
	if len(a.Buttons) == 0 {
		a.Buttons = make(ButtonPres, 0)
	} else {
		for index, buttonTree := range a.Buttons {
			a.Buttons[index] = buttonTree.Init()
		}
	}
	if len(a.SonMenus) == 0 {
		a.SonMenus = make(MenuTrees, 0)
	} else {
		for index, sonMenuTree := range a.SonMenus {
			a.SonMenus[index] = sonMenuTree.Init()
		}
	}
	return &a
}

func (a MenuTrees) Init() *MenuTrees {
	for index, item := range a {
		*a[index] = *item.Init()
	}
	return &a
}

// SetCreator 设置菜单树操作者
func (a MenuTree) SetCreator(creator uint64) *MenuTree {
	if creator != 0 {
		a.Creator = creator
		if len(a.SonMenus) != 0 {
			for _, sonMenu := range a.SonMenus {
				sonMenu.SetCreator(creator)
			}
		}
		if len(a.Buttons) != 0 {
			a.Buttons.SetCreator(creator)
		}
	}
	return &a
}

// SetCreator 设置菜单树操作者
func (a MenuTrees) SetCreator(creator uint64) *MenuTrees {
	if creator != 0 {
		for _, menuTree := range a {
			menuTree.SetCreator(creator)
		}
	}
	return &a
}

func (a MenuTree) ToSchemaMenu() *Menu {
	item := new(Menu)
	_ = common.Copy(a, item)
	return item
}

// GetMenuIDsByMenuID 获取指定菜单的菜单ID集合
func (a MenuTrees) GetMenuIDsByMenuID(menuID uint64, menuIDs *[]uint64) {
	for _, menuInfo := range a {
		if menuInfo.ID == menuID {
			menuInfo.GetMenuIDs(menuIDs)
			break
		} else if len(menuInfo.SonMenus) != 0 {
			menuInfo.SonMenus.GetMenuIDsByMenuID(menuID, menuIDs)
		}
	}
	return
}

// GetMenuIDs 获取菜单树的ID集合
func (a MenuTree) GetMenuIDs(menuIDs *[]uint64) {
	*menuIDs = append(*menuIDs, a.ID)
	if len(a.SonMenus) != 0 {
		for _, menuTreeInfo := range a.SonMenus {
			menuTreeInfo.GetMenuIDs(menuIDs)
		}
	}
}

// GetButtonIDsByMenuID 获取指定菜单的按钮ID集合
func (a MenuTrees) GetButtonIDsByMenuID(menuID uint64, buttonIDs *[]uint64) {
	for _, menuInfo := range a {
		if menuInfo.ID == menuID {
			menuInfo.GetButtonIDs(buttonIDs)
			return
		} else if len(menuInfo.SonMenus) != 0 {
			menuInfo.SonMenus.GetButtonIDsByMenuID(menuID, buttonIDs)
		}
	}
	return
}

func (a MenuTree) GetButtonIDs(buttonIDs *[]uint64) {
	if len(a.Buttons) != 0 {
		a.Buttons.GetIDs(buttonIDs)
	}
	if len(a.SonMenus) != 0 {
		for _, sonMenu := range a.SonMenus {
			sonMenu.GetButtonIDs(buttonIDs)
		}
	}
}

// GetSonButtonIDs 获取按钮的子按钮ID集合
func (a MenuTrees) GetSonButtonIDs(menuID, buttonID uint64) []uint64 {
	var sonButtonIDs []uint64
	for _, menuInfo := range a {
		if menuInfo.ID == menuID {
			menuInfo.Buttons.GetButtonIDsByButtonID(buttonID, &sonButtonIDs)
			break
		} else if len(menuInfo.SonMenus) != 0 {
			menuInfo.SonMenus.GetSonButtonIDs(menuID, buttonID)
		}
	}
	return sonButtonIDs
}
