package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
	"xiaoweishu/internal/repository"
	"xiaoweishu/internal/repository/cache"
	"xiaoweishu/internal/repository/dao"
	"xiaoweishu/internal/service"
	"xiaoweishu/internal/service/sms/memory"
	"xiaoweishu/internal/web"
	"xiaoweishu/internal/web/middleware"
)

func main() {

	db := initDB()

	server := initServer()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
	})
	u := initUser(db, rdb)
	u.RegisterRoutes(server)

	//server := gin.Default()
	server.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})
	server.Run(":8080")
}

func initServer() *gin.Engine {
	server := gin.Default()

	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: "localhost:16379",
	//})
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

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

	//store := memstore.NewStore([]byte("KntbYH88cXPKDRdFrXrQjh5yZpA7c5QQXKh3MHwYFnt2v43wGCy2d8XCSpmwPjFy"),
	//	[]byte("SAmc4oHXzZXPd2Q5tr7A2COHHB0rEk3wrLqfPiwxCDZw5jnNzCahyXxiCafRqkYN"))

	//store, err := redis.NewStore(16, "tcp", "localhost:16379", "", "",
	//	[]byte("KntbYH88cXPKDRdFrXrQjh5yZpA7c5QQXKh3MHwYFnt2v43wGCy2d8XCSpmwPjFy"),
	//	[]byte("SAmc4oHXzZXPd2Q5tr7A2COHHB0rEk3wrLqfPiwxCDZw5jnNzCahyXxiCafRqkYN"))
	//if err != nil {
	//	panic(err)
	//}

	//server.Use(sessions.Sessions("mysession", store))

	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/login").
	//	IgnorePaths("/users/signup").Build())

	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/login").
		IgnorePaths("/users/signup").
		IgnorePaths("/users/sms/login/send").
		IgnorePaths("/users/sms/login/verify").Build())
	return server
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	ud := dao.NewUserDao(db)
	uc := cache.NewUserCache(redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
	}), time.Minute*15)
	repo := repository.NewUserRepository(ud, uc)
	svc := service.NewUserService(repo)
	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	//smsSvc := tencent.NewService(
	//	sms.NewClient(),
	//	"14005555555",
	//	"测试短信",
	//)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)
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
