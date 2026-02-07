package ioc

import (
	"os"
	"xiaoweishu/internal/pkg/logger"
	"xiaoweishu/internal/service/oauth2/wechat"
	"xiaoweishu/internal/web"
)

func InitOauth2WechatService(l logger.LoggerV1) wechat.Service {
	appID, ok := os.LookupEnv("WECHAT_APPID")
	if !ok {
		panic("WECHAT_APPID not found")
	}
	appSecret, ok := os.LookupEnv("WECHAT_APPSECRET")
	if !ok {
		panic("WECHAT_APPSECRET not found")
	}
	return wechat.NewService(appID, appSecret, l)
}

func NewWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
