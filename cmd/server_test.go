package main

import (
	"bytes"
	"image_cache/internal/config"
	"io"
	"log/slog"
	"net/http"
	"testing"
)

func TestNewAPIFromConfigReturnsHandler(t *testing.T) {
	cfg := config.Config{
		Img: config.ImgConfig{
			Folder:  "images",
			Timeout: 7,
		},
		Cache: config.CacheConfig{
			Path:    "cache.json",
			MaxSize: 1024,
		},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	instance := newAPIFromConfig(cfg, logger)

	if instance == nil || instance.Handler() == nil {
		t.Fatal("ожидали настроенный API")
	}
}

func TestNewHTTPServerUsesConfigAddr(t *testing.T) {
	cfg := config.Config{
		HTTP: config.HTTPConfig{
			Addr: ":3001",
		},
	}

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	// http.NewServeMux() это функция из стандартной библиотеки Go, которая создаёт роутер.
	handler := http.NewServeMux()

	server := newHTTPServer(cfg, logger, handler)

	if server == nil {
		t.Fatal("ожидали, что server будет создан")
	}

	if server.Addr != ":3001" {
		t.Fatalf("ожидали Addr :3001, получили %s", server.Addr)
	}
}
