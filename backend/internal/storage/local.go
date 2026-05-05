package storage

import (
	"io"
	"os"
	"path/filepath"
)

// LocalStore stores files on the local filesystem.
type LocalStore struct {
	baseDir string
	baseURL string
}

func NewLocalStore(baseDir string, baseURL string) *LocalStore {
	return &LocalStore{baseDir: baseDir, baseURL: baseURL}
}

func (s *LocalStore) Save(path string, r io.Reader, _ int64) (string, error) {
	fullPath := filepath.Join(s.baseDir, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		os.Remove(fullPath)
		return "", err
	}

	return path, nil
}

func (s *LocalStore) Delete(path string) error {
	return os.Remove(filepath.Join(s.baseDir, path))
}

func (s *LocalStore) BaseURL() string {
	return s.baseURL
}

func (s *LocalStore) Subdir() string {
	return s.baseDir
}
