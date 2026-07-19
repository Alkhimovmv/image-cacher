package lru

import "image_cache/internal/linkedlist"

type CacheItem[V any] struct {
	Key   string
	Value V
	Size  uint
}

type LRU[V any] struct {
	maxSize uint
	size    uint
	items   map[string]*linkedlist.Item[CacheItem[V]]
	list    *linkedlist.List[CacheItem[V]]
}

func New[V any](maxSize uint) *LRU[V] {
	return &LRU[V]{
		maxSize: maxSize,
		items:   make(map[string]*linkedlist.Item[CacheItem[V]]),
		list:    linkedlist.New[CacheItem[V]](),
	}
}

func (c *LRU[V]) Get(key string) (V, bool) {
	item, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}

	c.list.MoveToFront(item)

	return item.Value.Value, true
}

func (c *LRU[V]) Set(key string, value V, size uint) []V {
	if existing, ok := c.items[key]; ok {
		c.size -= existing.Value.Size

		existing.Value.Key = key
		existing.Value.Value = value
		existing.Value.Size = size

		c.list.MoveToFront(existing)
		c.size += size

		return c.evictIfNeeded()
	}
	item := CacheItem[V]{
		Key:   key,
		Value: value,
		Size:  size,
	}

	listItem := c.list.PushFront(item)
	c.items[key] = listItem
	c.size += size

	return c.evictIfNeeded()
}

func (c *LRU[V]) Flush() []CacheItem[V] {
	items := make([]CacheItem[V], 0, c.list.Len())

	current := c.list.Front()
	for current != nil {
		items = append(items, current.Value)
		current = current.Next()
	}

	return items
}

func (c *LRU[V]) evictIfNeeded() []V {
	var evicted []V

	for c.size > c.maxSize {
		tail := c.list.Back()
		if tail == nil {
			break
		}

		evicted = append(evicted, tail.Value.Value)
		c.size -= tail.Value.Size
		delete(c.items, tail.Value.Key)
		c.list.Remove(tail)
	}

	return evicted
}

func (c *LRU[V]) RestoreData(items []CacheItem[V]) {
	c.items = make(map[string]*linkedlist.Item[CacheItem[V]])
	c.list = linkedlist.New[CacheItem[V]]()
	c.size = 0

	for _, item := range items {
		listItem := c.list.PushBack(item)
		c.items[item.Key] = listItem
		c.size += item.Size
	}
}
