package server

import (
	"blog/global"
	"blog/model/commond/request"
	"blog/model/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
)

type MenuService struct{}

var MenuServiceApp = new(MenuService)

//@function: GetMenuTree
//@description: 获取动态菜单树
//@param: authorityId string
//@return: menus []system.SysMenu, err error
func (menuService *MenuService) GetMenuTree(authorityId uint) (menus []system.SysMenu, err error) {
	if authorityId == 0 {
		return nil, errors.New("认证不通过")
	}
	menuTree, err := menuService.getMenuTreeMap(authorityId)
	menus = menuTree["0"]
	for i := 0; i < len(menus); i++ {
		//递归把树都东西都递归进去
		err = menuService.getChildrenList(&menus[i], menuTree)
	}
	return menus, err
}
func (menuService *MenuService) getChildrenList(menu *system.SysMenu, treeMap map[string][]system.SysMenu) (err error) {
	menu.Children = treeMap[menu.MenuId]
	for i := 0; i < len(menu.Children); i++ {
		err = menuService.getChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

func (menuService *MenuService) getMenuTreeMap(authorityId uint) (treeMap map[string][]system.SysMenu, err error) {
	var allMenus []system.SysMenu
	var baseMenu []system.SysBaseMenu
	var btns []system.SysAuthorityBtn
	treeMap = make(map[string][]system.SysMenu)
	var SysAuthorityMenus []system.SysAuthorityMenu
	//找对应权限
	err = global.GM_DB.Where("sys_authority_authority_id = ?", authorityId).Find(&SysAuthorityMenus).Error
	if err != nil {
		return
	}
	var MenuIds []string
	for i := range SysAuthorityMenus {
		//收集权限表
		MenuIds = append(MenuIds, SysAuthorityMenus[i].MenuId)
	}
	err = global.GM_DB.Where("id in (?)", MenuIds).Order("sort").Preload("Parameters").Find(&baseMenu).Error
	if err != nil {
		return
	}
	for i := range baseMenu {
		allMenus = append(allMenus, system.SysMenu{
			SysBaseMenu: baseMenu[i],
			AuthorityId: authorityId,
			MenuId:      strconv.Itoa(int(baseMenu[i].ID)),
			Parameters:  baseMenu[i].Parameters,
		})
	}
	err = global.GM_DB.Where("authority_id = ?", authorityId).Preload("SysBaseMenuBtn").Find(&btns).Error
	if err != nil {
		return
	}
	//
	var btnMap = make(map[uint]map[string]uint)
	for _, v := range btns {
		if btnMap[v.SysMenuID] == nil {
			btnMap[v.SysMenuID] = make(map[string]uint)
		}
		btnMap[v.SysMenuID][v.SysBaseMenuBtn.Name] = authorityId
	}
	for _, v := range allMenus {
		v.Btns = btnMap[v.ID]
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return treeMap, err
}

//@function: AddMenuAuthority
//@description: 为角色增加menu树
//@param: menus []model.SysBaseMenu, authorityId string
//@return: err error
func (menuService *MenuService) AddMenuAuthority(menus []system.SysBaseMenu, authorityId uint) (err error) {
	var auth system.SysAuthority
	auth.AuthorityId = authorityId
	auth.SysBaseMenus = menus
	err = AuthorityServiceApp.SetMenuAuthority(&auth)
	return err
}

//@function: GetBaseMenuTree
//@description: 获取基础路由树
//@return: menus []system.SysBaseMenu, err error

func (menuService *MenuService) GetBaseMenuTree() (menus []system.SysBaseMenu, err error) {
	treeMap, err := menuService.getBaseMenuTreeMap()
	menus = treeMap["0"]
	for i := 0; i < len(menus); i++ {
		err = menuService.getBaseChildrenList(&menus[i], treeMap)
	}
	return menus, err
}

//@function: getBaseMenuTreeMap
//@description: 获取路由总树mapGetMenuAuthority
//@return: treeMap map[string][]system.SysBaseMenu, err error

func (menuService *MenuService) getBaseMenuTreeMap() (treeMap map[string][]system.SysBaseMenu, err error) {
	var allMenus []system.SysBaseMenu
	treeMap = make(map[string][]system.SysBaseMenu)
	//err = global.GM_DB.Debug().Order("sort").Preload("MenuBtn").Preload("Parameters").Find(&allMenus).Error
	err = global.GM_DB.Order("sort").Preload("MenuBtn").Preload("Parameters").Find(&allMenus).Error
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return treeMap, err
}

//@function: getBaseChildrenList
//@description: 获取菜单的子菜单
//@param: menu *model.SysBaseMenu, treeMap map[string][]model.SysBaseMenu
//@return: err error

func (menuService *MenuService) getBaseChildrenList(menu *system.SysBaseMenu, treeMap map[string][]system.SysBaseMenu) (err error) {
	menu.Children = treeMap[strconv.Itoa(int(menu.ID))]
	for i := 0; i < len(menu.Children); i++ {
		err = menuService.getBaseChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

//@function: GetMenuAuthority
//@description: 查看当前角色树
//@param: info *request.GetAuthorityId
//@return: menus []system.SysMenu, err error

func (menuService *MenuService) GetMenuAuthority(info *request.GetAuthorityId) (menus []system.SysMenu, err error) {

	var baseMenu []system.SysBaseMenu
	//链接表
	var SysAuthorityMenus []system.SysAuthorityMenu
	//根据id 获取 SysAuthorityMenus
	err = global.GM_DB.Where("sys_authority_authority_id = ?", info.AuthorityId).Find(&SysAuthorityMenus).Error
	if err != nil {
		return
	}
	//剩下就是拿树

	var MenuIds []string
	for i := range SysAuthorityMenus {
		MenuIds = append(MenuIds, SysAuthorityMenus[i].MenuId)
	}
	err = global.GM_DB.Where("id in (?) ", MenuIds).Order("sort").Find(&baseMenu).Error
	//但是这里没有变树 是一位数组
	for i := range baseMenu {
		menus = append(menus, system.SysMenu{
			SysBaseMenu: baseMenu[i],
			AuthorityId: info.AuthorityId,
			MenuId:      strconv.Itoa(int(baseMenu[i].ID)),
			Parameters:  baseMenu[i].Parameters,
		})
	}

	// sql := "SELECT authority_menu.keep_alive,authority_menu.default_menu,authority_menu.created_at,authority_menu.updated_at,authority_menu.deleted_at,authority_menu.menu_level,authority_menu.parent_id,authority_menu.path,authority_menu.`name`,authority_menu.hidden,authority_menu.component,authority_menu.title,authority_menu.icon,authority_menu.sort,authority_menu.menu_id,authority_menu.authority_id FROM authority_menu WHERE authority_menu.authority_id = ? ORDER BY authority_menu.sort ASC"
	// err = global.GM_DB.Raw(sql, authorityId).Scan(&menus).Error
	return menus, err
}

//@function: GetInfoList
//@description: 获取路由分页
//@return: list interface{}, total int64,err error

func (menuService *MenuService) GetInfoList() (list interface{}, total int64, err error) {
	var menuList []system.SysBaseMenu
	treeMap, err := menuService.getBaseMenuTreeMap()
	menuList = treeMap["0"]
	for i := 0; i < len(menuList); i++ {
		err = menuService.getBaseChildrenList(&menuList[i], treeMap)
	}
	return menuList, total, err
}

//@function: AddBaseMenu
//@description: 添加基础路由
//@param: menu model.SysBaseMenu
//@return: error

func (menuService *MenuService) AddBaseMenu(menu system.SysBaseMenu) error {
	//直接插入
	// 拿如果是删除呢
	//记得要留意对应的属性
	if !errors.Is(global.GM_DB.Where("name = ?", menu.Name).First(&system.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	return global.GM_DB.Create(&menu).Error
}

//@function: DeleteBaseMenu
//@description: 删除基础路由
//@param: id float64
//@return: err error

func (baseMenuService *BaseMenuService) DeleteBaseMenu(id int) (err error) {
	//原来删除未必会吧全部删除
	//同时也找到子结点，但是为什么要找呢？ 如果有是不是像之前authiry 操作一样要删除子节点呢
	err = global.GM_DB.Preload("MenuBtn").Preload("Parameters").Where("parent_id = ?", id).First(&system.SysBaseMenu{}).Error
	if err != nil {
		//用一下事务
		var menu system.SysBaseMenu
		var db *gorm.DB
		err = global.GM_DB.Transaction(func(tx *gorm.DB) error {
			//删除 菜单的时候同时会牵连权限吗
			//并且删除响铃的
			db = tx.Debug().Preload("SysAuthoritys").Where("id = ?", id).First(&menu).Delete(&menu)
			err := tx.Debug().Delete(&system.SysBaseMenuParameter{}, "sys_base_menu_id = ?", id).Error
			err = tx.Debug().Delete(&system.SysBaseMenuBtn{}, "sys_base_menu_id = ?", id).Error
			err = tx.Debug().Delete(&system.SysAuthorityBtn{}, "sys_menu_id = ?", id).Error
			if err != nil {
				return err

			}
			global.GM_DB.Debug().Model(&menu).Association("SysAuthoritys").Delete(&menu.SysAuthoritys)
			if db.Error != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "并没有完全删除")
		}

	} else {
		return errors.New("此菜单存在子菜单不可删除")
	}
	return err
}
