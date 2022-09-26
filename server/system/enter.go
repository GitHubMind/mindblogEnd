package server

type ServerGroup struct {
	InitDBService
	UserServer
	MenuService
	JwtService
	ApiService
	AuthorityServer
	CasbinService
	BaseMenuService
	DictionaryService
	DictionaryDetailService
	OperationRecordService
	FileUploadAndDownloadService
}
