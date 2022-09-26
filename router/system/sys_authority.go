package system

import (
	"blog/api"
	"blog/middleware"
	"github.com/gin-gonic/gin"
)

type AuthorityRouter struct{}

func (s *AuthorityRouter) InitRouter(Router *gin.RouterGroup) {
	authorityRouterNoRecord := Router.Group("authority")
	authorityRouter := Router.Group("authority").Use(middleware.OperationRecord())
	authorityApi := api.ApiGroupApp.SystemApi.AuthorityApi
	{

		authorityRouter.POST("createAuthority", authorityApi.CreateAuthority)   // 创建角色
		authorityRouter.POST("deleteAuthority", authorityApi.DeleteAuthority)   // 删除角色
		authorityRouter.PUT("updateAuthority", authorityApi.UpdateAuthority)    // 更新角色
		authorityRouter.POST("copyAuthority", authorityApi.CopyAuthority)       // 拷贝角色
		authorityRouter.POST("setDataAuthority", authorityApi.SetDataAuthority) // 设置角色资源权限
	}
	{
		authorityRouterNoRecord.POST("getAuthorityList", authorityApi.GetAuthorityList) // 获取角色列表
	}
}
