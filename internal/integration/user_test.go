package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"xiaoweishu/internal/pkg/ginx"
	"xiaoweishu/ioc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name string
		// 需要考虑准备数据，验证数据。
		before func(t *testing.T)
		// 需要考虑数据库的数据对不对，redis的数据对不对
		after func(t *testing.T)
		// 输入
		reqBody string
		// 输出
		wantCode int
		wantBody ginx.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				// 不需要处理，redis 中无数据
			},
			after: func(t *testing.T) {
				// 需要清理时间
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
				val, err := rdb.GetDel(ctx, "phone_code:login:12345678901").Result()
				cancel()
				require.NoError(t, err)
				assert.True(t, 6 == len(val))
			},
			reqBody: `
{
	"phone":"12345678901"
}
`,
			wantCode: http.StatusOK,
			wantBody: ginx.Result{
				Code: 0,
				Msg:  "短信发送成功",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				// 此手机号已经有验证码了
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
				defer cancel()
				_, err := rdb.Set(ctx, "phone_code:login:12345678901", "123456", time.Minute*9+time.Second*30).Result()
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 需要清理时间
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
				val, err := rdb.GetDel(ctx, "phone_code:login:12345678901").Result()
				cancel()
				require.NoError(t, err)
				assert.True(t, "123456" == val)
			},
			reqBody: `
{
	"phone":"12345678901"
}
`,
			wantCode: http.StatusOK,
			wantBody: ginx.Result{
				Code: 4,
				Msg:  "短信发送太频繁，请稍后再试",
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
				defer cancel()
				// 有验证码，但是没有过期时间
				_, err := rdb.Set(ctx, "phone_code:login:12345678901", "123456", 0).Result()
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 需要清理时间
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
				val, err := rdb.GetDel(ctx, "phone_code:login:12345678901").Result()
				cancel()
				require.NoError(t, err)
				assert.True(t, 6 == len(val))
			},
			reqBody: `
{
	"phone":"12345678901"
}
`,
			wantCode: http.StatusOK,
			wantBody: ginx.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "手机号为空",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `
{
	"phone":""
}
`,
			wantCode: http.StatusOK,
			wantBody: ginx.Result{
				Code: 5,
				Msg:  "手机号不能为空",
			},
		},
		{
			name: "数据格式有误",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `
{
	"phone":,
}
`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/sms/login/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != http.StatusOK {
				return
			}
			var res ginx.Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)
			tc.after(t)
		})
	}
}
