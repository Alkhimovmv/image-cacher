package img

import (
	"fmt"
	domainimage "image_cache/internal/domain/image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Image struct {
	folder  string
	timeout time.Duration
}

func New(folder string, timeout time.Duration) *Image {
	return &Image{
		folder:  folder,
		timeout: timeout,
	}
}

func (i *Image) Download(originURL string, fileName string) (domainimage.Image, error) {
	return Download(originURL, i.folder, fileName, i.timeout)
}

func (i *Image) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (i *Image) Remove(path string) error {
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}

	return err
}

func Download(originURL string, folder string, fileName string, timeout time.Duration) (domainimage.Image, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(originURL)
	if err != nil {
		return domainimage.Image{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return domainimage.Image{}, &domainimage.StatusError{
			StatusCode: resp.StatusCode,
		}
	}

	contentType := resp.Header.Get("Content-Type")
	parts := strings.Split(contentType, "/")
	if len(parts) != 2 {
		return domainimage.Image{}, fmt.Errorf("invalid content-type: %s", contentType)
	}

	ext := parts[1]
	filePath := filepath.Join(folder, fileName+"."+ext)

	err = os.MkdirAll(folder, 0o755)
	if err != nil {
		return domainimage.Image{}, err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return domainimage.Image{}, err
	}
	defer file.Close()

	written, err := io.Copy(file, resp.Body)
	if err != nil {
		return domainimage.Image{}, err
	}

	return domainimage.Image{
		Path:        filePath,
		ContentType: contentType,
		Size:        uint(written),
	}, nil
}
