package api

import (
	"image_cache/internal/img"
	"image_cache/internal/lru"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

type Api struct {
	imgFolder string
	cachePath string
	timeout   time.Duration
	download  func(originURL string, folder string, fileName string, timeout time.Duration) (img.DownloadResult, error)
	cache     *lru.LRU[img.DownloadResult]
	logger    *slog.Logger
	mu        sync.Mutex
}

func New() *Api {
	return NewWithConfig(
		"app-data/images",
		"app-data/cache/data.json",
		5*time.Second,
		50*1024*1024,
	)
}

func NewWithConfig(imgFolder string, cachePath string, timeout time.Duration, maxCacheSize uint) *Api {
	return &Api{
		imgFolder: imgFolder,
		cachePath: cachePath,
		timeout:   timeout,
		download:  img.Download,
		cache:     lru.New[img.DownloadResult](maxCacheSize),
		logger:    slog.New(slog.NewTextHandler(os.Stdout, nil)),
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
