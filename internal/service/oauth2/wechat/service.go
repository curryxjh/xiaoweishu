package wechat

import (
	"context"
	"fmt"
	"net/url"
)

var redirectURI = url.PathEscape("https://wxw.xiaoweishu.com/oauth2/wechat/callback")

type Service interface {
	AuthUrl(ctx context.Context, state string) (string, error)
	Callback(ctx context.Context, code string) (string, error)
}

func NewService(appID string) Service {
	return &service{
		appID: appID,
	}
}

type service struct {
	appID string
}

func (s *service) AuthUrl(ctx context.Context, state string) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	return fmt.Sprintf(urlPattern, s.appID, redirectURI, state), nil
}

func (s *service) Callback(ctx context.Context, code string) (string, error) {
	return "", nil
}
