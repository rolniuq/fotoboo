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

// MockSessionRepository implements domain.SessionRepository for testing
type MockSessionRepository struct {
	sessions map[string]*domain.Session
}

func NewMockSessionRepository() *MockSessionRepository {
	return &MockSessionRepository{
		sessions: make(map[string]*domain.Session),
	}
}

func (m *MockSessionRepository) Save(session *domain.Session) error {
	m.sessions[session.ID] = session
	return nil
}

func (m *MockSessionRepository) FindByID(id string) (*domain.Session, error) {
	if s, ok := m.sessions[id]; ok {
		return s, nil
	}
	return nil, domain.ErrSessionNotFound
}

func (m *MockSessionRepository) Update(session *domain.Session) error {
	if _, ok := m.sessions[session.ID]; !ok {
		return domain.ErrSessionNotFound
	}
	m.sessions[session.ID] = session
	return nil
}

func (m *MockSessionRepository) ListAll() ([]*domain.Session, error) {
	result := make([]*domain.Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		result = append(result, s)
	}
	return result, nil
}

func (m *MockSessionRepository) CountAll() (int, error) {
	return len(m.sessions), nil
}

func (m *MockSessionRepository) CountByDate(date time.Time) (int, error) {
	count := 0
	dateStr := date.Format("2006-01-02")
	for _, s := range m.sessions {
		if s.CreatedAt.Format("2006-01-02") == dateStr {
			count++
		}
	}
	return count, nil
}

// SessionUseCaseTestSuite uses testify suite
type SessionUseCaseTestSuite struct {
	suite.Suite
	sessionRepo *MockSessionRepository
	photoRepo   *MockPhotoRepository
	uc          *usecase.SessionUseCase
}

func (s *SessionUseCaseTestSuite) SetupTest() {
	s.sessionRepo = NewMockSessionRepository()
	s.photoRepo = NewMockPhotoRepository()
	s.uc = usecase.NewSessionUseCase(s.sessionRepo, s.photoRepo)
}

func (s *SessionUseCaseTestSuite) TestStartSession_Success() {
	session, err := s.uc.StartSession("device-1")

	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), session.ID)
	assert.Equal(s.T(), "device-1", session.DeviceID)
	assert.Equal(s.T(), domain.SessionStatusActive, session.Status)
}

func (s *SessionUseCaseTestSuite) TestStartSession_NoDevice() {
	session, err := s.uc.StartSession("")

	require.NoError(s.T(), err)
	assert.Empty(s.T(), session.DeviceID)
}

func (s *SessionUseCaseTestSuite) TestGetSession_Success() {
	created, _ := s.uc.StartSession("")

	got, err := s.uc.GetSession(created.ID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), created.ID, got.ID)
}

func (s *SessionUseCaseTestSuite) TestGetSession_NotFound() {
	_, err := s.uc.GetSession("nonexistent")

	assert.ErrorIs(s.T(), err, domain.ErrSessionNotFound)
}

func (s *SessionUseCaseTestSuite) TestCompleteSession_Success() {
	session, _ := s.uc.StartSession("")
	assert.Equal(s.T(), domain.SessionStatusActive, session.Status)

	completed, err := s.uc.CompleteSession(session.ID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), domain.SessionStatusCompleted, completed.Status)
}

func (s *SessionUseCaseTestSuite) TestCompleteSession_NotFound() {
	_, err := s.uc.CompleteSession("nonexistent")

	assert.ErrorIs(s.T(), err, domain.ErrSessionNotFound)
}

func (s *SessionUseCaseTestSuite) TestListSessions() {
	s.uc.StartSession("d1")
	s.uc.StartSession("d2")
	s.uc.StartSession("d3")

	sessions, err := s.uc.ListSessions()

	require.NoError(s.T(), err)
	assert.Len(s.T(), sessions, 3)
}

func (s *SessionUseCaseTestSuite) TestCountSessions() {
	s.uc.StartSession("")
	s.uc.StartSession("")

	count, err := s.uc.CountSessions()

	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, count)
}

func (s *SessionUseCaseTestSuite) TestGetSessionPhotos_Success() {
	session, _ := s.uc.StartSession("")
	photoUC := usecase.NewPhotoUseCase(s.photoRepo)

	photoUC.UploadPhoto(session.ID, []byte("a"))
	photoUC.UploadPhoto(session.ID, []byte("b"))
	photoUC.UploadPhoto("other-session", []byte("c"))

	photos, err := s.uc.GetSessionPhotos(session.ID)

	require.NoError(s.T(), err)
	assert.Len(s.T(), photos, 2)
}

func (s *SessionUseCaseTestSuite) TestGetSessionPhotos_SessionNotFound() {
	_, err := s.uc.GetSessionPhotos("nonexistent")

	assert.ErrorIs(s.T(), err, domain.ErrSessionNotFound)
}

func TestSessionUseCaseSuite(t *testing.T) {
	suite.Run(t, new(SessionUseCaseTestSuite))
}
