package domain

import "time"

// User 领域对象，是 DDD 中的 entity
type User struct {
	Id       int64
	Email    string
	Phone    string
	Password string
	NickName string
	Birthday string
	AboutMe  string
	Ctime    time.Time
}
