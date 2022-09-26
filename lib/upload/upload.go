package upload

import (
	"blog/global"
	"mime/multipart"
)

type OSS interface {
	//上传文件
	UploadFile(file *multipart.FileHeader) (string, string, error)
	//珊瑚文件
	DeleteFile(key string) error
}

func NewOss() OSS {
	switch global.GM_CONFIG.System.OssType {
	//case "local":
	//	return &Local{}
	//case "qiniu":
	//	return &Qiniu{}
	//case "tencent-cos":
	//	return &TencentCOS{}
	//case "aliyun-oss":
	//	return &AliyunOSS{}
	//case "huawei-obs":
	//	return HuaWeiObs
	//case "aws-s3":
	//	return &AwsS3{}
	default:
		return &Local{}
	}
}
