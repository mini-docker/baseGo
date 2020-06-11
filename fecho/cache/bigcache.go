// 对bigcache的基本封装
package cache

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	defaultShards = 16
)

var (
	bcache       *cache.Cache
	bMutex       sync.RWMutex
	ErrNotInit   = errors.New(`bcache not initialized`)
	ErrExpire    = errors.New(`cache expired`)
	cacheControl bool
)

// 过期时间.
func InitCache(exire time.Duration, debug bool) error {
	bcache = cache.New(exire*time.Minute, 10*time.Minute)
	if debug {
		cacheControl = true
	}
	return nil
}

// 设置key  cacheControl 为true 的时候就不缓存.
func BSet(key string, val []byte) error {
	bMutex.RLock()
	defer bMutex.RUnlock()

	if cacheControl {
		return nil
	}
	if bcache == nil {
		return ErrNotInit
	}

	bcache.SetDefault(key, val)

	return nil
}

// 获取.
func BGet(key string) ([]byte, error) {
	bMutex.RLock()
	defer bMutex.RUnlock()

	if bcache == nil {
		return nil, ErrNotInit
	}

	b, has := bcache.Get(key)
	if !has {
		return nil, ErrExpire
	}

	bys, ok := b.([]byte)
	if !ok {
		return nil, errors.New(`bug: cache not []byte `)
	}

	return bys, nil
}

// 清除缓存.
func Reset(siteId string, siteIndexId string) error {
	bMutex.RLock()
	defer bMutex.RUnlock()
	if bcache == nil {
		return ErrNotInit
	}

	for k := range bcache.Items() {
		if strings.HasPrefix(k, siteId+siteIndexId) {
			bcache.Delete(k)
		}
	}

	return nil
}

func ResetAll(siteids []string) error {
	bMutex.RLock()
	defer bMutex.RUnlock()

	if bcache == nil {
		return ErrNotInit
	}

	for _, v := range siteids {
		for k := range bcache.Items() {
			if strings.HasPrefix(k, v) {
				bcache.Delete(k)
			}
		}
	}
	return nil
}
