package ioc

import (
	"xiaoweishu/internal/pkg/logger"
	"xiaoweishu/internal/repository"
	"xiaoweishu/internal/service"

	"go.uber.org/zap"
)

func InitUserHandler(repo repository.UserRepository) service.UserService {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return service.NewUserService(repo, logger.NewZapLogger(l))
}
