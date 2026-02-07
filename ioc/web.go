package ioc

import (
	"context"
	"strings"
	"time"
	"xiaoweishu/internal/pkg/ginx/middlewares/logger"
	"xiaoweishu/internal/pkg/ginx/middlewares/ratelimit"
	lg "xiaoweishu/internal/pkg/logger"
	"xiaoweishu/internal/web"
	ijwt "xiaoweishu/internal/web/jwt"
	"xiaoweishu/internal/web/middleware"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler, oauth2Hdl *web.Oauth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	oauth2Hdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable, jwtHdl ijwt.Handler, l lg.LoggerV1) []gin.HandlerFunc {
	bd := logger.NewBuilder(func(c context.Context, al *logger.AccessLog) {
		l.Debug("HTTP 请求", lg.Field{Key: "al", Value: al})
	}).AllowReqBody(true).AllowRespBody()
	viper.OnConfigChange(func(in fsnotify.Event) {
		ok := viper.GetBool("web.logger")
		bd.AllowReqBody(ok)
	})
	return []gin.HandlerFunc{
		corsHdl(),
		bd.Build(),
		ratelimit.NewBuilder(NewRateLimiter(redisClient, time.Second, 100)).Build(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePaths("/users/login").
			IgnorePaths("/users/signup").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/users/sms/login/send").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("oauth2/wechat/callback").
			IgnorePaths("/users/sms/login/verify").Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods: []string{},
		AllowHeaders: []string{"authorization", "Content-Type"},
		//
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		// 是否允许携带用户认证信息，如cookie
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "https://") {
				return true
			}
			return strings.Contains(origin, "http://youcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
