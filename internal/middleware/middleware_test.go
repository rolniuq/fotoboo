package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/fotoboo/fotoboo/internal/middleware"
)

// MetricsTestSuite tests Metrics middleware
type MetricsTestSuite struct {
	suite.Suite
	metrics *middleware.Metrics
}

func (s *MetricsTestSuite) SetupTest() {
	s.metrics = middleware.NewMetrics()
}

func (s *MetricsTestSuite) TestMiddleware_TracksRequests() {
	handler := s.metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	snap := s.metrics.Snapshot()
	assert.Equal(s.T(), int64(5), snap.TotalRequests)
	assert.Equal(s.T(), int64(0), snap.TotalErrors)
	assert.Equal(s.T(), int64(0), snap.ActiveRequests)
}

func (s *MetricsTestSuite) TestMiddleware_TracksErrors() {
	handler := s.metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	// Success
	req := httptest.NewRequest(http.MethodGet, "/success", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	// Error
	req = httptest.NewRequest(http.MethodGet, "/error", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	snap := s.metrics.Snapshot()
	assert.Equal(s.T(), int64(2), snap.TotalRequests)
	assert.Equal(s.T(), int64(1), snap.TotalErrors)
}

func (s *MetricsTestSuite) TestHandleMetrics_ReturnsJSON() {
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()

	s.metrics.HandleMetrics(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var snap middleware.MetricsSnapshot
	err := json.NewDecoder(rr.Body).Decode(&snap)
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), snap.Uptime)
}

func TestMetricsSuite(t *testing.T) {
	suite.Run(t, new(MetricsTestSuite))
}

// RateLimiterTestSuite tests RateLimiter middleware
type RateLimiterTestSuite struct {
	suite.Suite
}

func (s *RateLimiterTestSuite) TestAllowsWithinLimit() {
	rl := middleware.NewRateLimiter(5, time.Minute)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(s.T(), http.StatusOK, rr.Code, "Request %d should succeed", i+1)
	}
}

func (s *RateLimiterTestSuite) TestBlocksExcessRequests() {
	rl := middleware.NewRateLimiter(3, time.Minute)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Use up quota
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.2:12345"
		handler.ServeHTTP(httptest.NewRecorder(), req)
	}

	// 4th should be blocked
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.2:12345"
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusTooManyRequests, rr.Code)
}

func (s *RateLimiterTestSuite) TestDifferentIPsAreIndependent() {
	rl := middleware.NewRateLimiter(2, time.Minute)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Exhaust IP1
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		handler.ServeHTTP(httptest.NewRecorder(), req)
	}

	// IP1 blocked
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(s.T(), http.StatusTooManyRequests, rr.Code)

	// IP2 still allowed
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "10.0.0.2:12345"
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(s.T(), http.StatusOK, rr.Code)
}

func TestRateLimiterSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterTestSuite))
}

// SessionLimiterTestSuite tests SessionLimiter
type SessionLimiterTestSuite struct {
	suite.Suite
}

func (s *SessionLimiterTestSuite) TestAcquire_Success() {
	sl := middleware.NewSessionLimiter(3)

	assert.True(s.T(), sl.Acquire())
	assert.True(s.T(), sl.Acquire())
	assert.True(s.T(), sl.Acquire())
	assert.False(s.T(), sl.Acquire(), "4th acquire should fail")
	assert.Equal(s.T(), int64(3), sl.ActiveCount())
}

func (s *SessionLimiterTestSuite) TestRelease_FreesSlot() {
	sl := middleware.NewSessionLimiter(2)

	sl.Acquire()
	sl.Acquire()
	assert.False(s.T(), sl.Acquire())

	sl.Release()

	assert.True(s.T(), sl.Acquire(), "Should acquire after release")
}

func TestSessionLimiterSuite(t *testing.T) {
	suite.Run(t, new(SessionLimiterTestSuite))
}

// LoggerTestSuite tests Logger middleware
type LoggerTestSuite struct {
	suite.Suite
}

func (s *LoggerTestSuite) TestCapturesStatusCode() {
	logger := middleware.NewStructuredLogger()
	logMiddleware := middleware.Logger(logger)

	testCases := []struct {
		status int
	}{
		{http.StatusOK},
		{http.StatusNotFound},
		{http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		handler := logMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.status)
		}))

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(s.T(), tc.status, rr.Code)
	}
}

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
