package cache

import (
	"sync"

	"github.com/basedalex/nats-service/model"
)

type Cache struct {
	mutex sync.RWMutex
	data map[string]model.OrderData
}

func New() *Cache {
	return &Cache{
		data: map[string]model.OrderData{},
	}
}

func (c *Cache) Set(id string, data model.OrderData) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[id] = data
}

func (c *Cache) Get(id string) (model.OrderData, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	data, ok := c.data[id]
	return data, ok 
}
