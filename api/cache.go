package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"image_cache/internal/img"
	"image_cache/internal/lru"
	"os"
	"path/filepath"
)

func (a *Api) SaveCache() error {
	a.mu.Lock()
	items := a.cache.Flush()
	a.mu.Unlock()

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}

	dir := filepath.Dir(a.cachePath)
	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}

	return os.WriteFile(a.cachePath, data, 0o644)
}

func (a *Api) RestoreCache() error {
	data, err := os.ReadFile(a.cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	var items []lru.CacheItem[img.DownloadResult]
	err = json.Unmarshal(data, &items)
	if err != nil {
		return err
	}

	a.mu.Lock()
	a.cache.RestoreData(items)
	a.mu.Unlock()

	return nil
}

func fileNameFromOrigin(origin string) string {
	sum := md5.Sum([]byte(origin))
	return hex.EncodeToString(sum[:])
}

func (a *Api) removeEvictedFiles(items []img.DownloadResult) error {
	for _, item := range items {
		err := os.Remove(item.Path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}
