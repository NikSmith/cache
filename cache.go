package cache

import (
	"io"
	"sync"
	"time"
)

type (
	cacheItem struct {
		value  interface{}
		expire time.Time
	}

	cache struct {
		io.Closer

		data           map[string]*cacheItem
		duration       time.Duration
		tickerDuration time.Duration
		lock           sync.RWMutex
	}
)

func New(duration time.Duration, tickerDuration time.Duration) *cache {
	c := cache{
		duration:       duration,
		tickerDuration: tickerDuration,
		data:           make(map[string]*cacheItem),
	}

	go c.cleaner()

	return &c
}

func (c *cache) cleaner() {

	for {
		c.lock.RLock()
		duration := c.tickerDuration
		c.lock.RUnlock()

		if duration == 0 {
			break
		}

		time.Sleep(duration)

		c.lock.Lock()

		for key, item := range c.data {
			if item.expire.Before(time.Now()) {
				delete(c.data, key)
			}
		}

		c.lock.Unlock()
	}
}

func (c *cache) Set(key string, value interface{}) {
	expire := time.Now().Add(c.duration)
	c.SetWithExpire(key, value, expire)
}

func (c *cache) SetWithExpire(key string, value interface{}, expire time.Time) {

	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[key] = &cacheItem{
		value:  value,
		expire: expire,
	}
}

func (c *cache) Get(key string) (value interface{}) {

	c.lock.RLock()

	item, exist := c.data[key]
	valid := exist && !item.expire.Before(time.Now())
	if valid {
		value = item.value
	}

	c.lock.RUnlock()

	if exist && !valid {
		go c.Del(key)
	}

	return
}

func (c *cache) Del(key string) {

	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.data, key)
}

func (c *cache) Close() error {

	c.lock.Lock()
	defer c.lock.Unlock()

	c.tickerDuration = 0
	c.data = nil

	return nil
}
