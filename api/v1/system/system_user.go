package system

import (
	"blog/global"
	"blog/lib"
	jwtLib "blog/lib/jwt"
	"blog/model/commond/response"
	"blog/model/system"
	"blog/model/system/request"
	sysmRep "blog/model/system/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"strconv"
)

type BaseApi struct {
}

//// @Tags SysUser
//// @Summary 用户注册账号
//// @Produce  application/json
//// @Param data body request.Register true "用户名, 昵称, 密码, 角色ID"
//// @Success 200 {object} response.Response{data=response.SysUserResponse,msg=string} "用户注册账号,返回包括用户信息"
//// @Router /user/register [post]
func (b *BaseApi) Register(c *gin.Context) {
	var v request.Register
	//	验证数据
	_ = c.ShouldBindJSON(&v)
	if err := lib.Verify(v, lib.RegisterVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		log.Println(err)
		return
	}
	//  这个还是很奇怪 可能以后做权限的会碰到
	var authorities []system.SysAuthority
	for _, v := range v.AuthorityIds {
		authorities = append(authorities, system.SysAuthority{
			AuthorityId: v,
		})
	}
	user := &system.SysUser{Username: v.Username, NickName: v.NickName, Password: v.Password, HeaderImg: v.HeaderImg, AuthorityId: v.AuthorityId, Authorities: authorities, Enable: v.Enable}
	userReturn, err := userService.Register(*user)
	if err != nil {
		global.GM_LOG.Error("注册失败!", zap.Error(err))
		response.FailWithDetailed(sysmRep.SysUserResponse{User: userReturn}, "注册失败", c)
	} else {
		response.OkWithDetailed(sysmRep.SysUserResponse{User: userReturn}, "注册成功", c)
	}
}

//// @Tags SysUser
//// @Summary  用户登陆
//// @Produce  application/json
//// @Param data body request.Login true "用户名, 昵称, 密码, 角色ID"
//// @Success 200 {object} response.Response{data=sysmRep.LoginResponse,Token=string} "用户注册账号,返回包括用户信息"
//// @Router /base/login [post]
func (b *BaseApi) Login(c *gin.Context) {
	var v request.Login
	//	验证数据
	_ = c.ShouldBindJSON(&v)
	if err := lib.Verify(v, lib.RegisterVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := &system.SysUser{Username: v.Username, Password: v.Password}
	//debug
	if store.Verify(v.CaptchaId, v.Captcha, true) || global.GM_CONFIG.GinConfig.RunMode == "debug" {
		userReturn, err := userService.Login(user)
		if err != nil {
			global.GM_LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
			response.FailWithMessage("用户名不存在或者密码错误", c)
			return
		} else {
			//这里有问题 回来的时候吹
			log.Println(userReturn)
			if userReturn.Enable != 1 {
				global.GM_LOG.Error("登陆失败! 用户被禁止登录!")
				response.FailWithMessage("用户被禁止登录", c)
				return
			}
			//jwt

			b.TokenNext(c, userReturn)
		}
	} else {
		response.FailWithMessage("jwt验证码错误", c)
	}
}

// 登录以后签发jwt
func (b *BaseApi) TokenNext(c *gin.Context, user *system.SysUser) {
	//签发jwt

	clian := request.BaseClaims{UUID: user.UUID, BaseID: user.ID, Username: user.Username, NickName: user.NickName, AuthorityId: user.AuthorityId}
	jwt := jwtLib.NewJWT()
	claims := jwt.CreateClaims(clian)
	token, err := jwt.CreateToken(claims)
	if err != nil {
		global.GM_LOG.Error("登陆失败! 用户被禁止登录!")
		response.FailWithMessage("用户被禁止登录", c)
		return
	}
	if !global.GM_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(sysmRep.LoginResponse{
			User:  *user,
			Token: token,
			//到时候看看
			ExpiresAt: claims.RegisteredClaims.ExpiresAt,
		}, "登录成功", c)
		return
	}
}

// @Tags SysUser
// @Summary 获取用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取用户信息"
// @Router /user/getUserInfo [get]
func (b *BaseApi) GetUserInfo(c *gin.Context) {
	uuid := lib.GetUserUuid(c)
	if ReqUser, err := userService.GetUserInfo(uuid); err != nil {
		global.GM_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(gin.H{"userInfo": ReqUser}, "获取成功", c)
	}
}

// @Tags SysUser
// @Summary 更改用户权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param x-token header string true "Insert your access token"
// @Param data body request.SetUserAuth true "用户UUID, 角色ID"
// @Success 200 {object} response.Response{msg=string} "设置用户权限"
// @Router /user/setUserAuthority [post]
func (b *BaseApi) SetUserAuthority(c *gin.Context) {
	var sua request.SetUserAuth
	_ = c.ShouldBindJSON(&sua)
	if UserVerifyErr := lib.Verify(sua, lib.SetUserAuthorityVerify); UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}
	//从jtw 解析did
	userID := lib.GetUserID(c)
	if err := userService.SetUserAuthority(userID, sua.AuthorityId); err != nil {
		global.GM_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		claims := lib.GetUserInfo(c)
		//解析信息
		j := &jwtLib.JWT{SigningKey: []byte(global.GM_CONFIG.JWT.SigningKey)} // 唯一签名
		claims.AuthorityId = sua.AuthorityId
		if token, err := j.CreateToken(*claims); err != nil {
			global.GM_LOG.Error("修改失败!", zap.Error(err))
			response.FailWithMessage(err.Error(), c)
		} else {
			//但是其实也没更新，是因为天数是一样的吗
			c.Header("new-token", token)
			//修改了哪里头部
			//
			c.Header("new-expires-at", strconv.FormatInt(claims.ExpiresAt.Unix(), 10))
			response.OkWithMessage("修改成功", c)
		}

	}
}
