package api

import (
	"image_cache/internal/img"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestStatusHandlerReturnsOK(t *testing.T) {
	api := New()

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()

	api.Status(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ожидали статус 200, получили %d", resp.StatusCode)
	}

	if rec.Body.String() != "Ok" {
		t.Fatalf("ожидали тело Ok, получили %s", rec.Body.String())
	}
}

func TestSaveCacheWritesCacheStateToFile(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "cache.json")

	api := New()
	api.cachePath = cachePath

	api.cache.Set("test-key", img.DownloadResult{
		Path:        "app-data/images/test.jpeg",
		ContentType: "image/jpeg",
		Size:        123,
	}, 123)

	err := api.SaveCache()
	if err != nil {
		t.Fatalf("не ожидали ошибку при сохранении кеша, получили %v", err)
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatalf("не удалось прочитать файл кеша: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "test-key") {
		t.Fatalf("ожидали, что файл кеша содержит ключ test-key, получили %s", content)
	}

	if !strings.Contains(content, "app-data/images/test.jpeg") {
		t.Fatalf("ожидали, что файл кеша содержит путь к файлу, получили %s", content)
	}
}

func TestRestoreCacheLoadsCacheStateFromFile(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "cache.json")

	api := New()
	api.cachePath = cachePath

	data := `[{"Key":"test-key","Value":{"Path":"app-data/images/test.jpeg","ContentType":"image/jpeg","Size":123},"Size":123}]`

	err := os.WriteFile(cachePath, []byte(data), 0o644)
	if err != nil {
		t.Fatalf("не удалось подготовить файл кеша: %v", err)
	}

	err = api.RestoreCache()
	if err != nil {
		t.Fatalf("не ожидали ошибку при восстановлении кеша, получили %v", err)
	}

	value, ok := api.cache.Get("test-key")
	if !ok {
		t.Fatal("ожидали, что ключ test-key появится в кеше после восстановления")
	}

	if value.Path != "app-data/images/test.jpeg" {
		t.Fatalf("ожидали путь app-data/images/test.jpeg, получили %s", value.Path)
	}

	if value.ContentType != "image/jpeg" {
		t.Fatalf("ожидали Content-Type image/jpeg, получили %s", value.ContentType)
	}
}

func TestRestoreCacheReturnsNilWhenFileDoesNotExist(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "missing-cache.json")

	api := New()
	api.cachePath = cachePath

	err := api.RestoreCache()
	if err != nil {
		t.Fatalf("не ожидали ошибку, если файл кеша не существует, получили %v", err)
	}
}

func TestHandlerRoutesStatusRequest(t *testing.T) {
	api := New()

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()

	api.Handler().ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ожидали статус 200, получили %d", resp.StatusCode)
	}

	if rec.Body.String() != "Ok" {
		t.Fatalf("ожидали тело Ok, получили %s", rec.Body.String())
	}
}

func TestHandlerRoutesImageRequest(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test-file.jpeg")

	err := os.WriteFile(filePath, []byte("image-bytes"), 0o644)
	if err != nil {
		t.Fatalf("не удалось подготовить тестовый файл: %v", err)
	}

	api := New()
	api.download = func(originURL string, folder string, fileName string, timeout time.Duration) (img.DownloadResult, error) {
		return img.DownloadResult{
			Path:        filePath,
			ContentType: "image/jpeg",
			Size:        uint(len("image-bytes")),
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin=https://example.com/cat.jpg", nil)
	rec := httptest.NewRecorder()

	api.Handler().ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ожидали статус 200, получили %d", resp.StatusCode)
	}

	if rec.Body.String() != "image-bytes" {
		t.Fatalf("ожидали тело image-bytes, получили %s", rec.Body.String())
	}
}

func TestNewWithConfigUsesProvidedValues(t *testing.T) {
	api := NewWithConfig("custom-images", "custom-cache.json", 7*time.Second, 12345)

	if api.imgFolder != "custom-images" {
		t.Fatalf("ожидали imgFolder custom-images, получили %s", api.imgFolder)
	}

	if api.cachePath != "custom-cache.json" {
		t.Fatalf("ожидали cachePath custom-cache.json, получили %s", api.cachePath)
	}

	if api.timeout != 7*time.Second {
		t.Fatalf("ожидали timeout 7s, получили %v", api.timeout)
	}

	if api.cache == nil {
		t.Fatal("ожидали, что cache будет создан")
	}
}