package sms

import "context"

type NameArg struct {
	Val  string
	Name string
}

type Service interface {
	Send(ctx context.Context, tpl string, args []NameArg, numbers ...string) error
}
