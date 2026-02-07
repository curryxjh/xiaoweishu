package logger

import (
	"bytes"
	"context"
	"io"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessLog struct {
	// HTTP 请求的方法
	Method string
	// URL 整个请求
	Url        string
	Duration   string
	ReqBody    string
	RespBody   string
	StatusCode int
}

type LoggerMiddlewareBuilder struct {
	allowReqBody  atomic.Bool
	allowRespBody atomic.Bool
	loggerFunc    func(c context.Context, al *AccessLog)
}

func NewBuilder(fn func(c context.Context, al *AccessLog)) *LoggerMiddlewareBuilder {
	return &LoggerMiddlewareBuilder{
		loggerFunc:    fn,
		allowReqBody:  atomic.Bool{},
		allowRespBody: atomic.Bool{},
	}
}

func (l *LoggerMiddlewareBuilder) AllowReqBody(ok bool) *LoggerMiddlewareBuilder {
	l.allowReqBody.Store(ok)
	return l
}

func (l *LoggerMiddlewareBuilder) AllowRespBody() *LoggerMiddlewareBuilder {
	l.allowRespBody.Store(true)
	return l
}

func (l *LoggerMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		url := c.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: c.Request.Method,
			// 也可能很长
			Url: url,
		}
		if l.allowReqBody.Load() && c.Request.Body != nil {
			// 此处 Body 读完就没了, 它是一个流
			body, _ := c.GetRawData()
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			// 比较消耗 CPU 和内存的操作
			// 会引起复制
			al.ReqBody = string(body)

		}
		if l.allowRespBody.Load() {
			c.Writer = responseWriter{al: al, ResponseWriter: c.Writer}
		}
		defer func() {
			al.Duration = time.Since(start).String()
			//al.Duration = time.Now().Sub(start)
			l.loggerFunc(c, al)
		}()
		// 执行业务逻辑
		c.Next()
	}
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.al.RespBody = string(b)
	return w.ResponseWriter.Write(b)
}

func (w responseWriter) WriteString(s string) (int, error) {
	w.al.RespBody = s
	return w.ResponseWriter.WriteString(s)
}

func (w responseWriter) WriteHeader(statusCode int) {
	w.al.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
