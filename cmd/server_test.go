package main

import (
	"bytes"
	"image_cache/internal/config"
	"log/slog"
	"net/http"
	"testing"
)

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
