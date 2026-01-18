package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"xiaoweishu/internal/domain"
)

var redirectURI = url.PathEscape("https://wxw.xiaoweishu.com/oauth2/wechat/callback")

type Service interface {
	AuthUrl(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

func NewService(appID, appSecret string) Service {
	return &service{
		appID:     appID,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

type service struct {
	appID     string
	appSecret string
	client    *http.Client
}

func (s *service) AuthUrl(ctx context.Context, state string) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	return fmt.Sprintf(urlPattern, s.appID, redirectURI, state), nil
}

func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, s.appID, s.appSecret, code)
	//resp, err := http.Get(target)
	//if err != nil {
	//	return domain.WechatInfo{}, err
	//}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(req.Body)

	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("wechat verify code error: %s", res.ErrMsg)
	}
	return domain.WechatInfo{
		OpenID:  res.OpenID,
		UnionID: res.UnionID,
	}, nil
}

type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	OpenID  string `json:"openid"`
	Scope   string `json:"scope"`
	UnionID string `json:"unionid"`
}
