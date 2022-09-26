package initapp

import (
	"blog/config"
	"blog/global"
	"blog/lib"
	"blog/server/system"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GormMysql() *gorm.DB {
	m := global.GM_CONFIG.Mysql
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	//如果我单独来用会有什么问题
	init := lib.GormConfig{}
	config := config.GeneralDB{}
	config.Path = m.Path
	config.Port = m.Port
	config.Config = m.Config
	config.Dbname = m.Dbname
	config.Username = m.Username
	config.Password = m.Password
	config.Name = m.Name
	config.Dsn = m.Dsn()
	config.MaxIdleConns = m.MaxIdleConns
	config.MaxOpenConns = m.MaxOpenConns
	config.LogMode = m.LogMode
	config.LogZap = m.LogZap
	//创建gva的表
	//system.createDatabase()
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", config.Dbname)
	err := server.CreateDatabase(m.DsnNoDB(), "mysql", createSql)
	if err != nil {
		global.GM_LOG.Error("自动创建表出错了", err)
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), init.Config(&config)); err != nil {

		return nil
	} else {
		sqlDB, _ := db.DB()
		//设置连接数
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
