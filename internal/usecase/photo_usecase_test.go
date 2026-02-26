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

// MockPhotoRepository implements domain.PhotoRepository for testing
type MockPhotoRepository struct {
	photos  map[string]*domain.Photo
	data    map[string][]byte
	saveErr error
}

func NewMockPhotoRepository() *MockPhotoRepository {
	return &MockPhotoRepository{
		photos: make(map[string]*domain.Photo),
		data:   make(map[string][]byte),
	}
}

func (m *MockPhotoRepository) Save(photo *domain.Photo, data []byte) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	photo.FilePath = "/mock/" + photo.ID + ".jpg"
	m.photos[photo.ID] = photo
	m.data[photo.ID] = data
	return nil
}

func (m *MockPhotoRepository) FindByID(id string) (*domain.Photo, error) {
	if photo, ok := m.photos[id]; ok {
		return photo, nil
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
	if data, ok := m.data[photo.ID]; ok {
		return data, nil
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
	if _, ok := m.photos[id]; !ok {
		return domain.ErrPhotoNotFound
	}
	delete(m.photos, id)
	delete(m.data, id)
	return nil
}

func (m *MockPhotoRepository) CountAll() (int, error) {
	return len(m.photos), nil
}

func (m *MockPhotoRepository) CountByDate(date time.Time) (int, error) {
	count := 0
	dateStr := date.Format("2006-01-02")
	for _, p := range m.photos {
		if p.CreatedAt.Format("2006-01-02") == dateStr {
			count++
		}
	}
	return count, nil
}

func (m *MockPhotoRepository) TotalStorageBytes() (int64, error) {
	var total int64
	for _, d := range m.data {
		total += int64(len(d))
	}
	return total, nil
}

// PhotoUseCaseTestSuite uses testify suite
type PhotoUseCaseTestSuite struct {
	suite.Suite
	repo *MockPhotoRepository
	uc   *usecase.PhotoUseCase
}

func (s *PhotoUseCaseTestSuite) SetupTest() {
	s.repo = NewMockPhotoRepository()
	s.uc = usecase.NewPhotoUseCase(s.repo)
}

func (s *PhotoUseCaseTestSuite) TestUploadPhoto_Success() {
	data := []byte("fake image data")

	photo, err := s.uc.UploadPhoto("session-123", data)

	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), photo.ID)
	assert.Equal(s.T(), "session-123", photo.SessionID)
	assert.NotEmpty(s.T(), photo.FilePath)
}

func (s *PhotoUseCaseTestSuite) TestUploadPhoto_EmptyData() {
	_, err := s.uc.UploadPhoto("session-123", []byte{})

	assert.ErrorIs(s.T(), err, domain.ErrInvalidPhoto)
}

func (s *PhotoUseCaseTestSuite) TestUploadPhoto_NoSession() {
	photo, err := s.uc.UploadPhoto("", []byte("test"))

	require.NoError(s.T(), err)
	assert.Empty(s.T(), photo.SessionID)
}

func (s *PhotoUseCaseTestSuite) TestGetPhoto_Success() {
	photo, _ := s.uc.UploadPhoto("", []byte("test"))

	got, err := s.uc.GetPhoto(photo.ID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), photo.ID, got.ID)
}

func (s *PhotoUseCaseTestSuite) TestGetPhoto_NotFound() {
	_, err := s.uc.GetPhoto("nonexistent")

	assert.ErrorIs(s.T(), err, domain.ErrPhotoNotFound)
}

func (s *PhotoUseCaseTestSuite) TestGetPhotoData_Success() {
	originalData := []byte("image content here")
	photo, _ := s.uc.UploadPhoto("", originalData)

	got, data, err := s.uc.GetPhotoData(photo.ID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), photo.ID, got.ID)
	assert.Equal(s.T(), originalData, data)
}

func (s *PhotoUseCaseTestSuite) TestDeletePhoto_Success() {
	photo, _ := s.uc.UploadPhoto("", []byte("test"))

	err := s.uc.DeletePhoto(photo.ID)

	require.NoError(s.T(), err)

	_, err = s.uc.GetPhoto(photo.ID)
	assert.ErrorIs(s.T(), err, domain.ErrPhotoNotFound)
}

func (s *PhotoUseCaseTestSuite) TestListPhotos() {
	s.uc.UploadPhoto("s1", []byte("a"))
	s.uc.UploadPhoto("s2", []byte("b"))
	s.uc.UploadPhoto("s3", []byte("c"))

	photos, err := s.uc.ListPhotos()

	require.NoError(s.T(), err)
	assert.Len(s.T(), photos, 3)
}

func (s *PhotoUseCaseTestSuite) TestCountPhotos() {
	s.uc.UploadPhoto("", []byte("a"))
	s.uc.UploadPhoto("", []byte("b"))

	count, err := s.uc.CountPhotos()

	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, count)
}

func (s *PhotoUseCaseTestSuite) TestTotalStorageBytes() {
	s.uc.UploadPhoto("", []byte("12345"))      // 5 bytes
	s.uc.UploadPhoto("", []byte("1234567890")) // 10 bytes

	total, err := s.uc.TotalStorageBytes()

	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(15), total)
}

func (s *PhotoUseCaseTestSuite) TestGetPhotosBySession() {
	s.uc.UploadPhoto("session-A", []byte("a"))
	s.uc.UploadPhoto("session-A", []byte("b"))
	s.uc.UploadPhoto("session-B", []byte("c"))

	photos, err := s.uc.GetPhotosBySession("session-A")

	require.NoError(s.T(), err)
	assert.Len(s.T(), photos, 2)
}

func TestPhotoUseCaseSuite(t *testing.T) {
	suite.Run(t, new(PhotoUseCaseTestSuite))
}
