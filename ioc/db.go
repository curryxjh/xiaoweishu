package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"xiaoweishu/internal/repository/dao"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13306)/webook?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}

	if err := dao.InitTable(db); err != nil {
		panic(err)
	}
	return db
}
