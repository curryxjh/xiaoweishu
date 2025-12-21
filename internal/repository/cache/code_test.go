package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"xiaoweishu/internal/repository/cache/redismocks"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) redis.Cmdable
		// 输入
		ctx   context.Context
		biz   string
		phone string
		code  string
		// 输出
		wantErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(nil)
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:123456789"}, "123456").Return(res)
				return cmd
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "123456789",
			code:    "123456",
			wantErr: nil,
		},
		{
			name: "redis 错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errors.New("mock redis error"))
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:123456789"}, "123456").Return(res)
				return cmd
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "123456789",
			code:    "123456",
			wantErr: errors.New("mock redis error"),
		},
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:123456789"}, "123456").Return(res)
				return cmd
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "123456789",
			code:    "123456",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(nil)
				res.SetVal(int64(-2))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:123456789"}, "123456").Return(res)
				return cmd
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "123456789",
			code:    "123456",
			wantErr: errors.New("系统错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := tc.mock(ctrl)
			cache := NewCodeCache(client)
			err := cache.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
