package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type Item struct {
	Key   Key
	Value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if itemFromCache, ok := c.items[key]; ok == true {
		c.queue.MoveToFront(itemFromCache)
		itemFromCache.Value.(*Item).Value = value
		return true
	}
	if c.queue.Len() == c.capacity {
		itemLast := c.queue.Back()
		c.queue.Remove(itemLast)
		delete(c.items, itemLast.Value.(*Item).Key)
	}
	item := &Item{
		Key:   key,
		Value: value,
	}
	itemNew := c.queue.PushFront(item)
	c.items[item.Key] = itemNew
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	itemFromCache, ok := c.items[key]
	if ok == true {
		c.queue.MoveToFront(itemFromCache)
		return itemFromCache.Value.(*Item).Value, true
	} else {
		return nil, false
	}
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
