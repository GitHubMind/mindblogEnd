package server

import (
	"blog/global"
	"blog/model/commond/request"
	"blog/model/system"
	"blog/model/system/response"
	"errors"
	"gorm.io/gorm"
	"strconv"
)

type AuthorityServer struct {
}

var ErrRoleExistence = errors.New("存在相同角色id")
var AuthorityServiceApp = new(AuthorityServer)

func (authorityService *AuthorityServer) CreateAuthority(auth system.SysAuthority) (authority system.SysAuthority, err error) {
	var authorityBox system.SysAuthority
	if !errors.Is(global.GM_DB.Where("authority_id = ?", auth.AuthorityId).First(&authorityBox).Error, gorm.ErrRecordNotFound) {
		return auth, ErrRoleExistence
	}
	err = global.GM_DB.Create(&auth).Error
	return auth, err
}

//@parame  request.PageInfo
//@return interface{} int error
func (authorityService *AuthorityServer) GetAuthorityInfoList(info request.PageInfo) (list interface{}, total int64, err error) {
	//偏差
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	//这是直接拿全部？返回的是null
	var a system.SysAuthority
	db := global.GM_DB.Model(&a)
	//获取结点
	err = db.Where("parent_id = ?", "0").Count(&total).Error
	var authority []system.SysAuthority
	//分页拼接
	err = db.Limit(limit).Offset(offset).Preload("DataAuthorityId").Where("parent_id = ?", "0").Find(&authority).Error
	//
	if len(authority) > 0 {
		for k := range authority {
			err = authorityService.findChildrenAuthority(&authority[k])
		}
	}
	return authority, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: findChildrenAuthority
//@description: 查询子角色
//@param: authority *model.SysAuthority
//@return: err error
func (authorityService *AuthorityServer) findChildrenAuthority(authority *system.SysAuthority) (err error) {
	err = global.GM_DB.Preload("DataAuthorityId").Where("parent_id = ?", authority.AuthorityId).Find(&authority.Children).Error
	if len(authority.Children) > 0 {
		for k := range authority.Children {
			//优化处，没必要一个递归查询一次吧,应该一次性吧全部都拿出来再拼装
			err = authorityService.findChildrenAuthority(&authority.Children[k])
		}
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetMenuAuthority
//@description: 菜单与角色绑定
//@param: auth *system.SysAuthority
//@return: error
func (authorityService *AuthorityServer) SetMenuAuthority(auth *system.SysAuthority) error {
	var s system.SysAuthority
	global.GM_DB.Preload("SysBaseMenus").First(&s, "authority_id = ?", auth.AuthorityId)
	// INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`,`sys_base_menu_id`) VALUES (1,1),(1,3),(1,4) ON DUPLICATE KEY UPDATE `sys_authority_authority_id`=`sys_authority_authority_id`
	// 插入数据，如果冲突来就 键值就相等于
	// [rows:1] UPDATE `sys_authorities` SET `updated_at`='2022-08-18 10:51:16.526' WHERE `authority_id` = 1
	//  修改了什么，这是Association 造成的吗
	// DELETE FROM `sys_authority_menus` WHERE `sys_authority_menus`.`sys_authority_authority_id` = 1 AND `sys_authority_menus`.`sys_base_menu_id` NOT IN (1,3,4)
	//  删除掉id 其他对东西
	err := global.GM_DB.Debug().Model(&s).Association("SysBaseMenus").Replace(&auth.SysBaseMenus)
	return err
}

//@function: UpdateAuthority
//@description: 更改一个角色
//@param: auth model.SysAuthority
//@return: authority system.SysAuthority, err error

func (authorityService *AuthorityServer) UpdateAuthority(auth system.SysAuthority) (authority system.SysAuthority, err error) {
	//为什么用first呢
	//err = global.GM_DB.Debug().Where("authority_id = ?", auth.AuthorityId).First(&system.SysAuthority{}).Updates(&auth).Error
	err = global.GM_DB.Model(&system.SysAuthority{}).Where("authority_id = ?", auth.AuthorityId).Updates(&auth).Error

	return auth, err
}

//@function: DeleteAuthority
//@description: 删除角色
//@param: auth *model.SysAuthority
//@return: err error

func (authorityService *AuthorityServer) DeleteAuthority(auth *system.SysAuthority) (err error) {
	//这个一般来说 剩下的数据就只有一个而已 为什么要加 first
	if errors.Is(global.GM_DB.Debug().Preload("Users").First(&auth).Error, gorm.ErrRecordNotFound) {

		return errors.New("该角色不存在")
	}
	//核对2个东西
	// 在使用的
	// 子节点的角色

	if len(auth.Users) != 0 {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(global.GM_DB.Where("authority_id = ?", auth.AuthorityId).First(&system.SysUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(global.GM_DB.Where("parent_id = ?", auth.AuthorityId).First(&system.SysAuthority{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色存在子角色不允许删除")
	}
	db := global.GM_DB.Preload("SysBaseMenus").Where("authority_id = ?", auth.AuthorityId).First(auth)
	err = db.Unscoped().Delete(auth).Error
	if err != nil {
		return
	}
	if len(auth.SysBaseMenus) > 0 {
		err = global.GM_DB.Model(auth).Association("SysBaseMenus").Delete(auth.SysBaseMenus)
		if err != nil {
			return
		}
		// err = db.Association("SysBaseMenus").Delete(&auth)
	} else {
		err = db.Error
		if err != nil {
			return
		}
	}
	//删除两个关联表
	//SysUserAuthority  关联对应的页面
	err = global.GM_DB.Delete(&[]system.SysUserAuthority{}, "sys_authority_authority_id = ?", auth.AuthorityId).Error
	err = global.GM_DB.Delete(&[]system.SysAuthorityBtn{}, "authority_id = ?", auth.AuthorityId).Error
	authorityId := strconv.Itoa(int(auth.AuthorityId))
	CasbinServiceApp.ClearCasbin(0, authorityId)
	return err
}

//@function: CopyAuthority
//@description: 复制一个角色
//@param: copyInfo response.SysAuthorityCopyResponse
//@return: authority system.SysAuthority, err error
//优化项，可以递归复制吗
func (authorityService *AuthorityServer) CopyAuthority(copyInfo response.SysAuthorityCopyResponse) (authority system.SysAuthority, err error) {
	var authorityBox system.SysAuthority
	//判断是否为空
	if !errors.Is(global.GM_DB.Where("authority_id = ?", copyInfo.Authority.AuthorityId).First(&authorityBox).Error, gorm.ErrRecordNotFound) {
		return authority, ErrRoleExistence
	}
	//

	copyInfo.Authority.Children = []system.SysAuthority{}
	menus, err := MenuServiceApp.GetMenuAuthority(&request.GetAuthorityId{AuthorityId: copyInfo.OldAuthorityId})
	if err != nil {
		return
	}
	var baseMenu []system.SysBaseMenu
	for _, v := range menus {
		intNum, _ := strconv.Atoi(v.MenuId)
		v.SysBaseMenu.ID = uint(intNum)
		baseMenu = append(baseMenu, v.SysBaseMenu)
	}
	copyInfo.Authority.SysBaseMenus = baseMenu
	//通过这个可以同时整合很多many2many的东西
	err = global.GM_DB.Create(&copyInfo.Authority).Error
	if err != nil {
		return
	}

	var btns []system.SysAuthorityBtn

	err = global.GM_DB.Find(&btns, "authority_id = ?", copyInfo.OldAuthorityId).Error
	if err != nil {
		return
	}
	if len(btns) > 0 {
		for i := range btns {
			btns[i].AuthorityId = copyInfo.Authority.AuthorityId
		}
		//创建对应按钮
		err = global.GM_DB.Create(&btns).Error

		if err != nil {
			return
		}
	}
	//对接这里
	paths := CasbinServiceApp.GetPolicyPathByAuthorityId(copyInfo.OldAuthorityId)
	err = CasbinServiceApp.UpdateCasbin(copyInfo.Authority.AuthorityId, paths)
	if err != nil {
		_ = authorityService.DeleteAuthority(&copyInfo.Authority)
	}
	return copyInfo.Authority, err
}

//@function: SetDataAuthority
//@description: 设置角色资源权限
//@param: auth model.SysAuthority
//@return: error

func (authorityService *AuthorityServer) SetDataAuthority(auth system.SysAuthority) error {
	var s system.SysAuthority
	global.GM_DB.Preload("DataAuthorityId").First(&s, "authority_id = ?", auth.AuthorityId)
	//
	err := global.GM_DB.Model(&s).Association("DataAuthorityId").Replace(&auth.DataAuthorityId)
	return err
}
