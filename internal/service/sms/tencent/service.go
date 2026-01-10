package tencent

import (
	"context"
	"fmt"
	"xiaoweishu/internal/pkg/ratelimit"
	mysms "xiaoweishu/internal/service/sms"

	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
	limiter  ratelimit.Limiter
}

func NewService(client *sms.Client, appId, signName string, limiter ratelimit.Limiter) *Service {
	return &Service{
		appId:    &appId,
		signName: &signName,
		client:   client,
		limiter:  limiter,
	}
}

// biz 代表的就是templateID
func (s *Service) Send(ctx context.Context, biz string, args []mysms.NameArg, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SenderId = s.appId
	req.SignName = s.signName
	req.TemplateId = &biz
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = s.toNameArgPtrSlice(args)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *status.Code != "Ok" {
			return fmt.Errorf("短信发送失败, code: %s, msg: %s\n", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	var dst []*string
	for idx, _ := range src {
		dst = append(dst, &src[idx])
	}
	return dst
}
func (s *Service) toNameArgPtrSlice(src []mysms.NameArg) []*string {
	var dst []*string
	for idx, _ := range src {
		dst = append(dst, &src[idx].Val)
	}
	return dst
}
