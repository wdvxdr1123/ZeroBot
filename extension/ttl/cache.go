package ttl

import (
	"sync"
	"time"
)

// Cache is a synchronised map of items that auto-expire once stale
type Cache struct {
	sync.RWMutex
	ttl   time.Duration
	items map[interface{}]*Item
}

// NewCache 创建指定生命周期的 Cache
func NewCache(ttl time.Duration) *Cache {
	cache := &Cache{
		ttl:   ttl,
		items: map[interface{}]*Item{},
	}
	go cache.gc() // async gc
	return cache
}

func (c *Cache) gc() {
	ticker := time.NewTicker(time.Minute)
	for {
		<-ticker.C
		c.Lock()
		for key, item := range c.items {
			if item.expired() {
				delete(c.items, key)
			}
		}
		c.Unlock()
	}
}

// Get 通过 key 获取指定的元素
func (c *Cache) Get(key interface{}) interface{} {
	c.RLock()
	item, ok := c.items[key]
	c.RUnlock()
	if ok && item.expired() {
		c.Delete(key)
		return nil
	}
	if item == nil {
		return nil
	}
	item.exp = time.Now().Add(c.ttl) // reset the expired time
	return item.value
}

// Set 设置指定 key 的值
func (c *Cache) Set(key interface{}, val interface{}) {
	c.Lock()
	defer c.Unlock()
	item := &Item{
		exp:   time.Now().Add(c.ttl),
		value: val,
	}
	c.items[key] = item
}

// Delete 删除指定key
func (c *Cache) Delete(key interface{}) {
	c.Lock()
	defer c.Unlock()
	delete(c.items, key)
}

// Touch 为指定key添加一定生命周期
func (c *Cache) Touch(key interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()
	if c.items[key] != nil {
		c.items[key].exp = c.items[key].exp.Add(ttl)
	}
}
