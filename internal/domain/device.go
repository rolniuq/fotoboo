package domain

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // "webcam", "dslr", "phone"
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewDevice(name string, deviceType string) *Device {
	now := time.Now()
	return &Device{
		ID:        uuid.New().String(),
		Name:      name,
		Type:      deviceType,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type DeviceRepository interface {
	Save(device *Device) error
	FindByID(id string) (*Device, error)
	Update(device *Device) error
	ListAll() ([]*Device, error)
	Delete(id string) error
}
