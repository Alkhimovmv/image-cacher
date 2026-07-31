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

func Download(originURL string, folder string, fileName string, timeout time.Duration) (DownloadResult, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(originURL)
	if err != nil {
		return DownloadResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DownloadResult{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	parts := strings.Split(contentType, "/")
	if len(parts) != 2 {
		return DownloadResult{}, fmt.Errorf("invalid content-type: %s", contentType)
	}

	ext := parts[1]
	filePath := filepath.Join(folder, fileName+"."+ext)

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
