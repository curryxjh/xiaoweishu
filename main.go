package main

import (
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"xiaoweishu/internal/web"
	"github.com/gin-contrib/cors"
)

func main() {
	server := gin.Default()

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

	u := web.NewUserHandler()
	u.RegisterRoutes(server)
	server.Run(":8080")
}
