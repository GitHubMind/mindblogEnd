package blog

import service "blog/server"

type ApiGroup struct {
	blogCurd
}

var (
	blogService = service.ServiceGroupApp.BlogService
)
