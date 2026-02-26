package storage

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// LocalStorage implements Storage interface using local filesystem
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new local filesystem storage backend
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	// Ensure directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
	}, nil
}

func (s *LocalStorage) Save(ctx context.Context, key string, data []byte) (string, error) {
	filePath := filepath.Join(s.basePath, key)

	// Ensure parent directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

func (s *LocalStorage) Get(ctx context.Context, key string) ([]byte, error) {
	filePath := filepath.Join(s.basePath, key)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %w", err)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	filePath := filepath.Join(s.basePath, key)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	filePath := filepath.Join(s.basePath, key)
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}
	return true, nil
}

func (s *LocalStorage) Size(ctx context.Context, key string) (int64, error) {
	filePath := filepath.Join(s.basePath, key)
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, fmt.Errorf("file not found: %w", err)
		}
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}
	return info.Size(), nil
}

func (s *LocalStorage) ListKeys(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	searchPath := filepath.Join(s.basePath, prefix)

	err := filepath.WalkDir(searchPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Skip directories that don't exist
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}

		if !d.IsDir() {
			// Get relative path from basePath
			relPath, err := filepath.Rel(s.basePath, path)
			if err != nil {
				return err
			}
			keys = append(keys, relPath)
		}
		return nil
	})

	if err != nil {
		// If the prefix directory doesn't exist, return empty list
		if os.IsNotExist(err) {
			return keys, nil
		}
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return keys, nil
}

func (s *LocalStorage) TotalSize(ctx context.Context, prefix string) (int64, error) {
	var totalSize int64
	searchPath := filepath.Join(s.basePath, prefix)

	err := filepath.WalkDir(searchPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to calculate total size: %w", err)
	}

	return totalSize, nil
}

func (s *LocalStorage) Close() error {
	// No resources to release for local storage
	return nil
}
