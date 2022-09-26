package system

import (
	"blog/global"
	com "blog/model/commond/response"
	"blog/model/system/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)
import "github.com/mojocn/base64Captcha"

var store = base64Captcha.DefaultMemStore

// Captcha
// @Tags Base
// @Summary 生成验证码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=response.SysCaptchaResponse,msg=string} "生成验证码,返回包括随机数id,base64,验证码长度"
// @Router /base/captcha [post]
func (b BaseApi) Captcha(c *gin.Context) {
	//生成配置
	driver := base64Captcha.NewDriverDigit(global.GM_CONFIG.Captcha.ImgHeight, global.GM_CONFIG.Captcha.ImgWidth, global.GM_CONFIG.Captcha.KeyLong, 0.7, 80)
	//to
	cp := base64Captcha.NewCaptcha(driver, store)
	if id, b64s, err := cp.Generate(); err != nil {
		global.GM_LOG.Error("验证码获取失败!", zap.Error(err))
		com.FailWithMessage("验证码获取失败", c)
	} else {
		com.OkWithDetailed(response.SysCaptchaResponse{
			CaptchaId:     id,
			PicPath:       b64s,
			CaptchaLength: global.GM_CONFIG.Captcha.KeyLong,
		}, "验证码获取成功", c)
	}
}
