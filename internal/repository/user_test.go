package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"xiaoweishu/internal/domain"
	"xiaoweishu/internal/repository/cache"
	cachemocks "xiaoweishu/internal/repository/cache/mocks"
	"xiaoweishu/internal/repository/dao"
	daomocks "xiaoweishu/internal/repository/dao/mocks"
)

// 非异步的测试
//func TestCachedUserRepository_FindById(t *testing.T) {
//	now := time.Now()
//	// 去除毫秒以外的时间精度
//	now = time.UnixMilli(now.UnixMilli())
//	testCases := []struct {
//		name     string
//		ctx      context.Context
//		id       int64
//		mock     func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
//		wantUser domain.User
//		wantErr  error
//	}{
//		{
//			name: "缓存未命中, 查询成功",
//			ctx:  context.Background(),
//			id:   1,
//			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
//				uc := cachemocks.NewMockUserCache(ctrl)
//				uc.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, cache.ErrKeyNotFound).Times(1)
//				ud := daomocks.NewMockUserDao(ctrl)
//				ud.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.User{
//					Id: 1,
//					Email: sql.NullString{
//						String: "1234@qq.com",
//						Valid:  true,
//					},
//					Phone: sql.NullString{
//						String: "15972887248",
//						Valid:  true,
//					},
//					Password: "123456qwe",
//					NicName:  "xiaoweishu",
//					BirthDay: now,
//					AboutMe:  "hello world",
//					Ctime:    now.UnixMilli(),
//					Utime:    now.UnixMilli(),
//				}, nil).Times(1)
//
//				uc.EXPECT().Set(gomock.Any(), domain.User{
//					Id:       1,
//					Email:    "1234@qq.com",
//					Phone:    "15972887248",
//					Password: "123456qwe",
//					NickName: "xiaoweishu",
//					Birthday: now.String(),
//					AboutMe:  "hello world",
//					Ctime:    now,
//				}).Return(nil).Times(1)
//				return ud, uc
//			},
//			wantUser: domain.User{
//				Id:       1,
//				Email:    "1234@qq.com",
//				Phone:    "15972887248",
//				Password: "123456qwe",
//				NickName: "xiaoweishu",
//				Birthday: now.String(),
//				AboutMe:  "hello world",
//				Ctime:    now,
//			},
//			wantErr: nil,
//		},
//		{
//			name: "缓存命中",
//			ctx:  context.Background(),
//			id:   1,
//			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
//				uc := cachemocks.NewMockUserCache(ctrl)
//				uc.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{
//					Id:       1,
//					Email:    "1234@qq.com",
//					Phone:    "15972887248",
//					Password: "123456qwe",
//					NickName: "xiaoweishu",
//					Birthday: now.String(),
//					AboutMe:  "hello world",
//					Ctime:    now,
//				}, nil).Times(1)
//				ud := daomocks.NewMockUserDao(ctrl)
//				return ud, uc
//			},
//			wantUser: domain.User{
//				Id:       1,
//				Email:    "1234@qq.com",
//				Phone:    "15972887248",
//				Password: "123456qwe",
//				NickName: "xiaoweishu",
//				Birthday: now.String(),
//				AboutMe:  "hello world",
//				Ctime:    now,
//			},
//			wantErr: nil,
//		},
//		{
//			name: "缓存未命中, 数据库查询失败",
//			ctx:  context.Background(),
//			id:   1,
//			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
//				uc := cachemocks.NewMockUserCache(ctrl)
//				uc.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, cache.ErrKeyNotFound).Times(1)
//				ud := daomocks.NewMockUserDao(ctrl)
//				ud.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.User{}, errors.New("mock error")).Times(1)
//				return ud, uc
//			},
//			wantUser: domain.User{},
//			wantErr:  errors.New("mock error"),
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//			ud, uc := tc.mock(ctrl)
//			repo := NewUserRepository(ud, uc)
//			u, err := repo.FindById(tc.ctx, tc.id)
//			assert.Equal(t, tc.wantErr, err)
//			assert.Equal(t, tc.wantUser, u)
//		})
//	}
//}

// 异步测试
func TestCachedUserRepository_FindById(t *testing.T) {
	now := time.Now()
	// 去除毫秒以外的时间精度
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name     string
		ctx      context.Context
		id       int64
		mock     func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中, 查询成功",
			ctx:  context.Background(),
			id:   1,
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, cache.ErrKeyNotFound).Times(1)
				ud := daomocks.NewMockUserDao(ctrl)
				ud.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.User{
					Id: 1,
					Email: sql.NullString{
						String: "1234@qq.com",
						Valid:  true,
					},
					Phone: sql.NullString{
						String: "15972887248",
						Valid:  true,
					},
					Password: "123456qwe",
					NicName:  "xiaoweishu",
					BirthDay: now,
					AboutMe:  "hello world",
					Ctime:    now.UnixMilli(),
					Utime:    now.UnixMilli(),
				}, nil).Times(1)

				uc.EXPECT().Set(gomock.Any(), domain.User{
					Id:       1,
					Email:    "1234@qq.com",
					Phone:    "15972887248",
					Password: "123456qwe",
					NickName: "xiaoweishu",
					Birthday: now.String(),
					AboutMe:  "hello world",
					Ctime:    now,
				}).Return(nil).Times(1)
				return ud, uc
			},
			wantUser: domain.User{
				Id:       1,
				Email:    "1234@qq.com",
				Phone:    "15972887248",
				Password: "123456qwe",
				NickName: "xiaoweishu",
				Birthday: now.String(),
				AboutMe:  "hello world",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			ctx:  context.Background(),
			id:   1,
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{
					Id:       1,
					Email:    "1234@qq.com",
					Phone:    "15972887248",
					Password: "123456qwe",
					NickName: "xiaoweishu",
					Birthday: now.String(),
					AboutMe:  "hello world",
					Ctime:    now,
				}, nil).Times(1)
				ud := daomocks.NewMockUserDao(ctrl)
				return ud, uc
			},
			wantUser: domain.User{
				Id:       1,
				Email:    "1234@qq.com",
				Phone:    "15972887248",
				Password: "123456qwe",
				NickName: "xiaoweishu",
				Birthday: now.String(),
				AboutMe:  "hello world",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "缓存未命中, 数据库查询失败",
			ctx:  context.Background(),
			id:   1,
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, cache.ErrKeyNotFound).Times(1)
				ud := daomocks.NewMockUserDao(ctrl)
				ud.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.User{}, errors.New("mock error")).Times(1)
				return ud, uc
			},
			wantUser: domain.User{},
			wantErr:  errors.New("mock error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			u, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
			time.Sleep(time.Second)
		})
	}
}
