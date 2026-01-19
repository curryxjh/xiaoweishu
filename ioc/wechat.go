package ioc

import (
	"os"
	"xiaoweishu/internal/service/oauth2/wechat"
	"xiaoweishu/internal/web"
)

func InitOauth2WechatService() wechat.Service {
	appID, ok := os.LookupEnv("WECHAT_APPID")
	if !ok {
		panic("WECHAT_APPID not found")
	}
	appSecret, ok := os.LookupEnv("WECHAT_APPSECRET")
	if !ok {
		panic("WECHAT_APPSECRET not found")
	}
	return wechat.NewService(appID, appSecret)
}

func NewWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
