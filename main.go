package main

import (
	_ "blog/docs"
	"blog/global"
	"blog/initapp"
	"blog/router"
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	ValidationErrorMalformed        uint32 = 1 << iota // Token is malformed
	ValidationErrorUnverifiable                        // Token could not be verified because of signing problems
	ValidationErrorSignatureInvalid                    // Signature validation failed

	// Standard Claim validation errors
	ValidationErrorAudience      // AUD validation failed
	ValidationErrorExpired       // EXP validation failed
	ValidationErrorIssuedAt      // IAT validation failed
	ValidationErrorIssuer        // ISS validation failed
	ValidationErrorNotValidYet   // NBF validation failed
	ValidationErrorId            // JTI validation failed
	ValidationErrorClaimsInvalid // Generic claims validation error
)

// @title blog后台
// @version 1.0 版本
// @description blog后台
// @host	go.gzyezi.top
// @BasePath /
func main() {
	global.GlobalInit()
	initapp.MysqlRegister()
	gin.SetMode(global.GM_CONFIG.GinConfig.RunMode)
	global.GM_DB = initapp.GormMysql() // gorm连接数据库

	if global.GM_DB != nil {
		initapp.AutoTableInit(global.GM_DB) // 初始化表
		db, _ := global.GM_DB.DB()
		defer db.Close()
	}
	//如果是多服务器
	//初始化
	if err := router.InitRouter().Run(
		fmt.Sprintf(":%d", global.GM_CONFIG.GinConfig.Port)); err != nil {
		panic(err)
	}
}
