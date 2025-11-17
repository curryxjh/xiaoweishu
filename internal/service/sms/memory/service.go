package memory

import (
	"context"
	"fmt"
	"xiaoweishu/internal/service/sms"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tpl string, args []sms.NameArg, numbers ...string) error {
	fmt.Println(args)
	return nil

}
