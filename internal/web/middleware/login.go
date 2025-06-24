package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		//if c.Request.URL.Path == "/users/login" ||
		//	c.Request.URL.Path == "/users/signup" {
		//	return
		//}
		sess := sessions.Default(c)
		id := sess.Get("userId")

		if id == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")

		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})

		now := time.Now().UnixMilli()
		// 此时还没有刷新, 刚登录，还没有刷新
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}

		// updateTime是有的，存在上次刷新
		updateTimeVal, _ := updateTime.(int64)
		//if !ok {
		//	// 出问题了
		//	c.AbortWithStatus(http.StatusInternalServerError)
		//	return
		//}

		if now-updateTimeVal > 10*1000 {
			sess.Set("update_time", now)
			sess.Save()
		}

	}
}
