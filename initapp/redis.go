package initapp

import (
	"blog/global"
	server "blog/server/system"
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func RedisInit() {
	redisCfg := global.GM_CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.GM_LOG.Error("redis connect ping failed, err:", zap.Error(err))
	} else {
		global.GM_LOG.Info("redis connect ping response:", zap.String("pong", pong))
		global.GM_REDIS = client
	}
	// 从db加载jwt数据
	if global.GM_DB != nil {
		server.LoadAll()
	}
	//并且加载jwt数据
}
