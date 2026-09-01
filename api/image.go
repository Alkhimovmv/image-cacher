package api

import (
	"errors"
	"image_cache/internal/img"
	"net/http"
	"os"
)

func (a *Api) Image(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	if origin == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fileName := fileNameFromOrigin(origin)

	a.mu.Lock()
	cached, ok := a.cache.Get(fileName)
	a.mu.Unlock()

	if ok {
		a.writeCachedImage(w, r, cached)
		return
	}

	result, err := a.download(origin, a.imgFolder, fileName, a.timeout)
	if err != nil {
		a.handleDownloadError(w, r, err)
		return
	}

	a.mu.Lock()
	evicted := a.cache.Set(fileName, result, result.Size)
	a.mu.Unlock()

	err = a.removeEvictedFiles(evicted)
	if err != nil {
		a.logger.Error(
			"failed to remove evicted image",
			"method", r.Method,
			"uri", r.RequestURI,
			"host", r.Host,
			"error", err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.writeDownloadedImage(w, r, result)
}

func (a *Api) writeCachedImage(w http.ResponseWriter, r *http.Request, cached img.DownloadResult) {
	data, err := os.ReadFile(cached.Path)
	if err != nil {
		a.logger.Error(
			"failed to read cached image",
			"method", r.Method,
			"uri", r.RequestURI,
			"host", r.Host,
			"error", err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", cached.ContentType)
	w.Header().Set("X-IMAGE-CACHE-HIT", "true")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (a *Api) handleDownloadError(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Error(
		"failed to download image",
		"method", r.Method,
		"uri", r.RequestURI,
		"host", r.Host,
		"error", err,
	)

	var statusErr *img.StatusError
	if errors.As(err, &statusErr) {
		w.WriteHeader(statusErr.StatusCode)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func (a *Api) writeDownloadedImage(w http.ResponseWriter, r *http.Request, result img.DownloadResult) {
	data, err := os.ReadFile(result.Path)
	if err != nil {
		a.logger.Error(
			"failed to read downloaded image",
			"method", r.Method,
			"uri", r.RequestURI,
			"host", r.Host,
			"error", err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("X-IMAGE-CACHE-HIT", "false")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
