package api

import (
	blog "blog/api/v1/blog"
	"blog/api/v1/system"
)

type ApiGroup struct {
	SystemApi system.ApiGroup
	BlogApi   blog.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
