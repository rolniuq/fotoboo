package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/fotoboo/fotoboo/internal/handler"
)

// QRHandlerTestSuite tests QR handler
type QRHandlerTestSuite struct {
	suite.Suite
}

func (s *QRHandlerTestSuite) TestGenerateQR_WithTextParam() {
	req := httptest.NewRequest(http.MethodGet, "/qr?text=https://example.com", nil)
	rr := httptest.NewRecorder()

	handler.GenerateQR(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), "image/png", rr.Header().Get("Content-Type"))
}

func (s *QRHandlerTestSuite) TestGenerateQR_WithPathText() {
	req := httptest.NewRequest(http.MethodGet, "/qr/https://example.com", nil)
	rr := httptest.NewRecorder()

	handler.GenerateQR(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), "image/png", rr.Header().Get("Content-Type"))
}

func (s *QRHandlerTestSuite) TestGenerateQR_MissingText_Returns400() {
	// Test exact /qr path without text parameter
	req := httptest.NewRequest(http.MethodGet, "/qr", nil)
	rr := httptest.NewRecorder()

	handler.GenerateQR(rr, req)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Contains(s.T(), rr.Body.String(), "text parameter is required")
}

func (s *QRHandlerTestSuite) TestGenerateQR_MissingText_WithTrailingSlash() {
	req := httptest.NewRequest(http.MethodGet, "/qr/", nil)
	rr := httptest.NewRecorder()

	handler.GenerateQR(rr, req)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Contains(s.T(), rr.Body.String(), "text parameter is required")
}

func (s *QRHandlerTestSuite) TestGenerateQR_WrongMethod() {
	req := httptest.NewRequest(http.MethodPost, "/qr?text=test", nil)
	rr := httptest.NewRecorder()

	handler.GenerateQR(rr, req)

	assert.Equal(s.T(), http.StatusMethodNotAllowed, rr.Code)
}

func TestQRHandlerSuite(t *testing.T) {
	suite.Run(t, new(QRHandlerTestSuite))
}
