package storage

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStore stores files in MinIO (S3-compatible object storage).
type MinioStore struct {
	client  *minio.Client
	bucket  string
	baseURL string
}

// NewMinioStore creates a MinIO-backed file store.
// endpoint: MinIO server address (e.g., "localhost:9000")
// accessKey, secretKey: MinIO credentials
// bucket: bucket name (created if not exists)
// baseURL: public base URL for accessing files (e.g., "https://cdn.example.com")
// useSSL: whether to use HTTPS for the MinIO connection
func NewMinioStore(endpoint string, accessKey string, secretKey string, bucket string, baseURL string, useSSL bool) (*MinioStore, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("minio bucket check: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio create bucket: %w", err)
		}
		log.Printf("storage: created minio bucket %s", bucket)
	}

	log.Printf("storage: minio connected to %s, bucket=%s", endpoint, bucket)
	return &MinioStore{client: client, bucket: bucket, baseURL: baseURL}, nil
}

func (s *MinioStore) Save(path string, r io.Reader, size int64) (string, error) {
	ctx := context.Background()
	_, err := s.client.PutObject(ctx, s.bucket, path, r, size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", fmt.Errorf("minio put: %w", err)
	}
	return path, nil
}

func (s *MinioStore) Delete(path string) error {
	ctx := context.Background()
	return s.client.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
}

func (s *MinioStore) BaseURL() string {
	return s.baseURL
}

func (s *MinioStore) Subdir() string {
	return s.bucket
}
