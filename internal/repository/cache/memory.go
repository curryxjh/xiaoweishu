package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

type LocalCodeCache struct {
	lock       sync.Mutex
	cache      *lru.Cache[string, any]
	rwLock     sync.RWMutex
	expiration time.Duration
}

type codeItem struct {
	code   string
	cnt    int
	expire time.Time
}

func NewLocalCodeCache(c *lru.Cache[string, any], expiration time.Duration) CodeCache {
	return &LocalCodeCache{
		cache:      c,
		expiration: expiration,
	}
}

func (l *LocalCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	key := l.key(biz, phone)
	now := time.Now()
	val, ok := l.cache.Get(key)
	if !ok {
		l.cache.Add(key, codeItem{
			code:   code,
			cnt:    3,
			expire: now.Add(l.expiration),
		})
		return nil
	}
	item, ok := val.(codeItem)
	if !ok {
		return errors.New("server error")
	}
	if item.expire.Sub(now) > time.Minute*9 {
		return ErrCodeSendTooMany
	}
	l.cache.Add(key, codeItem{
		code:   code,
		cnt:    3,
		expire: now.Add(l.expiration),
	})
	return nil
}

func (l *LocalCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	key := l.key(biz, phone)
	val, ok := l.cache.Get(key)
	if !ok {
		return false, ErrKeyNotFound
	}
	item, ok := val.(codeItem)
	if !ok {
		return false, errors.New("server error")
	}
	if item.cnt <= 0 {
		return false, ErrCodeVerifyTooMany
	}
	item.cnt--
	l.cache.Add(key, item)
	if item.code != inputCode {
		return false, errors.New("code not match")
	}
	l.cache.Remove(key)
	return true, nil
}

func (l *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (l *LocalCodeCache) cntKey(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s:cnt", biz, phone)
}
