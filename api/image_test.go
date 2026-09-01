package api

import (
	"errors"
	"image_cache/internal/img"
	"image_cache/internal/lru"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestImageReturnsBadRequestWhenOriginIsMissing(t *testing.T) {
	api := New()

	req := httptest.NewRequest(http.MethodGet, "/image", nil)
	rec := httptest.NewRecorder()

	api.Image(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("ожидали статус 400, получили %d", resp.StatusCode)
	}
}

func TestImageReturnsBadRequestWhenOriginIsEmpty(t *testing.T) {
	api := New()

	req := httptest.NewRequest(http.MethodGet, "/image?origin=", nil)
	rec := httptest.NewRecorder()

	api.Image(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("ожидали статус 400, получили %d", resp.StatusCode)
	}
}

func TestImageCallsDownloadWhenOriginIsProvided(t *testing.T) {
	called := false

	api := New()
	api.download = func(originURL string, folder string, fileName string, timeout time.Duration) (img.DownloadResult, error) {
		called = true
		return img.DownloadResult{}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin=https://example.com/cat.jpg", nil)
	rec := httptest.NewRecorder()

	api.Image(rec, req)

	if !called {
		t.Fatal("ожидали, что обработчик вызовет скачивание изображения")
	}
}

func TestImageReturnsInternalServerErrorWhenDownloadFails(t *testing.T) {
	api := New()
	api.download = func(originURL string, folder string, fileName string, timeout time.Duration) (img.DownloadResult, error) {
		return img.DownloadResult{}, errors.New("download failed")
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin=https://example.com/cat.jpg", nil)
	rec := httptest.NewRecorder()

	api.Image(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("ожидали статус 500, получили %d", resp.StatusCode)
	}
}

func TestImageSetsContentTypeFromDownloadResult(t *testing.T) {
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

	api.Image(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ожидали статус 200, получили %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "image/jpeg" {
		t.Fatalf("ожидали Content-Type image/jpeg, получили %s", resp.Header.Get("Content-Type"))
	}
}

func TestImageSetsCacheHitHeaderToFalseForDownloadedImage(t *testing.T) {
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

	api.Image(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ожидали статус 200, получили %d", resp.StatusCode)
	}

	if resp.Header.Get("X-IMAGE-CACHE-HIT") != "false" {
		t.Fatalf("ожидали X-IMAGE-CACHE-HIT=false, получили %s", resp.Header.Get("X-IMAGE-CACHE-HIT"))
	}
}

func TestImageWritesDownloadedFileToResponse(t *testing.T) {
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

	api.Image(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ожидали статус 200, получили %d", resp.StatusCode)
	}

	if rec.Body.String() != "image-bytes" {
		t.Fatalf("ожидали тело image-bytes, получили %s", rec.Body.String())
	}
}

func TestImagePassesFolderAndTimeoutToDownload(t *testing.T) {
	expectedFolder := "app-data/images"
	expectedTimeout := 7 * time.Second

	api := New()
	api.imgFolder = expectedFolder
	api.timeout = expectedTimeout

	api.download = func(originURL string, folder string, fileName string, timeout time.Duration) (img.DownloadResult, error) {
		if folder != expectedFolder {
			t.Fatalf("ожидали folder %s, получили %s", expectedFolder, folder)
		}

		if timeout != expectedTimeout {
			t.Fatalf("ожидали timeout %v, получили %v", expectedTimeout, timeout)
		}

		return img.DownloadResult{
			Path:        "/tmp/test-file.jpeg",
			ContentType: "image/jpeg",
			Size:        10,
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin=https://example.com/cat.jpg", nil)
	rec := httptest.NewRecorder()

	api.Image(rec, req)
}

func TestImagePassesNonEmptyFileNameToDownload(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test-file.jpeg")

	err := os.WriteFile(filePath, []byte("image-bytes"), 0o644)
	if err != nil {
		t.Fatalf("не удалось подготовить тестовый файл: %v", err)
	}

	api := New()
	api.download = func(originURL string, folder string, fileName string, timeout time.Duration) (img.DownloadResult, error) {
		if fileName == "" {
			t.Fatal("ожидали, что обработчик передаст непустое имя файла")
		}

		if fileName == "temp-file" {
			t.Fatal("ожидали, что обработчик перестанет использовать захардкоженное имя temp-file")
		}

		return img.DownloadResult{
			Path:        filePath,
			ContentType: "image/jpeg",
			Size:        uint(len("image-bytes")),
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin=https://example.com/cat.jpg", nil)
	rec := httptest.NewRecorder()

	api.Image(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ожидали статус 200, получили %d", resp.StatusCode)
	}
}

func TestImagePassesMD5FileNameToDownload(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test-file.jpeg")

	err := os.WriteFile(filePath, []byte("image-bytes"), 0o644)
	if err != nil {
		t.Fatalf("не удалось подготовить тестовый файл: %v", err)
	}

	origin := "https://example.com/cat.jpg"
	expectedFileName := "c52000873ea33e72fb89929688eee8cb"

	api := New()
	api.download = func(originURL string, folder string, fileName string, timeout time.Duration) (img.DownloadResult, error) {
		if fileName != expectedFileName {
			t.Fatalf("ожидали имя файла %s, получили %s", expectedFileName, fileName)
		}

		return img.DownloadResult{
			Path:        filePath,
			ContentType: "image/jpeg",
			Size:        uint(len("image-bytes")),
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin="+origin, nil)
	rec := httptest.NewRecorder()

	api.Image(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ожидали статус 200, получили %d", resp.StatusCode)
	}
}

func TestImageDoesNotDownloadWhenFileIsAlreadyCached(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "cached-file.jpeg")

	err := os.WriteFile(filePath, []byte("cached-image"), 0o644)
	if err != nil {
		t.Fatalf("не удалось подготовить кешированный файл: %v", err)
	}

	api := New()
	api.download = func(originURL string, folder string, fileName string, timeout time.Duration) (img.DownloadResult, error) {
		t.Fatal("не ожидали вызов download при cache hit")
		return img.DownloadResult{}, nil
	}

	api.cache = lru.New[img.DownloadResult](50 * 1024 * 1024)
	api.cache.Set("c52000873ea33e72fb89929688eee8cb", img.DownloadResult{
		Path:        filePath,
		ContentType: "image/jpeg",
		Size:        uint(len("cached-image")),
	}, uint(len("cached-image")))

	req := httptest.NewRequest(http.MethodGet, "/image?origin=https://example.com/cat.jpg", nil)
	rec := httptest.NewRecorder()

	api.Image(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ожидали статус 200, получили %d", resp.StatusCode)
	}

	if resp.Header.Get("X-IMAGE-CACHE-HIT") != "true" {
		t.Fatalf("ожидали X-IMAGE-CACHE-HIT=true, получили %s", resp.Header.Get("X-IMAGE-CACHE-HIT"))
	}

	if rec.Body.String() != "cached-image" {
		t.Fatalf("ожидали тело cached-image, получили %s", rec.Body.String())
	}
}

func TestImageReturnsOriginStatusCodeWhenDownloadFailsWithStatusError(t *testing.T) {
	api := New()
	api.download = func(originURL string, folder string, fileName string, timeout time.Duration) (img.DownloadResult, error) {
		return img.DownloadResult{}, &img.StatusError{
			StatusCode: http.StatusNotFound,
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin=https://example.com/missing.jpg", nil)
	rec := httptest.NewRecorder()

	api.Image(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("ожидали статус 404, получили %d", resp.StatusCode)
	}
}

func TestImageReturnsAnyOriginStatusCodeWhenDownloadFailsWithStatusError(t *testing.T) {
	api := New()
	api.download = func(originURL string, folder string, fileName string, timeout time.Duration) (img.DownloadResult, error) {
		return img.DownloadResult{}, &img.StatusError{
			StatusCode: http.StatusForbidden,
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/image?origin=https://example.com/forbidden.jpg", nil)
	rec := httptest.NewRecorder()

	api.Image(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("ожидали статус 403, получили %d", resp.StatusCode)
	}
}

func TestImageRemovesEvictedFilesFromDisk(t *testing.T) {
	dir := t.TempDir()
	firstPath := filepath.Join(dir, "first.jpeg")
	secondPath := filepath.Join(dir, "second.jpeg")

	err := os.WriteFile(firstPath, []byte("first-image"), 0o644)
	if err != nil {
		t.Fatalf("не удалось подготовить первый файл: %v", err)
	}

	err = os.WriteFile(secondPath, []byte("second-image"), 0o644)
	if err != nil {
		t.Fatalf("не удалось подготовить второй файл: %v", err)
	}

	api := NewWithConfig(dir, filepath.Join(dir, "cache.json"), 5*time.Second, 15)

	call := 0
	api.download = func(originURL string, folder string, fileName string, timeout time.Duration) (img.DownloadResult, error) {
		call++

		if call == 1 {
			return img.DownloadResult{
				Path:        firstPath,
				ContentType: "image/jpeg",
				Size:        uint(len("first-image")),
			}, nil
		}

		return img.DownloadResult{
			Path:        secondPath,
			ContentType: "image/jpeg",
			Size:        uint(len("second-image")),
		}, nil
	}

	firstReq := httptest.NewRequest(http.MethodGet, "/image?origin=https://example.com/first.jpg", nil)
	firstRec := httptest.NewRecorder()
	api.Image(firstRec, firstReq)

	secondReq := httptest.NewRequest(http.MethodGet, "/image?origin=https://example.com/second.jpg", nil)
	secondRec := httptest.NewRecorder()
	api.Image(secondRec, secondReq)

	_, err = os.Stat(firstPath)
	if !os.IsNotExist(err) {
		t.Fatalf("ожидали, что первый файл будет удалён после вытеснения, получили ошибку %v", err)
	}
}
