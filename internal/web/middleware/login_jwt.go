package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

// LoginJWTMiddlewareBuilder JWT登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		// 使用 JWT 校验
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			// 未登录
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.SplitN(tokenHeader, " ", 2)

		tokenStr := segs[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("KntbYH88cXPKDRdFrXrQjh5yZpA7c5QQXKh3MHwYFnt2v43wGCy2d8XCSpmwPjFy"), nil
		})

		if err != nil {
			// 未登录
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// err != nil, token != nil
		if token == nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
