package cache

import "time"

type LRUCache struct {
	*baseCache
}

func NewLRUCache(b *baseCache) *LRUCache {
	cache := &LRUCache{
		baseCache: b,
	}
	return cache
}

func (L LRUCache) Set(key string, value any) {
}

func (L LRUCache) SetWithExpiration(key string, value any, expiration time.Duration) {
}

func (L LRUCache) Get(key string) (any, error) {
	return nil, nil
}

func (L LRUCache) Remove(key string) bool {
	return false
}
