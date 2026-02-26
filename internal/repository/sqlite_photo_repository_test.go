package repository_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/internal/repository"
)

// PhotoRepositoryTestSuite tests SQLitePhotoRepository
type PhotoRepositoryTestSuite struct {
	suite.Suite
	db          *sql.DB
	storagePath string
	tempDir     string
	repo        *repository.SQLitePhotoRepository
}

func (s *PhotoRepositoryTestSuite) SetupTest() {
	tempDir, err := os.MkdirTemp("", "fotoboo-test-*")
	require.NoError(s.T(), err)
	s.tempDir = tempDir

	dbPath := filepath.Join(tempDir, "test.db")
	s.storagePath = filepath.Join(tempDir, "photos")

	db, err := repository.InitDB(dbPath)
	require.NoError(s.T(), err)
	s.db = db

	s.repo = repository.NewSQLitePhotoRepository(db, s.storagePath)
}

func (s *PhotoRepositoryTestSuite) TearDownTest() {
	if s.db != nil {
		s.db.Close()
	}
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *PhotoRepositoryTestSuite) TestSave_Success() {
	photo := domain.NewPhoto("test-session", "")
	data := []byte("fake image data")

	err := s.repo.Save(photo, data)

	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), photo.FilePath)
	assert.FileExists(s.T(), photo.FilePath)
}

func (s *PhotoRepositoryTestSuite) TestSave_EmptyData() {
	photo := domain.NewPhoto("", "")

	err := s.repo.Save(photo, []byte{})

	assert.ErrorIs(s.T(), err, domain.ErrInvalidPhoto)
}

func (s *PhotoRepositoryTestSuite) TestFindByID_Success() {
	photo := domain.NewPhoto("", "")
	s.repo.Save(photo, []byte("test"))

	found, err := s.repo.FindByID(photo.ID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), photo.ID, found.ID)
}

func (s *PhotoRepositoryTestSuite) TestFindByID_NotFound() {
	_, err := s.repo.FindByID("nonexistent")

	assert.ErrorIs(s.T(), err, domain.ErrPhotoNotFound)
}

func (s *PhotoRepositoryTestSuite) TestFindBySessionID() {
	p1 := domain.NewPhoto("session-A", "")
	s.repo.Save(p1, []byte("a"))

	p2 := domain.NewPhoto("session-A", "")
	s.repo.Save(p2, []byte("b"))

	p3 := domain.NewPhoto("session-B", "")
	s.repo.Save(p3, []byte("c"))

	photos, err := s.repo.FindBySessionID("session-A")

	require.NoError(s.T(), err)
	assert.Len(s.T(), photos, 2)
}

func (s *PhotoRepositoryTestSuite) TestGetFileData_Success() {
	originalData := []byte("photo content here")
	photo := domain.NewPhoto("", "")
	s.repo.Save(photo, originalData)

	data, err := s.repo.GetFileData(photo)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), originalData, data)
}

func (s *PhotoRepositoryTestSuite) TestListAll() {
	for i := 0; i < 3; i++ {
		p := domain.NewPhoto("", "")
		s.repo.Save(p, []byte("test"))
	}

	photos, err := s.repo.ListAll()

	require.NoError(s.T(), err)
	assert.Len(s.T(), photos, 3)
}

func (s *PhotoRepositoryTestSuite) TestDelete_Success() {
	photo := domain.NewPhoto("", "")
	s.repo.Save(photo, []byte("test"))
	filePath := photo.FilePath

	err := s.repo.Delete(photo.ID)

	require.NoError(s.T(), err)

	_, err = s.repo.FindByID(photo.ID)
	assert.ErrorIs(s.T(), err, domain.ErrPhotoNotFound)

	_, err = os.Stat(filePath)
	assert.True(s.T(), os.IsNotExist(err), "File should be deleted")
}

func (s *PhotoRepositoryTestSuite) TestCountAll() {
	for i := 0; i < 5; i++ {
		p := domain.NewPhoto("", "")
		s.repo.Save(p, []byte("test"))
	}

	count, err := s.repo.CountAll()

	require.NoError(s.T(), err)
	assert.Equal(s.T(), 5, count)
}

func (s *PhotoRepositoryTestSuite) TestCountByDate() {
	for i := 0; i < 3; i++ {
		p := domain.NewPhoto("", "")
		s.repo.Save(p, []byte("test"))
	}

	count, err := s.repo.CountByDate(time.Now().UTC())

	require.NoError(s.T(), err)
	assert.Equal(s.T(), 3, count)
}

func (s *PhotoRepositoryTestSuite) TestTotalStorageBytes() {
	p1 := domain.NewPhoto("", "")
	s.repo.Save(p1, []byte("12345")) // 5 bytes

	p2 := domain.NewPhoto("", "")
	s.repo.Save(p2, []byte("1234567890")) // 10 bytes

	total, err := s.repo.TotalStorageBytes()

	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(15), total)
}

func TestPhotoRepositorySuite(t *testing.T) {
	suite.Run(t, new(PhotoRepositoryTestSuite))
}
