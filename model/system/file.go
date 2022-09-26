package system

import (
	"blog/global"
)

type FileUploadAndDownload struct {
	global.GM_MODEL
	Name string `json:"name" gorm:"comment:文件名"` // 文件名
	Url  string `json:"url" gorm:"comment:文件地址"` // 文件地址
	Tag  string `json:"tag" gorm:"comment:文件标签"` // 文件标签
	Key  string `json:"key" gorm:"comment:编号"`   // 编号
}

func (FileUploadAndDownload) TableName() string {
	return "file_upload_and_downloads"
}