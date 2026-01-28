package ioc

import (
	"time"
	"xiaoweishu/internal/pkg/ratelimit"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	if err := viper.UnmarshalKey("redis", &cfg); err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})

	return client
}

func NewRateLimiter(redisClient redis.Cmdable, interval time.Duration, rate int) ratelimit.Limiter {
	return ratelimit.NewRedisSlidingWindowLimiter(redisClient, interval, rate)
}
