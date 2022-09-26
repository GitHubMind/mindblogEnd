package jwtLib

import (
	"blog/global"
	"blog/model/system/request"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.GM_CONFIG.JWT.SigningKey),
	}
}

//创建jwt
func (j *JWT) CreateClaims(baseClaims request.BaseClaims) request.CustomClaims {
	start := time.Now()
	claims := request.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: global.GM_CONFIG.JWT.BufferTime, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			//(t Time) AddDate 只能操作到天数，只能通过这个转换来达到秒数
			NotBefore: jwt.NewNumericDate(time.Unix(time.Now().Unix()-1000, 0)),                             // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Unix(time.Now().Unix()+global.GM_CONFIG.JWT.ExpiresTime, 0)), // 过期时间 7天  配置文件
			Issuer:    global.GM_CONFIG.JWT.Issuer,                                                          // 签名的发行者

		},
	}
	log.Println("jtw 签发耗时:", time.Since(start))
	return claims
}

// 创建一个token
func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		//按照例子丢进去
		return j.SigningKey, nil
	})
	if err != nil {
		//疑惑 为什么不直接用  jwt.Is 就直接判断所有错误了，不过这样可以细分大相对的错误
		if ve, ok := err.(*jwt.ValidationError); ok {
			//token
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
				//超时
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired

				//	no before 在这个前面签发的都无效
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			//唔错 顺利返回对应的参数
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (j *JWT) CreateTokenByOldToken(oldToken string, claims request.CustomClaims) (string, error) {
	v, err, _ := global.GVA_Concurrency_Control.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}
