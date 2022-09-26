package global

import (
	"blog/config"
	logger "blog/core/log"
	viperLib "blog/lib/viper"
	"github.com/go-redis/redis/v8"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	GM_DB *gorm.DB
	//多db处理
	GM_DBList map[string]*gorm.DB
	GM_REDIS  *redis.Client
	GM_CONFIG config.Server
	GM_VP     *viper.Viper
	//GM_LOG    *oplogging.Logger
	GM_LOG *zap.SugaredLogger
	//GVA_Timer               timer.Timer = timer.NewTimerTask()
	//缓存库，防止内存击穿
	GVA_Concurrency_Control = &singleflight.Group{}
	//这是一个很有趣的库
	BlackCache local_cache.Cache
	lock       sync.RWMutex
)

func GlobalInit() {
	GM_CONFIG, GM_VP = viperLib.InitViper()
	GM_LOG = logger.GetLogger()
	BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(time.Second * time.Duration(GM_CONFIG.JWT.ExpiresTime)),
	)
}
