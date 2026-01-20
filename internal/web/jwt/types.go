package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	ExtractToken(c *gin.Context) string
	SetJwtToken(c *gin.Context, uid int64, ssid string) error
	SetRefreshToken(c *gin.Context, uid int64, ssid string) error
	ClearToken(c *gin.Context) error
	CheckSession(c *gin.Context, ssid string) error
	SetLoginToken(c *gin.Context, uid int64) error
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
	Ssid      string
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}
