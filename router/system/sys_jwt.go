package system

import (
	"blog/api"
	"blog/middleware"
	"github.com/gin-gonic/gin"
)

type JwtRouter struct{}

func (s *JwtRouter) InitRouter(Router *gin.RouterGroup) {
	jwtRouter := Router.Group("jwt").Use(middleware.OperationRecord())
	jwtApi := api.ApiGroupApp.SystemApi.JwtApi
	{
		jwtRouter.POST("jsonInBlacklist", jwtApi.JsonInBlacklist) // jwt加入黑名单
	}
}
