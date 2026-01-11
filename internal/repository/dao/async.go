package dao

import (
	"context"
	"gorm.io/gorm"
)

var ErrWaitingSMSNotFound = gorm.ErrRecordNotFound

type AsyncSmsDao interface {
	Insert(ctx context.Context, s AsyncSms) error
}

type AsyncSms struct {
	Id int64
	// 重试次数
	RetryCnt int
	// 最大重试次数
	RetryMax int
	Status   uint8
	Ctime    int64
	Utime    int64 `gorm:"index"`
}

type SmsConfig struct {
	TplId   string
	Args    []string
	Numbers []string
}
