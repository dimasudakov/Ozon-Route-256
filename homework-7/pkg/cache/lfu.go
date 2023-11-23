package cache

import (
	"fmt"
	"math"
	"time"
)

type LFUCache struct {
	storage map[string]*lfuItem
	length  uint32
	*baseCache
}

type lfuItem struct {
	value      any
	frequency  uint32
	expiration time.Time
}

func NewLFUCache(b *baseCache) *LFUCache {
	cache := &LFUCache{
		storage:   make(map[string]*lfuItem),
		length:    0,
		baseCache: b,
	}
	return cache
}

func (c *LFUCache) Set(key string, value any) {
	c.SetWithExpiration(key, value, c.defaultExpiration)
}

func (c *LFUCache) SetWithExpiration(key string, value any, expiration time.Duration) {
	c.Lock()
	defer c.Unlock()

	if _, exist := c.storage[key]; exist {
		c.storage[key].value = value
	} else {
		if c.length == uint32(c.capacity) {
			c.evict()
		}
		c.storage[key] = &lfuItem{
			value: value,
		}
		c.length++
	}
	c.storage[key].expiration = time.Now().Add(expiration)
	c.storage[key].frequency++
}

func (c *LFUCache) Remove(key string) bool {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.storage[key]; ok {
		delete(c.storage, key)
		return true
	}
	return false
}

func (c *LFUCache) Get(key string) (any, error) {
	c.RLock()
	item, ok := c.storage[key]
	c.RUnlock()

	if ok {
		if item.expiration.Before(time.Now()) {
			c.Lock()
			delete(c.storage, key)
			c.Unlock()
			return nil, fmt.Errorf("the element with key: %s has expired", key)
		}
		c.Lock()
		item.frequency++
		c.Unlock()
		return item.value, nil
	} else {
		return nil, fmt.Errorf("no cached element with key: %s", key)
	}
}

func (c *LFUCache) evict() {
	minFrequency := uint32(math.MaxUint32)
	var evictKey string

	for key, item := range c.storage {
		if item.expiration.Before(time.Now()) {
			delete(c.storage, key)
			return
		}
		if item.frequency < minFrequency {
			minFrequency = item.frequency
			evictKey = key
		}
	}

	delete(c.storage, evictKey)
	c.length--
}
