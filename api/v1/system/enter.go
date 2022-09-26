package system

import service "blog/server"

type ApiGroup struct {
	DBApi
	JwtApi
	BaseApi
	//SystemApi
	CasbinApi
	//AutoCodeApi
	SystemApiApi
	AuthorityApi
	DictionaryApi
	AuthorityMenuApi
	OperationRecordApi
	//AutoCodeHistoryApi
	DictionaryDetailApi
	AuthorityInit
	FileUploadAndDownloadApi
}

var (
	fileUploadAndDownloadService = service.ServiceGroupApp.FileUploadAndDownloadService
	apiService                   = service.ServiceGroupApp.ApiService
	jwtService                   = service.ServiceGroupApp.JwtService
	menuService                  = service.ServiceGroupApp.MenuService
	userService                  = service.ServiceGroupApp.UserServer
	initDBService                = service.ServiceGroupApp.InitDBService
	casbinService                = service.ServiceGroupApp.CasbinService
	//autoCodeService         = service.ServiceGroupApp.SystemServiceGroup.AutoCodeService
	baseMenuService   = service.ServiceGroupApp.BaseMenuService
	authorityService  = service.ServiceGroupApp.AuthorityServer
	dictionaryService = service.ServiceGroupApp.DictionaryService
	//systemConfigService     = service.ServiceGroupApp.SystemServiceGroup.SystemConfigService
	operationRecordService = service.ServiceGroupApp.OperationRecordService
	//autoCodeHistoryService  = service.ServiceGroupApp.SystemServiceGroup.AutoCodeHistoryService
	dictionaryDetailService = service.ServiceGroupApp.DictionaryDetailService
	//authorityBtnService     = service.ServiceGroupApp.SystemServiceGroup.AuthorityBtnService
)
