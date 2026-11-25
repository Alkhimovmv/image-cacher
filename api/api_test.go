package api

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	domainimage "image_cache/internal/domain/image"
)

func newTestAPI(t *testing.T) (*Api, *fakeImage, *fakeCache) {
	t.Helper()

	image := &fakeImage{t: t}
	cache := &fakeCache{t: t}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	return New(Config{
		Image:  image,
		Cache:  cache,
		Logger: logger,
	}), image, cache
}

func TestNewUsesInjectedDependencies(t *testing.T) {
	instance, image, cache := newTestAPI(t)

	if instance.image != image {
		t.Fatal("ожидали переданную зависимость Image")
	}

	if instance.cache != cache {
		t.Fatal("ожидали переданную зависимость Cache")
	}
}

func TestStatusHandlerReturnsOK(t *testing.T) {
	instance, _, _ := newTestAPI(t)
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()

	instance.Status(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ожидали статус 200, получили %d", rec.Code)
	}

	if rec.Body.String() != "Ok" {
		t.Fatalf("ожидали тело Ok, получили %s", rec.Body.String())
	}
}

func TestHandlerRoutesStatusRequest(t *testing.T) {
	instance, _, _ := newTestAPI(t)
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()

	instance.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK || rec.Body.String() != "Ok" {
		t.Fatalf("неожиданный ответ: status=%d body=%q", rec.Code, rec.Body.String())
	}
}

func TestHandlerRoutesImageRequest(t *testing.T) {
	instance, image, cache := newTestAPI(t)
	result := domainimage.Image{
		Path:        "virtual/cached.jpeg",
		ContentType: "image/jpeg",
		Size:        12,
	}

	cache.getFunc = func(key string) (domainimage.Image, bool) {
		if key != "c52000873ea33e72fb89929688eee8cb" {
			t.Fatalf("неожиданный ключ кеша: %s", key)
		}
		return result, true
	}
	image.readFunc = func(path string) ([]byte, error) {
		if path != result.Path {
			t.Fatalf("неожиданный путь изображения: %s", path)
		}
		return []byte("cached-image"), nil
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin=https://example.com/cat.jpg", nil)
	rec := httptest.NewRecorder()
	instance.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK || rec.Body.String() != "cached-image" {
		t.Fatalf("неожиданный ответ: status=%d body=%q", rec.Code, rec.Body.String())
	}
}

func TestSaveCacheDelegatesToCache(t *testing.T) {
	expectedErr := errors.New("save failed")
	tests := []struct {
		name string
		err  error
	}{
		{name: "success"},
		{name: "error", err: expectedErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance, _, cache := newTestAPI(t)
			cache.saveFunc = func() error { return tt.err }

			err := instance.SaveCache()
			if !errors.Is(err, tt.err) {
				t.Fatalf("ожидали ошибку %v, получили %v", tt.err, err)
			}
			if cache.saveCalls != 1 {
				t.Fatalf("ожидали один вызов Cache.Save, получили %d", cache.saveCalls)
			}
		})
	}
}

func TestRestoreCacheDelegatesToCache(t *testing.T) {
	expectedErr := errors.New("restore failed")
	tests := []struct {
		name string
		err  error
	}{
		{name: "success"},
		{name: "error", err: expectedErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance, _, cache := newTestAPI(t)
			cache.restoreFunc = func() error { return tt.err }

			err := instance.RestoreCache()
			if !errors.Is(err, tt.err) {
				t.Fatalf("ожидали ошибку %v, получили %v", tt.err, err)
			}
			if cache.restoreCalls != 1 {
				t.Fatalf("ожидали один вызов Cache.Restore, получили %d", cache.restoreCalls)
			}
		})
	}
}
