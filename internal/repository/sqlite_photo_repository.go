package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/fotoboo/fotoboo/internal/domain"
)

type SQLitePhotoRepository struct {
	db       *sql.DB
	basePath string
}

func NewSQLitePhotoRepository(db *sql.DB, basePath string) *SQLitePhotoRepository {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		panic(fmt.Sprintf("failed to create storage directory: %v", err))
	}

	repo := &SQLitePhotoRepository{
		db:       db,
		basePath: basePath,
	}

	return repo
}

func (r *SQLitePhotoRepository) Save(photo *domain.Photo, data []byte) error {
	if len(data) == 0 {
		return domain.ErrInvalidPhoto
	}

	filePath := filepath.Join(r.basePath, photo.ID+".jpg")

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write photo file: %w", err)
	}

	photo.FilePath = filePath

	_, err := r.db.Exec(
		`INSERT INTO photos (id, session_id, file_path, created_at) VALUES (?, ?, ?, ?)`,
		photo.ID, photo.SessionID, photo.FilePath, photo.CreatedAt.UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		// Clean up the file if DB insert fails
		os.Remove(filePath)
		return fmt.Errorf("failed to save photo metadata: %w", err)
	}

	return nil
}

func (r *SQLitePhotoRepository) FindByID(id string) (*domain.Photo, error) {
	row := r.db.QueryRow(`SELECT id, session_id, file_path, created_at FROM photos WHERE id = ?`, id)
	return r.scanPhoto(row)
}

func (r *SQLitePhotoRepository) FindBySessionID(sessionID string) ([]*domain.Photo, error) {
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

func (r *SQLitePhotoRepository) GetFileData(photo *domain.Photo) ([]byte, error) {
	data, err := os.ReadFile(photo.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read photo file: %w", err)
	}
	return data, nil
}

func (r *SQLitePhotoRepository) ListAll() ([]*domain.Photo, error) {
	rows, err := r.db.Query(`SELECT id, session_id, file_path, created_at FROM photos ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list photos: %w", err)
	}
	defer rows.Close()

	return r.scanPhotos(rows)
}

func (r *SQLitePhotoRepository) Delete(id string) error {
	photo, err := r.FindByID(id)
	if err != nil {
		return err
	}

	// Delete the file
	if err := os.Remove(photo.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete photo file: %w", err)
	}

	// Delete from database
	_, err = r.db.Exec(`DELETE FROM photos WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete photo record: %w", err)
	}

	return nil
}

func (r *SQLitePhotoRepository) CountAll() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM photos`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count photos: %w", err)
	}
	return count, nil
}

func (r *SQLitePhotoRepository) CountByDate(date time.Time) (int, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, time.UTC).Format(time.RFC3339)

	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM photos WHERE created_at >= ? AND created_at <= ?`, startOfDay, endOfDay).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count photos by date: %w", err)
	}
	return count, nil
}

func (r *SQLitePhotoRepository) TotalStorageBytes() (int64, error) {
	rows, err := r.db.Query(`SELECT file_path FROM photos`)
	if err != nil {
		return 0, fmt.Errorf("failed to query photo paths: %w", err)
	}
	defer rows.Close()

	var totalSize int64
	for rows.Next() {
		var filePath string
		if err := rows.Scan(&filePath); err != nil {
			continue
		}
		info, err := os.Stat(filePath)
		if err != nil {
			continue
		}
		totalSize += info.Size()
	}

	return totalSize, nil
}

func (r *SQLitePhotoRepository) scanPhoto(row *sql.Row) (*domain.Photo, error) {
	var photo domain.Photo
	var createdAt string

	err := row.Scan(&photo.ID, &photo.SessionID, &photo.FilePath, &createdAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrPhotoNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan photo: %w", err)
	}

	photo.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
	return &photo, nil
}

func (r *SQLitePhotoRepository) scanPhotos(rows *sql.Rows) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	for rows.Next() {
		var photo domain.Photo
		var createdAt string

		err := rows.Scan(&photo.ID, &photo.SessionID, &photo.FilePath, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan photo row: %w", err)
		}

		photo.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
		photos = append(photos, &photo)
	}

	return photos, nil
}
