package repository_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/internal/repository"
)

// DeviceRepositoryTestSuite tests SQLiteDeviceRepository
type DeviceRepositoryTestSuite struct {
	suite.Suite
	db      *sql.DB
	tempDir string
	repo    *repository.SQLiteDeviceRepository
}

func (s *DeviceRepositoryTestSuite) SetupTest() {
	tempDir, err := os.MkdirTemp("", "fotoboo-device-test-*")
	require.NoError(s.T(), err)
	s.tempDir = tempDir

	dbPath := filepath.Join(tempDir, "test.db")

	db, err := repository.InitDB(dbPath)
	require.NoError(s.T(), err)
	s.db = db

	s.repo = repository.NewSQLiteDeviceRepository(db)
}

func (s *DeviceRepositoryTestSuite) TearDownTest() {
	if s.db != nil {
		s.db.Close()
	}
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *DeviceRepositoryTestSuite) TestSave_Success() {
	device := domain.NewDevice("Test Camera", "webcam")

	err := s.repo.Save(device)

	require.NoError(s.T(), err)

	// Verify it was saved
	found, err := s.repo.FindByID(device.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), device.ID, found.ID)
	assert.Equal(s.T(), "Test Camera", found.Name)
}

func (s *DeviceRepositoryTestSuite) TestFindByID_NotFound() {
	_, err := s.repo.FindByID("nonexistent")

	assert.ErrorIs(s.T(), err, domain.ErrDeviceNotFound)
}

func (s *DeviceRepositoryTestSuite) TestUpdate_Success() {
	device := domain.NewDevice("Original", "webcam")
	s.repo.Save(device)

	device.Name = "Updated"
	device.Type = "dslr"
	device.Active = false

	err := s.repo.Update(device)

	require.NoError(s.T(), err)

	found, _ := s.repo.FindByID(device.ID)
	assert.Equal(s.T(), "Updated", found.Name)
	assert.Equal(s.T(), "dslr", found.Type)
	assert.False(s.T(), found.Active)
}

func (s *DeviceRepositoryTestSuite) TestListAll() {
	s.repo.Save(domain.NewDevice("D1", "webcam"))
	s.repo.Save(domain.NewDevice("D2", "dslr"))
	s.repo.Save(domain.NewDevice("D3", "phone"))

	devices, err := s.repo.ListAll()

	require.NoError(s.T(), err)
	assert.Len(s.T(), devices, 3)
}

func (s *DeviceRepositoryTestSuite) TestDelete_Success() {
	device := domain.NewDevice("ToDelete", "webcam")
	s.repo.Save(device)

	err := s.repo.Delete(device.ID)

	require.NoError(s.T(), err)

	_, err = s.repo.FindByID(device.ID)
	assert.ErrorIs(s.T(), err, domain.ErrDeviceNotFound)
}

func (s *DeviceRepositoryTestSuite) TestDelete_NonexistentReturns404() {
	// This is the fix for API-018
	err := s.repo.Delete("nonexistent-device-id")

	assert.ErrorIs(s.T(), err, domain.ErrDeviceNotFound)
}

func TestDeviceRepositorySuite(t *testing.T) {
	suite.Run(t, new(DeviceRepositoryTestSuite))
}
