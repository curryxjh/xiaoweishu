package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"xiaoweishu/internal/pkg/ginx/middlewares/ratelimit"
	"xiaoweishu/internal/repository"
	"xiaoweishu/internal/repository/dao"
	"xiaoweishu/internal/service"
	"xiaoweishu/internal/web"
	"xiaoweishu/internal/web/middleware"
)

func main() {

	db := initDB()

	server := initServer()

	u := initUser(db)
	u.RegisterRoutes(server)

	//server := gin.Default()
	//server.GET("/hello", func(c *gin.Context) {
	//	c.String(http.StatusOK, "hello world")
	//})
	server.Run(":8080")
}

func initServer() *gin.Engine {
	server := gin.Default()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods: []string{},
		AllowHeaders: []string{"authorization", "Content-Type"},
		//
		ExposeHeaders: []string{"x-jwt-token"},
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

	//store := cookie.NewStore([]byte("secret"))

	store := memstore.NewStore([]byte("KntbYH88cXPKDRdFrXrQjh5yZpA7c5QQXKh3MHwYFnt2v43wGCy2d8XCSpmwPjFy"),
		[]byte("SAmc4oHXzZXPd2Q5tr7A2COHHB0rEk3wrLqfPiwxCDZw5jnNzCahyXxiCafRqkYN"))

	//store, err := redis.NewStore(16, "tcp", "localhost:16379", "", "",
	//	[]byte("KntbYH88cXPKDRdFrXrQjh5yZpA7c5QQXKh3MHwYFnt2v43wGCy2d8XCSpmwPjFy"),
	//	[]byte("SAmc4oHXzZXPd2Q5tr7A2COHHB0rEk3wrLqfPiwxCDZw5jnNzCahyXxiCafRqkYN"))
	//if err != nil {
	//	panic(err)
	//}

	server.Use(sessions.Sessions("mysession", store))

	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/login").
	//	IgnorePaths("/users/signup").Build())

	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/login").
		IgnorePaths("/users/signup").Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13306)/webook?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}

	if err := dao.InitTable(db); err != nil {
		panic(err)
	}
	return db
}
