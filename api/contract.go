package api

import domainimage "image_cache/internal/domain/image"

type Image interface {
	Download(originURL string, fileName string) (domainimage.Image, error)
	Read(path string) ([]byte, error)
	Remove(path string) error
}

type Cache interface {
	Get(key string) (domainimage.Image, bool)
	Set(key string, value domainimage.Image, size uint) []domainimage.Image
	Save() error
	Restore() error
}
