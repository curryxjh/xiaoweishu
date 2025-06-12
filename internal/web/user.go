package web

import "github.com/gin-gonic/gin"

// UserHandler 定义所有和user相关的路由
type UserHandler struct{}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/users/signup", u.SignUp)

	server.POST("/users/signin", u.LogIn)

	server.POST("/users/edit", u.Edit)

	server.GET("/users/profile", u.Profile)
}

func (u *UserHandler) SignUp(c *gin.Context) {}

func (u *UserHandler) LogIn(c *gin.Context) {}

func (u *UserHandler) Edit(c *gin.Context) {}

func (u *UserHandler) Profile(c *gin.Context) {}
