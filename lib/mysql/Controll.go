package mysqlGorm

import (
	"blog/config"
	"blog/lib"
	"errors"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

var (
	DB *gorm.DB
)

func gormMysql() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       Dsn(), // DSN data resource name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	return db, err
}
func Dsn() string {
	log.Print(viper.GetString("GormMysql"))
	return viper.GetString("GormMysql")
}
func Create() *gorm.DB {
	if DB != nil {
		sqlDB, err := DB.DB()
		err = sqlDB.Ping()
		defer sqlDB.Close()
		//如果没有问题
		if err == nil {
			return DB
		}
	}
	//否则每3秒重试
	DB, err := gormMysql()
	for err != nil {
		DB, err = gormMysql() //先修改再延迟
		time.Sleep(1 * time.Second)
	}
	return DB
}

/*
 *  用来配置vip的,难道是单线程的？
 */
var mutx sync.Mutex

func GormMysqlByConfig(m *config.GeneralDB) (*gorm.DB, error) {

	if m.Name == "" {
		return nil, errors.New("need name by viper.config ")
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn, // DSN data resource name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	var gc lib.GormConfig
	if db, err := gorm.Open(mysql.New(mysqlConfig), gc.Config(m)); err != nil {
		log.Print("mysql:")
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db, nil
	}
}
