package cache

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
	"sync"
	"time"
)

type LocalCodeCache struct {
	lock  sync.Mutex
	cache *lru.Cache[string, any]
}

type codeItem struct {
	code   string
	cnt    int
	expire time.Time
}

func (l *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (l *LocalCodeCache) cntKey(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s:cnt", biz, phone)
}
