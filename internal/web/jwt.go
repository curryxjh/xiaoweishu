package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type jwtHandler struct {
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明自己要放进 token 的数据
	Uid int64

	UserAgent string
}

func (h jwtHandler) SetJWTToken(c *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uid,
		UserAgent: c.Request.UserAgent(),
	}

	// 使用 JWT 设置登录状态
	// 生成一个 JWT
	//token := jwt.New(jwt.SigningMethodHS512)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("KntbYH88cXPKDRdFrXrQjh5yZpA7c5QQXKh3MHwYFnt2v43wGCy2d8XCSpmwPjFy"))
	if err != nil {
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}
