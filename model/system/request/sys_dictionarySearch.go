package request

import (
	"blog/model/commond/request"
	"blog/model/system"
)

type SysDictionarySearch struct {
	system.SysDictionary
	request.PageInfo
}
