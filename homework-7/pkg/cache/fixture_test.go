package cache

import (
	"math/rand"
	"sync"
	"time"
)

var (
	fixture *CacheFixture
	once    sync.Once
)

type CacheFixture struct {
	lfu Cache
	rnd *rand.Rand
}

func NewCacheFixture(capacity int, defaultExpiration time.Duration) *CacheFixture {
	once.Do(func() {
		lfu, err := New(capacity).LFU().Expiration(defaultExpiration).Build()
		if err != nil {
			panic("can't initialize lfu cache")
		}
		source := rand.NewSource(time.Now().UnixNano())
		rnd := rand.New(source)
		fixture = &CacheFixture{
			lfu: lfu,
			rnd: rnd,
		}
	})
	return fixture
}

func (f *CacheFixture) GenerateCacheItem(maxExpirationMs int) cacheValue {
	key := f.generateRandomString(10)
	expiration := time.Millisecond * time.Duration(f.rnd.Intn(maxExpirationMs))
	valueType := f.rnd.Intn(4)
	switch valueType {
	case 0:
		return cacheValue{
			key:        key,
			value:      f.rnd.Intn(100000),
			expiration: expiration,
		}
	case 1:
		return cacheValue{
			key:        key,
			value:      f.rnd.Float64() * 10000,
			expiration: expiration,
		}
	case 2:
		return cacheValue{
			key:        key,
			value:      f.generateRandomString(5),
			expiration: expiration,
		}
	default:
		return cacheValue{
			key:        key,
			value:      f.rnd.Intn(2)%2 == 0,
			expiration: expiration,
		}
	}
}

func (f *CacheFixture) generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[f.rnd.Intn(len(charset))]
	}
	return string(result)
}
