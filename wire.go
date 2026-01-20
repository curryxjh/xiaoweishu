//go:build wireinject

package main

import (
	"xiaoweishu/internal/repository"
	"xiaoweishu/internal/repository/cache"
	"xiaoweishu/internal/repository/dao"
	"xiaoweishu/internal/service"
	"xiaoweishu/internal/web"
	ijwt "xiaoweishu/internal/web/jwt"
	"xiaoweishu/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// DB
		ioc.InitDB,
		// Cache
		ioc.InitRedis,
		// DAO
		dao.NewUserDao,
		cache.NewUserCache,
		cache.NewCodeCache,
		// Repository
		repository.NewUserRepository,
		repository.NewCodeRepository,
		// Service
		service.NewUserService,
		service.NewCodeService,
		ioc.InitSmsService,
		ioc.InitOauth2WechatService,

		// Handler
		ijwt.NewRedisJwtHandler,
		web.NewUserHandler,
		ioc.NewWechatHandlerConfig,
		web.NewOauth2WechatHandler,

		// middleware
		ioc.InitMiddlewares,

		ioc.InitWebServer,
	)
	return new(gin.Engine)
}
