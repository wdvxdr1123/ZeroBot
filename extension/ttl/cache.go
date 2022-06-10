package ttl

import (
	"sync"
	"time"
)

// Cache is a synchronised map of items that auto-expire once stale
type Cache[K comparable, V any] struct {
	sync.RWMutex
	ttl   time.Duration
	items map[K]*Item[V]
	onset func(K, V)
	onget func(K, V)
	ondel func(K, V)
	ontch func(K, V)
	stop  func() // Stop stops the gc loop
}

// NewCache 创建指定生命周期的 Cache
func NewCache[K comparable, V any](ttl time.Duration) *Cache[K, V] {
	return NewCacheOn(ttl, [4]func(K, V){})
}

// NewCacheOn 创建指定生命周期的 Cache
//   on: [onset, onget, ondel, ontouch]
func NewCacheOn[K comparable, V any](ttl time.Duration, on [4]func(K, V)) *Cache[K, V] {
	cache := &Cache[K, V]{
		ttl:   ttl,
		items: map[K]*Item[V]{},
		onset: on[0],
		onget: on[1],
		ondel: on[2],
		ontch: on[3],
	}
	cache.stop = cache.gc() // async gc
	return cache
}

func (c *Cache[K, V]) gc() (stop func()) {
	ticker := time.NewTicker(time.Minute)
	stopchan := make(chan struct{})
	go func() {
	loop:
		for {
			select {
			case <-ticker.C:
				c.Lock()
				for key, item := range c.items {
					if item.expired() {
						if c.ondel != nil {
							c.ondel(key, c.items[key].value)
						}
						delete(c.items, key)
					}
				}
				c.Unlock()
			case <-stopchan:
				break loop
			}
		}
	}()
	return func() {
		ticker.Stop()
		stopchan <- struct{}{}
	}
}

// Destroy 销毁 chahe, 不可再使用, 否则 panic
func (c *Cache[K, V]) Destroy() {
	c.stop()
	c.items = nil
	c.stop = nil
}

// Get 通过 key 获取指定的元素
func (c *Cache[K, V]) Get(key K) (v V) {
	c.RLock()
	item, ok := c.items[key]
	c.RUnlock()
	if ok && item.expired() {
		c.Delete(key)
		return
	}
	if item == nil {
		return
	}
	item.exp = time.Now().Add(c.ttl) // reset the expired time
	if c.onget != nil {
		c.onget(key, item.value)
	}
	return item.value
}

// Set 设置指定 key 的值
func (c *Cache[K, V]) Set(key K, val V) {
	c.Lock()
	defer c.Unlock()
	item := &Item[V]{
		exp:   time.Now().Add(c.ttl),
		value: val,
	}
	c.items[key] = item
	if c.onset != nil {
		c.onset(key, val)
	}
}

// Delete 删除指定key
func (c *Cache[K, V]) Delete(key K) {
	c.Lock()
	defer c.Unlock()
	if c.ondel != nil {
		c.ondel(key, c.items[key].value)
	}
	delete(c.items, key)
}

// Touch 为指定key添加一定生命周期
func (c *Cache[K, V]) Touch(key K, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()
	if c.items[key] != nil {
		c.items[key].exp = c.items[key].exp.Add(ttl)
		if c.ontch != nil {
			c.ontch(key, c.items[key].value)
		}
	}
}
