package router

import (
	"blog/router/blog"
	"blog/router/system"
)

type RouterGroup struct {
	system.MenuRouter
	system.UserRouter
	system.BaseRouter
	system.JwtRouter
	system.AuthorityRouter
	system.ApiRouter
	system.CasbinRouter
	system.DictionaryRouter
	system.OpreationRouter
	system.FileUploadAndDownloadRouter
	blog.BlogRouter
}

var RouterGroupApp = new(RouterGroup)
