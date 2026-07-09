package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/internal/handler"
	"github.com/fotoboo/fotoboo/internal/usecase"
)

// MockDeviceRepository for handler tests
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

// DeviceHandlerTestSuite tests DeviceHandler
type DeviceHandlerTestSuite struct {
	suite.Suite
	repo *MockDeviceRepository
	uc   *usecase.DeviceUseCase
	h    *handler.DeviceHandler
}

func (s *DeviceHandlerTestSuite) SetupTest() {
	s.repo = NewMockDeviceRepository()
	s.uc = usecase.NewDeviceUseCase(s.repo)
	s.h = handler.NewDeviceHandler(s.uc)
}

func (s *DeviceHandlerTestSuite) TestHandleDevices_POST_Success() {
	body := bytes.NewReader([]byte(`{"name":"Booth Camera 1","type":"webcam"}`))
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.h.HandleDevices(rr, req)

	assert.Equal(s.T(), http.StatusCreated, rr.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), resp["id"])
	assert.Equal(s.T(), "Booth Camera 1", resp["name"])
	assert.Equal(s.T(), "webcam", resp["type"])
	assert.True(s.T(), resp["active"].(bool))
}

func (s *DeviceHandlerTestSuite) TestHandleDevices_POST_EmptyName() {
	body := bytes.NewReader([]byte(`{"name":"","type":"webcam"}`))
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.h.HandleDevices(rr, req)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
}

func (s *DeviceHandlerTestSuite) TestHandleDevices_POST_DefaultsType() {
	body := bytes.NewReader([]byte(`{"name":"Default Cam"}`))
	req := httptest.NewRequest(http.MethodPost, "/devices", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.h.HandleDevices(rr, req)

	assert.Equal(s.T(), http.StatusCreated, rr.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "webcam", resp["type"])
}

func (s *DeviceHandlerTestSuite) TestHandleDevices_GET_Success() {
	s.uc.RegisterDevice("Cam 1", "webcam")
	s.uc.RegisterDevice("Cam 2", "dslr")

	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	rr := httptest.NewRecorder()

	s.h.HandleDevices(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var devices []map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&devices)
	require.NoError(s.T(), err)
	assert.Len(s.T(), devices, 2)
}

func (s *DeviceHandlerTestSuite) TestHandleDevices_WrongMethod() {
	req := httptest.NewRequest(http.MethodPut, "/devices", nil)
	rr := httptest.NewRecorder()

	s.h.HandleDevices(rr, req)

	assert.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
}

func (s *DeviceHandlerTestSuite) TestHandleDevice_GET_Success() {
	device, _ := s.uc.RegisterDevice("Test Cam", "webcam")

	req := httptest.NewRequest(http.MethodGet, "/devices/"+device.ID, nil)
	rr := httptest.NewRecorder()

	s.h.HandleDevice(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), device.ID, resp["id"])
	assert.Equal(s.T(), "Test Cam", resp["name"])
}

func (s *DeviceHandlerTestSuite) TestHandleDevice_GET_NotFound() {
	req := httptest.NewRequest(http.MethodGet, "/devices/nonexistent", nil)
	rr := httptest.NewRecorder()

	s.h.HandleDevice(rr, req)

	assert.Equal(s.T(), http.StatusNotFound, rr.Code)
}

func (s *DeviceHandlerTestSuite) TestHandleDevice_GET_EmptyID() {
	req := httptest.NewRequest(http.MethodGet, "/devices/", nil)
	rr := httptest.NewRecorder()

	s.h.HandleDevice(rr, req)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
}

func (s *DeviceHandlerTestSuite) TestHandleDevice_PUT_Success() {
	device, _ := s.uc.RegisterDevice("Original", "webcam")

	body := bytes.NewReader([]byte(`{"name":"Updated Cam","type":"dslr","active":false}`))
	req := httptest.NewRequest(http.MethodPut, "/devices/"+device.ID, body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.h.HandleDevice(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated Cam", resp["name"])
	assert.Equal(s.T(), "dslr", resp["type"])
	assert.False(s.T(), resp["active"].(bool))
}

func (s *DeviceHandlerTestSuite) TestHandleDevice_PUT_NotFound() {
	body := bytes.NewReader([]byte(`{"name":"New","type":"webcam","active":true}`))
	req := httptest.NewRequest(http.MethodPut, "/devices/nonexistent", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.h.HandleDevice(rr, req)

	assert.Equal(s.T(), http.StatusNotFound, rr.Code)
}

func (s *DeviceHandlerTestSuite) TestHandleDevice_PUT_InvalidBody() {
	body := bytes.NewReader([]byte(`invalid json`))
	req := httptest.NewRequest(http.MethodPut, "/devices/some-id", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.h.HandleDevice(rr, req)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
}

func (s *DeviceHandlerTestSuite) TestHandleDevice_DELETE_Success() {
	device, _ := s.uc.RegisterDevice("To Delete", "webcam")

	req := httptest.NewRequest(http.MethodDelete, "/devices/"+device.ID, nil)
	rr := httptest.NewRecorder()

	s.h.HandleDevice(rr, req)

	assert.Equal(s.T(), http.StatusNoContent, rr.Code)
}

func (s *DeviceHandlerTestSuite) TestHandleDevice_DELETE_NotFound() {
	req := httptest.NewRequest(http.MethodDelete, "/devices/nonexistent", nil)
	rr := httptest.NewRecorder()

	s.h.HandleDevice(rr, req)

	assert.Equal(s.T(), http.StatusNotFound, rr.Code)
}

func (s *DeviceHandlerTestSuite) TestHandleDevice_WrongMethod() {
	req := httptest.NewRequest(http.MethodPost, "/devices/some-id", nil)
	rr := httptest.NewRecorder()

	s.h.HandleDevice(rr, req)

	assert.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
}

func TestDeviceHandlerSuite(t *testing.T) {
	suite.Run(t, new(DeviceHandlerTestSuite))
}
