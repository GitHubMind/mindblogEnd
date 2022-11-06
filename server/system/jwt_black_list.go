package server

import (
	"blog/global"
	"blog/model/blog"
	"blog/model/system"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type JwtService struct{}

//通过这个方法来定义无状态，那么如何设置存在时间,过多无状态的时间会不会导致数据存储大量无用的东西
func (jwtService *JwtService) JsonInBlacklist(jwtList system.JwtBlacklist) (err error) {
	err = global.GM_DB.Create(&jwtList).Error
	if err != nil {
		return
	}
	//这里报错了 我猜测是空地址 打个断点
	//  c.defaultExpire 估计变成了一个空指针，可能这个包需要初始化
	global.BlackCache.SetDefault(jwtList.Jwt, struct{}{})
	return
}

//@function: IsBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt string
//@return: bool
func (jwtService *JwtService) IsBlacklist(jwt string) bool {
	//一个是从内部 一个是从数据库
	//_, ok := global.BlackCache.Get(jwt)
	//return ok
	err := global.GM_DB.Where("jwt = ?", jwt).First(&system.JwtBlacklist{}).Error
	isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !isNotFound
}

//用于加载
func LoadAll() {
	var data []string
	err := global.GM_DB.Model(&system.JwtBlacklist{}).Select("jwt").Find(&data).Error
	if err != nil {
		global.GM_LOG.Error("加载数据库jwt黑名单失败!", zap.Error(err))
		return
	}
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	} // jwt黑名单 加入 BlackCache 中

}
func LoadLike() {
	var data []blog.LikeAndWatch
	_ = global.GM_DB.Model(&blog.LikeAndWatch{}).Find(&data).Error
	for _, value := range data {
		if value.Like == 1 {
			key := string(value.ArticleID) + value.Ip + "like"
			_, _ = global.GM_REDIS.Set(context.Background(), key, value.Like, 0).Result()
		}
	}

}
