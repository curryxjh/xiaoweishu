package auth

import (
	"context"
	"errors"
	"xiaoweishu/internal/service/sms"

	"github.com/golang-jwt/jwt/v5"
)

type SMSService struct {
	svc sms.Service
	key string
}

func (s *SMSService) Send(ctx context.Context, biz string, args []sms.NameArg, numbers ...string) error {
	// 权限校验
	var tc Claims
	// 此处能成功解析就说明是对应的业务，是有效的
	token, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.key), nil
	})
	if err != nil {
		return err
	}
	// 校验token是否有效
	if !token.Valid {
		return errors.New("invalid token")
	}
	return s.svc.Send(ctx, tc.Tpl, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string
}
