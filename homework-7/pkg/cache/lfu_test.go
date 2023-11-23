package cache

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		defaultExpiration time.Duration
		key               string
		value             interface{}
		requestKey        string
		waitTime          time.Duration
		expectedError     error
	}{
		{
			name:              "Success1",
			defaultExpiration: time.Second * 10,
			key:               "Key1",
			value:             "Dima",
			requestKey:        "Key1",
			waitTime:          0,
			expectedError:     nil,
		},
		{
			name:              "Success2",
			defaultExpiration: time.Second * 10,
			key:               "Key2",
			value:             322,
			requestKey:        "Key2",
			waitTime:          0,
			expectedError:     nil,
		},
		{
			name:              "Fail1 (expired)",
			key:               "Key3",
			defaultExpiration: time.Millisecond * 20,
			value:             1909.0,
			requestKey:        "Key3",
			waitTime:          time.Millisecond * 100,
			expectedError:     errors.New("the element with key: Key3 has expired"),
		},
		{
			name:              "Fail2 (no element with key)",
			key:               "Key4",
			defaultExpiration: time.Second * 2,
			value:             false,
			requestKey:        "Key2389",
			waitTime:          0,
			expectedError:     errors.New("no cached element with key: Key2389"),
		},
		{},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			lfu, err := New(1).LFU().Expiration(tc.defaultExpiration).Build()
			assert.NoError(t, err)

			lfu.Set(tc.key, tc.value)

			time.Sleep(tc.waitTime)

			result, err := lfu.Get(tc.requestKey)
			if tc.expectedError != nil {
				assert.Nil(t, result)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.value, result)
			}
		})
	}
}

func TestSetWithExpiration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		defaultExpiration time.Duration
		key               string
		value             interface{}
		expiration        time.Duration
		requestKey        string
		waitTime          time.Duration
		expectedError     error
	}{
		{
			name:              "Success1",
			defaultExpiration: 0,
			key:               "Key1",
			value:             "Dima",
			expiration:        time.Second * 100,
			requestKey:        "Key1",
			waitTime:          0,
			expectedError:     nil,
		},
		{
			name:              "Success2",
			defaultExpiration: time.Second * 10,
			key:               "Key2",
			value:             2434,
			expiration:        time.Second * 5,
			requestKey:        "Key2",
			waitTime:          0,
			expectedError:     nil,
		},
		{
			name:              "Fail 1",
			defaultExpiration: time.Second * 10,
			key:               "Key1",
			value:             true,
			expiration:        time.Millisecond * 10,
			requestKey:        "Key1",
			waitTime:          time.Millisecond * 100,
			expectedError:     errors.New("the element with key: Key1 has expired"),
		},
		{
			name:              "Fail 2",
			defaultExpiration: time.Second * 10,
			key:               "Key4",
			value:             2390.483,
			expiration:        time.Second,
			requestKey:        "Key3443",
			waitTime:          0,
			expectedError:     errors.New("no cached element with key: Key3443"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			lfu, err := New(1).LFU().Expiration(tc.defaultExpiration).Build()
			assert.NoError(t, err)

			lfu.SetWithExpiration(tc.key, tc.value, tc.expiration)

			time.Sleep(tc.waitTime)

			result, err := lfu.Get(tc.requestKey)
			if tc.expectedError != nil {
				assert.Nil(t, result)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.value, result)
			}
		})
	}
}

type cacheValue struct {
	key        string
	value      interface{}
	expiration time.Duration
}

func TestEviction(t *testing.T) {
	source := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(source)
	f := &CacheFixture{rnd: rnd}

	N := 100000
	capacity := 10

	lfu, err := New(capacity).LFU().Expiration(time.Second).Build()
	assert.NoError(t, err)

	values := make([]cacheValue, N)
	for i := 0; i < N; i++ {
		values[i] = f.GenerateCacheItem(400)
	}

	for i := 0; i < capacity-1; i++ {
		lfu.Set(values[i].key, values[i].value)

		// набиваем частоту использования для первых capacity - 1 эл-тов
		for j := 0; j < 100; j++ {
			val, err := lfu.Get(values[i].key)
			assert.NoError(t, err)
			assert.Equal(t, values[i].value, val)
		}
	}

	lfu.Set(values[capacity-1].key, values[capacity-1].value)

	for i := capacity; i < N; i++ {
		lfu.Set(values[i].key, values[i].value)

		val, err := lfu.Get(values[i].key)
		assert.NoError(t, err)
		assert.Equal(t, values[i].value, val)

		val, err = lfu.Get(values[i-1].key)
		assert.EqualError(t, err, fmt.Sprintf("no cached element with key: %s", values[i-1].key))
		assert.Nil(t, val)
	}

}

func TestDataRaces(t *testing.T) {

	var wg sync.WaitGroup
	N := 100000

	f := NewCacheFixture(N, time.Second*2)

	values := make([]cacheValue, N)
	for i := 0; i < N; i++ {
		values[i] = f.GenerateCacheItem(400)
	}

	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				values[i].expiration += time.Second * 10
			}
			f.lfu.SetWithExpiration(values[i].key, values[i].value, values[i].expiration)
		}(i)
	}

	wg.Wait()
	time.Sleep(400 * time.Millisecond)
	for i := 0; i < N; i++ {
		go func(i int, t *testing.T) {
			result, err := f.lfu.Get(values[i].key)
			if i%2 == 0 {
				assert.NoError(t, err)
				assert.Equal(t, values[i].value, result)

				f.lfu.Remove(values[i].key)
				result, err = f.lfu.Get(values[i].key)
				assert.EqualError(t, err, fmt.Sprintf("no cached element with key: %s", values[i].key))
				assert.Nil(t, result)

			} else {
				assert.Nil(t, result)
				assert.EqualError(t, err, fmt.Sprintf("the element with key: %s has expired", values[i].key))
			}
		}(i, t)
	}

}
