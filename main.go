package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"strings"
	"time"
	"xiaoweishu/internal/repository"
	"xiaoweishu/internal/repository/dao"
	"xiaoweishu/internal/service"
	"xiaoweishu/internal/web"
	"github.com/gin-contrib/cors"
)

func main() {
	server := gin.Default()

	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)

	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods: []string{},
		AllowHeaders: []string{"authorization", "Content-Type"},
		//ExposeHeaders: []string{"x-jwt-token"},
		// 是否允许携带用户认证信息，如cookie
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "https://") {
				return true
			}
			return strings.Contains(origin, "http://youcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	u.RegisterRoutes(server)
	server.Run(":8080")
}
