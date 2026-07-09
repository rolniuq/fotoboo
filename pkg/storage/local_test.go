package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// LocalStorageTestSuite tests LocalStorage
type LocalStorageTestSuite struct {
	suite.Suite
	storage *LocalStorage
	baseDir string
	ctx     context.Context
}

func (s *LocalStorageTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "fotoboo-storage-test-*")
	require.NoError(s.T(), err)
	s.baseDir = dir
	s.ctx = context.Background()
	s.storage, err = NewLocalStorage(dir)
	require.NoError(s.T(), err)
}

func (s *LocalStorageTestSuite) TearDownTest() {
	os.RemoveAll(s.baseDir)
}

func (s *LocalStorageTestSuite) TestSave_And_Get() {
	key := "photos/test.jpg"
	data := []byte("test image data")

	path, err := s.storage.Save(s.ctx, key, data)
	require.NoError(s.T(), err)
	assert.Contains(s.T(), path, key)

	retrieved, err := s.storage.Get(s.ctx, key)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), data, retrieved)
}

func (s *LocalStorageTestSuite) TestGet_NotFound() {
	_, err := s.storage.Get(s.ctx, "nonexistent.jpg")
	assert.Error(s.T(), err)
}

func (s *LocalStorageTestSuite) TestDelete() {
	key := "photos/to-delete.jpg"
	s.storage.Save(s.ctx, key, []byte("data"))

	err := s.storage.Delete(s.ctx, key)
	require.NoError(s.T(), err)

	_, err = s.storage.Get(s.ctx, key)
	assert.Error(s.T(), err)
}

func (s *LocalStorageTestSuite) TestDelete_NonExistent() {
	// Deleting a non-existent file should not error
	err := s.storage.Delete(s.ctx, "nonexistent.jpg")
	assert.NoError(s.T(), err)
}

func (s *LocalStorageTestSuite) TestExists() {
	s.storage.Save(s.ctx, "photos/exists.jpg", []byte("data"))

	exists, err := s.storage.Exists(s.ctx, "photos/exists.jpg")
	require.NoError(s.T(), err)
	assert.True(s.T(), exists)

	notExists, err := s.storage.Exists(s.ctx, "photos/nope.jpg")
	require.NoError(s.T(), err)
	assert.False(s.T(), notExists)
}

func (s *LocalStorageTestSuite) TestSize() {
	s.storage.Save(s.ctx, "test-size.jpg", []byte("12345"))

	size, err := s.storage.Size(s.ctx, "test-size.jpg")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(5), size)
}

func (s *LocalStorageTestSuite) TestSize_NotFound() {
	_, err := s.storage.Size(s.ctx, "nonexistent.jpg")
	assert.Error(s.T(), err)
}

func (s *LocalStorageTestSuite) TestListKeys() {
	s.storage.Save(s.ctx, "photos/a.jpg", []byte("1"))
	s.storage.Save(s.ctx, "photos/b.jpg", []byte("2"))
	s.storage.Save(s.ctx, "other/c.jpg", []byte("3"))

	photoKeys, err := s.storage.ListKeys(s.ctx, "photos/")
	require.NoError(s.T(), err)
	assert.Len(s.T(), photoKeys, 2)

	allKeys, err := s.storage.ListKeys(s.ctx, "")
	require.NoError(s.T(), err)
	assert.Len(s.T(), allKeys, 3)
}

func (s *LocalStorageTestSuite) TestListKeys_EmptyPrefix() {
	keys, err := s.storage.ListKeys(s.ctx, "nonexistent/")
	require.NoError(s.T(), err)
	assert.Empty(s.T(), keys)
}

func (s *LocalStorageTestSuite) TestTotalSize() {
	s.storage.Save(s.ctx, "photos/a.jpg", []byte("12345"))     // 5 bytes
	s.storage.Save(s.ctx, "photos/b.jpg", []byte("1234567890")) // 10 bytes

	total, err := s.storage.TotalSize(s.ctx, "photos/")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(15), total)
}

func (s *LocalStorageTestSuite) TestTotalSize_Empty() {
	total, err := s.storage.TotalSize(s.ctx, "empty/")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(0), total)
}

func (s *LocalStorageTestSuite) TestClose() {
	err := s.storage.Close()
	assert.NoError(s.T(), err)
}

func (s *LocalStorageTestSuite) TestNestedDirectories() {
	key := "a/b/c/d/e/test.jpg"
	data := []byte("nested test")

	_, err := s.storage.Save(s.ctx, key, data)
	require.NoError(s.T(), err)

	// Verify the file is at the right path
	info, err := os.Stat(filepath.Join(s.baseDir, key))
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(len(data)), info.Size())
}

func TestLocalStorageSuite(t *testing.T) {
	suite.Run(t, new(LocalStorageTestSuite))
}
