package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	domainimage "image_cache/internal/domain/image"
)

const (
	testOrigin = "https://example.com/cat.jpg"
	testKey    = "c52000873ea33e72fb89929688eee8cb"
)

func TestImageReturnsBadRequestWithoutOrigin(t *testing.T) {
	for _, target := range []string{"/image", "/image?origin="} {
		t.Run(target, func(t *testing.T) {
			instance, _, _ := newTestAPI(t)
			req := httptest.NewRequest(http.MethodGet, target, nil)
			rec := httptest.NewRecorder()

			instance.Image(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("ожидали статус 400, получили %d", rec.Code)
			}
		})
	}
}

func TestImageReturnsCachedImage(t *testing.T) {
	instance, image, cache := newTestAPI(t)
	result := domainimage.Image{Path: "virtual/cached.jpeg", ContentType: "image/jpeg", Size: 12}

	cache.getFunc = func(key string) (domainimage.Image, bool) {
		if key != testKey {
			t.Fatalf("ожидали ключ %s, получили %s", testKey, key)
		}
		return result, true
	}
	image.readFunc = func(path string) ([]byte, error) {
		if path != result.Path {
			t.Fatalf("ожидали путь %s, получили %s", result.Path, path)
		}
		return []byte("cached-image"), nil
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin="+testOrigin, nil)
	rec := httptest.NewRecorder()
	instance.Image(rec, req)

	if rec.Code != http.StatusOK || rec.Body.String() != "cached-image" {
		t.Fatalf("неожиданный ответ: status=%d body=%q", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("Content-Type") != "image/jpeg" {
		t.Fatalf("ожидали Content-Type image/jpeg, получили %q", rec.Header().Get("Content-Type"))
	}
	if rec.Header().Get("X-IMAGE-CACHE-HIT") != "true" {
		t.Fatalf("ожидали X-IMAGE-CACHE-HIT=true, получили %q", rec.Header().Get("X-IMAGE-CACHE-HIT"))
	}
}

func TestImageDownloadsAndReturnsCacheMiss(t *testing.T) {
	instance, image, cache := newTestAPI(t)
	result := domainimage.Image{Path: "virtual/downloaded.jpeg", ContentType: "image/jpeg", Size: 11}

	cache.getFunc = func(key string) (domainimage.Image, bool) {
		if key != testKey {
			t.Fatalf("ожидали ключ %s, получили %s", testKey, key)
		}
		return domainimage.Image{}, false
	}
	image.downloadFunc = func(originURL string, fileName string) (domainimage.Image, error) {
		if originURL != testOrigin || fileName != testKey {
			t.Fatalf("неожиданные параметры Download: origin=%s fileName=%s", originURL, fileName)
		}
		return result, nil
	}
	cache.setFunc = func(key string, value domainimage.Image, size uint) []domainimage.Image {
		if key != testKey || value != result || size != result.Size {
			t.Fatalf("неожиданные параметры Set: key=%s value=%#v size=%d", key, value, size)
		}
		return nil
	}
	image.readFunc = func(path string) ([]byte, error) {
		if path != result.Path {
			t.Fatalf("ожидали путь %s, получили %s", result.Path, path)
		}
		return []byte("image-bytes"), nil
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin="+testOrigin, nil)
	rec := httptest.NewRecorder()
	instance.Image(rec, req)

	if rec.Code != http.StatusOK || rec.Body.String() != "image-bytes" {
		t.Fatalf("неожиданный ответ: status=%d body=%q", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("Content-Type") != "image/jpeg" {
		t.Fatalf("ожидали Content-Type image/jpeg, получили %q", rec.Header().Get("Content-Type"))
	}
	if rec.Header().Get("X-IMAGE-CACHE-HIT") != "false" {
		t.Fatalf("ожидали X-IMAGE-CACHE-HIT=false, получили %q", rec.Header().Get("X-IMAGE-CACHE-HIT"))
	}
}

func TestImageReturnsInternalServerErrorWhenDownloadFails(t *testing.T) {
	instance, image, cache := newTestAPI(t)
	cache.getFunc = func(string) (domainimage.Image, bool) { return domainimage.Image{}, false }
	image.downloadFunc = func(string, string) (domainimage.Image, error) {
		return domainimage.Image{}, errors.New("download failed")
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin="+testOrigin, nil)
	rec := httptest.NewRecorder()
	instance.Image(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("ожидали статус 500, получили %d", rec.Code)
	}
}

func TestImageReturnsOriginStatusCode(t *testing.T) {
	for _, statusCode := range []int{http.StatusNotFound, http.StatusCreated} {
		t.Run(http.StatusText(statusCode), func(t *testing.T) {
			instance, image, cache := newTestAPI(t)
			cache.getFunc = func(string) (domainimage.Image, bool) { return domainimage.Image{}, false }
			image.downloadFunc = func(string, string) (domainimage.Image, error) {
				return domainimage.Image{}, &domainimage.StatusError{StatusCode: statusCode}
			}

			req := httptest.NewRequest(http.MethodGet, "/image?origin="+testOrigin, nil)
			rec := httptest.NewRecorder()
			instance.Image(rec, req)

			if rec.Code != statusCode {
				t.Fatalf("ожидали статус %d, получили %d", statusCode, rec.Code)
			}
		})
	}
}

func TestImageReturnsInternalServerErrorWhenReadFails(t *testing.T) {
	for _, cacheHit := range []bool{true, false} {
		t.Run(map[bool]string{true: "cached image", false: "downloaded image"}[cacheHit], func(t *testing.T) {
			instance, image, cache := newTestAPI(t)
			result := domainimage.Image{Path: "virtual/cat.jpeg", ContentType: "image/jpeg", Size: 5}
			cache.getFunc = func(string) (domainimage.Image, bool) { return result, cacheHit }
			if !cacheHit {
				image.downloadFunc = func(string, string) (domainimage.Image, error) { return result, nil }
				cache.setFunc = func(string, domainimage.Image, uint) []domainimage.Image { return nil }
			}
			image.readFunc = func(string) ([]byte, error) { return nil, errors.New("read failed") }

			req := httptest.NewRequest(http.MethodGet, "/image?origin="+testOrigin, nil)
			rec := httptest.NewRecorder()
			instance.Image(rec, req)

			if rec.Code != http.StatusInternalServerError {
				t.Fatalf("ожидали статус 500, получили %d", rec.Code)
			}
		})
	}
}

func TestImageRemovesEvictedImages(t *testing.T) {
	instance, image, cache := newTestAPI(t)
	result := domainimage.Image{Path: "virtual/new.jpeg", ContentType: "image/jpeg", Size: 3}
	evicted := []domainimage.Image{{Path: "virtual/old-1.jpeg"}, {Path: "virtual/old-2.jpeg"}}
	var calls []string

	cache.getFunc = func(string) (domainimage.Image, bool) { return domainimage.Image{}, false }
	image.downloadFunc = func(string, string) (domainimage.Image, error) { return result, nil }
	cache.setFunc = func(string, domainimage.Image, uint) []domainimage.Image { return evicted }
	image.removeFunc = func(path string) error {
		calls = append(calls, "remove:"+path)
		return nil
	}
	image.readFunc = func(path string) ([]byte, error) {
		calls = append(calls, "read:"+path)
		return []byte("new"), nil
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin="+testOrigin, nil)
	rec := httptest.NewRecorder()
	instance.Image(rec, req)

	if rec.Code != http.StatusOK || rec.Body.String() != "new" {
		t.Fatalf("неожиданный ответ: status=%d body=%q", rec.Code, rec.Body.String())
	}
	expectedCalls := []string{
		"remove:virtual/old-1.jpeg",
		"remove:virtual/old-2.jpeg",
		"read:virtual/new.jpeg",
	}
	if !reflect.DeepEqual(calls, expectedCalls) {
		t.Fatalf("неожиданный порядок вызовов: %#v", calls)
	}
}

func TestImageReturnsInternalServerErrorWhenEvictedImageRemovalFails(t *testing.T) {
	instance, image, cache := newTestAPI(t)
	result := domainimage.Image{Path: "virtual/new.jpeg", ContentType: "image/jpeg", Size: 3}
	evicted := []domainimage.Image{{Path: "virtual/old.jpeg"}}

	cache.getFunc = func(string) (domainimage.Image, bool) { return domainimage.Image{}, false }
	image.downloadFunc = func(string, string) (domainimage.Image, error) { return result, nil }
	cache.setFunc = func(string, domainimage.Image, uint) []domainimage.Image { return evicted }
	image.removeFunc = func(path string) error {
		if path != evicted[0].Path {
			t.Fatalf("ожидали путь %s, получили %s", evicted[0].Path, path)
		}
		return errors.New("remove failed")
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin="+testOrigin, nil)
	rec := httptest.NewRecorder()
	instance.Image(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("ожидали статус 500, получили %d", rec.Code)
	}
}

func TestFileNameFromOriginUsesMD5(t *testing.T) {
	if got := fileNameFromOrigin(testOrigin); got != testKey {
		t.Fatalf("ожидали ключ %s, получили %s", testKey, got)
	}
}
