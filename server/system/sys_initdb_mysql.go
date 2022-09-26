package server

import (
	"blog/config"
	"blog/global"
	"blog/lib"
	"blog/model/system/request"
	"context"
	"errors"
	"fmt"
	"github.com/gookit/color"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type MysqlInitHandler struct{}

func NewMysqlInitHandler() *MysqlInitHandler {
	return &MysqlInitHandler{}
}

// WriteConfig mysql回写配置,再写入文件
func (h MysqlInitHandler) WriteConfig(ctx context.Context) error {
	c, ok := ctx.Value("config").(config.Mysql)
	if !ok {
		return errors.New("mysql config invalid")
	}
	global.GM_CONFIG.System.DbType = "mysql"
	global.GM_CONFIG.Mysql = c
	global.GM_CONFIG.JWT.SigningKey = uuid.NewV4().String()
	cs := lib.StructToMap(global.GM_CONFIG)
	//是这个的问题吗？
	//这个覆盖是可以改变文件的吗？
	for k, v := range cs {
		//这里报错了
		global.GM_VP.Set(k, v)
	}
	return global.GM_VP.WriteConfig()
}

// EnsureDB 创建数据库并初始化 mysql
func (h MysqlInitHandler) EnsureDB(ctx context.Context, conf *request.InitDB) (next context.Context, err error) {
	//选择 携程变量中国内地  dbtype
	if s, ok := ctx.Value("dbtype").(string); !ok || s != "mysql" {

		//如果找不到就报错
		return ctx, ErrDBTypeMismatch
	}
	//开始初始化,拿到该参数
	c := conf.ToMysqlConfig()
	//他打算用着哪里
	next = context.WithValue(ctx, "config", c)
	if c.Dbname == "" {
		return ctx, nil
	} // 如果没有数据库名, 则跳出初始化数据
	//dsn := conf.MysqlEmptyDsn()
	//操作数据库
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", c.Dbname)
	if err = CreateDatabase(c.Dsn(), "mysql", createSql); err != nil {
		//return nil, err
		log.Println("数据库执行失败，因为有可能已经有了", err)
	} // 创建数据库
	var db *gorm.DB
	if db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       c.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: true,    // 根据版本自动配置
	}), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}); err != nil {
		return ctx, err
	}
	//这是用来干嘛的 用来自动生成代码的吗  应该是插件的
	//global.GM_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	next = context.WithValue(next, "db", db)
	return next, err
}

//初始化表
func (h MysqlInitHandler) InitTables(ctx context.Context, inits initSlice) error {
	return createTables(ctx, inits)
}

//初始化数据吗
func (h MysqlInitHandler) InitData(ctx context.Context, inits initSlice) error {
	next, cancel := context.WithCancel(ctx)
	//函数运行完就结束该线程
	defer func(c func()) { c() }(cancel)
	for _, init := range inits {
		//判断是否数据以及初始化了
		if init.DataInserted(next) {
			color.Info.Printf(InitDataExist, Mysql, init.InitializerName())
			continue
		}
		//注入数据
		if n, err := init.InitializeData(next); err != nil {
			color.Info.Printf(InitDataFailed, Mysql, init.InitializerName(), err)
			return err
		} else {
			next = n
			color.Info.Printf(InitDataSuccess, Mysql, init.InitializerName())
		}
	}
	color.Info.Printf(InitSuccess, Mysql)
	return nil
}
