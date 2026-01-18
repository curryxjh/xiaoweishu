package ioc

import (
	"github.com/redis/go-redis/v9"
	"time"
	"xiaoweishu/internal/service/sms"
	"xiaoweishu/internal/service/sms/memory"
	"xiaoweishu/internal/service/sms/ratelimit"
	limiter "xiaoweishu/internal/pkg/ratelimit"
	"xiaoweishu/internal/service/sms/retryable"
)

func InitSmsService(cmd redis.Cmdable) sms.Service {
	svc := ratelimit.NewRatelimitSMSService(memory.NewService(), limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 100))
	return retryable.NewService(svc, 3)
}
