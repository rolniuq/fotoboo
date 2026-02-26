package handler_test

import (
	"bytes"
	"encoding/json"
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

// MockSessionRepository for handler tests
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

func (m *MockSessionRepository) CountAll() (int, error)                  { return len(m.sessions), nil }
func (m *MockSessionRepository) CountByDate(date time.Time) (int, error) { return 0, nil }

// SessionHandlerTestSuite uses testify suite
type SessionHandlerTestSuite struct {
	suite.Suite
	sessionRepo *MockSessionRepository
	photoRepo   *MockPhotoRepository
	uc          *usecase.SessionUseCase
	h           *handler.SessionHandler
}

func (s *SessionHandlerTestSuite) SetupTest() {
	s.sessionRepo = NewMockSessionRepository()
	s.photoRepo = NewMockPhotoRepository()
	s.uc = usecase.NewSessionUseCase(s.sessionRepo, s.photoRepo)
	s.h = handler.NewSessionHandler(s.uc)
}

func (s *SessionHandlerTestSuite) TestHandleSessions_POST_Success() {
	body := bytes.NewReader([]byte(`{"device_id":"dev-123"}`))
	req := httptest.NewRequest(http.MethodPost, "/sessions", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.h.HandleSessions(rr, req)

	assert.Equal(s.T(), http.StatusCreated, rr.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), resp["id"])
	assert.Equal(s.T(), "dev-123", resp["device_id"])
	assert.Equal(s.T(), "active", resp["status"])
}

func (s *SessionHandlerTestSuite) TestHandleSessions_GET_Success() {
	s.uc.StartSession("d1")
	s.uc.StartSession("d2")

	req := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	rr := httptest.NewRecorder()

	s.h.HandleSessions(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var sessions []map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&sessions)
	require.NoError(s.T(), err)
	assert.Len(s.T(), sessions, 2)
}

func (s *SessionHandlerTestSuite) TestHandleSession_GET_Success() {
	session, _ := s.uc.StartSession("")

	req := httptest.NewRequest(http.MethodGet, "/sessions/"+session.ID, nil)
	rr := httptest.NewRecorder()

	s.h.HandleSession(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), session.ID, resp["id"])
}

func (s *SessionHandlerTestSuite) TestHandleSession_Complete_Success() {
	session, _ := s.uc.StartSession("")

	req := httptest.NewRequest(http.MethodPost, "/sessions/"+session.ID+"/complete", nil)
	rr := httptest.NewRecorder()

	s.h.HandleSession(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "completed", resp["status"])
}

func (s *SessionHandlerTestSuite) TestHandleSession_NotFound() {
	req := httptest.NewRequest(http.MethodGet, "/sessions/nonexistent", nil)
	rr := httptest.NewRecorder()

	s.h.HandleSession(rr, req)

	assert.Equal(s.T(), http.StatusNotFound, rr.Code)
}

func (s *SessionHandlerTestSuite) TestHandleSession_Photos_Success() {
	session, _ := s.uc.StartSession("")
	photoUC := usecase.NewPhotoUseCase(s.photoRepo)
	photoUC.UploadPhoto(session.ID, []byte("a"))
	photoUC.UploadPhoto(session.ID, []byte("b"))

	req := httptest.NewRequest(http.MethodGet, "/sessions/"+session.ID+"/photos", nil)
	rr := httptest.NewRecorder()

	s.h.HandleSession(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var photos []map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&photos)
	require.NoError(s.T(), err)
	assert.Len(s.T(), photos, 2)
}

func TestSessionHandlerSuite(t *testing.T) {
	suite.Run(t, new(SessionHandlerTestSuite))
}
