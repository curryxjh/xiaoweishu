package ratelimit

import "context"

type Limiter interface {
	// Limited 有没有触发限流, key 就是限流对象, 比如手机号, IP 等
	// 返回值:
	// 	- bool: 是否触发限流
	// 	- error: 错误信息(限流器是否有错误)
	Limit(ctx context.Context, key string) (bool, error)
}
