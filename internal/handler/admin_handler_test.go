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

// AdminHandlerTestSuite tests AdminHandler
type AdminHandlerTestSuite struct {
	suite.Suite
	photoRepo   *MockPhotoRepository
	sessionRepo *MockSessionRepository
	deviceRepo  *MockDeviceRepository
	configStore *domain.ConfigStore
	uc          *usecase.AdminUseCase
	h           *handler.AdminHandler
}

func (s *AdminHandlerTestSuite) SetupTest() {
	s.photoRepo = NewMockPhotoRepository()
	s.sessionRepo = NewMockSessionRepository()
	s.deviceRepo = NewMockDeviceRepository()
	s.configStore = domain.NewConfigStore()
	s.uc = usecase.NewAdminUseCase(s.photoRepo, s.sessionRepo, s.deviceRepo, s.configStore)
	s.h = handler.NewAdminHandler(s.uc)
}

func (s *AdminHandlerTestSuite) TestHandleStats_GET_Success() {
	req := httptest.NewRequest(http.MethodGet, "/admin/stats", nil)
	rr := httptest.NewRecorder()

	s.h.HandleStats(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var stats map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&stats)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), float64(0), stats["total_photos"])
	assert.Equal(s.T(), float64(0), stats["total_sessions"])
	assert.Equal(s.T(), float64(0), stats["total_devices"])
}

func (s *AdminHandlerTestSuite) TestHandleStats_GET_WrongMethod() {
	req := httptest.NewRequest(http.MethodPost, "/admin/stats", nil)
	rr := httptest.NewRecorder()

	s.h.HandleStats(rr, req)

	assert.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
}

func (s *AdminHandlerTestSuite) TestHandleConfig_GET_Success() {
	req := httptest.NewRequest(http.MethodGet, "/admin/config", nil)
	rr := httptest.NewRecorder()

	s.h.HandleConfig(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var cfg map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&cfg)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "FotoBoo Event", cfg["event_name"])
	assert.Equal(s.T(), float64(3), cfg["countdown_duration"])
}

func (s *AdminHandlerTestSuite) TestHandleConfig_PUT_Success() {
	newCfg := `{
		"event_name":"Custom Event",
		"countdown_duration":5,
		"available_frames":["none","simple"],
		"available_filters":["none","grayscale"],
		"max_upload_size_mb":20,
		"photo_retention_days":60
	}`
	body := bytes.NewReader([]byte(newCfg))
	req := httptest.NewRequest(http.MethodPut, "/admin/config", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.h.HandleConfig(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var cfg map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&cfg)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Custom Event", cfg["event_name"])
	assert.Equal(s.T(), float64(5), cfg["countdown_duration"])
}

func (s *AdminHandlerTestSuite) TestHandleConfig_PUT_InvalidCountdown() {
	body := bytes.NewReader([]byte(`{"countdown_duration":0}`))
	req := httptest.NewRequest(http.MethodPut, "/admin/config", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.h.HandleConfig(rr, req)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
}

func (s *AdminHandlerTestSuite) TestHandleConfig_PUT_InvalidUploadSize() {
	body := bytes.NewReader([]byte(`{"countdown_duration":3,"max_upload_size_mb":0}`))
	req := httptest.NewRequest(http.MethodPut, "/admin/config", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.h.HandleConfig(rr, req)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
}

func (s *AdminHandlerTestSuite) TestHandleConfig_PUT_InvalidRetention() {
	body := bytes.NewReader([]byte(`{"countdown_duration":3,"max_upload_size_mb":10,"photo_retention_days":0}`))
	req := httptest.NewRequest(http.MethodPut, "/admin/config", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.h.HandleConfig(rr, req)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
}

func (s *AdminHandlerTestSuite) TestHandleConfig_PUT_InvalidBody() {
	body := bytes.NewReader([]byte(`invalid`))
	req := httptest.NewRequest(http.MethodPut, "/admin/config", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.h.HandleConfig(rr, req)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
}

func (s *AdminHandlerTestSuite) TestHandleConfig_WrongMethod() {
	req := httptest.NewRequest(http.MethodDelete, "/admin/config", nil)
	rr := httptest.NewRecorder()

	s.h.HandleConfig(rr, req)

	assert.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
}

func TestAdminHandlerSuite(t *testing.T) {
	suite.Run(t, new(AdminHandlerTestSuite))
}
