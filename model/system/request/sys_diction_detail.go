package request

import (
	"blog/model/commond/request"
	"blog/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}

