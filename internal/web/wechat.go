package web

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"xiaoweishu/internal/pkg/ginx"
	"xiaoweishu/internal/service"
	"xiaoweishu/internal/service/oauth2/wechat"
	ijwt "xiaoweishu/internal/web/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
)

type WechatHandlerConfig struct {
	Secure bool
}

type Oauth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	ijwt.Handler
	stateKey []byte
	cfg      WechatHandlerConfig
}

func NewOauth2WechatHandler(svc wechat.Service, userSvc service.UserService, cfg WechatHandlerConfig, jwtHdl ijwt.Handler) *Oauth2WechatHandler {
	return &Oauth2WechatHandler{
		svc:      svc,
		userSvc:  userSvc,
		stateKey: []byte("KntbYH88cXJHKDRdFrXrQjh5yZp7c5QQXKh3MXJHwYFnt2v43wGCy2d8XCSpmwPjFy"),
		cfg:      cfg,
		Handler:  jwtHdl,
	}
}

func (h *Oauth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/oauth2/wechat")
	ug.GET("/authurl", h.AuthURL)
	ug.Any("/callback", h.Callback)
}

func (h *Oauth2WechatHandler) AuthURL(c *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthUrl(c, state)
	if err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "获取授权URL失败",
		})
		return
	}
	if err := h.SetStateCookie(c, state); err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	c.JSON(http.StatusOK, ginx.Result{Data: gin.H{"url": url}})
}

func (h *Oauth2WechatHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	err := h.verifyState(c)
	if err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: http.StatusInternalServerError,
			Msg:  "登陆失败",
		})
	}
	info, err := h.svc.VerifyCode(c, code)
	if err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: http.StatusInternalServerError,
			Msg:  "系统错误",
		})
	}
	// 这里需要设置为登录态, 需要 set token, 还需要拿到 uid
	u, err := h.userSvc.FindOrCreateByWechat(c, info)
	if err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: http.StatusInternalServerError,
			Msg:  "系统错误",
		})
	}

	if err := h.SetLoginToken(c, u.Id); err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: http.StatusInternalServerError,
			Msg:  "系统错误",
		})
		return
	}

	c.JSON(http.StatusOK, ginx.Result{
		Code: http.StatusOK,
		Data: info,
	})
}

type StateClaims struct {
	jwt.RegisteredClaims
	state string
}

func (h *Oauth2WechatHandler) SetStateCookie(c *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		state: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return err
	}
	c.SetCookie("jwt-state", tokenStr, 6000, "/oauth2/wechat/callback", "", h.cfg.Secure, true)
	return nil
}

func (h *Oauth2WechatHandler) verifyState(c *gin.Context) error {
	state := c.Query("state")
	ck, err := c.Cookie("jwt-state")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginx.Result{
			Code: http.StatusInternalServerError,
			Msg:  "系统错误",
		})
		return fmt.Errorf("拿不到 state 的 cookie, %w", err)
	}
	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusInternalServerError, ginx.Result{
			Code: http.StatusInternalServerError,
			Msg:  "系统错误",
		})
		return fmt.Errorf("token 已过期, %w", err)
	}
	if sc.state != state {
		return errors.New("state 不匹配")
	}
	return nil
}
