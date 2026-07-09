package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/pkg/storage"
)

type PhotoRepository struct {
	db      *sql.DB
	storage storage.Storage
}

func NewPhotoRepository(db *sql.DB, storage storage.Storage) *PhotoRepository {
	return &PhotoRepository{
		db:      db,
		storage: storage,
	}
}

func (r *PhotoRepository) Save(photo *domain.Photo, data []byte) error {
	if len(data) == 0 {
		return domain.ErrInvalidPhoto
	}

	ctx := context.Background()
	objectKey := fmt.Sprintf("photos/%s.jpg", photo.ID)

	// Save to storage backend
	_, err := r.storage.Save(ctx, objectKey, data)
	if err != nil {
		return fmt.Errorf("failed to save photo to storage: %w", err)
	}

	// Store the object key (not the full path) for consistency across storage backends
	photo.FilePath = objectKey

	// Save metadata to database
	_, err = r.db.Exec(
		`INSERT INTO photos (id, session_id, file_path, created_at) VALUES (?, ?, ?, ?)`,
		photo.ID, photo.SessionID, photo.FilePath, photo.CreatedAt.UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		// Clean up the stored file if DB insert fails
		r.storage.Delete(ctx, objectKey)
		return fmt.Errorf("failed to save photo metadata: %w", err)
	}

	return nil
}

func (r *PhotoRepository) FindByID(id string) (*domain.Photo, error) {
	row := r.db.QueryRow(`SELECT id, session_id, file_path, created_at FROM photos WHERE id = ?`, id)
	return r.scanPhoto(row)
}

func (r *PhotoRepository) FindBySessionID(sessionID string) ([]*domain.Photo, error) {
	rows, err := r.db.Query(
		`SELECT id, session_id, file_path, created_at FROM photos WHERE session_id = ? ORDER BY created_at ASC`,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query photos by session: %w", err)
	}
	defer rows.Close()

	return r.scanPhotos(rows)
}

func (r *PhotoRepository) GetFileData(photo *domain.Photo) ([]byte, error) {
	ctx := context.Background()

	// Use the stored key directly (it's already just the object key)
	data, err := r.storage.Get(ctx, photo.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get photo data: %w", err)
	}

	return data, nil
}

func (r *PhotoRepository) ListAll() ([]*domain.Photo, error) {
	rows, err := r.db.Query(`SELECT id, session_id, file_path, created_at FROM photos ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list photos: %w", err)
	}
	defer rows.Close()

	return r.scanPhotos(rows)
}

func (r *PhotoRepository) Delete(id string) error {
	photo, err := r.FindByID(id)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Delete from storage using the stored key
	if err := r.storage.Delete(ctx, photo.FilePath); err != nil {
		return fmt.Errorf("failed to delete photo from storage: %w", err)
	}

	// Delete from database
	_, err = r.db.Exec(`DELETE FROM photos WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete photo record: %w", err)
	}

	return nil
}

func (r *PhotoRepository) CountAll() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM photos`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count photos: %w", err)
	}
	return count, nil
}

func (r *PhotoRepository) CountByDate(date time.Time) (int, error) {
	return countByDate(r.db, "photos", date)
}

func (r *PhotoRepository) TotalStorageBytes() (int64, error) {
	ctx := context.Background()
	totalSize, err := r.storage.TotalSize(ctx, "photos/")
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total storage: %w", err)
	}
	return totalSize, nil
}

func (r *PhotoRepository) scanPhoto(row *sql.Row) (*domain.Photo, error) {
	var photo domain.Photo
	var createdAt string

	err := row.Scan(&photo.ID, &photo.SessionID, &photo.FilePath, &createdAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrPhotoNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan photo: %w", err)
	}

	photo.CreatedAt = parseTime(createdAt)
	return &photo, nil
}

func (r *PhotoRepository) scanPhotos(rows *sql.Rows) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	for rows.Next() {
		var photo domain.Photo
		var createdAt string

		err := rows.Scan(&photo.ID, &photo.SessionID, &photo.FilePath, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan photo row: %w", err)
		}

		photo.CreatedAt = parseTime(createdAt)
		photos = append(photos, &photo)
	}

	return photos, nil
}
