package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/fotoboo/fotoboo/internal/domain"
)

type FilePhotoRepository struct {
	basePath string
	photos   map[string]*domain.Photo
	mu       sync.RWMutex
}

func NewFilePhotoRepository(basePath string) *FilePhotoRepository {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		panic(err)
	}

	repo := &FilePhotoRepository{
		basePath: basePath,
		photos:   make(map[string]*domain.Photo),
	}

	repo.loadMetadata()
	return repo
}

func (r *FilePhotoRepository) metadataPath() string {
	return filepath.Join(r.basePath, "metadata.json")
}

func (r *FilePhotoRepository) loadMetadata() {
	data, err := os.ReadFile(r.metadataPath())
	if err != nil {
		return
	}

	var photos []*domain.Photo
	if err := json.Unmarshal(data, &photos); err != nil {
		return
	}

	for _, p := range photos {
		r.photos[p.ID] = p
	}
}

func (r *FilePhotoRepository) saveMetadata() error {
	photos := make([]*domain.Photo, 0, len(r.photos))
	for _, p := range r.photos {
		photos = append(photos, p)
	}

	data, err := json.MarshalIndent(photos, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.metadataPath(), data, 0644)
}

func (r *FilePhotoRepository) Save(photo *domain.Photo, data []byte) error {
	if len(data) == 0 {
		return domain.ErrInvalidPhoto
	}

	filePath := filepath.Join(r.basePath, photo.ID+".jpg")

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	photo.FilePath = filePath

	r.mu.Lock()
	r.photos[photo.ID] = photo
	err := r.saveMetadata()
	r.mu.Unlock()

	return err
}

func (r *FilePhotoRepository) FindByID(id string) (*domain.Photo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	photo, exists := r.photos[id]
	if !exists {
		return nil, domain.ErrPhotoNotFound
	}

	return photo, nil
}

func (r *FilePhotoRepository) GetFileData(photo *domain.Photo) ([]byte, error) {
	return os.ReadFile(photo.FilePath)
}
