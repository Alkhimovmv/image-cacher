package api

import (
	"log/slog"
	"net/http"
)

type Config struct {
	Image  Image
	Cache  Cache
	Logger *slog.Logger
}

type Api struct {
	image  Image
	cache  Cache
	logger *slog.Logger
}

func New(cfg Config) *Api {
	return &Api{
		image:  cfg.Image,
		cache:  cfg.Cache,
		logger: cfg.Logger,
	}
}

func (a *Api) Status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Ok"))
}

func (a *Api) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", a.Status)
	mux.HandleFunc("/image", a.Image)

	return mux
}
