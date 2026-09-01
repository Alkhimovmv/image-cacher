package img

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDownloadSavesImageToDisk(t *testing.T) {
	dir := t.TempDir()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test-image-data"))
	}))
	defer server.Close()

	result, err := Download(server.URL, dir, "test-file", 5*time.Second)
	if err != nil {
		t.Fatalf("ожидали успешное скачивание, получили ошибку %v", err)
	}

	if result.ContentType != "image/jpeg" {
		t.Fatalf("ожидали content-type image/jpeg, получили %s", result.ContentType)
	}

	if filepath.Ext(result.Path) != ".jpeg" {
		t.Fatalf("ожидали расширение .jpeg, получили %s", filepath.Ext(result.Path))
	}

	data, err := os.ReadFile(result.Path)
	if err != nil {
		t.Fatalf("ожидали, что файл сохранён на диск, получили ошибку чтения %v", err)
	}

	if string(data) != "test-image-data" {
		t.Fatalf("ожидали содержимое test-image-data, получили %s", string(data))
	}

	if filepath.Base(result.Path) != "test-file.jpeg" {
		t.Fatalf("ожидали имя файла test-file.jpeg, получили %s", filepath.Base(result.Path))
	}

	if result.Size != uint(len("test-image-data")) {
		t.Fatalf("ожидали размер %d, получили %d", len("test-image-data"), result.Size)
	}
}

func TestDownloadReturnsErrorForInvalidContentType(t *testing.T) {
	dir := t.TempDir()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "invalid")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test-image-data"))
	}))
	defer server.Close()

	_, err := Download(server.URL, dir, "test-file", 5*time.Second)
	if err == nil {
		t.Fatal("ожидали ошибку для некорректного content-type, получили nil")
	}
}

func TestDownloadReturnsStatusErrorForNonOKStatus(t *testing.T) {
	dir := t.TempDir()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not-found"))
	}))
	defer server.Close()

	_, err := Download(server.URL, dir, "test-file", 5*time.Second)
	if err == nil {
		t.Fatal("ожидали ошибку для ответа не 200 OK, получили nil")
	}

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("ожидали ошибку типа StatusError, получили %v", err)
	}

	if statusErr.StatusCode != http.StatusNotFound {
		t.Fatalf("ожидали статус 404 внутри ошибки, получили %d", statusErr.StatusCode)
	}
}

func TestDownloadAcceptsAny2xxStatus(t *testing.T) {
	dir := t.TempDir()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("created-image"))
	}))
	defer server.Close()

	result, err := Download(server.URL, dir, "created-file", 5*time.Second)
	if err != nil {
		t.Fatalf("ожидали успешное скачивание для 201 Created, получили ошибку %v", err)
	}

	data, err := os.ReadFile(result.Path)
	if err != nil {
		t.Fatalf("не удалось прочитать сохранённый файл: %v", err)
	}

	if string(data) != "created-image" {
		t.Fatalf("ожидали содержимое created-image, получили %s", string(data))
	}
}

func TestDownloadReturnsErrorOnTimeout(t *testing.T) {
	dir := t.TempDir()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("slow-image"))
	}))
	defer server.Close()

	_, err := Download(server.URL, dir, "slow-file", 10*time.Millisecond)
	if err == nil {
		t.Fatal("ожидали ошибку timeout, получили nil")
	}
}
