package lib

import (
	"blog/global"
	jwtLib "blog/lib/jwt"
	systemReq "blog/model/system/request"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
)

func GetClaims(c *gin.Context) (*systemReq.CustomClaims, error) {
	token := c.Request.Header.Get("x-token")
	j := jwtLib.NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.GM_LOG.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) uint {
	//get 就是加了个锁拿map
	if claims, exists := c.Get("claims"); !exists {
		//不存在就去x-token 拿
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.BaseID
		}
	} else {
		//存在就直接翻译过来
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.BaseID
	}
}

// ？这个没有好一点设计模式写一下吗,如此多的沉淀代码
// GetUserUuid 从Gin的Context中获取从jwt解析出来的用户UUID
func GetUserUuid(c *gin.Context) uuid.UUID {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return uuid.UUID{}
		} else {
			return cl.UUID
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.UUID
	}
}

// GetUserAuthorityId 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserAuthorityId(c *gin.Context) uint {
	claims, exists := c.Get("claims")
	log.Println(claims)
	if !exists {
		//如果不存在，就是如果没有传入
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.AuthorityId
		}
	} else {

		waitUse := claims.(*systemReq.CustomClaims)
		//初始化，估计也是0啦
		return waitUse.AuthorityId
	}
}

// GetUserInfo 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserInfo(c *gin.Context) *systemReq.CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse
	}
}
