package ioc

import (
	"time"
	"xiaoweishu/internal/pkg/logger"
	"xiaoweishu/internal/repository/dao"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitDB(l logger.LoggerV1) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config = Config{
		DSN: "root:root@tcp(127.0.0.1:13306)/webook_default",
	}
	if err := viper.UnmarshalKey("db", &cfg); err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询阈值, 只有执行时间超过这个阈值, 才会使 用
			// 50ms, 100ms 是一般的
			// SQL 查询必然要求命中索引, 最后就是直走一次磁盘 IO
			// 一次磁盘 IO 一般是不到 10ms
			SlowThreshold:             50 * time.Millisecond,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  glogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}

	if err := dao.InitTable(db); err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}
