package api

import (
	"errors"
	domainimage "image_cache/internal/domain/image"
	"net/http"
)

func (a *Api) Image(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	if origin == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fileName := fileNameFromOrigin(origin)

	cached, ok := a.cache.Get(fileName)

	if ok {
		a.writeImage(w, r, cached, true)
		return
	}

	result, err := a.image.Download(origin, fileName)
	if err != nil {
		a.handleDownloadError(w, r, err)
		return
	}

	evicted := a.cache.Set(fileName, result, result.Size)
	for _, item := range evicted {
		if err := a.image.Remove(item.Path); err != nil {
			a.logInternalError("failed to remove evicted image", r, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	a.writeImage(w, r, result, false)
}

func (a *Api) writeImage(w http.ResponseWriter, r *http.Request, image domainimage.Image, cacheHit bool) {
	data, err := a.image.Read(image.Path)
	if err != nil {
		a.logInternalError("failed to read image", r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", image.ContentType)
	if cacheHit {
		w.Header().Set("X-IMAGE-CACHE-HIT", "true")
	} else {
		w.Header().Set("X-IMAGE-CACHE-HIT", "false")
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (a *Api) handleDownloadError(w http.ResponseWriter, r *http.Request, err error) {
	a.logInternalError("failed to download image", r, err)

	var statusErr *domainimage.StatusError
	if errors.As(err, &statusErr) {
		w.WriteHeader(statusErr.StatusCode)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func (a *Api) logInternalError(message string, r *http.Request, err error) {
	a.logger.Error(
		message,
		"method", r.Method,
		"uri", r.RequestURI,
		"host", r.Host,
		"error", err,
	)
}
