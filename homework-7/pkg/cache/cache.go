package cache

import (
	"fmt"
	"sync"
	"time"
)

const (
	TYPE_LRU = "lru"
	TYPE_LFU = "lfu"

	maxDuration time.Duration = 1<<63 - 1
)

type (
	DeserializeFunc func(interface{}, interface{}) (interface{}, error)
	SerializeFunc   func(interface{}, interface{}) (interface{}, error)
)

type Cache interface {
	Set(key string, value any)
	SetWithExpiration(key string, value any, expiration time.Duration)
	Get(key string) (any, error)
	Remove(key string) bool
}

type CacheBuilder struct {
	baseCache baseCache
}

type baseCache struct {
	evictType         string
	capacity          int
	defaultExpiration time.Duration
	serializeFunc     SerializeFunc
	deserializeFunc   DeserializeFunc
	sync.RWMutex
}

func New(capacity int) *CacheBuilder {
	return &CacheBuilder{
		baseCache{
			evictType:         TYPE_LFU,
			capacity:          capacity,
			defaultExpiration: maxDuration,
		},
	}
}

func (cb *CacheBuilder) LRU() *CacheBuilder {
	return cb.evictType(TYPE_LRU)
}

func (cb *CacheBuilder) LFU() *CacheBuilder {
	return cb.evictType(TYPE_LFU)
}

func (cb *CacheBuilder) evictType(tp string) *CacheBuilder {
	cb.baseCache.evictType = tp
	return cb
}

func (cb *CacheBuilder) Expiration(expiration time.Duration) *CacheBuilder {
	cb.baseCache.defaultExpiration = expiration
	return cb
}

func (cb *CacheBuilder) SerializeFunc(serializeFunc SerializeFunc) *CacheBuilder {
	cb.baseCache.serializeFunc = serializeFunc
	return cb
}

func (cb *CacheBuilder) DeserializeFunc(deserializeFunc DeserializeFunc) *CacheBuilder {
	cb.baseCache.deserializeFunc = deserializeFunc
	return cb
}

func (cb *CacheBuilder) Build() (Cache, error) {
	if cb.baseCache.capacity <= 0 {
		return nil, fmt.Errorf("invalid cache capacity: %d", cb.baseCache.capacity)
	}
	switch cb.baseCache.evictType {
	case TYPE_LRU:
		return NewLRUCache(&cb.baseCache), nil
	case TYPE_LFU:
		return NewLFUCache(&cb.baseCache), nil
	default:
		return nil, fmt.Errorf("unknown cache type: %s", cb.baseCache.evictType)
	}
}
