package image

import "fmt"

type Image struct {
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
