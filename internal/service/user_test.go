package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
	"xiaoweishu/internal/domain"
	"xiaoweishu/internal/repository"
	repomocks "xiaoweishu/internal/repository/mocks"
)

func Test_userService_Login(t *testing.T) {
	// 做成一个测试用例使用的公共时间
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository
		// 输入
		ctx      context.Context
		email    string
		password string

		// 输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登陆成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRepo := repomocks.NewMockUserRepository(ctrl)
				userRepo.EXPECT().FindByEmail(gomock.Any(), "1234@qq.com").
					Return(domain.User{
						Email:    "1234@qq.com",
						Password: "$2a$10$PMr9A/ENZZNqJm8oRUZVl.SWd.XZ3pquYS0GjvcW7Dv98Lt02R0LO",
						Phone:    "19852897878",
						Ctime:    now,
					}, nil).Times(1)
				return userRepo
			},
			email:    "1234@qq.com",
			password: "hello@world123",
			wantUser: domain.User{
				Email:    "1234@qq.com",
				Password: "$2a$10$PMr9A/ENZZNqJm8oRUZVl.SWd.XZ3pquYS0GjvcW7Dv98Lt02R0LO",
				Phone:    "19852897878",
				Ctime:    now,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

func Test_EncryptPassword(t *testing.T) {
	password := "hello@world123"
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
