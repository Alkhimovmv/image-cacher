package main

import (
	"context"
	"image_cache/api"
	"image_cache/internal/cache"
	"image_cache/internal/config"
	"image_cache/internal/img"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func loadConfig() (config.Config, error) {
	return config.Load("config/config.yml")
}

func newAPIFromConfig(cfg config.Config, logger *slog.Logger) *api.Api {
	imageAdapter := img.New(
		cfg.Img.Folder,
		time.Duration(cfg.Img.Timeout)*time.Second,
	)
	cacheAdapter := cache.New(cfg.Cache.Path, cfg.Cache.MaxSize)

	return api.New(api.Config{
		Image:  imageAdapter,
		Cache:  cacheAdapter,
		Logger: logger,
	})
}

func newHTTPServer(cfg config.Config, logger *slog.Logger, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    cfg.HTTP.Addr,
		Handler: accessLogMiddleware(logger, handler),
	}
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func accessLogMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		logger.Info(
			"request completed",
			"method", r.Method,
			"uri", r.RequestURI,
			"status", recorder.statusCode,
			"duration", time.Since(start),
		)
	})
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	apiInstance := newAPIFromConfig(cfg, logger)

	err = apiInstance.RestoreCache()
	if err != nil {
		return err
	}

	server := newHTTPServer(cfg, logger, apiInstance.Handler())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case <-stop:
	case err := <-serverErr:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = apiInstance.SaveCache()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := run()
	if err != nil {
		logger.Error("service stopped with error", "error", err)
		os.Exit(1)
	}
}
