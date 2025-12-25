package domain

import (
	"time"

	"github.com/google/uuid"
)

type Photo struct {
	ID        string    `json:"id"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
}

func NewPhoto(filePath string) *Photo {
	return &Photo{
		ID:        uuid.New().String(),
		FilePath:  filePath,
		CreatedAt: time.Now(),
	}
}

type PhotoRepository interface {
	Save(photo *Photo, data []byte) error
	FindByID(id string) (*Photo, error)
	GetFileData(photo *Photo) ([]byte, error)
}
