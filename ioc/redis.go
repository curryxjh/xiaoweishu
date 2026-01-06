package ioc

import "github.com/redis/go-redis/v9"

func InitRedis() redis.Cmdable {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
	})

	return client
}
