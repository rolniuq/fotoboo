package usecase

import (
	"github.com/fotoboo/fotoboo/internal/domain"
)

type PhotoUseCase struct {
	repo domain.PhotoRepository
}

func NewPhotoUseCase(repo domain.PhotoRepository) *PhotoUseCase {
	return &PhotoUseCase{
		repo: repo,
	}
}

func (uc *PhotoUseCase) UploadPhoto(data []byte) (*domain.Photo, error) {
	if len(data) == 0 {
		return nil, domain.ErrInvalidPhoto
	}

	photo := domain.NewPhoto("")

	if err := uc.repo.Save(photo, data); err != nil {
		return nil, err
	}

	return photo, nil
}

func (uc *PhotoUseCase) GetPhoto(id string) (*domain.Photo, error) {
	return uc.repo.FindByID(id)
}

func (uc *PhotoUseCase) GetPhotoData(id string) (*domain.Photo, []byte, error) {
	photo, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}

	data, err := uc.repo.GetFileData(photo)
	if err != nil {
		return nil, nil, err
	}

	return photo, data, nil
}
