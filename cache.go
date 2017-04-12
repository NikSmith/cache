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

	Cache struct {
		io.Closer

		Duration       time.Duration
		TickerDuration time.Duration

		data map[string]*cacheItem
		lock sync.RWMutex
		once sync.Once
	}
)

func New(duration time.Duration, tickerDuration time.Duration) *Cache {
	return &Cache{
		Duration:       duration,
		TickerDuration: tickerDuration,
	}
}

func (c *Cache) cleaner() {
	go func() {
		for {
			c.lock.RLock()
			duration := c.TickerDuration
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
	}()
}

func (c *Cache) Set(key string, value interface{}) {
	expire := time.Now().Add(c.Duration)
	c.SetWithExpire(key, value, expire)
}

func (c *Cache) SetWithExpire(key string, value interface{}, expire time.Time) {
	c.internalInit()

	c.lock.Lock()
	c.data[key] = &cacheItem{value, expire}
	c.lock.Unlock()
}

func (c *Cache) Get(key string) (value interface{}) {
	c.internalInit()

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

func (c *Cache) Del(key string) {
	c.internalInit()

	c.lock.Lock()
	delete(c.data, key)
	c.lock.Unlock()
}

func (c *Cache) Close() error {

	c.lock.Lock()
	defer c.lock.Unlock()

	c.TickerDuration = 0
	c.data = nil

	return nil
}

func (c *Cache) internalInit() {
	c.once.Do(func() {
		c.data = make(map[string]*cacheItem)
		c.cleaner()
	})
}
