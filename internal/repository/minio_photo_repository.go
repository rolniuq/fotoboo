package repository

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/fotoboo/fotoboo/internal/domain"
)

type MinioPhotoRepository struct {
	db          *sql.DB
	minioClient *minio.Client
	bucketName  string
}

func NewMinioPhotoRepository(db *sql.DB, config *domain.MinioConfig) (*MinioPhotoRepository, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	repo := &MinioPhotoRepository{
		db:          db,
		minioClient: minioClient,
		bucketName:  config.BucketName,
	}

	// Ensure bucket exists
	if err := repo.ensureBucket(); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	return repo, nil
}

func (r *MinioPhotoRepository) ensureBucket() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := r.minioClient.BucketExists(ctx, r.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = r.minioClient.MakeBucket(ctx, r.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return nil
}

func (r *MinioPhotoRepository) Save(photo *domain.Photo, data []byte) error {
	if len(data) == 0 {
		return domain.ErrInvalidPhoto
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	objectName := fmt.Sprintf("photos/%s.jpg", photo.ID)

	// Detect content type from data
	contentType := http.DetectContentType(data)

	// Upload to MinIO
	_, err := r.minioClient.PutObject(
		ctx,
		r.bucketName,
		objectName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upload photo to MinIO: %w", err)
	}

	// Store object name as file path in database
	photo.FilePath = objectName

	// Save metadata to database
	_, err = r.db.Exec(
		`INSERT INTO photos (id, session_id, file_path, created_at) VALUES (?, ?, ?, ?)`,
		photo.ID, photo.SessionID, photo.FilePath, photo.CreatedAt.UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		// Clean up the uploaded object if DB insert fails
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		r.minioClient.RemoveObject(cleanupCtx, r.bucketName, objectName, minio.RemoveObjectOptions{})
		return fmt.Errorf("failed to save photo metadata: %w", err)
	}

	return nil
}

func (r *MinioPhotoRepository) FindByID(id string) (*domain.Photo, error) {
	row := r.db.QueryRow(`SELECT id, session_id, file_path, created_at FROM photos WHERE id = ?`, id)
	return r.scanPhoto(row)
}

func (r *MinioPhotoRepository) FindBySessionID(sessionID string) ([]*domain.Photo, error) {
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

func (r *MinioPhotoRepository) GetFileData(photo *domain.Photo) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Download from MinIO
	object, err := r.minioClient.GetObject(ctx, r.bucketName, photo.FilePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from MinIO: %w", err)
	}
	defer object.Close()

	// Read all data
	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read photo data: %w", err)
	}

	return data, nil
}

func (r *MinioPhotoRepository) ListAll() ([]*domain.Photo, error) {
	rows, err := r.db.Query(`SELECT id, session_id, file_path, created_at FROM photos ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list photos: %w", err)
	}
	defer rows.Close()

	return r.scanPhotos(rows)
}

func (r *MinioPhotoRepository) Delete(id string) error {
	photo, err := r.FindByID(id)
	if err != nil {
		return err
	}

	// Delete from database first (can be retried if MinIO delete fails)
	_, err = r.db.Exec(`DELETE FROM photos WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete photo record: %w", err)
	}

	// Delete from MinIO (best-effort; orphaned objects can be cleaned up later)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = r.minioClient.RemoveObject(ctx, r.bucketName, photo.FilePath, minio.RemoveObjectOptions{})

	return nil
}

func (r *MinioPhotoRepository) CountAll() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM photos`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count photos: %w", err)
	}
	return count, nil
}

func (r *MinioPhotoRepository) CountByDate(date time.Time) (int, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, time.UTC).Format(time.RFC3339)

	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM photos WHERE created_at >= ? AND created_at <= ?`, startOfDay, endOfDay).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count photos by date: %w", err)
	}
	return count, nil
}

func (r *MinioPhotoRepository) TotalStorageBytes() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var totalSize int64
	objectCh := r.minioClient.ListObjects(ctx, r.bucketName, minio.ListObjectsOptions{
		Prefix:    "photos/",
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return 0, fmt.Errorf("failed to list objects for storage calculation: %w", object.Err)
		}
		totalSize += object.Size
	}

	return totalSize, nil
}

func (r *MinioPhotoRepository) scanPhoto(row *sql.Row) (*domain.Photo, error) {
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

func (r *MinioPhotoRepository) scanPhotos(rows *sql.Rows) ([]*domain.Photo, error) {
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
