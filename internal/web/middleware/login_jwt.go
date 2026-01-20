package middleware

import (
	"net/http"
	ijwt "xiaoweishu/internal/web/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// LoginJWTMiddlewareBuilder JWT登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: jwtHdl,
	}
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
		tokenStr := l.ExtractToken(c)

		claims := &ijwt.UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("KntbYH88cXPKDRdFrXrQjh5yZpA7c5QQXKh3MHwYFnt2v43wGCy2d8XCSpmwPjFy"), nil
		})

		//token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//	return []byte("KntbYH88cXPKDRdFrXrQjh5yZpA7c5QQXKh3MHwYFnt2v43wGCy2d8XCSpmwPjFy"), nil
		//})

		if err != nil {
			// 未登录
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// err != nil, token != nil
		if token == nil || !token.Valid || claims.Uid == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserAgent != c.Request.UserAgent() {
			// 严重的安全问题
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = l.CheckSession(c, claims.Ssid)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//// 每10s刷新一次
		//now := time.Now()
		//if claims.ExpiresAt.Sub(now) < time.Second*50 {
		//	claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Minute))
		//	tokenStr, err := token.SignedString([]byte("ntbYH88cXPKDRdFrXrQjh5yZpA7c5QQXKh3MHwYFnt2v43wGCy2d8XCSpmwPjFy"))
		//	if err != nil {
		//		log.Println("jwt 续约失败  ", err)
		//	}
		//	c.Header("x-jwt-token", tokenStr)
		//}
		c.Set("claims", claims)
	}
}
