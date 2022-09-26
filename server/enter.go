package service

import (
	"blog/server/blog"
	server "blog/server/system"
)

type ServiceGroup struct {
	server.ServerGroup
	blog.BlogService
}

var ServiceGroupApp = new(ServiceGroup)
