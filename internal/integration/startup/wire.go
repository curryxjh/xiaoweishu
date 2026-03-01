//go:build wireinject

package startup

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

var thirdPartySet = wire.NewSet(
	ioc.InitDB, ioc.InitRedis, ioc.InitLogger,
)

var userSvcProvider = wire.NewSet(
	dao.NewUserDao,
	cache.NewUserCache,
	repository.NewUserRepository,
	service.NewUserService,
)

var articleSvcProvider = wire.NewSet(
	repository.NewArticleRepository,
	dao.NewGormArticleDao,
	service.NewArticleService,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdPartySet,
		userSvcProvider,
		articleSvcProvider,
		// DAO
		cache.NewCodeCache,
		// Repository
		repository.NewCodeRepository,
		// Service
		service.NewCodeService,
		ioc.InitSmsService,
		ioc.InitOauth2WechatService,
		// Handler
		ijwt.NewRedisJwtHandler,
		web.NewUserHandler,
		web.NewArticleHandler,
		ioc.NewWechatHandlerConfig,
		web.NewOauth2WechatHandler,

		// middlewares
		ioc.InitMiddlewares,

		ioc.InitWebServer,
	)
	return new(gin.Engine)
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(thirdPartySet, dao.NewGormArticleDao, repository.NewArticleRepository, service.NewArticleService, web.NewArticleHandler)
	return &web.ArticleHandler{}
}
