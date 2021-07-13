package hw04lrucache

import (
	"errors"
	"sync"
)

var ErrCacheWithZeroCapacity error = errors.New("Could not initalize cache with zero capacity")

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
	cacheItems
}

type cacheItem struct {
	key   Key
	value interface{}
}

type cacheItems []*cacheItem

func NewCache(capacity int) Cache {
	if capacity == 0 {
		panic(ErrCacheWithZeroCapacity)
	}
	return &lruCache{
		capacity:   capacity,
		queue:      NewList(),
		items:      make(map[Key]*ListItem, capacity),
		cacheItems: cacheItems(make([]*cacheItem, 0, capacity)),
	}
}

func (ci cacheItems) findIndex(key Key) int {
	index := -1
	for i, cacheItem := range ci {
		if cacheItem.key == key {
			index = i
		}
	}
	return index
}

func (ci cacheItems) moveFront(key Key, value interface{}) {
	currCacheItem := &cacheItem{key, value}
	indexForReplace := ci.findIndex(key)
	// Нужно сначала удалить элемент по найденному индексу, а затем добавить его в конец
	ci = append(ci[:indexForReplace], ci[indexForReplace+1:]...)
	ci = append(ci, currCacheItem)
}

func (c *lruCache) Set(key Key, value interface{}) bool {

	val, ok := c.items[key]
	if ok {
		val.Value = value
		c.Lock()
		c.items[key] = val
		c.Unlock()
		c.queue.PushFront(value)
		c.cacheItems.moveFront(key, value)
	} else {
		c.queue.PushFront(value)
		c.Lock()
		c.items[key] = c.queue.Front()
		c.Unlock()
		c.cacheItems = append(c.cacheItems, &cacheItem{key, value})
	}

	if c.queue.Len() > c.capacity {
		lastElement := c.queue.Back()
		c.queue.Remove(lastElement)
		cacheItemToDelete := c.cacheItems[0]
		delete(c.items, cacheItemToDelete.key)
		c.cacheItems = c.cacheItems[1:]
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {

	val, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.PushFront(val.Value)
	c.cacheItems.moveFront(key, val.Value)
	return val.Value, true
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.Lock()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.Lock()
	c.cacheItems = cacheItems(make([]*cacheItem, 0, c.capacity))
}
