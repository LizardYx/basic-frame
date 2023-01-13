package bll

import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/ginx"
	"basic-frame/util/ginx/errors"
	"basic-frame/util/logger"
	"basic-frame/util/mysql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var MenuBll = &Menu{
	MenuModel: model.MenuModel,
}

type Menu struct {
	MenuModel          *model.Menu
	ButtonModel        *model.Button
	RestfulApiModel    *model.RestfulApi
	DisabledFieldModel *model.DisabledField
	SystemConfigModel  *model.SystemConfig
}

// InitData 初始化菜单数据
func (a *Menu) InitData() {
	loadFileFailed := false
	sleepTime := 5
	for {
		// 读取前端页面菜单配置文件
		data, err := a.readData(common.SysConfig.MenuFile)
		if err != nil {
			fmt.Println("读取前端页面菜单配置文件失败")
			if !loadFileFailed {
				logger.Log.Warning("读取前端页面菜单配置文件失败")
				loadFileFailed = true
			}
			if sleepTime < 10 {
				sleepTime += 5
			}
			time.Sleep(time.Duration(sleepTime) * time.Second)
			continue
		}
		if loadFileFailed {
			loadFileFailed = false
			logger.Log.Warning("读取前端页面菜单配置文件成功")
		}
		// Json文件中的版本号大于DB中的版本号。则表明Json文件有更新
		if data.MenuVersion > common.SysConfig.MenuVersion {
			// 更新菜单树和可禁用字段
			if err = a.UpdatePermissionTree(&gin.Context{}, schema.PermissionTree{
				MenuTrees:      data.MenuTrees,
				DisabledFields: data.DisabledFields,
			}); err != nil {
				if sleepTime < 10 {
					sleepTime += 5
				}
				time.Sleep(time.Duration(sleepTime) * time.Second)
				continue
			} else {
				// 获取系统配置ID
				var SystemConfigID uint64
				if SystemConfigQueryResult, err := a.SystemConfigModel.Query(schema.SystemConfigQueryParam{}); err != nil {
					logger.Log.Warningf("获取系统配置失败: %s", err.Error())
					fmt.Println(fmt.Sprintf("获取系统配置失败: %s", err.Error()))
					return
				} else if len(SystemConfigQueryResult.Data) == 0 {
					logger.Log.Warningf("获取系统配置失败: %s", err.Error())
					fmt.Println(fmt.Sprintf("获取系统配置失败: %s", err.Error()))
					return
				} else {
					SystemConfigID = SystemConfigQueryResult.Data[0].ID
				}

				// 更新系统配置中的菜单版本号
				if err = a.SystemConfigModel.UpdateByID(SystemConfigID, map[string]interface{}{
					"menu_version": data.MenuVersion,
				}); err != nil {
					logger.Log.Warningf("更新系统配置中的菜单版本号: %s", err.Error())
					fmt.Println(fmt.Sprintf("更新系统配置中的菜单版本号: %s", err.Error()))
					return
				} else {
					common.SysConfig.MenuVersion = data.MenuVersion
				}
			}
		}
		sleepTime = 5
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}

// readData 读取前端页面菜单配置文件
func (a *Menu) readData(name string) (data schema.MenuTreeJson, err error) {
	file, err := os.Open(name)
	if err != nil {
		return data, err
	}
	defer file.Close()
	// 创建json解码器
	decoder := json.NewDecoder(file)
	// 解析前端页面路由配置文件
	if err = decoder.Decode(&data); err != nil {
		logger.Log.Warningf("解析前端页面路由配置文件失败: %s", err.Error())
		fmt.Println(fmt.Sprintf("解析前端页面路由配置文件失败: %s", err.Error()))
		return
	}
	return
}

func (a *Menu) Query(c *gin.Context, params schema.MenuQueryParam) (*schema.MenuQueryResult, error) {
	return a.MenuModel.Query(params)
}

func (a *Menu) Get(c *gin.Context, id uint64) (*schema.Menu, error) {
	item, err := a.MenuModel.Get(id)
	if err != nil {
		return nil, errors.WithMessage(err, "获取菜单信息失败")
	} else if item == nil {
		return nil, errors.New("未找到菜单信息")
	}

	return item, nil
}

func (a *Menu) Create(c *gin.Context, item schema.Menu) (*common.IDResult, error) {
	item.UUID = common.GetUUID()
	return a.MenuModel.Create(item)
}

// Update 更新菜单基本信息
func (a *Menu) Update(c *gin.Context, id uint64, item schema.Menu) error {
	// 参数检查
	if err := a.MenuParamsCheck(c, &item); err != nil {
		return err
	}

	// 检查菜单是否存在
	if _, err := a.Get(c, id); err != nil {
		return err
	}

	// 更新菜单基本信息
	return a.MenuModel.UpdateByID(id, map[string]interface{}{
		"select":      item.Select,
		"name":        item.Name,
		"icon":        item.Icon,
		"class":       item.Class,
		"router":      item.Router,
		"sequence":    item.Sequence,
		"parent_id":   item.ParentID,
		"show_status": item.ShowStatus,
		"status":      item.Status,
		"memo":        item.Memo,
	})
}

// BatchUpdateMenus 批量更新菜单基本信息
func (a *Menu) BatchUpdateMenus(c *gin.Context, items schema.Menus) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := a.Update(c, item.ID, *item); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// Delete 删除菜单、菜单调用的Api关联以及菜单的按钮
func (a *Menu) Delete(c *gin.Context, id uint64) error {
	// 查询该菜单是否存在
	if _, err := a.Get(c, id); err != nil {
		return err
	}

	// 查询菜单是否有子项
	if MenuQueryResult, err := a.MenuModel.Query(schema.MenuQueryParam{
		PaginationParam: common.PaginationParam{
			OnlyCount: true,
		},
		ParentID: id,
	}); err != nil {
		return errors.WithMessage(err, "查询菜单是否有子项失败")
	} else if MenuQueryResult.PageResult.Total != 0 {
		return errors.New("该菜单含有子项，请先删除子项")
	}

	// 删除菜单、菜单调用的Api关联以及菜单的按钮
	if err := a.MenuModel.Delete(id); err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// ----------------------------------------MenuTrees--------------------------------------

// UpdateMenuTrees 更新前端菜单及按钮(包括菜单和按钮调用的api)
func (a *Menu) UpdateMenuTrees(c *gin.Context, parentID *uint64, list schema.MenuTrees) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range list {
			item.ParentID = parentID
			ButtonTrees := item.Buttons
			sonMenus := item.SonMenus
			item.Buttons = nil
			item.SonMenus = nil
			item.Name = strings.TrimSpace(item.Name)
			if item.Name == "" {
				return errors.New("菜单名称不能为空")
			}
			if item.UUID == "" {
				// 创建菜单及菜单关联的Api
				fmt.Println("创建菜单及菜单关联的Api:", item.Name)
				if menuIDResult, err := a.CreateMenuTree(c, *item); err != nil {
					return errors.New(fmt.Sprintf("创建菜单: %s 及菜单关联的Api失败", item.Name))
				} else {
					item.ID = menuIDResult.ID
				}
			} else {
				// 查询该UUID是否存在
				MenuQueryResult, err := a.Query(c, schema.MenuQueryParam{
					UUID: item.UUID,
				})
				if err != nil {
					return errors.New("查询菜单UUID是否存在失败")
				} else if len(MenuQueryResult.Data) != 0 {
					// 菜单存在，更新菜单及菜单关联的Api
					fmt.Println("更新菜单及菜单关联的Api:", item.Name)
					item.ID = MenuQueryResult.Data[0].ID
					if err = a.UpdateMenuTree(c, *item); err != nil {
						return errors.New(fmt.Sprintf("更新菜单: %s 及菜单关联的Api失败", item.Name))
					}
				} else {
					fmt.Println("创建菜单及菜单关联的Api:", item.Name)
					if menuIDResult, err := a.CreateMenuTree(c, *item); err != nil {
						return errors.New(fmt.Sprintf("创建菜单: %s 及菜单关联的Api失败", item.Name))
					} else {
						item.ID = menuIDResult.ID
					}
				}
			}
			// 创建菜单的按钮
			if len(ButtonTrees) != 0 {
				if err := a.createButtons(c, item.ID, nil, ButtonTrees); err != nil {
					return err
				}
			}
			// 创建子菜单
			if len(sonMenus) != 0 {
				if err := a.UpdateMenuTrees(c, &item.ID, sonMenus); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return err
}

// createButtons 更新前端按钮及按钮关联的Api
func (a *Menu) createButtons(c *gin.Context, MenuID uint64, parentID *uint64, list schema.ButtonPres) error {
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range list {
			item.MenuID = MenuID
			item.ParentID = parentID
			sonButtons := item.SonButtons
			item.SonButtons = nil
			item.Name = strings.TrimSpace(item.Name)
			if item.Name == "" {
				return errors.New("按钮名称不能为空")
			}
			if item.UUID == "" {
				// 创建按钮及按钮关联的Api
				fmt.Println("创建按钮及按钮关联的Api:", item.Name)
				item.ID = 0
				if item.UUID == "" {
					item.UUID = common.GetUUID()
				}
				for _, restfulApi := range item.RestfulApis {
					restfulApi.ID = 0
					if restfulApi.UUID == "" {
						restfulApi.UUID = common.GetUUID()
					}
				}
				if BtnIDResult, err := a.ButtonModel.CreateButtonPre(*item); err != nil {
					return errors.New(fmt.Sprintf("创建按钮: %s 及按钮关联的Api失败", item.Name))
				} else {
					item.ID = BtnIDResult.ID
				}
			} else {
				// 查询按钮是否存在
				ButtonQueryResult, err := a.ButtonModel.Query(schema.ButtonQueryParam{
					UUID: item.UUID,
				})
				if err != nil {
					return errors.New("查询按钮是否存在失败")
				} else if len(ButtonQueryResult.Data) != 0 {
					// 更新按钮的Api信息
					if len(item.RestfulApis) != 0 {
						if err = a.UpdateRestfulApis(c, &item.RestfulApis); err != nil {
							return errors.New(fmt.Sprintf("更新按钮: %s 关联的Api失败", item.Name))
						}
					}
					// 更新按钮及按钮关联的Api
					fmt.Println("更新按钮及按钮关联的Api:", item.Name)
					item.ID = ButtonQueryResult.Data[0].ID
					if err = a.ButtonModel.UpdateButtonPre(*item); err != nil {
						return errors.New(fmt.Sprintf("更新按钮: %s 及按钮关联的Api失败", item.Name))
					}
				} else {
					// 创建按钮及按钮关联的Api
					fmt.Println("创建按钮及按钮关联的Api:", item.Name)
					item.ID = 0
					if item.UUID == "" {
						item.UUID = common.GetUUID()
					}
					for _, restfulApi := range item.RestfulApis {
						restfulApi.ID = 0
						if restfulApi.UUID == "" {
							restfulApi.UUID = common.GetUUID()
						}
					}
					if BtnIDResult, err := a.ButtonModel.CreateButtonPre(*item); err != nil {
						return errors.New(fmt.Sprintf("创建按钮: %s 及按钮关联的Api失败", item.Name))
					} else {
						item.ID = BtnIDResult.ID
					}
				}
			}
			// 创建子按钮
			if len(sonButtons) != 0 {
				if err := a.createButtons(c, MenuID, &item.ID, sonButtons); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// QueryMenuTree 获取当前用户的菜单树
func (a *Menu) QueryMenuTree() (*schema.MenuTrees, error) {
	return a.MenuModel.QueryMenuTree()
}

func (a *Menu) CreateMenuTree(c *gin.Context, item schema.MenuTree) (*common.IDResult, error) {
	item.ID = 0
	if item.UUID == "" {
		item.UUID = common.GetUUID()
	}
	for _, restfulApi := range item.RestfulApis {
		restfulApi.ID = 0
		if restfulApi.UUID == "" {
			restfulApi.UUID = common.GetUUID()
		}
	}
	return a.MenuModel.CreateMenuTree(item)
}

// UpdateMenuRestfulApis 更新菜单关联的Api
func (a *Menu) UpdateMenuRestfulApis(c *gin.Context, id uint64, item schema.RestfulApis) error {
	// 检查菜单是否存在
	if _, err := a.Get(c, id); err != nil {
		return err
	}

	// 更新菜单关联的Api
	return a.MenuModel.UpdateMenuRestfulApis(id, item)
}

// UpdateMenuTree 更新菜单及菜单关联的Api
func (a *Menu) UpdateMenuTree(c *gin.Context, item schema.MenuTree) error {
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 更新菜单信息
		if err := a.Update(c, item.ID, *item.ToSchemaMenu()); err != nil {
			return err
		}

		// 更新RestfulApi信息
		if len(item.RestfulApis) != 0 {
			if err := a.UpdateRestfulApis(c, &item.RestfulApis); err != nil {
				return err
			}
		}

		// 更新菜单与RestfulApi关联信息
		if err := a.UpdateMenuRestfulApis(c, item.ID, item.RestfulApis); err != nil {
			return err
		}
		return nil
	})
}

// UpdateRestfulApis 更新或新增RestfulApi的信息
func (a *Menu) UpdateRestfulApis(c *gin.Context, items *schema.RestfulApis) error {
	for _, item := range *items {
		if item.UUID == "" {
			// 创建restfulApi
			item.ID = 0
			if IDResult, err := a.RestfulApiModel.Create(*item); err != nil {
				return err
			} else {
				item.ID = IDResult.ID
			}
		} else {
			// 查询restfulApi是否存在
			RestfulApiQueryResult, err := a.RestfulApiModel.Query(schema.RestfulApiQueryParam{
				UUID: item.UUID,
			})
			if err != nil {
				return err
			} else if len(RestfulApiQueryResult.Data) != 0 {
				// 更新restfulApi
				if err = a.RestfulApiModel.UpdateByUUID(item.UUID, map[string]interface{}{
					"method": item.Method,
					"path":   item.Path,
					"memo":   item.Memo,
				}); err != nil {
					return err
				}
			} else {
				// 创建restfulApi
				item.ID = 0
				if IDResult, err := a.RestfulApiModel.Create(*item); err != nil {
					return err
				} else {
					item.ID = IDResult.ID
				}
			}
		}
	}
	return nil
}

// ----------------------------------------PermissionTree--------------------------------------

// GetPermissionTree 编辑权限树时调用的接口
func (a *Menu) GetPermissionTree(c *gin.Context) (*schema.PermissionTree, error) {
	permissionTree := &schema.PermissionTree{
		MenuTrees:      make(schema.MenuTrees, 0),
		DisabledFields: make(schema.DisabledFields, 0),
	}
	if menuTrees, err := a.MenuModel.QueryMenuTree(); err != nil {
		return permissionTree, errors.New("获取菜单树失败")
	} else {
		permissionTree.MenuTrees = *menuTrees.SortMenuTrees().Init()
	}

	if fieldResult, err := a.DisabledFieldModel.Query(schema.DisabledFieldQueryParam{}); err != nil {
		return permissionTree, errors.New("获取特殊接口失败")
	} else {
		permissionTree.DisabledFields = *fieldResult.Data.Init()
	}
	return permissionTree, nil
}

func (a *Menu) UpdatePermissionTree(c *gin.Context, item schema.PermissionTree) error {
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 菜单树更新
		for _, menuTree := range item.MenuTrees {
			parentId := menuTree.ParentID
			if menuTree.ParentID != nil {
				if *parentId == 0 {
					parentId = nil
				}
			}
			if err := a.UpdateMenuTrees(c, parentId, schema.MenuTrees{menuTree}); err != nil {
				return err
			}
		}
		// 可禁用字段更新
		if len(item.DisabledFields) != 0 {
			for _, disabledField := range item.DisabledFields {
				disabledField.KeyName = strings.TrimSpace(disabledField.KeyName)
				disabledField.KeyValue = strings.TrimSpace(disabledField.KeyValue)
				if disabledField.KeyName == "" {
					return errors.New("可禁用字段名称不能为空")
				} else if disabledField.KeyValue == "" {
					return errors.New("可禁用字段值不能为空")
				}
				disabledField.ID = 0
				if disabledField.UUID == "" {
					// 新增可禁用字段
					disabledField.UUID = common.GetUUID()
					disabledField.Creator = ginx.GetUserID(c)
					if _, err := a.DisabledFieldModel.Create(*disabledField); err != nil {
						return errors.New(fmt.Sprintf("创建可禁用字段: %s 失败", disabledField.KeyName))
					}
				} else {
					// 检查可禁用字段是否存在
					DisabledFieldQueryResult, err := a.DisabledFieldModel.Query(schema.DisabledFieldQueryParam{
						UUID: disabledField.UUID,
					})
					if err != nil {
						return errors.New(fmt.Sprintf("创建可禁用字段: %s 失败", disabledField.KeyName))
					} else if len(DisabledFieldQueryResult.Data) != 0 {
						// 更新可禁用字段
						if err = a.DisabledFieldModel.UpdateByID(DisabledFieldQueryResult.Data[0].ID, map[string]interface{}{
							"key_name":  disabledField.KeyName,
							"key_value": disabledField.KeyValue,
							"memo":      disabledField.Memo,
						}); err != nil {
							return errors.New(fmt.Sprintf("更新可禁用字段: %s 失败", disabledField.KeyName))
						}
					} else {
						// 新增可禁用字段
						disabledField.Creator = ginx.GetUserID(c)
						if _, err = a.DisabledFieldModel.Create(*disabledField); err != nil {
							return errors.New(fmt.Sprintf("创建可禁用字段: %s 失败", disabledField.KeyName))
						}
					}
				}
			}
		}
		return nil
	})
}

// UpdatePermissionTreeFile 使用DB中的数据替换配置文件内容
func (a *Menu) UpdatePermissionTreeFile(c *gin.Context) error {
	// 判断文件是否存在
	filePath := common.SysConfig.MenuFile
	// 以重写方式打开配置文件，若不存在则创建该文件(不会创建目录)
	file, errByOpenFile := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if errByOpenFile != nil {
		return errors.New("打开权限树配置文件失败")
	}
	defer file.Close()
	// 获取Json文件内容
	menuTreeJson, err := a.GetMenuJson(c)
	if err != nil {
		return err
	}
	if fileContent, err := json.Marshal(menuTreeJson); err != nil {
		return errors.New("获取配置文件失败")
	} else {
		if err = ioutil.WriteFile(filePath, fileContent, 0666); err != nil {
			return errors.New("获取配置文件失败")
		}
	}
	return nil
}

// GetPermissionTreeForCreateRole 创建角色时调用的接口
func (a *Menu) GetPermissionTreeForCreateRole(c *gin.Context) (*schema.PermissionTree, error) {
	permissionTree := &schema.PermissionTree{
		MenuTrees:      make(schema.MenuTrees, 0),
		DisabledFields: make(schema.DisabledFields, 0),
	}
	if menuTrees, err := a.MenuModel.QueryMenuTreeForCreateRole(); err != nil {
		return permissionTree, errors.New("获取菜单树失败")
	} else {
		permissionTree.MenuTrees = *menuTrees.SortMenuTrees().Init()
	}

	if fieldResult, err := a.DisabledFieldModel.Query(schema.DisabledFieldQueryParam{}); err != nil {
		return permissionTree, errors.New("获取特殊接口失败")
	} else {
		permissionTree.DisabledFields = *fieldResult.Data.Init()
	}
	return permissionTree, nil
}

// ----------------------------------------MenuTreeJson--------------------------------------

func (a *Menu) GetMenuJson(c *gin.Context) (*schema.MenuTreeJson, error) {
	menuTreeJson := schema.MenuTreeJson{}
	// 获取权限树
	if permissionTree, err := a.GetPermissionTree(c); err != nil {
		return &menuTreeJson, err
	} else {
		menuTreeJson.MenuTrees = *permissionTree.MenuTrees.SortMenuTrees()
		menuTreeJson.DisabledFields = permissionTree.DisabledFields
	}

	// 获取菜单版本号
	if systemBaseConfig, err := a.SystemConfigModel.First(); err != nil {
		return &menuTreeJson, errors.New("获取系统版本号失败")
	} else {
		menuTreeJson.MenuVersion = systemBaseConfig.MenuVersion
	}
	return &menuTreeJson, nil
}

// ---------------------------------------- Params  Validate --------------------------------------

func (a *Menu) MenuParamsCheck(c *gin.Context, item *schema.Menu) error {
	// 参数检查
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return errors.New("菜单名称不能为空")
	}

	// 检查父菜单ID
	if item.ParentID != nil {
		parentID := uint64(0)
		if *item.ParentID == parentID {
			item.ParentID = nil
		}
	}

	return nil
}
