package async

import (
	"context"
	"xiaoweishu/internal/service/sms"
)

type Service struct {
	svc sms.Service
}

func NewService(svc sms.Service) *Service {
	return &Service{
		svc: svc,
	}
}

func (s Service) Send(ctx context.Context, biz string, args []sms.NameArg, numbers ...string) error {
	//TODO implement me
	panic("implement me")
}
