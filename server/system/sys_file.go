package server

import (
	"blog/global"
	"blog/lib/upload"
	"blog/model/commond/request"
	"blog/model/system"
	"mime/multipart"
	"strings"
)

type FileUploadAndDownloadService struct{}

//@function: UploadFile
//@description: 根据配置文件判断是文件上传到本地或者七牛云
//@param: header *multipart.FileHeader, tag string
//@return: file model.ExaFileUploadAndDownload, err error
func (e *FileUploadAndDownloadService) UploadFile(header *multipart.FileHeader, tag string) (file system.FileUploadAndDownload, err error) {
	oss := upload.NewOss()
	filePath, key, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		panic(err)
	}
	//已什么方式存储进去
	if tag == "0" {
		s := strings.Split(header.Filename, ".")
		f := system.FileUploadAndDownload{
			Url:  filePath,
			Name: header.Filename,
			//截取最后一段作为
			Tag: s[len(s)-1],
			Key: key,
		}
		return f, e.Upload(&f)
	}
	return
}
func (e *FileUploadAndDownloadService) Upload(file *system.FileUploadAndDownload) (err error) {
	return global.GM_DB.Create(file).Error
}
func (e *FileUploadAndDownloadService) GetFileRecordInfoList(info request.PageInfo) (list []system.FileUploadAndDownload, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	keyword := info.Keyword
	db := global.GM_DB.Model(&system.FileUploadAndDownload{})
	//关键词
	if len(keyword) > 0 {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&list).Error
	return
}
func (e *FileUploadAndDownloadService) FindFile(id uint) (system.FileUploadAndDownload, error) {
	var file system.FileUploadAndDownload
	err := global.GM_DB.Where("id = ?", id).First(&file).Error
	return file, err
}
func (e *FileUploadAndDownloadService) EditFileName(file *system.FileUploadAndDownload) error {
	var tmp system.FileUploadAndDownload
	//如果不加&会出问题，会找不到这个数
	return global.GM_DB.Where("id = ? ", file.ID).First(&tmp).Update("name", file.Name).Error
}
func (e *FileUploadAndDownloadService) DeleteFile(file *system.FileUploadAndDownload) (err error) {
	var fileFromDb system.FileUploadAndDownload
	//为什么要找一次呢 如果删除失败不直接报错就可以了吗，因为会担心前端有延时这个特征

	fileFromDb, err = e.FindFile(file.ID)
	if err != nil {
		return
	}

	oss := upload.NewOss()
	//删除加密之后的东西
	err = oss.DeleteFile(fileFromDb.Key)
	if err != nil {
		return
	}
	err = global.GM_DB.Where("id = ?", fileFromDb.ID).Unscoped().Delete(file).Error
	if err != nil {
		return
	}
	return
}
