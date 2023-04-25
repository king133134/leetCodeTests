package cache

import (
	"sync"
	"time"
)

var cacheInitOnce sync.Once

var cacheSingle *Cache

func NewCache() *Cache {
	cacheInitOnce.Do(cacheInit)
	return cacheSingle
}

func cacheInit() {
	cacheSingle = &Cache{data: map[uint32]*cacheItem{}}
}

func hashKey(key interface{}) uint32 {
	if val, ok := key.(uint32); ok {
		return val
	}
	return hash32(key.(string))
}

type cacheItem struct {
	value      interface{}
	expiration int
}

type Cache struct {
	data map[uint32]*cacheItem
	mu   sync.Mutex
}

func (_this *Cache) Set(key string, value interface{}, second int) {
	defer _this.mu.Unlock()
	_this.mu.Lock()
	_this.data[hash32(key)] = &cacheItem{value: value, expiration: second + time.Now().Second()}
}

func (_this *Cache) del(key interface{}, lock bool) {
	if lock {
		_this.mu.Lock()
		defer _this.mu.Unlock()
	}
	delete(_this.data, hashKey(key))
}

func (_this *Cache) Del(key string) {
	_this.del(key, true)
}

func (_this *Cache) Get(key string) (interface{}, bool) {
	defer _this.mu.Unlock()
	_this.mu.Lock()
	h := hash32(key)
	item, ok := _this.data[h]
	if !ok {
		return nil, false
	}
	if item.expiration < time.Now().Second() {
		_this.del(h, false)
		return nil, false
	}
	return item.value, true
}

type RememberFunc func() interface{}

func (_this *Cache) Remember(key string, f RememberFunc, second int) (value interface{}) {
	value, ok := _this.Get(key)
	if ok {
		return
	}
	value = f()
	_this.Set(key, value, second)
	return
}
