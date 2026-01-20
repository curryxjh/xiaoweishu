package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"xiaoweishu/internal/pkg/ginx"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	AccessTokenKey  = []byte("KntbYH88cXPKDRdFrXrQjh5yZpA7c5QQXKh3MHwYFnt2v43wGCy2d8XCSpmwPjFy")
	RefreshTokenKey = []byte("KntcTH88cXPKDRdFrXrQjh5yZpA7c5QQXKh3MHwYFnt2v43wGCy2d8JCSpmwPjFy")
)

type RedisJwtHandler struct {
	cmd redis.Cmdable
}

func NewRedisJwtHandler(cmd redis.Cmdable) Handler {
	return &RedisJwtHandler{cmd: cmd}
}

func (r *RedisJwtHandler) SetLoginToken(c *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := r.SetJwtToken(c, uid, ssid)
	if err != nil {
		return err
	}
	err = r.SetRefreshToken(c, uid, ssid)
	return err
}

func (r *RedisJwtHandler) SetJwtToken(c *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: c.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(AccessTokenKey)
	if err != nil {
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}

func (r *RedisJwtHandler) SetRefreshToken(c *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:  uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(RefreshTokenKey)
	if err != nil {
		return err
	}
	c.Header("x-refresh-token", tokenStr)
	return nil
}

func (r *RedisJwtHandler) ExtractToken(c *gin.Context) string {
	tokenHeader := c.GetHeader("Authorization")
	if tokenHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	segs := strings.SplitN(tokenHeader, " ", 2)
	if len(segs) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	return segs[1]
}

func (r *RedisJwtHandler) ClearToken(c *gin.Context) error {
	c.Header("x-jwt-token", "")
	c.Header("x-refresh-token", "")
	ca, _ := c.Get("claims")
	claims, ok := ca.(*UserClaims)
	if !ok {
		c.JSON(http.StatusOK, ginx.Result{
			Code: http.StatusInternalServerError,
			Msg:  "系统错误",
			Data: nil,
		})
		return nil
	}
	return r.cmd.Set(c, fmt.Sprintf("users:ssid:%s", claims.Ssid), "", time.Hour*24*7).Err()
}

func (r *RedisJwtHandler) CheckSession(c *gin.Context, ssid string) error {
	val, err := r.cmd.Exists(c, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return nil
	case errors.Is(err, nil):
		if val == 0 {
			return nil
		}
		return errors.New("session already expired")
	default:
		return err
	}
}
