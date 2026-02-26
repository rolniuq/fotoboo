package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/fotoboo/fotoboo/internal/domain"
)

// PhotoTestSuite tests Photo entity
type PhotoTestSuite struct {
	suite.Suite
}

func (s *PhotoTestSuite) TestNewPhoto_GeneratesID() {
	photo := domain.NewPhoto("session-1", "")

	assert.NotEmpty(s.T(), photo.ID)
	assert.Len(s.T(), photo.ID, 36) // UUID format
}

func (s *PhotoTestSuite) TestNewPhoto_SetsSessionID() {
	photo := domain.NewPhoto("my-session", "")

	assert.Equal(s.T(), "my-session", photo.SessionID)
}

func (s *PhotoTestSuite) TestNewPhoto_SetsFilePath() {
	photo := domain.NewPhoto("", "/path/to/file.jpg")

	assert.Equal(s.T(), "/path/to/file.jpg", photo.FilePath)
}

func (s *PhotoTestSuite) TestNewPhoto_SetsCreatedAt() {
	photo := domain.NewPhoto("", "")

	assert.False(s.T(), photo.CreatedAt.IsZero())
}

func TestPhotoSuite(t *testing.T) {
	suite.Run(t, new(PhotoTestSuite))
}

// SessionTestSuite tests Session entity
type SessionTestSuite struct {
	suite.Suite
}

func (s *SessionTestSuite) TestNewSession_GeneratesID() {
	session := domain.NewSession("device-1")

	assert.NotEmpty(s.T(), session.ID)
	assert.Len(s.T(), session.ID, 36)
}

func (s *SessionTestSuite) TestNewSession_SetsDeviceID() {
	session := domain.NewSession("my-device")

	assert.Equal(s.T(), "my-device", session.DeviceID)
}

func (s *SessionTestSuite) TestNewSession_StartsActive() {
	session := domain.NewSession("")

	assert.Equal(s.T(), domain.SessionStatusActive, session.Status)
}

func (s *SessionTestSuite) TestNewSession_SetsTimestamps() {
	session := domain.NewSession("")

	assert.False(s.T(), session.CreatedAt.IsZero())
	assert.False(s.T(), session.UpdatedAt.IsZero())
}

func (s *SessionTestSuite) TestSession_Complete() {
	session := domain.NewSession("")
	assert.Equal(s.T(), domain.SessionStatusActive, session.Status)

	session.Complete()

	assert.Equal(s.T(), domain.SessionStatusCompleted, session.Status)
}

func TestSessionSuite(t *testing.T) {
	suite.Run(t, new(SessionTestSuite))
}

// DeviceTestSuite tests Device entity
type DeviceTestSuite struct {
	suite.Suite
}

func (s *DeviceTestSuite) TestNewDevice_GeneratesID() {
	device := domain.NewDevice("Camera 1", "webcam")

	assert.NotEmpty(s.T(), device.ID)
	assert.Len(s.T(), device.ID, 36)
}

func (s *DeviceTestSuite) TestNewDevice_SetsFields() {
	device := domain.NewDevice("My Camera", "dslr")

	assert.Equal(s.T(), "My Camera", device.Name)
	assert.Equal(s.T(), "dslr", device.Type)
}

func (s *DeviceTestSuite) TestNewDevice_ActiveByDefault() {
	device := domain.NewDevice("Test", "webcam")

	assert.True(s.T(), device.Active)
}

func (s *DeviceTestSuite) TestNewDevice_SetsTimestamps() {
	device := domain.NewDevice("Test", "webcam")

	assert.False(s.T(), device.CreatedAt.IsZero())
	assert.False(s.T(), device.UpdatedAt.IsZero())
}

func TestDeviceSuite(t *testing.T) {
	suite.Run(t, new(DeviceTestSuite))
}

// ConfigTestSuite tests Config and ConfigStore
type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestDefaultConfig() {
	cfg := domain.DefaultConfig()

	assert.Equal(s.T(), "FotoBoo Event", cfg.EventName)
	assert.Equal(s.T(), 3, cfg.CountdownDuration)
	assert.Equal(s.T(), 10, cfg.MaxUploadSizeMB)
	assert.Equal(s.T(), 30, cfg.PhotoRetentionDay)
	assert.NotEmpty(s.T(), cfg.AvailableFrames)
	assert.NotEmpty(s.T(), cfg.AvailableFilters)
}

func (s *ConfigTestSuite) TestConfigStore_Get() {
	store := domain.NewConfigStore()

	cfg := store.Get()

	assert.Equal(s.T(), "FotoBoo Event", cfg.EventName)
}

func (s *ConfigTestSuite) TestConfigStore_Update() {
	store := domain.NewConfigStore()

	newCfg := domain.Config{
		EventName:         "Updated Event",
		CountdownDuration: 5,
		MaxUploadSizeMB:   20,
		PhotoRetentionDay: 60,
		AvailableFrames:   []string{"custom"},
		AvailableFilters:  []string{"custom"},
	}

	store.Update(newCfg)

	cfg := store.Get()
	assert.Equal(s.T(), "Updated Event", cfg.EventName)
	assert.Equal(s.T(), 5, cfg.CountdownDuration)
}

func (s *ConfigTestSuite) TestConfigStore_GetReturnsCopy() {
	store := domain.NewConfigStore()

	cfg1 := store.Get()
	cfg1.EventName = "Modified"

	cfg2 := store.Get()
	assert.Equal(s.T(), "FotoBoo Event", cfg2.EventName, "Original should be unchanged")
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

// ErrorsTestSuite tests domain errors
type ErrorsTestSuite struct {
	suite.Suite
}

func (s *ErrorsTestSuite) TestErrors_AreDefined() {
	assert.NotNil(s.T(), domain.ErrPhotoNotFound)
	assert.NotNil(s.T(), domain.ErrInvalidPhoto)
	assert.NotNil(s.T(), domain.ErrSessionNotFound)
	assert.NotNil(s.T(), domain.ErrDeviceNotFound)
	assert.NotNil(s.T(), domain.ErrInvalidDevice)
}

func (s *ErrorsTestSuite) TestErrors_HaveMessage() {
	assert.Contains(s.T(), domain.ErrPhotoNotFound.Error(), "not found")
	assert.Contains(s.T(), domain.ErrSessionNotFound.Error(), "not found")
	assert.Contains(s.T(), domain.ErrDeviceNotFound.Error(), "not found")
}

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsTestSuite))
}
