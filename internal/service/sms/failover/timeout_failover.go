package failover

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"xiaoweishu/internal/service/sms"
)

type TimeoutFailOverSMSService struct {
	svcs []sms.Service
	idx  int32
	// 连续超时个数
	cnt int32

	// 阈值
	// 连续超时超过阈值时, 切换到下一个短信服务
	threshold int32
}

func NewTimeoutFailOverSMSService(svcs []sms.Service, threshold int32) *TimeoutFailOverSMSService {
	return &TimeoutFailOverSMSService{
		svcs:      svcs,
		threshold: threshold,
	}
}

func (t *TimeoutFailOverSMSService) Send(ctx context.Context, tpl string, args []sms.NameArg, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt >= t.threshold {
		// 切换到下一个短信服务
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			atomic.StoreInt32(&t.cnt, 0)
		}
		// 出现并发了
		idx = atomic.LoadInt32(&t.idx)
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, context.DeadlineExceeded):
		atomic.AddInt32(&t.cnt, 1)
		return err
	default:
		// 不知道什么错误
		// 可以考虑换下一个
		// 超时错误可能是偶然发生, 可以在重试
		// 其他错误, 直接切换下一个
		log.Printf("sms failover, idx: %d, err: %v", idx, err)
		return err
	}
}
