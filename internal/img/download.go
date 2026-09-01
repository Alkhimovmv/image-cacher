package img

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DownloadResult struct {
	Path        string
	ContentType string
	Size        uint
}

type StatusError struct {
	StatusCode int
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("unexpected status code: %d", e.StatusCode)
}

func Download(originURL string, folder string, fileName string, timeout time.Duration) (DownloadResult, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(originURL)
	if err != nil {
		return DownloadResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return DownloadResult{}, &StatusError{
			StatusCode: resp.StatusCode,
		}
	}

	contentType := resp.Header.Get("Content-Type")
	parts := strings.Split(contentType, "/")
	if len(parts) != 2 {
		return DownloadResult{}, fmt.Errorf("invalid content-type: %s", contentType)
	}

	ext := parts[1]
	filePath := filepath.Join(folder, fileName+"."+ext)

	err = os.MkdirAll(folder, 0o755)
	if err != nil {
		return DownloadResult{}, err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return DownloadResult{}, err
	}
	defer file.Close()

	written, err := io.Copy(file, resp.Body)
	if err != nil {
		return DownloadResult{}, err
	}

	return DownloadResult{
		Path:        filePath,
		ContentType: contentType,
		Size:        uint(written),
	}, nil
}
