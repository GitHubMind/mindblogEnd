package system

import (
	"blog/api"
	"blog/middleware"
	"github.com/gin-gonic/gin"
)

type OpreationRouter struct {
}

func (s *OpreationRouter) InitRouter(Router *gin.RouterGroup) {
	operationRecordRouter := Router.Group("sysOperationRecord").Use(middleware.OperationRecord())
	opreationApi := api.ApiGroupApp.SystemApi.OperationRecordApi
	{
		operationRecordRouter.POST("createSysOperationRecord", opreationApi.CreateSysOperationRecord)             // 新建SysOperationRecord
		operationRecordRouter.DELETE("deleteSysOperationRecord", opreationApi.DeleteSysOperationRecord)           // 删除SysOperationRecord
		operationRecordRouter.DELETE("deleteSysOperationRecordByIds", opreationApi.DeleteSysOperationRecordByIds) // 批量删除SysOperationRecord
		operationRecordRouter.GET("findSysOperationRecord", opreationApi.FindSysOperationRecord)                  // 根据ID获取SysOperationRecord
		operationRecordRouter.GET("getSysOperationRecordList", opreationApi.GetSysOperationRecordList)            // 获取SysOperationRecord列表
	}
}
