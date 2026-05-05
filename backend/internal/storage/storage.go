package storage

import (
	"errors"
	"io"
)

var ErrNotFound = errors.New("file not found")

// FileStore abstracts file storage operations.
// Implementations: Local (disk), MinIO (object storage).
type FileStore interface {
	// Save stores a file and returns the relative path.
	Save(path string, r io.Reader, size int64) (string, error)

	// Delete removes a file at the given path.
	Delete(path string) error

	// BaseURL returns the public base URL for stored files.
	// E.g., "http://localhost:8081" for nginx, or MinIO endpoint.
	BaseURL() string

	// Subdir returns the storage subdirectory name ("raw", "hls", "thumbs").
	Subdir() string
}
