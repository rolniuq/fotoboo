package usecase

import (
	"github.com/fotoboo/fotoboo/internal/domain"
)

type DeviceUseCase struct {
	repo domain.DeviceRepository
}

func NewDeviceUseCase(repo domain.DeviceRepository) *DeviceUseCase {
	return &DeviceUseCase{repo: repo}
}

func (uc *DeviceUseCase) RegisterDevice(name string, deviceType string) (*domain.Device, error) {
	if name == "" {
		return nil, domain.ErrInvalidDevice
	}

	device := domain.NewDevice(name, deviceType)

	if err := uc.repo.Save(device); err != nil {
		return nil, err
	}

	return device, nil
}

func (uc *DeviceUseCase) GetDevice(id string) (*domain.Device, error) {
	return uc.repo.FindByID(id)
}

func (uc *DeviceUseCase) UpdateDevice(id string, name string, deviceType string, active bool) (*domain.Device, error) {
	device, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	device.Name = name
	device.Type = deviceType
	device.Active = active

	if err := uc.repo.Update(device); err != nil {
		return nil, err
	}

	return device, nil
}

func (uc *DeviceUseCase) ListDevices() ([]*domain.Device, error) {
	return uc.repo.ListAll()
}

func (uc *DeviceUseCase) DeleteDevice(id string) error {
	return uc.repo.Delete(id)
}
