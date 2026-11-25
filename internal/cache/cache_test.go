package cache

import (
	"os"
	"path/filepath"
	"testing"

	domainimage "image_cache/internal/domain/image"
)

func TestCacheSaveAndRestore(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "cache.json")
	value := domainimage.Image{
		Path:        "images/cat.jpeg",
		ContentType: "image/jpeg",
		Size:        11,
	}

	first := New(path, 1024)
	first.Set("cat", value, value.Size)
	if err := first.Save(); err != nil {
		t.Fatalf("не удалось сохранить кеш: %v", err)
	}

	second := New(path, 1024)
	if err := second.Restore(); err != nil {
		t.Fatalf("не удалось восстановить кеш: %v", err)
	}

	got, ok := second.Get("cat")
	if !ok || got != value {
		t.Fatalf("неожиданное восстановленное значение: %#v, найдено=%v", got, ok)
	}
}

func TestCacheRestoreMissingFileReturnsNil(t *testing.T) {
	cache := New(filepath.Join(t.TempDir(), "missing.json"), 1024)

	if err := cache.Restore(); err != nil {
		t.Fatalf("не ожидали ошибку для отсутствующего файла: %v", err)
	}
}

func TestCacheRestoreInvalidJSONReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.json")
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}

	cache := New(path, 1024)
	if err := cache.Restore(); err == nil {
		t.Fatal("ожидали ошибку для некорректного JSON")
	}
}

func TestCacheSetReturnsEvictedImages(t *testing.T) {
	cache := New(filepath.Join(t.TempDir(), "cache.json"), 5)
	first := domainimage.Image{Path: "first", Size: 3}
	second := domainimage.Image{Path: "second", Size: 3}

	if evicted := cache.Set("first", first, first.Size); len(evicted) != 0 {
		t.Fatalf("не ожидали вытеснение первого элемента: %#v", evicted)
	}

	evicted := cache.Set("second", second, second.Size)
	if len(evicted) != 1 || evicted[0] != first {
		t.Fatalf("ожидали вытеснение первого элемента, получили %#v", evicted)
	}
}
