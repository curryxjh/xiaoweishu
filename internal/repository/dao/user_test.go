package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDao_Insert(t *testing.T) {
	testCases := []struct {
		name string
		mock func(t *testing.T) *sql.DB
		// 输入
		ctx  context.Context
		user User
		// 输出
		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				res := sqlmock.NewResult(3, 1)
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnResult(res)
				require.NoError(t, err)
				return mockDB
			},
			ctx: context.Background(),
			user: User{
				Email: sql.NullString{
					String: "test@example.com",
					Valid:  true,
				},
			},
			wantErr: nil,
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnError(&mysql.MySQLError{
					Number: 1062,
				})
				require.NoError(t, err)
				return mockDB
			},
			ctx: context.Background(),
			user: User{
				Email: sql.NullString{
					String: "test@example.com",
					Valid:  true,
				},
			},
			wantErr: ErrUserDuplicate,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnError(errors.New("mysql error"))
				require.NoError(t, err)
				return mockDB
			},
			ctx: context.Background(),
			user: User{
				Email: sql.NullString{
					String: "test@example.com",
					Valid:  true,
				},
			},
			wantErr: errors.New("mysql error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db, err := gorm.Open(gormMysql.New(gormMysql.Config{
				Conn: tc.mock(t),
				// select version 跳过
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				// 不需要默认开启事物
				SkipDefaultTransaction: true,
				// 跳过自动ping数据库
				DisableAutomaticPing: true,
			})
			d := NewUserDao(db)
			err = d.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
