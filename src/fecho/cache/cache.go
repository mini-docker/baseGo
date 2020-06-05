package cache

import (
	"sync"
	"time"

	"errors"
)

// 简单缓存实现。可以迁移至实现了该接口的缓存.

var (
	ErrTimeExpire = errors.New("cache has been expire")
	ErrNotFound   = errors.New("key not found!")
	DefaultCache  Cache
)

type Cache interface {
	Get(string) (interface{}, error)        //  interface{}  -> 内存存储直接为数据结构，不需要序列化，兼容其他缓存格式，需要手动断言.
	Set(string, interface{}, time.Duration) // 0 代表不过期.
}

type kv struct {
	k      string
	v      interface{}
	expire time.Time
}

type memoryCache struct {
	kv map[string]kv
	sync.RWMutex

	tk *time.Ticker
}

func (m *memoryCache) Get(k string) (interface{}, error) {
	m.RLock()
	data, ok := m.kv[k]
	m.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	none := data.expire == time.Time{}
	if !none && data.expire.Sub(time.Now()) <= time.Duration(0) {
		return nil, ErrTimeExpire
	}
	return data.v, nil
}

func (m *memoryCache) Set(k string, v interface{}, t time.Duration) {
	var ex time.Time
	if t != 0 {
		ex = time.Now().Add(t)
	} else {
		ex = time.Time{}
	}

	m.Lock()
	m.kv[k] = kv{
		k:      k,
		v:      v,
		expire: ex,
	}
	m.Unlock()
}

func (m *memoryCache) run(closeSign chan struct{}) {
	for {
		select {
		case <-m.tk.C:
			m.Lock()
			now := time.Now()
			for k, v := range m.kv {
				none := v.expire == time.Time{}
				if !none && v.expire.Sub(now) <= time.Duration(0) {
					delete(m.kv, k)
				}
			}
			m.Unlock()
		case <-closeSign:
			m.tk.Stop()
			return
			//default:
		}
	}
}

func NewMemoryCache(closeSign chan struct{}) Cache {
	m := &memoryCache{
		tk: time.NewTicker(30 * time.Second),
		kv: make(map[string]kv),
	}
	DefaultCache = m
	go m.run(closeSign)
	return m
}

func Set(s string, v interface{}, t time.Duration) {
	DefaultCache.Set(s, v, t)
}

func Get(s string) (interface{}, error) {
	return DefaultCache.Get(s)
}
