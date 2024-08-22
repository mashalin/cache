package cache

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	data       map[string]interface{}
	expiration map[string]time.Time
	mu         sync.Mutex
}

func New() *Cache {
	return &Cache{
		data:       make(map[string]interface{}),
		expiration: make(map[string]time.Time),
	}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	if ttl > 0 {
		c.expiration[key] = time.Now().Add(ttl)
	} else {
		delete(c.expiration, key)
	}
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if expiration, ok := c.expiration[key]; ok {
		if time.Now().After(expiration) {
			c.deleteWithoutLock(key)
			return nil, errors.New("no such value")
		}
	}

	return c.data[key], nil
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.deleteWithoutLock(key)
}

func (c *Cache) deleteWithoutLock(key string) {
	delete(c.data, key)
	delete(c.expiration, key)
}
