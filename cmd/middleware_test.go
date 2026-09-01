package main

import (
	"bytes"
	"image_cache/internal/config"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAccessLogMiddlewareLogsMethodURIAndStatus(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	handler := accessLogMiddleware(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	logOutput := buf.String()

	if !strings.Contains(logOutput, "method=GET") {
		t.Fatalf("ожидали, что лог содержит method=GET, получили %s", logOutput)
	}

	if !strings.Contains(logOutput, "uri=/status") {
		t.Fatalf("ожидали, что лог содержит uri=/status, получили %s", logOutput)
	}

	if !strings.Contains(logOutput, "status=201") {
		t.Fatalf("ожидали, что лог содержит status=201, получили %s", logOutput)
	}
}

func TestNewHTTPServerWrapsHandlerWithAccessLogMiddleware(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	cfg := config.Config{
		HTTP: config.HTTPConfig{
			Addr: ":3000",
		},
	}

	apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	server := newHTTPServer(cfg, logger, apiHandler)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	logOutput := buf.String()

	if !strings.Contains(logOutput, "method=GET") {
		t.Fatalf("ожидали, что лог содержит method=GET, получили %s", logOutput)
	}

	if !strings.Contains(logOutput, "uri=/status") {
		t.Fatalf("ожидали, что лог содержит uri=/status, получили %s", logOutput)
	}

	if !strings.Contains(logOutput, "status=204") {
		t.Fatalf("ожидали, что лог содержит status=204, получили %s", logOutput)
	}
}
