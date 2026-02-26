package domain

import (
	"time"

	"github.com/google/uuid"
)

type Photo struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
}

func NewPhoto(sessionID string, filePath string) *Photo {
	return &Photo{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		FilePath:  filePath,
		CreatedAt: time.Now(),
	}
}

type PhotoRepository interface {
	Save(photo *Photo, data []byte) error
	FindByID(id string) (*Photo, error)
	FindBySessionID(sessionID string) ([]*Photo, error)
	GetFileData(photo *Photo) ([]byte, error)
	ListAll() ([]*Photo, error)
	Delete(id string) error
	CountAll() (int, error)
	CountByDate(date time.Time) (int, error)
	TotalStorageBytes() (int64, error)
}
