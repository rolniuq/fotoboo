package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/internal/handler"
	"github.com/fotoboo/fotoboo/internal/usecase"
)

// MockPhotoRepository for handler tests
type MockPhotoRepository struct {
	photos map[string]*domain.Photo
	data   map[string][]byte
}

func NewMockPhotoRepository() *MockPhotoRepository {
	return &MockPhotoRepository{
		photos: make(map[string]*domain.Photo),
		data:   make(map[string][]byte),
	}
}

func (m *MockPhotoRepository) Save(photo *domain.Photo, data []byte) error {
	photo.FilePath = "/mock/" + photo.ID + ".jpg"
	m.photos[photo.ID] = photo
	m.data[photo.ID] = data
	return nil
}

func (m *MockPhotoRepository) FindByID(id string) (*domain.Photo, error) {
	if p, ok := m.photos[id]; ok {
		return p, nil
	}
	return nil, domain.ErrPhotoNotFound
}

func (m *MockPhotoRepository) FindBySessionID(sessionID string) ([]*domain.Photo, error) {
	var result []*domain.Photo
	for _, p := range m.photos {
		if p.SessionID == sessionID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockPhotoRepository) GetFileData(photo *domain.Photo) ([]byte, error) {
	if d, ok := m.data[photo.ID]; ok {
		return d, nil
	}
	return nil, domain.ErrPhotoNotFound
}

func (m *MockPhotoRepository) ListAll() ([]*domain.Photo, error) {
	result := make([]*domain.Photo, 0, len(m.photos))
	for _, p := range m.photos {
		result = append(result, p)
	}
	return result, nil
}

func (m *MockPhotoRepository) Delete(id string) error {
	delete(m.photos, id)
	delete(m.data, id)
	return nil
}

func (m *MockPhotoRepository) CountAll() (int, error)                  { return len(m.photos), nil }
func (m *MockPhotoRepository) CountByDate(date time.Time) (int, error) { return 0, nil }
func (m *MockPhotoRepository) TotalStorageBytes() (int64, error)       { return 0, nil }

// PhotoHandlerTestSuite uses testify suite
type PhotoHandlerTestSuite struct {
	suite.Suite
	repo *MockPhotoRepository
	uc   *usecase.PhotoUseCase
	h    *handler.PhotoHandler
}

func (s *PhotoHandlerTestSuite) SetupTest() {
	s.repo = NewMockPhotoRepository()
	s.uc = usecase.NewPhotoUseCase(s.repo)
	s.h = handler.NewPhotoHandler(s.uc)
}

func (s *PhotoHandlerTestSuite) TestUploadPhoto_Success() {
	body := bytes.NewReader([]byte("fake image data"))
	req := httptest.NewRequest(http.MethodPost, "/photos", body)
	req.Header.Set("Content-Type", "application/octet-stream")
	rr := httptest.NewRecorder()

	s.h.UploadPhoto(rr, req)

	assert.Equal(s.T(), http.StatusCreated, rr.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), resp["id"])
	assert.NotNil(s.T(), resp["created_at"])
}

func (s *PhotoHandlerTestSuite) TestUploadPhoto_EmptyBody() {
	req := httptest.NewRequest(http.MethodPost, "/photos", bytes.NewReader([]byte{}))
	rr := httptest.NewRecorder()

	s.h.UploadPhoto(rr, req)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
}

func (s *PhotoHandlerTestSuite) TestUploadPhoto_WrongMethod() {
	req := httptest.NewRequest(http.MethodGet, "/photos", nil)
	rr := httptest.NewRecorder()

	s.h.UploadPhoto(rr, req)

	assert.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
}

func (s *PhotoHandlerTestSuite) TestGetPhoto_Success() {
	photo, _ := s.uc.UploadPhoto("", []byte("test image bytes"))

	req := httptest.NewRequest(http.MethodGet, "/photos/"+photo.ID, nil)
	rr := httptest.NewRecorder()

	s.h.GetPhoto(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), "image/jpeg", rr.Header().Get("Content-Type"))

	body, _ := io.ReadAll(rr.Body)
	assert.Equal(s.T(), "test image bytes", string(body))
}

func (s *PhotoHandlerTestSuite) TestGetPhoto_NotFound() {
	req := httptest.NewRequest(http.MethodGet, "/photos/nonexistent", nil)
	rr := httptest.NewRecorder()

	s.h.GetPhoto(rr, req)

	assert.Equal(s.T(), http.StatusNotFound, rr.Code)
}

func (s *PhotoHandlerTestSuite) TestListPhotos_Success() {
	s.uc.UploadPhoto("", []byte("a"))
	s.uc.UploadPhoto("", []byte("b"))

	req := httptest.NewRequest(http.MethodGet, "/photos", nil)
	rr := httptest.NewRecorder()

	s.h.ListPhotos(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var photos []map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&photos)
	require.NoError(s.T(), err)
	assert.Len(s.T(), photos, 2)
}

func (s *PhotoHandlerTestSuite) TestDeletePhoto_Success() {
	photo, _ := s.uc.UploadPhoto("", []byte("test"))

	req := httptest.NewRequest(http.MethodDelete, "/photos/"+photo.ID, nil)
	rr := httptest.NewRecorder()

	s.h.DeletePhoto(rr, req)

	assert.Equal(s.T(), http.StatusNoContent, rr.Code)

	// Verify deleted
	_, err := s.uc.GetPhoto(photo.ID)
	assert.ErrorIs(s.T(), err, domain.ErrPhotoNotFound)
}

func TestPhotoHandlerSuite(t *testing.T) {
	suite.Run(t, new(PhotoHandlerTestSuite))
}
