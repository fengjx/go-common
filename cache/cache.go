package cache

import (
	"context"
)

// Result 缓存结果
type Result[V any] struct {
	val V
	err error
}

// Val 获取缓存值
func (r *Result[V]) Val() V {
	return r.val
}

// Err 获取缓存错误
func (r *Result[V]) Err() error {
	return r.err
}

// Result 获取缓存结果
func (r *Result[V]) Result() (V, error) {
	return r.val, r.err
}

// Cache 缓存接口定义
type Cache[K comparable, V any] interface {

	// Get 获取缓存值
	Get(ctx context.Context, key K) *Result[V]

	// GetWithFallback 获取缓存值，当缓存未命中时，使用当前回源函数查询
	GetWithFallback(ctx context.Context, key K, fn Fallback[K, V]) *Result[V]

	// GetMulti 批量获取缓存值
	GetMulti(ctx context.Context, keys []K) *Result[map[K]V]

	// GetMultiWithFallback 批量获取缓存值，当缓存未命中时，使用当前回源函数查询
	GetMultiWithFallback(ctx context.Context, keys []K, fn FallbackMulti[K, V]) *Result[map[K]V]

	// Set 设置缓存值
	Set(ctx context.Context, key K, val V) *Result[bool]

	// SetMulti 批量设置缓存值
	SetMulti(ctx context.Context, values map[K]V) *Result[bool]

	// Del 删除缓存，返回被删除的数量
	Del(ctx context.Context, keys ...K) *Result[int]

	// Has 判断缓存是否存在
	Has(ctx context.Context, key K) *Result[bool]

	// Clear 清空缓存
	Clear(ctx context.Context) *Result[bool]
}

// Fallback 缓存回源查询函数
type Fallback[K comparable, V any] func(ctx context.Context, missKey K) (V, error)

// FallbackMulti 缓存回源查询函数，查询多个
type FallbackMulti[K comparable, V any] func(ctx context.Context, missKeys []K) (map[K]V, error)

// Options 缓存配置
type Options struct {
	cacheEmpty bool
}

// Option 配置项赋值函数
type Option func(*Options)

// WithCacheEmpty 设置是否缓存空值
func WithCacheEmpty(cacheEmpty bool) Option {
	return func(o *Options) {
		o.cacheEmpty = cacheEmpty
	}
}
