package repository_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/internal/repository"
)

// SQLiteSessionRepositoryTestSuite tests SQLiteSessionRepository
type SQLiteSessionRepositoryTestSuite struct {
	suite.Suite
	repo *repository.SQLiteSessionRepository
	dbPath string
}

func (s *SQLiteSessionRepositoryTestSuite) SetupTest() {
	dir, _ := os.MkdirTemp("", "fotoboo-session-test-*")
	s.dbPath = filepath.Join(dir, "test.db")
	db, err := repository.InitDB(s.dbPath)
	require.NoError(s.T(), err)
	s.repo = repository.NewSQLiteSessionRepository(db)
}

func (s *SQLiteSessionRepositoryTestSuite) TearDownTest() {
	os.RemoveAll(filepath.Dir(s.dbPath))
}

func (s *SQLiteSessionRepositoryTestSuite) TestSave_Success() {
	session := domain.NewSession("device-1")
	err := s.repo.Save(session)

	require.NoError(s.T(), err)
}

func (s *SQLiteSessionRepositoryTestSuite) TestFindByID_Success() {
	session := domain.NewSession("device-1")
	s.repo.Save(session)

	found, err := s.repo.FindByID(session.ID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), session.ID, found.ID)
	assert.Equal(s.T(), session.DeviceID, found.DeviceID)
	assert.Equal(s.T(), domain.SessionStatusActive, found.Status)
}

func (s *SQLiteSessionRepositoryTestSuite) TestFindByID_NotFound() {
	_, err := s.repo.FindByID("nonexistent")

	assert.ErrorIs(s.T(), err, domain.ErrSessionNotFound)
}

func (s *SQLiteSessionRepositoryTestSuite) TestUpdate_Success() {
	session := domain.NewSession("device-1")
	s.repo.Save(session)

	session.Complete()
	err := s.repo.Update(session)

	require.NoError(s.T(), err)

	found, _ := s.repo.FindByID(session.ID)
	assert.Equal(s.T(), domain.SessionStatusCompleted, found.Status)
}

func (s *SQLiteSessionRepositoryTestSuite) TestListAll() {
	s1 := domain.NewSession("d1")
	s2 := domain.NewSession("d2")
	s3 := domain.NewSession("d3")

	s.repo.Save(s1)
	s.repo.Save(s2)
	s.repo.Save(s3)

	sessions, err := s.repo.ListAll()

	require.NoError(s.T(), err)
	assert.Len(s.T(), sessions, 3)
}

func (s *SQLiteSessionRepositoryTestSuite) TestCountAll() {
	s.repo.Save(domain.NewSession("d1"))
	s.repo.Save(domain.NewSession("d2"))

	count, err := s.repo.CountAll()

	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, count)
}

func (s *SQLiteSessionRepositoryTestSuite) TestCountByDate() {
	s.repo.Save(domain.NewSession("d1"))
	s.repo.Save(domain.NewSession("d2"))

	count, err := s.repo.CountByDate(time.Now())

	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, count)
}

func (s *SQLiteSessionRepositoryTestSuite) TestCountByDate_PreviousDate() {
	s.repo.Save(domain.NewSession("d1"))

	yesterday := time.Now().Add(-24 * time.Hour)
	count, err := s.repo.CountByDate(yesterday)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), 0, count)
}

func TestSQLiteSessionRepositorySuite(t *testing.T) {
	suite.Run(t, new(SQLiteSessionRepositoryTestSuite))
}
