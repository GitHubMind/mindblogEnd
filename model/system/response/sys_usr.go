package response

import (
	"blog/model/system"
	"github.com/golang-jwt/jwt/v4"
)

type SysUserResponse struct {
	User system.SysUser `json:"user"`
}

//type LoginResponse struct {
//	User      system.SysUser `json:"user"`
//	Token     string         `json:"token"`
//	ExpiresAt int64          `json:"expiresAt"`
//}

type Register struct {
}
type LoginResponse struct {
	User      system.SysUser   `json:"user"`
	Token     string           `json:"token"`
	ExpiresAt *jwt.NumericDate `json:"expiresAt"`
}
