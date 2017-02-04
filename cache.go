package internalCache

import (
	"sync"
	"time"
)

type (
	cacheItem struct {
		value interface{}
		expire time.Time
	}

	cache struct {
		data map[string]cacheItem
		duration time.Duration
		ticker *time.Ticker
		lock sync.RWMutex
	}
)

func NewCache(duration time.Duration, tickerDuration time.Duration) *cache {
	c:= cache{
		duration: duration,
		ticker: time.NewTicker(tickerDuration),
		data: make(map[string]cacheItem),
	}

	go c.cleaner()

	return &c
}

func (c *cache) cleaner() {
	for t := range c.ticker.C {
		for key, item := range c.data {
			if item.expire.Before(t) {
				c.Del(key)
			}
		}
	}
}

func (c *cache) Set(key string, value interface{}) {
	expire := time.Now().Add(c.duration)
	c.SetWithExpire(key, value, expire)
}

func (c *cache) SetWithExpire(key string, value interface{}, expire time.Time) {
	c.lock.Lock()

	_, ok := c.data[key]
	if ok {
		delete(c.data, key)
	}

	c.data[key] = cacheItem{
		value: value,
		expire: expire,
	}

	c.lock.Unlock()
}

func (c *cache) Get(key string) interface{} {
	data, ok := c.data[key]
	if !ok {
		return nil
	}

	now := time.Now()
	if data.expire.Before(now) {
		c.Del(key)
		return nil
	}

	return data.value
}

func (c *cache) Del(key string) {
	c.lock.Lock()
	delete(c.data, key)
	c.lock.Unlock()
}

func (c *cache) Destroy() {
	c.ticker.Stop()

	c.lock.Lock()

	for  key := range c.data {
		delete(c.data, key)
	}

	c.lock.Unlock()
}