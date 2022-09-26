package initapp

import (
	"blog/global"
	mysqlGorm "blog/lib/mysql"
	"gorm.io/gorm"
)

func MysqlRegister() {
	global.GM_DBList = map[string]*gorm.DB{}
	for _, v := range global.GM_CONFIG.MysqlList {
		server, err := mysqlGorm.GormMysqlByConfig(&v)
		if err == nil {
			global.GM_DBList[v.Name] = server
		}

	}
}
