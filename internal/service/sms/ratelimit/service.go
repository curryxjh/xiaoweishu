package ratelimit

import (
	"context"
	"fmt"
	"xiaoweishu/internal/pkg/ratelimit"
	"xiaoweishu/internal/service/sms"
)

var errLimited = fmt.Errorf("触发了限流")

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (r *RatelimitSMSService) Send(ctx context.Context, tpl string, args []sms.NameArg, numbers ...string) error {
	// 这里可以加新特性

	limited, err := r.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		// 系统错误
		// 可以限流: 保守策略, 说明下游的承载量不行
		// 可以不限: 容错策略, 说明下游的承载量比较大, 可以暂时不限制
		return fmt.Errorf("短信服务判断是否限流出现错误: %w", err)
	}
	if limited {
		return errLimited
	}

	err = r.svc.Send(ctx, tpl, args, numbers...)
	// 这里可以加新特性
	return err
}
