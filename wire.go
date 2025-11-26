//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"xiaoweishu/internal/repository"
	"xiaoweishu/internal/repository/cache"
	"xiaoweishu/internal/repository/dao"
	"xiaoweishu/internal/service"
	"xiaoweishu/internal/web"
	"xiaoweishu/ioc"
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
		// Handler
		web.NewUserHandler,

		// middlewares
		ioc.InitMiddlewares,

		ioc.InitWebServer,
	)
	return new(gin.Engine)
}
