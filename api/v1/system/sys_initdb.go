package system

import (
	"blog/global"
	"blog/model/commond/response"
	"blog/model/system/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

type DBApi struct{}

// CheckDB
// @Tags CheckDB
// @Summary 测试数据库是否需要初始化
// @Produce  application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "初始化用户数据库"
// @Router /init/checkdb [post]
func (i *DBApi) CheckDB(c *gin.Context) {
	var (
		message  = "前往初始化数据库"
		needInit = true
	)

	if global.GM_DB != nil {
		message = "数据库无需初始化"
		needInit = false
	}
	global.GM_LOG.Info(message)
	//没想明白之前不改
	log.Println(needInit)
	//response.OkWithDetailed(gin.H{"needInit": needInit}, message, c)
	response.OkWithDetailed(nil, message, c)
	return
}

// InitDB
// @Tags InitDB
// @Summary 初始化用户数据库
// @Produce  application/json
// @Param data body request.InitDB true "初始化数据库参数"
// @Success 200 {object} response.Response{data=string} "初始化用户数据库"
// @Router /init/initdb [post]
func (i *DBApi) InitDB(c *gin.Context) {
	if global.GM_DB != nil {
		global.GM_LOG.Error("已存在数据库配置!")
		response.FailWithMessage("已存在数据库配置", c)
		return
	}
	var dbInfo request.InitDB
	if err := c.ShouldBindJSON(&dbInfo); err != nil {
		global.GM_LOG.Error("参数校验不通过!", zap.Error(err))
		response.FailWithMessage("参数校验不通过", c)
		return
	}
	if err := initDBService.InitDB(dbInfo); err != nil {
		global.GM_LOG.Error("自动创建数据库失败!", zap.Error(err))
		response.FailWithMessage("自动创建数据库失败，请查看后台日志，检查后在进行初始化", c)
		return
	}
	response.OkWithMessage("自动创建数据库成功", c)
}
