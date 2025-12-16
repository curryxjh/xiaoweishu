package web

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"xiaoweishu/internal/service"
	svcmocks "xiaoweishu/internal/service/mocks"
)

//func TestEncrypt(t *testing.T) {
//	password := "x@123456"
//	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
//	assert.NoError(t, err)
//}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				// 注册成功 return nil
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				return userSvc
			},
			reqBody: `
{
	"email":"1234@qq.com",
	"password":"hello@world123",
	"confirmPassword":"hello@world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功!",
		},
		{
			name: "参数不对, bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
{
	"email":"1234@qq.com",
	"password":"hello@world123"
	"confirmPassword":"hello@world123"
}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式有误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
{
	"email":"1234qq.com",
	"password":"hello@world123",
	"confirmPassword":"hello@world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "你的邮箱格式不正确",
		},
		{
			name: "两次输入的密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
{
	"email":"1234@qq.com",
	"password":"hello@world123",
	"confirmPassword":"hello@world12"
}
`,
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
{
	"email":"1234@qq.com",
	"password":"helloworld123",
	"confirmPassword":"helloworld123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位，且包含特殊字符",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				// 注册成功 return nil
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrUserDuplicate).Times(1)
				return userSvc
			},
			reqBody: `
{
	"email":"1234@qq.com",
	"password":"hello@world123",
	"confirmPassword":"hello@world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				// 注册成功 return nil
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("系统异常")).Times(1)
				return userSvc
			},
			reqBody: `
{
	"email":"1234@qq.com",
	"password":"hello@world123",
	"confirmPassword":"hello@world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			// SignUp 没有使用CodeService
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			// 数据是 JSON 格式
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()

			// HTTP 请求进入 GIN 框架的入口
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())

		})
	}
}
