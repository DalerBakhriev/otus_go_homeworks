package hw04lrucache

import (
	"errors"
	"sync"
)

var ErrCacheWithZeroOrNegativeCapacity error = errors.New("Could not initalize cache with zero capacity")

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

// Cache instance initializer. Panics in case of zero or negative capacity
func NewCache(capacity int) Cache {
	if capacity <= 0 {
		panic(ErrCacheWithZeroOrNegativeCapacity)
	}
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {

	val, ok := c.items[key]
	newValue := &cacheItem{key, value}
	if ok {
		val.Value = newValue
		c.Lock()
		c.items[key] = val
		c.Unlock()
		c.queue.MoveToFront(val)
	} else {
		c.queue.PushFront(newValue)
		c.Lock()
		c.items[key] = c.queue.Front()
		c.Unlock()
	}

	if c.queue.Len() > c.capacity {
		lastElement := c.queue.Back()
		c.queue.Remove(lastElement)
		delete(c.items, lastElement.Value.(*cacheItem).key)
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {

	val, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(val)
	return val.Value.(*cacheItem).value, true
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.Lock()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.Unlock()
}
