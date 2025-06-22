package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"xiaoweishu/internal/domain"
	"xiaoweishu/internal/service"
)

// UserHandler 定义所有和user相关的路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)

	return &UserHandler{
		svc:         svc,
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/signin", u.LogIn)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

func (u *UserHandler) SignUp(c *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析你的数据到 req 里面
	// 解析错了，就会返回一个 400 的错误
	if err := c.Bind(&req); err != nil {
		log.Println(err)
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		log.Println(err)
		return
	}

	if !ok {
		c.String(http.StatusOK, "你的邮箱格式不正确")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		log.Println(err)
		return
	}

	if req.ConfirmPassword != req.ConfirmPassword {
		c.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	if !ok {
		c.String(http.StatusOK, "密码必须大于8位，且包含特殊字符")
		return
	}

	err = u.svc.SignUp(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.String(http.StatusOK, "系统异常")
		return
	}

	c.String(http.StatusOK, "success")
}

func (u *UserHandler) LogIn(c *gin.Context) {}

func (u *UserHandler) Edit(c *gin.Context) {}

func (u *UserHandler) Profile(c *gin.Context) {}
