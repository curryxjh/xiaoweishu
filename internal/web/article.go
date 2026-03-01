package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaoweishu/internal/domain"
	"xiaoweishu/internal/pkg/ginx"
	"xiaoweishu/internal/pkg/logger"
	"xiaoweishu/internal/service"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.LoggerV1
}

func NewArticleHandler(svc service.ArticleService, l logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (a *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/articles")
	ug.POST("/edit", a.Edit)
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 检测输入

	// 调用 svc
	id, err := a.svc.Save(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg:  "OK",
		Data: id,
	})
}
