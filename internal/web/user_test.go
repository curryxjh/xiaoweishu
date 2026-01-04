package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"xiaoweishu/internal/domain"
	"xiaoweishu/internal/pkg/ginx"
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

// todo 未完成，需要修复
func TestUserHandler_LoginSMS(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		// 输入
		reqBody string
		// 输出
		wantCode int
		wantBody ginx.Result
	}{
		{
			name: "登陆成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				svc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), biz, gomock.Any(), "123456").Return(true, nil).Times(1)
				svc.EXPECT().FindOrCreate(gomock.Any(), "12345678901").Return(domain.User{
					Phone: "12345678901",
				}, nil).Times(1)
				return svc, codeSvc
			},
			reqBody: `
{
	"phone":"12345678901",
	"code":"123456"
}
`,
			wantCode: http.StatusOK,
			wantBody: ginx.Result{
				Code: http.StatusOK,
				Msg:  "短信验证成功",
			},
		},
		{
			name: "参数错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				svc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return svc, codeSvc
			},
			reqBody: `
{
	"phone":"101,
	"code":"123456"
}
`,
			wantCode: http.StatusBadRequest,
			wantBody: ginx.Result{
				Code: http.StatusBadRequest,
				Msg:  "参数错误",
			},
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				svc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return svc, codeSvc
			},
			reqBody: `
{
	"phone":"101,
	"code":"123456"
}
`,
			wantCode: http.StatusBadRequest,
			wantBody: ginx.Result{
				Code: http.StatusBadRequest,
				Msg:  "参数错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			h := NewUserHandler(tc.mock(ctrl))
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/login_sms", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			var res ginx.Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)
		})
	}
}
