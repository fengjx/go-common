package cache

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/fengjx/go-halo/cache/internal/lru"
)

// LRUCached 使用 LRU 实现缓存
type lruCache[K comparable, V any] struct {
	lru           *lru.Cache
	mu            sync.RWMutex
	ttl           time.Duration
	cacheEmpty    bool
	fallbackMulti FallbackMulti[K, V]
}

type lruItem[K comparable, V any] struct {
	key      K
	val      V
	expireAt int64
}

// IsExpire 判断是否过期
// 返回 true 表示已过期
func (i lruItem[K, V]) IsExpire() bool {
	return i.expireAt <= time.Now().Unix()
}

// NewLRUCache 创建一个 LRU 缓存
func NewLRUCache[K comparable, V any](capacity int, ttl time.Duration, fallback FallbackMulti[K, V], opts ...Option) Cache[K, V] {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	cache := &lruCache[K, V]{
		lru:           lru.New(capacity),
		ttl:           ttl,
		cacheEmpty:    opt.cacheEmpty,
		fallbackMulti: fallback,
	}
	return cache
}

func (c *lruCache[K, V]) Get(ctx context.Context, key K) *Result[V] {
	return c.GetWithFallback(ctx, key, func(ctx context.Context, missKey K) (v V, err error) {
		var m map[K]V
		m, err = c.fallbackMulti(ctx, []K{missKey})
		if err != nil {
			return
		}
		return m[missKey], nil
	})
}

func (c *lruCache[K, V]) GetWithFallback(ctx context.Context, key K, fn Fallback[K, V]) *Result[V] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.getCache(key)
	if ok {
		return &Result[V]{val: v, err: nil}
	}
	if c.fallbackMulti != nil {
		vf, err := fn(ctx, key)
		if err != nil {
			return &Result[V]{err: err}
		}
		if !c.cacheEmpty && isEmpty(vf) {
			return &Result[V]{err: nil}
		}
		v = vf
	}
	c.setCache(key, v)
	return &Result[V]{val: v, err: nil}
}

func (c *lruCache[K, V]) getCache(key K) (v V, ok bool) {
	val, ok := c.lru.Get(key)
	if ok {
		item := val.(*lruItem[K, V])
		if !item.IsExpire() {
			c.lru.Remove(key)
			return
		}
		return item.val, true
	}
	return
}

func (c *lruCache[K, V]) GetMulti(ctx context.Context, keys []K) *Result[map[K]V] {
	return c.GetMultiWithFallback(ctx, keys, c.fallbackMulti)
}

func (c *lruCache[K, V]) GetMultiWithFallback(ctx context.Context, keys []K, fn FallbackMulti[K, V]) *Result[map[K]V] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	vals := make(map[K]V)
	var missKeys []K
	for _, k := range keys {
		v, ok := c.getCache(k)
		if ok {
			vals[k] = v
		} else {
			missKeys = append(missKeys, k)
		}
	}
	if len(missKeys) == 0 {
		return &Result[map[K]V]{val: vals, err: nil}
	}
	if c.fallbackMulti != nil {
		m, _ := fn(ctx, missKeys)
		for _, k := range missKeys {
			v := m[k]
			if !c.cacheEmpty && isEmpty(v) {
				continue
			}
			c.setCache(k, v)
			vals[k] = v
		}
	}
	return &Result[map[K]V]{val: vals, err: nil}
}

func (c *lruCache[K, V]) Set(ctx context.Context, key K, val V) *Result[bool] {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setCache(key, val)
	return &Result[bool]{true, nil}
}

func (c *lruCache[K, V]) setCache(key K, val V) {
	item := &lruItem[K, V]{
		key:      key,
		val:      val,
		expireAt: time.Now().Add(c.ttl).Unix(),
	}
	c.lru.Add(key, item)
}

func (c *lruCache[K, V]) SetMulti(ctx context.Context, values map[K]V) *Result[bool] {
	for k, v := range values {
		c.Set(ctx, k, v)
	}
	return &Result[bool]{true, nil}
}

func (c *lruCache[K, V]) Del(ctx context.Context, keys ...K) *Result[int] {
	c.mu.Lock()
	defer c.mu.Unlock()
	cnt := 0
	for _, k := range keys {
		exist := c.lru.Remove(k)
		if exist {
			cnt++
		}
	}
	return &Result[int]{cnt, nil}
}

func (c *lruCache[K, V]) Has(ctx context.Context, key K) *Result[bool] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.getCache(key)
	if !ok {
		return &Result[bool]{false, nil}
	}
	return &Result[bool]{ok, nil}
}

func (c *lruCache[K, V]) Clear(ctx context.Context) *Result[bool] {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lru.Clear()
	return &Result[bool]{true, nil}
}

// isNotEmpty 判断一个值是否是nil
func isEmpty(v any) bool {
	vl := reflect.ValueOf(v)
	return vl.Kind() != reflect.Pointer || vl.IsNil()
}
