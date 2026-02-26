package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Logger creates a structured logging middleware
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"duration_ms", time.Since(start).Milliseconds(),
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

// Metrics tracks request metrics in memory
type Metrics struct {
	TotalRequests  atomic.Int64
	TotalErrors    atomic.Int64
	ActiveRequests atomic.Int64
	startTime      time.Time
	pathCounts     sync.Map // map[string]*atomic.Int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		startTime: time.Now(),
	}
}

func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.TotalRequests.Add(1)
		m.ActiveRequests.Add(1)

		// Track per-path
		path := r.URL.Path
		val, _ := m.pathCounts.LoadOrStore(path, &atomic.Int64{})
		val.(*atomic.Int64).Add(1)

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		m.ActiveRequests.Add(-1)

		if rw.statusCode >= 400 {
			m.TotalErrors.Add(1)
		}
	})
}

type MetricsSnapshot struct {
	Uptime         string `json:"uptime"`
	TotalRequests  int64  `json:"total_requests"`
	TotalErrors    int64  `json:"total_errors"`
	ActiveRequests int64  `json:"active_requests"`
	ErrorRate      string `json:"error_rate"`
}

func (m *Metrics) Snapshot() MetricsSnapshot {
	total := m.TotalRequests.Load()
	errors := m.TotalErrors.Load()

	errorRate := "0.00%"
	if total > 0 {
		errorRate = formatPercent(float64(errors) / float64(total) * 100)
	}

	return MetricsSnapshot{
		Uptime:         time.Since(m.startTime).Round(time.Second).String(),
		TotalRequests:  total,
		TotalErrors:    errors,
		ActiveRequests: m.ActiveRequests.Load(),
		ErrorRate:      errorRate,
	}
}

func formatPercent(f float64) string {
	return fmt.Sprintf("%.2f%%", f)
}

// HandleMetrics returns the metrics endpoint handler
func (m *Metrics) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m.Snapshot())
}

// RateLimiter implements per-IP token bucket rate limiting
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int           // requests per window
	window   time.Duration // time window
}

type visitor struct {
	tokens    int
	lastSeen  time.Time
	resetTime time.Time
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}

	// Cleanup old entries every minute
	go func() {
		for {
			time.Sleep(time.Minute)
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if !rl.allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "rate limit exceeded",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	v, exists := rl.visitors[ip]
	if !exists || now.After(v.resetTime) {
		rl.visitors[ip] = &visitor{
			tokens:    rl.rate - 1,
			lastSeen:  now,
			resetTime: now.Add(rl.window),
		}
		return true
	}

	v.lastSeen = now
	if v.tokens > 0 {
		v.tokens--
		return true
	}

	return false
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, v := range rl.visitors {
		if now.Sub(v.lastSeen) > 5*time.Minute {
			delete(rl.visitors, ip)
		}
	}
}

// SessionLimiter limits the number of concurrent active sessions
type SessionLimiter struct {
	maxSessions int
	active      atomic.Int64
}

func NewSessionLimiter(max int) *SessionLimiter {
	return &SessionLimiter{maxSessions: max}
}

func (sl *SessionLimiter) Acquire() bool {
	current := sl.active.Load()
	if int(current) >= sl.maxSessions {
		return false
	}
	sl.active.Add(1)
	return true
}

func (sl *SessionLimiter) Release() {
	sl.active.Add(-1)
}

func (sl *SessionLimiter) ActiveCount() int64 {
	return sl.active.Load()
}

// EnhancedHealthCheck returns a detailed health check handler
func EnhancedHealthCheck(db *sql.DB, storagePath string, startTime time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := "ok"
		checks := make(map[string]string)

		// DB check
		if err := db.Ping(); err != nil {
			status = "degraded"
			checks["database"] = "error: " + err.Error()
		} else {
			checks["database"] = "ok"
		}

		// Storage check
		if _, err := os.Stat(storagePath); err != nil {
			status = "degraded"
			checks["storage"] = "error: " + err.Error()
		} else {
			checks["storage"] = "ok"
		}

		checks["uptime"] = time.Since(startTime).Round(time.Second).String()

		result := map[string]interface{}{
			"status": status,
			"checks": checks,
		}

		w.Header().Set("Content-Type", "application/json")
		if status != "ok" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		json.NewEncoder(w).Encode(result)
	}
}

// NewStructuredLogger creates a slog.Logger with JSON output
func NewStructuredLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
