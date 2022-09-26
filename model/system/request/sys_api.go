package request

import (
	"blog/model/commond/request"
	"blog/model/system"
)

// api分页条件查询及排序结构体
type SearchApiParams struct {
	//具体方法
	system.SysApi
	//分页以及关键值
	request.PageInfo
	//下面两个还可以封装
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
