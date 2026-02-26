package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/internal/usecase"
)

// MockDeviceRepository implements domain.DeviceRepository for testing
type MockDeviceRepository struct {
	devices map[string]*domain.Device
}

func NewMockDeviceRepository() *MockDeviceRepository {
	return &MockDeviceRepository{
		devices: make(map[string]*domain.Device),
	}
}

func (m *MockDeviceRepository) Save(device *domain.Device) error {
	m.devices[device.ID] = device
	return nil
}

func (m *MockDeviceRepository) FindByID(id string) (*domain.Device, error) {
	if d, ok := m.devices[id]; ok {
		return d, nil
	}
	return nil, domain.ErrDeviceNotFound
}

func (m *MockDeviceRepository) Update(device *domain.Device) error {
	if _, ok := m.devices[device.ID]; !ok {
		return domain.ErrDeviceNotFound
	}
	m.devices[device.ID] = device
	return nil
}

func (m *MockDeviceRepository) ListAll() ([]*domain.Device, error) {
	result := make([]*domain.Device, 0, len(m.devices))
	for _, d := range m.devices {
		result = append(result, d)
	}
	return result, nil
}

func (m *MockDeviceRepository) Delete(id string) error {
	if _, ok := m.devices[id]; !ok {
		return domain.ErrDeviceNotFound
	}
	delete(m.devices, id)
	return nil
}

// DeviceUseCaseTestSuite uses testify suite
type DeviceUseCaseTestSuite struct {
	suite.Suite
	repo *MockDeviceRepository
	uc   *usecase.DeviceUseCase
}

func (s *DeviceUseCaseTestSuite) SetupTest() {
	s.repo = NewMockDeviceRepository()
	s.uc = usecase.NewDeviceUseCase(s.repo)
}

func (s *DeviceUseCaseTestSuite) TestRegisterDevice_Success() {
	device, err := s.uc.RegisterDevice("Booth Camera 1", "webcam")

	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), device.ID)
	assert.Equal(s.T(), "Booth Camera 1", device.Name)
	assert.Equal(s.T(), "webcam", device.Type)
	assert.True(s.T(), device.Active)
}

func (s *DeviceUseCaseTestSuite) TestRegisterDevice_EmptyName() {
	_, err := s.uc.RegisterDevice("", "webcam")

	assert.ErrorIs(s.T(), err, domain.ErrInvalidDevice)
}

func (s *DeviceUseCaseTestSuite) TestGetDevice_Success() {
	created, _ := s.uc.RegisterDevice("Test", "dslr")

	got, err := s.uc.GetDevice(created.ID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), created.ID, got.ID)
}

func (s *DeviceUseCaseTestSuite) TestGetDevice_NotFound() {
	_, err := s.uc.GetDevice("nonexistent")

	assert.ErrorIs(s.T(), err, domain.ErrDeviceNotFound)
}

func (s *DeviceUseCaseTestSuite) TestUpdateDevice_Success() {
	device, _ := s.uc.RegisterDevice("Original", "webcam")

	updated, err := s.uc.UpdateDevice(device.ID, "Updated Name", "dslr", false)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated Name", updated.Name)
	assert.Equal(s.T(), "dslr", updated.Type)
	assert.False(s.T(), updated.Active)
}

func (s *DeviceUseCaseTestSuite) TestUpdateDevice_NotFound() {
	_, err := s.uc.UpdateDevice("nonexistent", "Name", "type", true)

	assert.ErrorIs(s.T(), err, domain.ErrDeviceNotFound)
}

func (s *DeviceUseCaseTestSuite) TestListDevices() {
	s.uc.RegisterDevice("D1", "webcam")
	s.uc.RegisterDevice("D2", "dslr")

	devices, err := s.uc.ListDevices()

	require.NoError(s.T(), err)
	assert.Len(s.T(), devices, 2)
}

func (s *DeviceUseCaseTestSuite) TestDeleteDevice_Success() {
	device, _ := s.uc.RegisterDevice("ToDelete", "webcam")

	err := s.uc.DeleteDevice(device.ID)

	require.NoError(s.T(), err)

	_, err = s.uc.GetDevice(device.ID)
	assert.ErrorIs(s.T(), err, domain.ErrDeviceNotFound)
}

func (s *DeviceUseCaseTestSuite) TestDeleteDevice_NotFound() {
	err := s.uc.DeleteDevice("nonexistent")

	assert.ErrorIs(s.T(), err, domain.ErrDeviceNotFound)
}

func TestDeviceUseCaseSuite(t *testing.T) {
	suite.Run(t, new(DeviceUseCaseTestSuite))
}
