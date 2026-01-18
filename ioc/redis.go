package ioc

import (
	"time"
	"xiaoweishu/internal/pkg/ratelimit"

	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
	})

	return client
}

func NewRateLimiter(redisClient redis.Cmdable, interval time.Duration, rate int) ratelimit.Limiter {
	return ratelimit.NewRedisSlidingWindowLimiter(redisClient, interval, rate)
}
