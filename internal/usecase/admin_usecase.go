package usecase

import (
	"fmt"
	"time"

	"github.com/fotoboo/fotoboo/internal/domain"
)

// AdminStats holds dashboard statistics
type AdminStats struct {
	TotalPhotos      int    `json:"total_photos"`
	TotalSessions    int    `json:"total_sessions"`
	TotalDevices     int    `json:"total_devices"`
	PhotosToday      int    `json:"photos_today"`
	SessionsToday    int    `json:"sessions_today"`
	StorageBytes     int64  `json:"storage_bytes"`
	StorageFormatted string `json:"storage_formatted"`
}

type AdminUseCase struct {
	photoRepo   domain.PhotoRepository
	sessionRepo domain.SessionRepository
	deviceRepo  domain.DeviceRepository
	config      *domain.ConfigStore
}

func NewAdminUseCase(
	photoRepo domain.PhotoRepository,
	sessionRepo domain.SessionRepository,
	deviceRepo domain.DeviceRepository,
	config *domain.ConfigStore,
) *AdminUseCase {
	return &AdminUseCase{
		photoRepo:   photoRepo,
		sessionRepo: sessionRepo,
		deviceRepo:  deviceRepo,
		config:      config,
	}
}

func (uc *AdminUseCase) GetStats() (*AdminStats, error) {
	today := time.Now().UTC()

	totalPhotos, err := uc.photoRepo.CountAll()
	if err != nil {
		return nil, err
	}

	totalSessions, err := uc.sessionRepo.CountAll()
	if err != nil {
		return nil, err
	}

	devices, err := uc.deviceRepo.ListAll()
	if err != nil {
		return nil, err
	}

	photosToday, err := uc.photoRepo.CountByDate(today)
	if err != nil {
		return nil, err
	}

	sessionsToday, err := uc.sessionRepo.CountByDate(today)
	if err != nil {
		return nil, err
	}

	storageBytes, err := uc.photoRepo.TotalStorageBytes()
	if err != nil {
		return nil, err
	}

	return &AdminStats{
		TotalPhotos:      totalPhotos,
		TotalSessions:    totalSessions,
		TotalDevices:     len(devices),
		PhotosToday:      photosToday,
		SessionsToday:    sessionsToday,
		StorageBytes:     storageBytes,
		StorageFormatted: formatBytes(storageBytes),
	}, nil
}

func (uc *AdminUseCase) GetConfig() domain.Config {
	return uc.config.Get()
}

func (uc *AdminUseCase) UpdateConfig(cfg domain.Config) {
	uc.config.Update(cfg)
}

func formatBytes(b int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case b >= GB:
		return fmt.Sprintf("%.2f GB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.2f MB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.2f KB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
