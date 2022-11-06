package system

import (
	server "blog/server/system"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)
import sys "blog/model/system"

const initOrderMenu = initOrderAuthority + 1

type initMenu struct{}

func (i initMenu) MigrateTable(ctx context.Context) (next context.Context, err error) {
	//TODO implement me
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, server.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&sys.SysBaseMenu{},
		&sys.SysBaseMenuParameter{},
		&sys.SysBaseMenuBtn{},
	)
}

func (i initMenu) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, server.ErrMissingDBContext
	}
	entities := []sys.SysBaseMenu{
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "dashboard", Name: "dashboard", Component: "view/dashboard/index.vue", Sort: 1, Meta: sys.Meta{Title: "仪表盘", Icon: "odometer"}},
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "about", Name: "about", Component: "view/about/index.vue", Sort: 9, Meta: sys.Meta{Title: "关于我们", Icon: "info-filled"}},
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "admin", Name: "superAdmin", Component: "view/superAdmin/index.vue", Sort: 3, Meta: sys.Meta{Title: "超级管理员", Icon: "user"}},
		{MenuLevel: 0, Hidden: false, ParentId: "3", Path: "authority", Name: "authority", Component: "view/superAdmin/authority/authority.vue", Sort: 1, Meta: sys.Meta{Title: "角色管理", Icon: "avatar"}},
		{MenuLevel: 0, Hidden: false, ParentId: "3", Path: "menu", Name: "menu", Component: "view/superAdmin/menu/menu.vue", Sort: 2, Meta: sys.Meta{Title: "菜单管理", Icon: "tickets", KeepAlive: true}},
		{MenuLevel: 0, Hidden: false, ParentId: "3", Path: "api", Name: "api", Component: "view/superAdmin/api/api.vue", Sort: 3, Meta: sys.Meta{Title: "api管理", Icon: "platform", KeepAlive: true}},
		{MenuLevel: 0, Hidden: false, ParentId: "3", Path: "user", Name: "user", Component: "view/superAdmin/user/user.vue", Sort: 4, Meta: sys.Meta{Title: "用户管理", Icon: "coordinate"}},
		{MenuLevel: 0, Hidden: true, ParentId: "0", Path: "person", Name: "person", Component: "view/person/person.vue", Sort: 4, Meta: sys.Meta{Title: "个人信息", Icon: "message"}},
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "example", Name: "example", Component: "view/example/index.vue", Sort: 7, Meta: sys.Meta{Title: "示例文件", Icon: "management"}},
		{MenuLevel: 0, Hidden: false, ParentId: "9", Path: "excel", Name: "excel", Component: "view/example/excel/excel.vue", Sort: 4, Meta: sys.Meta{Title: "excel导入导出", Icon: "takeaway-box"}},
		{MenuLevel: 0, Hidden: false, ParentId: "9", Path: "upload", Name: "upload", Component: "view/example/upload/upload.vue", Sort: 5, Meta: sys.Meta{Title: "媒体库（上传下载）", Icon: "upload"}},
		//{MenuLevel: 0, Hidden: false, ParentId: "9", Path: "breakpoint", Name: "breakpoint", Component: "view/example/breakpoint/breakpoint.vue", Sort: 6, Meta: sys.Meta{Title: "断点续传", Icon: "upload-filled"}},
		//{MenuLevel: 0, Hidden: false, ParentId: "9", Path: "customer", Name: "customer", Component: "view/example/customer/customer.vue", Sort: 7, Meta: sys.Meta{Title: "客户列表（资源示例）", Icon: "avatar"}},
		//{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "systemTools", Name: "systemTools", Component: "view/systemTools/index.vue", Sort: 5, Meta: sys.Meta{Title: "系统工具", Icon: "tools"}},
		//{MenuLevel: 0, Hidden: false, ParentId: "14", Path: "autoCode", Name: "autoCode", Component: "view/systemTools/autoCode/index.vue", Sort: 1, Meta: sys.Meta{Title: "代码生成器", Icon: "cpu", KeepAlive: true}},
		//{MenuLevel: 0, Hidden: false, ParentId: "14", Path: "formCreate", Name: "formCreate", Component: "view/systemTools/formCreate/index.vue", Sort: 2, Meta: sys.Meta{Title: "表单生成器", Icon: "magic-stick", KeepAlive: true}},
		{MenuLevel: 0, Hidden: false, ParentId: "14", Path: "system", Name: "system", Component: "view/systemTools/system/system.vue", Sort: 3, Meta: sys.Meta{Title: "系统配置", Icon: "operation"}},
		{MenuLevel: 0, Hidden: false, ParentId: "3", Path: "dictionary", Name: "dictionary", Component: "view/superAdmin/dictionary/sysDictionary.vue", Sort: 5, Meta: sys.Meta{Title: "字典管理", Icon: "notebook"}},
		{MenuLevel: 0, Hidden: true, ParentId: "3", Path: "dictionaryDetail/:id", Name: "dictionaryDetail", Component: "view/superAdmin/dictionary/sysDictionaryDetail.vue", Sort: 1, Meta: sys.Meta{Title: "字典详情-${id}", Icon: "order"}},
		{MenuLevel: 0, Hidden: false, ParentId: "3", Path: "operation", Name: "operation", Component: "view/superAdmin/operation/sysOperationRecord.vue", Sort: 6, Meta: sys.Meta{Title: "操作历史", Icon: "pie-chart"}},
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "state", Name: "state", Component: "view/system/state.vue", Sort: 8, Meta: sys.Meta{Title: "服务器状态", Icon: "cloudy"}},
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "state", Name: "state", Component: "view/system/state.vue", Sort: 8, Meta: sys.Meta{Title: "服务器状态", Icon: "cloudy"}},
		{MenuLevel: 0, Hidden: false, ParentId: "31", Path: "blog", Name: "blog", Component: "view/blog/routerHolder.vue", Sort: 0, Meta: sys.Meta{Title: "博客", Icon: "aim"}},
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "addArticle", Name: "addArticle", Component: "view/blog/addArticleList/index.vue", Sort: 0, Meta: sys.Meta{Title: "文章列表", Icon: "arrow-left-bold"}},
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "tag", Name: "tag", Component: "view/blog/tagCaregroy/index.vue", Sort: 0, Meta: sys.Meta{Title: "标签和分类", Icon: "aim"}},
		{MenuLevel: 0, Hidden: false, ParentId: "0", Path: "articleDetail/:id", Name: "articleDetail", Component: "view/blog/addArticle/index.vue", Sort: 0, Meta: sys.Meta{Title: "文章详情", Icon: "apple"}},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sys.SysBaseMenu{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i initMenu) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	m := db.Migrator()
	return m.HasTable(&sys.SysBaseMenu{}) &&
		m.HasTable(&sys.SysBaseMenuParameter{}) &&
		m.HasTable(&sys.SysBaseMenuBtn{})
}

func (i initMenu) DataInserted(ctx context.Context) bool {
	//TODO implement me
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}

	if errors.Is(db.Where("path = ?", "autoPkg").First(&sys.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}

func init() {
	server.RegisterInit(initOrderMenu, &initMenu{})
}
func (i initMenu) InitializerName() string {
	return sys.SysBaseMenu{}.TableName()
}
