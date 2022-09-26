package request

import "github.com/golang-jwt/jwt/v4"
import uuid "github.com/satori/go.uuid"

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	//开
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UUID        uuid.UUID
	BaseID      uint
	Username    string
	NickName    string
	AuthorityId uint
}
