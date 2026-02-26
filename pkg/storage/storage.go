package storage

import "context"

// Storage is the interface that all storage backends must implement
type Storage interface {
	// Save stores data with the given key and returns the final path/key
	Save(ctx context.Context, key string, data []byte) (string, error)

	// Get retrieves data by key
	Get(ctx context.Context, key string) ([]byte, error)

	// Delete removes data by key
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists
	Exists(ctx context.Context, key string) (bool, error)

	// Size returns the size in bytes for a given key
	Size(ctx context.Context, key string) (int64, error)

	// ListKeys returns all keys with the given prefix
	ListKeys(ctx context.Context, prefix string) ([]string, error)

	// TotalSize returns the total size of all objects with the given prefix
	TotalSize(ctx context.Context, prefix string) (int64, error)

	// Close releases any resources held by the storage backend
	Close() error
}

// Config holds common configuration for all storage backends
type Config struct {
	Type string // "local", "minio", "s3"
}
