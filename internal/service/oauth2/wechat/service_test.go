//go:build manual

package wechat

import (
	"context"
	"os"
	"testing"
)

func Test_service_manual_VerifyCode(t *testing.T) {
	appID, ok := os.LookupEnv("WECHAT_APPID")
	if !ok {
		panic("WECHAT_APPID not found")
	}
	appSecret, ok := os.LookupEnv("WECHAT_APPSECRET")
	if !ok {
		panic("WECHAT_APPSECRET not found")
	}
	svc := NewService(appID, appSecret)
	res, err := svc.VerifyCode(context.Background(), "")
	if err != nil {
		return
	}
	t.Log(res)
}
