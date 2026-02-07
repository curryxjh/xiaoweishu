package web

import (
	"errors"
	"log"
	"net/http"
	"xiaoweishu/internal/domain"
	"xiaoweishu/internal/pkg/ginx"
	"xiaoweishu/internal/service"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	ijwt "xiaoweishu/internal/web/jwt"
)

const biz = "login"

var _ handler = (*UserHandler)(nil)

// UserHandler 定义所有和user相关的路由
type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	ijwt.Handler
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, cmd redis.Cmdable) *UserHandler {
	const (
		emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)

	return &UserHandler{
		svc:         svc,
		codeSvc:     codeSvc,
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		Handler:     ijwt.NewRedisJwtHandler(cmd),
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.LogIn)
	ug.POST("/login", u.LogInJWT)
	ug.POST("login_sms", u.LoginSMS)
	ug.POST("/edit", u.Edit)
	//ug.GET("/profile", u.Profile)
	//ug.POST("/logout", u.Logout)
	ug.POST("/logout", u.LogoutJWT)
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("/sms/login/send", u.SendLoginSMSCode)
	ug.POST("/sms/login/verify", u.VerifyLoginSMSCode)
	ug.POST("/refresh_token", u.RefreshToekn)
}

func (u *UserHandler) SendLoginSMSCode(c *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}

	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}

	// 校验是否是合法手机号
	// 可以使用正则表达式
	if req.Phone == "" {
		c.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "手机号不能为空",
			Data: nil,
		})
		return
	}
	err := u.codeSvc.Send(c, biz, req.Phone)
	switch err {
	case nil:
		c.JSON(http.StatusOK, ginx.Result{
			Msg:  "短信发送成功",
			Data: nil,
		})
		return
	case service.ErrCodeSendTooMany:
		c.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
			Data: nil,
		})
	default:
		c.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}

}

func (u *UserHandler) VerifyLoginSMSCode(c *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "参数错误",
			Data: nil,
		})
		return
	}
	ok, err := u.codeSvc.Verify(c, biz, req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "短信验证失败",
			Data: nil,
		})
		return
	}
	user, err := u.svc.FindOrCreate(c, req.Phone)
	if err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}
	if err := u.SetLoginToken(c, user.Id); err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK, ginx.Result{
		Code: 4,
		Msg:  "短信验证成功",
		Data: nil,
	})
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

	if req.Password != req.ConfirmPassword {
		c.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		//log.Println(err)
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
	if err == service.ErrUserDuplicate {
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
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		c.String(http.StatusOK, "用户名/密码错误")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}

	if err := u.SetLoginToken(c, user.Id); err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}
	c.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) LoginSMS(c *gin.Context) {
	type SmsLoginReq struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req SmsLoginReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: http.StatusBadRequest,
			Msg:  "参数错误",
			Data: nil,
		})
		return
	}
	ok, err := u.codeSvc.Verify(c, biz, req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: http.StatusInternalServerError,
			Msg:  "系统错误",
			Data: nil,
		})
		zap.L().Error("短信验证失败", zap.Error(err))
		return
	}
	if !ok {
		c.JSON(http.StatusOK, ginx.Result{
			Code: http.StatusUnauthorized,
			Msg:  "手机号/验证码错误",
			Data: nil,
		})
		return
	}
	user, err := u.svc.FindOrCreate(c, req.Phone)
	if err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: http.StatusInternalServerError,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}

	if err := u.SetLoginToken(c, user.Id); err != nil {
		c.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, ginx.Result{
		Code: http.StatusOK,
		Msg:  "短信验证成功",
		Data: nil,
	})
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

func (u *UserHandler) LogoutJWT(c *gin.Context) {
	err := u.ClearToken(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	c.JSON(http.StatusOK, ginx.Result{
		Code: 0,
		Msg:  "退出登录成功",
	})
}

func (u *UserHandler) Edit(c *gin.Context) {
	type EditReq struct {
		NickName string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"about_me"`
	}

	var req EditReq
	if err := c.Bind(&req); err != nil {
		return
	}
	if req.NickName == "" {
		return
	}

}

func (u *UserHandler) Profile(c *gin.Context) {
	type ProfileReq struct {
		Email string `json:"email"`
	}
	sess := sessions.Default(c)
	id := sess.Get("userId").(int64)
	res, err := u.svc.Profile(c, id)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	c.JSON(http.StatusOK, ginx.Result{
		Code: 0,
		Msg:  "成功",
		Data: res.Email,
	})
}

func (u *UserHandler) ProfileJWT(c *gin.Context) {
	type ProfileResp struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		NicName  string `json:"nicName"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"about_me"`
	}
	uc := c.MustGet("claims").(*ijwt.UserClaims)
	res, err := u.svc.Profile(c, uc.Uid)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	c.JSON(http.StatusOK, ProfileResp{
		Email:    res.Email,
		Phone:    res.Phone,
		NicName:  res.NickName,
		Birthday: res.Birthday,
		AboutMe:  res.AboutMe,
	})
}

func (u *UserHandler) RefreshToekn(c *gin.Context) {
	// 只有这个接口, 拿出来的才是 refresh-token, 其余的都是 access-token
	refresh_token := u.ExtractToken(c)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(refresh_token, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RefreshTokenKey, nil
	})
	if err != nil {
	}
	if err != nil || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = u.CheckSession(c, rc.Ssid)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := u.SetJwtToken(c, rc.Uid, rc.Ssid); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		zap.L().Error("设置 JWT token 出现错误", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, ginx.Result{
		Msg: "success",
	})
}
