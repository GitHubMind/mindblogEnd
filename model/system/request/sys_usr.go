package request

import "blog/model/system"

//gva 是按照登陆进去然后注册的。
// User register structure
type Register struct {
	Username     string `json:"userName"example:"123123"`
	Password     string `json:"passWord"example:"123123"`
	NickName     string `json:"nickName"example:"QMPlusUser"gorm:"default:'QMPlusUser'"`
	HeaderImg    string `json:"headerImg"example:"应该是md5加密的"gorm:"default:'https://qmplusimg.henrongyi.top/gva_header.jpg'"`
	AuthorityId  uint   `json:"authorityId"example"888" gorm:"default:888"`
	Enable       int    `json:"enable"` // 账户可行度
	AuthorityIds []uint `json:"authorityIds"`
}

// User login structure
type Login struct {
	Username  string `json:"username" example:"123123"` // 用户名
	Password  string `json:"password" example:"123123"` // 密码
	Captcha   string `json:"captcha"`                   // 验证码
	CaptchaId string `json:"captchaId"`                 // 验证码ID
}

// Modify password structure
type ChangePasswordReq struct {
	ID          uint   `json:"-"`           // 从 JWT 中提取 user id，避免越权
	Password    string `json:"password"`    // 密码
	NewPassword string `json:"newPassword"` // 新密码
}

// Modify  user's auth structure
type SetUserAuth struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
}

// Modify  user's auth structure
type SetUserAuthorities struct {
	ID           uint
	AuthorityIds []uint `json:"authorityIds"` // 角色ID
}

type ChangeUserInfo struct {
	ID            uint                  `gorm:"primarykey"`                                                                           // 主键ID
	NickName      string                `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                            // 用户昵称
	Phone         string                `json:"phone"  gorm:"comment:用户手机号"`                                                          // 用户手机号
	AuthorityIds  []uint                `json:"authorityIds" gorm:"-"`                                                                // 角色ID
	Email         string                `json:"email"  gorm:"comment:用户邮箱"`                                                           // 用户邮箱
	HeaderImg     string                `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	SideMode      string                `json:"sideMode"  gorm:"comment:用户侧边主题"`                                                      // 用户侧边主题
	Enable        int                   `json:"enable" gorm:"comment:冻结用户"`                                                           //冻结用户
	Authorities   []system.SysAuthority `json:"-" gorm:"many2many:sys_user_authority;"`
	GitHubAddress string                `json:"github_address"  gorm:"comment:github地址"` // 用户邮箱
}
