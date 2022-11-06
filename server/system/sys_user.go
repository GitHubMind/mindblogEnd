package server

import (
	"blog/global"
	"blog/lib"
	"blog/model/commond/request"
	"blog/model/system"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type UserServer struct {
}

func (receiver UserServer) Register(u system.SysUser) (userInter system.SysUser, err error) {

	if !errors.Is(global.GM_DB.Where("username = ?", u.Username).First(&userInter).Error, gorm.ErrRecordNotFound) {
		return userInter, errors.New("用户名已注册")
	}
	u.Password = lib.BcryptHash(u.Password)
	u.UUID = uuid.NewV4()
	err = global.GM_DB.Create(&u).Error
	return
}

//关联其他东西
func getInfo(user *system.SysUser) (err error) {
	var SysAuthorityMenus []system.SysAuthorityMenu
	//这里这里直接复制
	err = global.GM_DB.Where("sys_authority_authority_id = ?", user.AuthorityId).Find(&SysAuthorityMenus).Error
	if err != nil {
		return
	}
	var MenuIds []string
	//获取菜单对应id
	for i := range SysAuthorityMenus {
		MenuIds = append(MenuIds, SysAuthorityMenus[i].MenuId)
	}
	////对应的路由处理
	var am system.SysBaseMenu
	//内链接的感觉
	ferr := global.GM_DB.First(&am, "name = ? and id in (?)", user.Authority.DefaultRouter, MenuIds).Error
	// 找到
	if errors.Is(ferr, gorm.ErrRecordNotFound) {
		user.Authority.DefaultRouter = "404"
	}
	return
}

func (receiver UserServer) Login(u *system.SysUser) (user *system.SysUser, err error) {

	//判断数据库是否还能操作
	if nil == global.GM_DB {
		return nil, fmt.Errorf("db not init")
	}
	//对应的
	if !errors.Is(global.GM_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error, gorm.ErrRecordNotFound) {
		if ok := lib.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("密码错误")
		}
		//获取详细信息
		err = getInfo(user)

	} else {
		return user, errors.New("没有找到对应的用户")
	}
	return
}

//@function: GetUserInfo
//@description: 获取用户信息
//@param: uuid uuid.UUID
//@return: err error, user system.SysUser
func (receiver UserServer) GetUserInfo(uuid uuid.UUID) (reqUser *system.SysUser, err error) {
	err = global.GM_DB.Preload("Authorities").Preload("Authority").First(&reqUser, "uuid = ?", uuid).Error
	if err != nil {
		return
	}
	err = getInfo(reqUser)
	return
}

//@function: SetUserAuthority
//@description: 设置一个用户的权限
//@param: uuid uuid.UUID, authorityId string
//@return: err error

func (userService *UserServer) SetUserAuthority(id uint, authorityId uint) (err error) {
	assignErr := global.GM_DB.Where("sys_user_id = ? AND sys_authority_authority_id = ?", id, authorityId).First(&system.SysUserAuthority{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("该用户无此角色")
	}
	err = global.GM_DB.Where("id = ?", id).First(&system.SysUser{}).Update("authority_id", authorityId).Error
	return err
}

func (userService *UserServer) SetUserInfo(req system.SysUser) error {
	return global.GM_DB.Updates(&req).Error
}

//@function: SetUserAuthorities
//@description: 更具id来改变权限
//@param: uuid uuid.UUID
//@return: err error, user system.SysUser
func (userService *UserServer) SetUserAuthorities(id uint, authorityIds []uint) (err error) {
	//事务使用的案例
	return global.GM_DB.Transaction(func(tx *gorm.DB) error {
		//删除 软删除
		TxErr := tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error
		if TxErr != nil {
			return TxErr
		}

		var useAuthority []system.SysUserAuthority
		for _, v := range authorityIds {
			useAuthority = append(useAuthority, system.SysUserAuthority{
				SysUserId: id, SysAuthorityAuthorityId: v,
			})
		}
		//添加新数据
		TxErr = tx.Create(&useAuthority).Error
		if TxErr != nil {
			return TxErr
		}
		//去修改链接存放链接 authority_id
		TxErr = tx.Where("id = ?", id).First(&system.SysUser{}).Update("authority_id", authorityIds[0]).Error
		if TxErr != nil {
			return TxErr
		}
		// 返回 nil 提交事务
		return nil
	})
}

//@function: GetUserInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserServer) GetUserInfoList(info request.PageInfo) (list interface{}, total int64, err error) {
	//分页
	//总数
	//拿到的数据做匹配
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GM_DB.Model(&system.SysUser{})
	var userList []system.SysUser
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Preload("Authorities").Preload("Authority").Find(&userList).Error
	return userList, total, err
}

//@function: DeleteUser
//@description: 删除用户
//@param: id float64
//@return: err error

func (userService *UserServer) DeleteUser(id int) (err error) {
	var user system.SysUser
	//删除 软删除
	err = global.GM_DB.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return err
	}
	err = global.GM_DB.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error
	return err
}

//@function: ChangePassword
//@description: 修改用户密码
//@param: u *model.SysUser, newPassword string
//@return: userInter *model.SysUser,err error

func (userService *UserServer) ChangePassword(u *system.SysUser, newPassword string) (userInter *system.SysUser, err error) {
	var user system.SysUser
	//
	if err = global.GM_DB.Where("id = ?", u.ID).First(&user).Error; err != nil {
		return nil, err
	}
	//
	if ok := lib.BcryptCheck(u.Password, user.Password); !ok {
		return nil, errors.New("原密码错误")
	}
	user.Password = lib.BcryptHash(newPassword)
	err = global.GM_DB.Save(&user).Error
	return &user, err
}

//@function: resetPassword
//@description: 修改用户密码
//@param: ID uint
//@return: err error

func (userService *UserServer) ResetPassword(ID uint) (err error) {
	//grom 正常curd 修改流程
	err = global.GM_DB.Model(&system.SysUser{}).Where("id = ?", ID).Update("password", lib.BcryptHash("123456")).Error
	return err
}
