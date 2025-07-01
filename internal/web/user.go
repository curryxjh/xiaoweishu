package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
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
	//ug.POST("/login", u.LogIn)
	ug.POST("/login", u.LogInJWT)
	ug.POST("/edit", u.Edit)
	//ug.GET("/profile", u.Profile)
	ug.GET("/profile", u.ProfileJWT)
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
	if err == service.ErrUserDuplicateEmail {
		c.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统异常")
		return
	}

	c.String(http.StatusOK, "注册成功!")
	return
}

func (u *UserHandler) LogIn(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(c, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		c.String(http.StatusOK, "用户名/密码错误")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}

	// 登录成功后，把 session 拿出来
	// 设置session
	sess := sessions.Default(c)
	// 我需要放在 session 中的值
	sess.Set("userId", user.Id)

	sess.Options(sessions.Options{
		//Secure: true,
		//HttpOnly: true,
		MaxAge: 60,
	})
	sess.Save()

	c.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) LogInJWT(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(c, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		c.String(http.StatusOK, "用户名/密码错误")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       user.Id,
		UserAgent: c.Request.UserAgent(),
	}

	// 使用 JWT 设置登录状态
	// 生成一个 JWT
	//token := jwt.New(jwt.SigningMethodHS512)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("KntbYH88cXPKDRdFrXrQjh5yZpA7c5QQXKh3MHwYFnt2v43wGCy2d8XCSpmwPjFy"))

	c.Header("x-jwt-token", tokenStr)

	c.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Logout(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()

	c.String(http.StatusOK, "退出登录成功")
}

func (u *UserHandler) Edit(c *gin.Context) {
	c.String(http.StatusOK, "Edit")
}

func (u *UserHandler) Profile(c *gin.Context) {
	c.String(http.StatusOK, "Profile")
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, ok := ctx.Get("claims")
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	fmt.Println(claims.Uid)
	ctx.String(http.StatusOK, "Profile")
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明自己要放进 token 的数据
	Uid int64

	UserAgent string
}
