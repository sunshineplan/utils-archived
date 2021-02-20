package cache

import (
	"sync"
	"time"
)

type item struct {
	Value      interface{}
	Duration   time.Duration
	Expiration int64
	Regenerate func() interface{}
}

func (i item) Expired() bool {
	if i.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.Expiration
}

// Cache is cache struct.
type Cache struct {
	cache sync.Map
}

// New creates a new cache with auto clean or not.
func New(autoClean bool) *Cache {
	c := &Cache{}

	if autoClean {
		go c.check()
	}

	return c
}

// Set sets cache value for a key, if f is presented, this value will regenerate when expired.
func (c *Cache) Set(key, value interface{}, d time.Duration, f func() interface{}) {
	c.cache.Store(key, item{
		Value:      value,
		Duration:   d,
		Expiration: time.Now().Add(d).UnixNano(),
		Regenerate: f,
	})
}

func (c *Cache) regenerate(key interface{}, i item) {
	i.Value = i.Regenerate()
	i.Expiration = time.Now().Add(i.Duration).UnixNano()
	c.cache.Store(key, i)
}

// Get gets cache value by key and whether value was found.
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	value, ok := c.cache.Load(key)
	if !ok {
		return nil, false
	}

	i := value.(item)
	if i.Expired() {
		if i.Regenerate == nil {
			c.cache.Delete(key)
			return nil, false
		}
		defer c.regenerate(key, i)
	}

	return i.Value, true
}

// Delete deletes the value for a key.
func (c *Cache) Delete(key interface{}) {
	c.cache.Delete(key)
}

// Empty deletes all values in cache.
func (c *Cache) Empty() {
	c.cache.Range(func(key, _ interface{}) bool {
		c.cache.Delete(key)
		return true
	})
}

func (c *Cache) check() {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			c.cache.Range(func(key, value interface{}) bool {
				i := value.(item)
				if i.Expired() {
					c.cache.Delete(key)
					if i.Regenerate != nil {
						defer c.regenerate(key, i)
					}
				}

				return true
			})
		}
	}
}
