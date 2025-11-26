package ioc

import (
	"xiaoweishu/internal/service/sms"
	"xiaoweishu/internal/service/sms/memory"
)

func InitSmsService() sms.Service {
	return memory.NewService()
}
