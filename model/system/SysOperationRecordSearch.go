package system

import "blog/model/commond/request"

type SysOperationRecordSearch struct {
	SysOperationRecord
	request.PageInfo
}
