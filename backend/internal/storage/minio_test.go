//go:build minio
// +build minio

package storage

import (
	"bytes"
	"os"
	"testing"
)

func TestMinioStoreIntegration(t *testing.T) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		t.Skip("MINIO_ENDPOINT not set")
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	store, err := NewMinioStore(endpoint, accessKey, secretKey, bucket, "", useSSL)
	if err != nil {
		t.Fatalf("NewMinioStore: %v", err)
	}

	// Test Save
	data := []byte("Hello, MinIO integration test!")
	path, err := store.Save("test/integration.txt", bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if path != "test/integration.txt" {
		t.Fatalf("expected path test/integration.txt, got %s", path)
	}

	// Test BaseURL
	if store.BaseURL() != "" {
		t.Logf("BaseURL: %s", store.BaseURL())
	}

	// Test Delete
	if err := store.Delete("test/integration.txt"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	// Test Delete non-existent (should not error)
	if err := store.Delete("test/nonexistent.txt"); err != nil {
		t.Logf("Delete non-existent: %v (expected)", err)
	}
}
