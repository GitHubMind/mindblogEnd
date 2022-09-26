package response

import "blog/model/system"

type FileResponse struct {
	File system.FileUploadAndDownload `json:"file"`
}
