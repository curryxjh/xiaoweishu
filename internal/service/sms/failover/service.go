package failover

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"xiaoweishu/internal/service/sms"
)

type FailoverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSMSService(svcs ...sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

func (f *FailoverSMSService) Send(ctx context.Context, tpl string, args []sms.NameArg, numbers ...string) error {
	// 缺点: 每次都从头开始遍历所有的短信服务, 绝大多数请求都是svcs[0]就能成功, 负载不均衡
	// 如果 len(svcs) 是几十个, 那么轮训的速度很慢
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		// 记录错误日志
		// 需要做好监控
		log.Printf("failover sms send failed, err: %v\n", err)
	}
	return errors.New("all sms send failed")
}

func (f *FailoverSMSService) SendV1(ctx context.Context, tpl string, args []sms.NameArg, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		default:
			// 记录错误日志
			// 需要做好监控
			log.Printf("failover sms send failed, err: %v\n", err)
		}
	}
	return errors.New("all sms send failed")
}
