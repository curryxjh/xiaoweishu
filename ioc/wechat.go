package ioc

import (
	"os"
	"xiaoweishu/internal/service/oauth2/wechat"
)

func InitOauth2WechatService() wechat.Service {
	appID, ok := os.LookupEnv("WECHAT_APPID")
	if !ok {
		panic("WECHAT_APPID not found")
	}
	return wechat.NewService(appID)

}
