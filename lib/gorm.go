package lib

import (
	"blog/config"
	"blog/global"
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type GormConfig struct {
}

func (GC *GormConfig) Config(m *config.GeneralDB) *gorm.Config {

	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}

	_default := logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags), m.LogZap), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})

	switch m.LogMode {
	case "silent", "Silent":
		config.Logger = _default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = _default.LogMode(logger.Info)
	default:
		config.Logger = _default.LogMode(logger.Info)
	}
	return config
}

type writer struct {
	Writer logger.Writer
	isZap  bool
}

func NewWriter(w logger.Writer, isZap bool) *writer {
	//初始化
	return &writer{Writer: w, isZap: isZap}
}

func (w *writer) Printf(message string, data ...interface{}) {
	//var logZap bool
	//switch global.GVA_CONFIG.System.DbType {
	//case "mysql":
	//	logZap = global.GVA_CONFIG.Mysql.LogZap
	//case "pgsql":
	//	logZap = global.GVA_CONFIG.Pgsql.LogZap
	//}
	//但是我想每个设置都不一样
	//log.Println(global.GM_VP)
	if w.isZap {
		global.GM_LOG.Info(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}
