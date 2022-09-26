package system

import "blog/global"

type JwtBlacklist struct {
	global.GM_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
