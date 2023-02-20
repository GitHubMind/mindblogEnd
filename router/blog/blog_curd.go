package blog

import (
	"blog/api"
	"blog/middleware"
	"github.com/gin-gonic/gin"
)

type BlogRouter struct{}

func (s *BlogRouter) InitRouter(Router *gin.RouterGroup) {
	//blog 后台使用的
	blogEndApi := Router.Group("blog").Use(middleware.OperationRecord()).Use(middleware.JWTAuth())

	blogFrontApi := Router.Group("blog").Use(middleware.JWTAuth())
	//apiRouterWithoutRecord := Router.Group("blog")

	//authorityRouter := Router.Group("authority").Use(middleware.OperationRecord())
	//apiRouterApi := api.ApiGroupApp.SystemApi.Blog

	//tag
	{
		blogEndApi.POST("CreateTag", api.ApiGroupApp.BlogApi.CreateTag)   // 创建Api
		blogEndApi.GET("GetTaglist", api.ApiGroupApp.BlogApi.GetTag)      // 创建Api
		blogEndApi.DELETE("DeleteTag", api.ApiGroupApp.BlogApi.DeleteTag) // 创建Api
	}
	//Category
	{
		blogEndApi.POST("CreateCategory", api.ApiGroupApp.BlogApi.CreateCategory)   // 创建Api
		blogEndApi.GET("GetCategorylist", api.ApiGroupApp.BlogApi.GetCategory)      // 创建Api
		blogEndApi.DELETE("DeleteCategory", api.ApiGroupApp.BlogApi.DeleteCategory) // 创建Api
	}
	//article
	{
		blogEndApi.POST("CreateArticle", api.ApiGroupApp.BlogApi.CreateArticle)                           // 创建Api
		blogEndApi.POST("GetSearchArticleList", api.ApiGroupApp.BlogApi.GetSearchArticleList)             // 获取所有文章
		blogEndApi.DELETE("DeleteArticle", api.ApiGroupApp.BlogApi.DeleteArticle)                         // 删除
		blogEndApi.POST("UpdateArticle", api.ApiGroupApp.BlogApi.UpdateArticle)                           // 更新api
		blogEndApi.POST("FindArticle", api.ApiGroupApp.BlogApi.FindArticle)                               // 通过id寻找他api
		blogEndApi.DELETE("DeleteMulArticle", api.ApiGroupApp.BlogApi.DeleteMulArticle)                   // 删除多个文章
		blogEndApi.POST("UpdateArticleContent", api.ApiGroupApp.BlogApi.UpdateArticleContent)             // 修改文章id
		blogEndApi.POST("UpdateArticleContentOnLine", api.ApiGroupApp.BlogApi.UpdateArticleContentOnLine) // 发布文章
	}
	//首页展示(后台）
	{
		blogEndApi.GET("GetRateNumber", api.ApiGroupApp.BlogApi.GetRateNumber)         // 获得好评率
		blogEndApi.GET("GetRateLikeNumber", api.ApiGroupApp.BlogApi.GetRateLikeNumber) // 发布文章
	}
	//前端 不用加密
	{
		blogFrontApi.GET("GetBlogInfoByName", api.ApiGroupApp.BlogApi.GetBlogInfoByName)                   // 通过title来获取用户文章信息
		blogFrontApi.POST("GetBlogSearchArticleList", api.ApiGroupApp.BlogApi.GetBlogSearchArticleList)    // 获取所有文章
		blogFrontApi.POST("FindBlogArticle", api.ApiGroupApp.BlogApi.FindArticle)                          // 通过id寻找文章
		blogFrontApi.GET("GetBlogInfoById", api.ApiGroupApp.BlogApi.GetBlogInfoById)                       // 通过id来获取信息 包括名字Ï
		blogFrontApi.GET("GetBlogCategoryTaglistById", api.ApiGroupApp.BlogApi.GetBlogCategoryTaglistById) // 通过id 来获取用户的信息(标签和类别）
		blogFrontApi.GET("GETLike", api.ApiGroupApp.BlogApi.ClickBlogLike)                                 // 通过title来获取所有信息
		blogFrontApi.DELETE("CancelLike", api.ApiGroupApp.BlogApi.CancelLike)                              // 通过title来获取所有信息
		blogFrontApi.GET("ClickBlog", api.ApiGroupApp.BlogApi.ClickBlog)                                   //  游览量++
	}

}
