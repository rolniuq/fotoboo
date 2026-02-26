# Storage Refactoring Summary

## Overview

The photo storage system has been refactored to follow **Clean Architecture** and **Strategy Pattern** principles. Storage implementation is now abstracted behind an interface in the `pkg/storage` package, allowing you to easily switch between different storage backends (Local, MinIO, S3, etc.) without modifying repository code.

---

## Architecture Changes

### Before (Tightly Coupled)
```
Repository Layer
├── SQLitePhotoRepository (local filesystem logic)
└── MinioPhotoRepository (MinIO logic)
```

**Problems:**
- Duplicated database logic in both repositories
- Switching storage requires rewriting entire repository
- Cannot easily add new storage backends (S3, Azure Blob, etc.)

### After (Clean Separation)
```
pkg/storage (Strategy Pattern)
├── Storage interface
├── LocalStorage implementation
└── MinioStorage implementation

Repository Layer
└── PhotoRepository (uses Storage interface)
```

**Benefits:**
- Single repository implementation
- Storage backend is pluggable via dependency injection
- Easy to add new storage types (just implement Storage interface)
- Repository focuses on database logic only

---

## File Structure

### New Files

**`pkg/storage/storage.go`**
- Defines `Storage` interface with methods: `Save()`, `Get()`, `Delete()`, `Exists()`, `Size()`, `ListKeys()`, `TotalSize()`, `Close()`
- All storage backends implement this interface

**`pkg/storage/local.go`**
- `LocalStorage` implementation using filesystem
- Stores files in configurable base path
- Same behavior as old `SQLitePhotoRepository`

**`pkg/storage/minio.go`**
- `MinioStorage` implementation using MinIO SDK
- Compatible with MinIO and AWS S3
- Manages bucket creation and object storage

**`internal/repository/photo_repository.go`**
- Unified `PhotoRepository` that accepts any `Storage` implementation
- Handles database operations (metadata)
- Delegates file operations to storage backend

### Modified Files

**`cmd/api/main.go`**
- Creates storage backend based on `USE_MINIO` env var
- Injects storage into `PhotoRepository`
- Clean separation of concerns

### Deleted Files

- ~~`internal/repository/sqlite_photo_repository.go`~~ (merged into `photo_repository.go`)
- ~~`internal/repository/minio_photo_repository.go`~~ (logic moved to `pkg/storage/minio.go`)
- ~~`internal/domain/minio_config.go`~~ (moved to `pkg/storage/minio.go`)

---

## How to Switch Storage Backends

### Local Storage (Default)
```bash
# No configuration needed
./bin/fotoboo-api
```

### MinIO Storage
```bash
export USE_MINIO=true
export MINIO_ENDPOINT=localhost:9000
export MINIO_ACCESS_KEY=minioadmin
export MINIO_SECRET_KEY=minioadmin
export MINIO_BUCKET=fotoboo
./bin/fotoboo-api
```

### AWS S3 (Future)
Just set MinIO config to S3:
```bash
export USE_MINIO=true
export MINIO_ENDPOINT=s3.amazonaws.com
export MINIO_ACCESS_KEY=your-aws-key
export MINIO_SECRET_KEY=your-aws-secret
export MINIO_BUCKET=your-bucket
export MINIO_USE_SSL=true
./bin/fotoboo-api
```

---

## Adding New Storage Backends

To add a new storage backend (e.g., Azure Blob, Google Cloud Storage):

1. Create `pkg/storage/azure.go` (or `gcs.go`)
2. Implement the `Storage` interface:
   ```go
   type AzureStorage struct { /* ... */ }
   func (s *AzureStorage) Save(ctx, key, data) (string, error) { /* ... */ }
   func (s *AzureStorage) Get(ctx, key) ([]byte, error) { /* ... */ }
   // ... implement other methods
   ```
3. Update `cmd/api/main.go` to instantiate your storage
4. **No changes needed in repository layer!**

---

## Storage Interface Contract

```go
type Storage interface {
    Save(ctx context.Context, key string, data []byte) (string, error)
    Get(ctx context.Context, key string) ([]byte, error)
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    Size(ctx context.Context, key string) (int64, error)
    ListKeys(ctx context.Context, prefix string) ([]string, error)
    TotalSize(ctx context.Context, prefix string) (int64, error)
    Close() error
}
```

### Key Design Decisions

1. **Keys are relative paths**: `photos/{id}.jpg`
   - Local storage: converted to absolute filesystem path
   - MinIO: used directly as object key
   - Database stores only the key (not full path)

2. **Context-aware**: All methods accept `context.Context` for:
   - Cancellation support
   - Timeout handling
   - Request tracing (future)

3. **Error handling**: All errors are wrapped with context using `fmt.Errorf("...: %w", err)`

---

## Benefits of This Refactoring

### 1. **Single Responsibility Principle**
- Repository handles database logic
- Storage handles file/object storage
- Clean separation

### 2. **Open/Closed Principle**
- Open for extension (add new storage backends)
- Closed for modification (no changes to repository)

### 3. **Dependency Inversion Principle**
- Repository depends on `Storage` interface (abstraction)
- Not on concrete implementations (LocalStorage, MinioStorage)

### 4. **Testability**
- Easy to mock `Storage` interface in tests
- Can test repository without real storage
- Can test storage implementations independently

### 5. **Maintainability**
- Storage logic is isolated in `pkg/storage`
- Repository code is simpler and focused
- Easy to understand and modify

---

## Testing

All existing tests continue to work:
```bash
go test ./internal/repository/... -v
```

Tests have been updated to:
- Use `storage.NewLocalStorage()` instead of direct repository
- Inject storage into `repository.NewPhotoRepository()`
- Verify behavior through the abstraction

---

## Migration Guide

If you're upgrading from the old code:

1. **No data migration needed** - file paths in database remain unchanged
2. **Environment variables** - same as before
3. **API behavior** - completely backward compatible
4. **Tests** - all passing with updated implementation

---

## Future Enhancements

With this architecture, you can easily add:

- **Cloud Storage**: Azure Blob, Google Cloud Storage, DigitalOcean Spaces
- **Caching Layer**: Redis + storage backend
- **CDN Integration**: CloudFront, Fastly
- **Compression**: Automatic image compression on save
- **Encryption**: At-rest encryption for sensitive photos
- **Multi-backend**: Primary + backup storage simultaneously

All without changing the repository layer!
