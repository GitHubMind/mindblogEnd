package system

import (
	"blog/model/system"
	"blog/server/system"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

//用这个来做牵引
const initOrderAuthority = initOrderCasbin + 1

type initAuthority struct{}

func (i initAuthority) InitializerName() string {

	return system.SysAuthority{}.TableName()
}

//如果是多服务器呢？
func (i initAuthority) MigrateTable(ctx context.Context) (next context.Context, err error) {
	//TODO implement me
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, server.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&system.SysAuthority{})
}

func (i initAuthority) InitializeData(ctx context.Context) (next context.Context, err error) {
	//初始化数据
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, server.ErrMissingDBContext
	}
	entities := []system.SysAuthority{
		{AuthorityId: 888, AuthorityName: "普通用户", ParentId: 0, DefaultRouter: "dashboard"},
		{AuthorityId: 9528, AuthorityName: "测试角色", ParentId: 0, DefaultRouter: "dashboard"},
		{AuthorityId: 8881, AuthorityName: "普通用户子角色", ParentId: 888, DefaultRouter: "dashboard"},
	}
	if err := db.Create(&entities).Error; err != nil {
		//确实是一个又有趣的做法
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!", system.SysAuthority{}.TableName())
	}
	//Association
	if err := db.Model(&entities[1]).Association("DataAuthorityId").Replace(
		[]*system.SysAuthority{
			{AuthorityId: 9528},
			{AuthorityId: 8881},
		}); err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
			db.Model(&entities[1]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	}
	//绑定数据 初始化完成之后数据就会绑定在这里
	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i initAuthority) TableCreated(ctx context.Context) bool {
	//TODO implement me
	//导入服务器 未必注册和业务的服务器不必是在一起的
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&system.SysAuthority{})
}

func (i initAuthority) DataInserted(ctx context.Context) bool {
	//TODO implement me
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	//为什么要判断 属于ping的于是吗
	if errors.Is(db.Where("authority_id = ?", "8881").
		First(&system.SysAuthority{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}

func init() {
	//想用俺这个函数 必须的规范  ! 经典接口引用
	server.RegisterInit(initOrderAuthority, &initAuthority{})
}
