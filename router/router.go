package router

import (
	"blog/api"
	"blog/global"
	"blog/middleware"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerfiles "github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

//initrouter 约束
type IRouter interface {
	InitRouter(*gin.RouterGroup)
}

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.StaticFS(global.GM_CONFIG.Local.Path, http.Dir(global.GM_CONFIG.Local.StorePath)) // 为用户头像和文件提供静态地址
	if global.GM_CONFIG.GinConfig.IsOpenSwager {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	r.Use(middleware.Cors())
	r.Use(middleware.LoggerHandler())
	systemRouter := RouterGroupApp.BaseRouter
	menuRouter := RouterGroupApp.MenuRouter
	userRouter := RouterGroupApp.UserRouter
	jwtRouter := RouterGroupApp.JwtRouter
	authorityRouter := RouterGroupApp.AuthorityRouter
	apiRouter := RouterGroupApp.ApiRouter
	casbinRouter := RouterGroupApp.CasbinRouter
	dictoryRouter := RouterGroupApp.DictionaryRouter
	opreationRouter := RouterGroupApp.OpreationRouter
	FileUploadAndDownloadRouter := RouterGroupApp.FileUploadAndDownloadRouter
	BlogRouter := RouterGroupApp.BlogRouter
	PublicGroup := r.Group("")
	{
		PublicGroup.GET("/upload", func(c *gin.Context) {
			test := make(map[string]any)
			test["errno"] = 1 // 只要不等于 0 就行
			test["message"] = "失败信息"
			c.JSON(200, test)
		}).Use(middleware.Cors())
		PublicGroup.GET("/health", func(c *gin.Context) {
			panic("我就是要报错")
			c.JSON(200, "ok")
		})
		PublicGroup.GET("/test", api.ApiGroupApp.SystemApi.AuthorityInit.InitTable)
		PublicGroup.POST("/init/initdb", api.ApiGroupApp.SystemApi.DBApi.InitDB)
		PublicGroup.POST("/init/checkdb", api.ApiGroupApp.SystemApi.DBApi.CheckDB)
		PublicGroup.POST("/user/register", api.ApiGroupApp.SystemApi.BaseApi.Register)
	}
	systemRouter.InitRouter(PublicGroup)
	userRouter.InitRouter(PublicGroup)
	menuRouter.InitRouter(PublicGroup)
	jwtRouter.InitRouter(PublicGroup)
	authorityRouter.InitRouter(PublicGroup)
	apiRouter.InitRouter(PublicGroup)
	casbinRouter.InitRouter(PublicGroup)
	dictoryRouter.InitRouter(PublicGroup)
	opreationRouter.InitRouter(PublicGroup)
	BlogRouter.InitRouter(PublicGroup)
	FileUploadAndDownloadRouter.InitRouter(PublicGroup)
	return r
}
