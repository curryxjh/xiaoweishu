package web

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"xiaoweishu/internal/pkg/ginx"
	"xiaoweishu/internal/service/oauth2/wechat"
)

type Oauth2WechatHandler struct {
	svc wechat.Service
}

func NewOauth2WechatHandler(svc wechat.Service) *Oauth2WechatHandler {
	return &Oauth2WechatHandler{
		svc: svc,
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

	c.JSON(http.StatusOK, ginx.Result{
		Code: http.StatusOK,
		Data: info,
	})
}

func (h *Oauth2WechatHandler) verifyState(c *gin.Context) error {
	return nil
}
