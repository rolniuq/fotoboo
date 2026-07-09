package handler_test

import (
	"encoding/json"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/fotoboo/fotoboo/internal/handler"
	"github.com/fotoboo/fotoboo/internal/usecase"
)

// PrintHandlerTestSuite tests PrintHandler
type PrintHandlerTestSuite struct {
	suite.Suite
	photoRepo *MockPhotoRepository
	uc        *usecase.PhotoUseCase
	h         *handler.PrintHandler
}

func (s *PrintHandlerTestSuite) SetupTest() {
	s.photoRepo = NewMockPhotoRepository()
	s.uc = usecase.NewPhotoUseCase(s.photoRepo)
	s.h = handler.NewPrintHandler(s.uc)
}

// Helper to create a valid JPEG for print tests
func validJPEG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	pr, pw := io.Pipe()
	go func() {
		_ = jpeg.Encode(pw, img, &jpeg.Options{Quality: 90})
		pw.Close()
	}()
	data, _ := io.ReadAll(pr)
	return data
}

func (s *PrintHandlerTestSuite) TestHandlePrintSizes_GET_Success() {
	req := httptest.NewRequest(http.MethodGet, "/print-sizes", nil)
	rr := httptest.NewRecorder()

	handler.HandlePrintSizes(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var sizes []map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&sizes)
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), sizes)

	keys := make([]string, len(sizes))
	for i, s := range sizes {
		keys[i] = s["key"].(string)
	}
	assert.Contains(s.T(), keys, "4x6")
	assert.Contains(s.T(), keys, "5x7")
	assert.Contains(s.T(), keys, "6x8")
	assert.Contains(s.T(), keys, "2x6")
}

func (s *PrintHandlerTestSuite) TestHandlePrintSizes_WrongMethod() {
	req := httptest.NewRequest(http.MethodPost, "/print-sizes", nil)
	rr := httptest.NewRecorder()

	handler.HandlePrintSizes(rr, req)

	assert.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
}

func (s *PrintHandlerTestSuite) TestHandlePrint_PhotoNotFound() {
	req := httptest.NewRequest(http.MethodGet, "/photos/nonexistent/print?size=4x6", nil)
	rr := httptest.NewRecorder()

	s.h.HandlePrint(rr, req, "nonexistent")

	assert.Equal(s.T(), http.StatusNotFound, rr.Code)
}

func (s *PrintHandlerTestSuite) TestHandlePrint_InvalidSize() {
	photo, _ := s.uc.UploadPhoto("", []byte("test"))
	req := httptest.NewRequest(http.MethodGet, "/photos/"+photo.ID+"/print?size=invalid", nil)
	rr := httptest.NewRecorder()

	s.h.HandlePrint(rr, req, photo.ID)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
}

func (s *PrintHandlerTestSuite) TestHandlePrint_DefaultsTo4x6() {
	photo, _ := s.uc.UploadPhoto("", validJPEG())
	req := httptest.NewRequest(http.MethodGet, "/photos/"+photo.ID+"/print", nil)
	rr := httptest.NewRecorder()

	s.h.HandlePrint(rr, req, photo.ID)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), "image/jpeg", rr.Header().Get("Content-Type"))
}

func (s *PrintHandlerTestSuite) TestHandlePrint_ReturnsAttachment() {
	photo, _ := s.uc.UploadPhoto("", validJPEG())
	req := httptest.NewRequest(http.MethodGet, "/photos/"+photo.ID+"/print?size=4x6", nil)
	rr := httptest.NewRecorder()

	s.h.HandlePrint(rr, req, photo.ID)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Contains(s.T(), rr.Header().Get("Content-Disposition"), "attachment")
}

func (s *PrintHandlerTestSuite) TestHandlePrint_WrongMethod() {
	req := httptest.NewRequest(http.MethodPost, "/photos/some-id/print", nil)
	rr := httptest.NewRecorder()

	s.h.HandlePrint(rr, req, "some-id")

	assert.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
}

func TestPrintHandlerSuite(t *testing.T) {
	suite.Run(t, new(PrintHandlerTestSuite))
}
