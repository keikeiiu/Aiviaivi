package storage

import (
	"io"
	"log"
)

// MinioStore is a placeholder for future MinIO integration.
// To enable:
//   1. Add github.com/minio/minio-go/v7 dependency
//   2. Set STORAGE=minio and MINIO_ENDPOINT/MINIO_ACCESS_KEY/MINIO_SECRET_KEY env vars
//   3. Uncomment the implementation below
type MinioStore struct {
	endpoint string
	bucket   string
	baseURL  string
}

// NewMinioStore creates a MinIO-backed file store.
// Currently returns a stub — MinIO is not yet wired.
func NewMinioStore(endpoint string, bucket string, baseURL string) *MinioStore {
	log.Printf("storage: MinIO stub created (endpoint=%s, bucket=%s) — not yet functional", endpoint, bucket)
	return &MinioStore{endpoint: endpoint, bucket: bucket, baseURL: baseURL}
}

func (s *MinioStore) Save(path string, r io.Reader, size int64) (string, error) {
	// TODO: implement with minio-go client.PutObject
	return path, nil
}

func (s *MinioStore) Delete(path string) error {
	// TODO: implement with minio-go client.RemoveObject
	return nil
}

func (s *MinioStore) BaseURL() string {
	return s.baseURL
}

func (s *MinioStore) Subdir() string {
	return s.bucket
}
