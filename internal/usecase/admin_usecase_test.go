package usecase_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/internal/usecase"
)

// AdminUseCaseTestSuite tests AdminUseCase
type AdminUseCaseTestSuite struct {
	suite.Suite
	photoRepo   *MockPhotoRepository
	sessionRepo *MockSessionRepository
	deviceRepo  *MockDeviceRepository
	configStore *domain.ConfigStore
	uc          *usecase.AdminUseCase
}

func (s *AdminUseCaseTestSuite) SetupTest() {
	s.photoRepo = NewMockPhotoRepository()
	s.sessionRepo = NewMockSessionRepository()
	s.deviceRepo = NewMockDeviceRepository()
	s.configStore = domain.NewConfigStore()
	s.uc = usecase.NewAdminUseCase(s.photoRepo, s.sessionRepo, s.deviceRepo, s.configStore)
}

func (s *AdminUseCaseTestSuite) TestGetStats_Empty() {
	stats, err := s.uc.GetStats()

	require.NoError(s.T(), err)
	assert.Equal(s.T(), 0, stats.TotalPhotos)
	assert.Equal(s.T(), 0, stats.TotalSessions)
	assert.Equal(s.T(), 0, stats.TotalDevices)
	assert.Equal(s.T(), int64(0), stats.StorageBytes)
}

func (s *AdminUseCaseTestSuite) TestGetStats_WithData() {
	// Add photos
	photoUC := usecase.NewPhotoUseCase(s.photoRepo)
	photoUC.UploadPhoto("s1", []byte("12345"))      // 5 bytes
	photoUC.UploadPhoto("s1", []byte("1234567890")) // 10 bytes

	// Add sessions
	sessionUC := usecase.NewSessionUseCase(s.sessionRepo, s.photoRepo)
	sessionUC.StartSession("d1")
	sessionUC.StartSession("d2")

	// Add devices
	deviceUC := usecase.NewDeviceUseCase(s.deviceRepo)
	deviceUC.RegisterDevice("Dev1", "webcam")

	stats, err := s.uc.GetStats()

	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, stats.TotalPhotos)
	assert.Equal(s.T(), 2, stats.TotalSessions)
	assert.Equal(s.T(), 1, stats.TotalDevices)
	assert.Equal(s.T(), int64(15), stats.StorageBytes)
	assert.Equal(s.T(), "15 B", stats.StorageFormatted)
}

func (s *AdminUseCaseTestSuite) TestGetStats_PhotosToday() {
	photoUC := usecase.NewPhotoUseCase(s.photoRepo)
	photoUC.UploadPhoto("", []byte("a"))
	photoUC.UploadPhoto("", []byte("b"))

	stats, err := s.uc.GetStats()

	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, stats.PhotosToday)
}

func (s *AdminUseCaseTestSuite) TestGetStats_SessionsToday() {
	sessionUC := usecase.NewSessionUseCase(s.sessionRepo, s.photoRepo)
	sessionUC.StartSession("")
	sessionUC.StartSession("")
	sessionUC.StartSession("")

	stats, err := s.uc.GetStats()

	require.NoError(s.T(), err)
	assert.Equal(s.T(), 3, stats.SessionsToday)
}

func (s *AdminUseCaseTestSuite) TestGetConfig_Default() {
	cfg := s.uc.GetConfig()

	assert.Equal(s.T(), "FotoBoo Event", cfg.EventName)
	assert.Equal(s.T(), 3, cfg.CountdownDuration)
	assert.Equal(s.T(), 10, cfg.MaxUploadSizeMB)
	assert.Equal(s.T(), 30, cfg.PhotoRetentionDay)
	assert.Contains(s.T(), cfg.AvailableFrames, "simple")
	assert.Contains(s.T(), cfg.AvailableFilters, "grayscale")
}

func (s *AdminUseCaseTestSuite) TestUpdateConfig() {
	newCfg := domain.Config{
		EventName:         "Custom Event",
		CountdownDuration: 5,
		AvailableFrames:   []string{"none", "custom"},
		AvailableFilters:  []string{"none"},
		MaxUploadSizeMB:   20,
		PhotoRetentionDay: 60,
	}

	s.uc.UpdateConfig(newCfg)

	cfg := s.uc.GetConfig()
	assert.Equal(s.T(), "Custom Event", cfg.EventName)
	assert.Equal(s.T(), 5, cfg.CountdownDuration)
	assert.Equal(s.T(), 20, cfg.MaxUploadSizeMB)
	assert.Equal(s.T(), 60, cfg.PhotoRetentionDay)
	assert.Equal(s.T(), []string{"none", "custom"}, cfg.AvailableFrames)
}

func (s *AdminUseCaseTestSuite) TestStorageFormatted_KB() {
	photoUC := usecase.NewPhotoUseCase(s.photoRepo)
	// Add ~2KB of data
	data := make([]byte, 2048)
	photoUC.UploadPhoto("", data)

	stats, _ := s.uc.GetStats()
	assert.Equal(s.T(), "2.00 KB", stats.StorageFormatted)
}

func TestAdminUseCaseSuite(t *testing.T) {
	suite.Run(t, new(AdminUseCaseTestSuite))
}

// MockSessionRepository with CountByDate for today
func (m *MockSessionRepository) countToday() int {
	count := 0
	today := time.Now().Format("2006-01-02")
	for _, s := range m.sessions {
		if s.CreatedAt.Format("2006-01-02") == today {
			count++
		}
	}
	return count
}
