package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	domainimage "image_cache/internal/domain/image"
	"image_cache/internal/lru"
)

type Cache struct {
	path string
	lru  *lru.LRU[domainimage.Image]
	mu   sync.Mutex
}

func New(path string, maxSize uint) *Cache {
	return &Cache{
		path: path,
		lru:  lru.New[domainimage.Image](maxSize),
	}
}

func (c *Cache) Get(key string) (domainimage.Image, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.lru.Get(key)
}

func (c *Cache) Set(key string, value domainimage.Image, size uint) []domainimage.Image {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.lru.Set(key, value, size)
}

func (c *Cache) Save() error {
	c.mu.Lock()
	items := c.lru.Flush()
	c.mu.Unlock()

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(c.path, data, 0o644)
}

func (c *Cache) Restore() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	var items []lru.CacheItem[domainimage.Image]
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	c.mu.Lock()
	c.lru.RestoreData(items)
	c.mu.Unlock()

	return nil
}
