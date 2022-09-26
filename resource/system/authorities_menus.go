package system

import (
	"blog/model/system"
	server "blog/server/system"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderMenuAuthority = initOrderMenu + initOrderAuthority

type initMenuAuthority struct{}

// auto run
func init() {
	server.RegisterInit(initOrderMenuAuthority, &initMenuAuthority{})
}

func (i *initMenuAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil // do nothing
}

func (i *initMenuAuthority) TableCreated(ctx context.Context) bool {
	return false // always replace
}

func (i initMenuAuthority) InitializerName() string {
	return "sys_menu_authorities"
}

func (i *initMenuAuthority) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, server.ErrMissingDBContext
	}
	//把表缓存进contxnt 了
	authorities, ok := ctx.Value(initAuthority{}.InitializerName()).([]system.SysAuthority)
	if !ok {
		return ctx, errors.Wrap(server.ErrMissingDependentContext, "创建 [菜单-权限] 关联失败, 未找到权限表初始化数据")
	}
	//context 用表的名字去收缩
	menus, ok := ctx.Value(initMenu{}.InitializerName()).([]system.SysBaseMenu)
	if !ok {
		return next, errors.Wrap(errors.New(""), "创建 [菜单-权限] 关联失败, 未找到菜单表初始化数据")
	}
	next = ctx

	// 888
	//
	if err = db.Model(&authorities[0]).Association("SysBaseMenus").Replace(menus[:20]); err != nil {
		return next, err
	}
	if err = db.Model(&authorities[0]).Association("SysBaseMenus").Append(menus[21:]); err != nil {
		return next, err
	}

	// 8881
	menu8881 := menus[:2]
	menu8881 = append(menu8881, menus[7])
	if err = db.Model(&authorities[1]).Association("SysBaseMenus").Replace(menu8881); err != nil {
		return next, err
	}

	// 9528
	if err = db.Model(&authorities[2]).Association("SysBaseMenus").Replace(menus[:12]); err != nil {
		return next, err
	}
	if err = db.Model(&authorities[2]).Association("SysBaseMenus").Append(menus[13:17]); err != nil {
		return next, err
	}
	return next, nil
}

func (i *initMenuAuthority) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	auth := &system.SysAuthority{}
	if ret := db.Model(auth).
		Where("authority_id = ?", 9528).Preload("SysBaseMenus").Find(auth); ret != nil {
		if ret.Error != nil {
			return false
		}
		return len(auth.SysBaseMenus) > 0
	}
	return false
}
