package system

import (
	"blog/global"
	"blog/model/commond/request"
	"blog/model/commond/response"
	"blog/model/system"
	sysRes "blog/model/system/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"time"
)

type FileUploadAndDownloadApi struct{}

//用于记录 time.After的
type MarkUploadType struct {
	Ticker    *time.Ticker
	Id        uint
	CloseChan chan bool
}

var MarkUpload map[string]MarkUploadType

func init() {
	MarkUpload = make(map[string]MarkUploadType)
}

// @Tags FileUploadAndDownload
// @Summary 上传文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "上传文件示例"
// @Success 200 {object} response.Response{data=response.FileResponse,msg=string} "上传文件示例,返回包括文件详情"
// @Router /fileUploadAndDownload/upload [post]
func (receiver *FileUploadAndDownloadApi) UploadFile(c *gin.Context) {
	//var file system.FileUploadAndDownload
	//预留标志，用于存储不同方式的
	tag := c.DefaultQuery("tag", "0")
	//core的原因 应该去排查一下
	_, header, err := c.Request.FormFile("file")
	//header.Content 里面的东西了
	if err != nil {
		global.GM_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	file, err := fileUploadAndDownloadService.UploadFile(header, tag) // 文件上传后拿到文件路径
	if err != nil {
		global.GM_LOG.Error("修改数据库链接失败!", zap.Error(err))
		response.FailWithMessage("修改数据库链接失败", c)
		return
	}
	response.OkWithDetailed(sysRes.FileResponse{File: file}, "上传成功", c)
}

// @Tags FileUploadAndDownload
// @Summary 上传文件但会一定时间删除
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "上传文件示例"
// @Success 200 {object} response.Response{data=response.FileResponse,msg=string} "上传文件示例,返回包括文件详情"
// @Router /fileUploadAndDownload/UploadWillDelete [post]
func (receiver *FileUploadAndDownloadApi) UploadWillDelete(c *gin.Context) {
	//var file system.FileUploadAndDownload
	//预留标志，用于存储不同方式的
	tag := c.DefaultQuery("tag", "0")
	//core的原因 应该去排查一下
	_, header, err := c.Request.FormFile("file")
	//header.Content 里面的东西了
	if err != nil {
		global.GM_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	file, err := fileUploadAndDownloadService.UploadFile(header, tag) // 文件上传后拿到文件路径
	if err != nil {
		global.GM_LOG.Error("修改数据库链接失败!", zap.Error(err))
		response.FailWithMessage("修改数据库链接失败", c)
		return
	}
	//添加进入
	go func() {
		//不用锁 因为每一个url都会不一样
		MarkUpload[file.Url] = MarkUploadType{Ticker: time.NewTicker(2 * 60 * time.Second), Id: file.ID, CloseChan: make(chan bool)}
		//MarkUpload[file.Url].Ticker.Stop()
		for {
			select {
			case <-MarkUpload[file.Url].Ticker.C:
				var del system.FileUploadAndDownload
				del.ID = MarkUpload[file.Url].Id
				err := fileUploadAndDownloadService.DeleteFile(&del)
				if err != nil {
					global.GM_LOG.Error("删除照片!", zap.Error(err))
				}

				MarkUpload[file.Url].Ticker.Stop()
			case stop := <-MarkUpload[file.Url].CloseChan:
				if stop {
					return
				}
			}

		}

	}()

	response.OkWithDetailed(sysRes.FileResponse{File: file}, "上传成功", c)
}

// @Tags ExaFileUploadAndDownload
// @Summary 分页文件列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页文件列表,返回包括列表,总数,页码,每页数量"
// @Router /fileUploadAndDownload/getFileList [post]
func (b *FileUploadAndDownloadApi) GetFileList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	list, total, err := fileUploadAndDownloadService.GetFileRecordInfoList(pageInfo)
	log.Println(list)
	if err != nil {
		global.GM_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags ExaFileUploadAndDownload
// @Summary 删除文件
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body system.FileUploadAndDownload true "传入文件里面id即可"
// @Success 200 {object} response.Response{msg=string} "删除文件"
// @Router /fileUploadAndDownload/deleteFile [post]
func (b *FileUploadAndDownloadApi) DeleteFile(c *gin.Context) {
	var file system.FileUploadAndDownload
	_ = c.ShouldBindJSON(&file)
	if err := fileUploadAndDownloadService.DeleteFile(&file); err != nil {
		global.GM_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// @Tags ExaFileUploadAndDownload
// @Summary 修改文件名字
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body system.FileUploadAndDownload true "传入要修改的名字"
// @Param file formData file true "上传文件示例"
// @Success 200 {object} response.Response{data=sysRes.FileResponse,msg=string} "上传文件示例,返回包括文件详情"
// @Router /fileUploadAndDownload/editFileName [post]
// EditFileName 编辑文件名或者备注
func (b *FileUploadAndDownloadApi) EditFileName(c *gin.Context) {
	var file system.FileUploadAndDownload
	_ = c.ShouldBindJSON(&file)
	if err := fileUploadAndDownloadService.EditFileName(&file); err != nil {
		global.GM_LOG.Error("编辑失败!", zap.Error(err))
		response.FailWithMessage("编辑失败", c)
		return
	}
	response.OkWithMessage("编辑成功", c)
}
