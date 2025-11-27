package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"xiaoweishu/internal/repository"
	"xiaoweishu/internal/service/sms"
)

const (
	codeTplId = "1877556"
)

var (
	ErrCodeSendTooMany   = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = repository.ErrCodeSendTooMany
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type codeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// Send 发送验证码 biz 区分业务场景
func (c *codeService) Send(ctx context.Context, biz string, phone string) error {
	// 生成验证码
	vcode := c.generateCode()
	// 放入redis
	err := c.repo.Store(ctx, biz, phone, vcode)
	if err != nil {
		return err
	}
	// 发送验证码
	err = c.smsSvc.Send(ctx, codeTplId, []sms.NameArg{{Val: vcode, Name: "code"}}, phone)
	if err != nil {
		// 意味着redis存在该验证码，但是发送失败
		return err
	}
	return nil
}

// Verify 验证验证码 biz 区分业务场景
func (c *codeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return c.repo.Verify(ctx, biz, phone, inputCode)
}

func (c *codeService) generateCode() string {
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06d", randSeed.Int31n(1000000))
	return vcode
}
