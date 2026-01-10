package sms

import "context"

type NameArg struct {
	Val  string
	Name string
}

type Service interface {
	// Send 发送短信 biz 含糊的业务
	Send(ctx context.Context, biz string, args []NameArg, numbers ...string) error
}
