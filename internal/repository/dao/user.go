package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

// User 直接对应数据库表结构
type User struct {
	Id       int64
	Email    string
	Password string

	Ctime int64 // 创新时间, 毫秒数
	Utime int64 // 更新时间, 毫秒数
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	// 存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}
