package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func InitDB(dbPath string) (*sql.DB, error) {
	if err := ensureDBDir(dbPath); err != nil {
		return nil, fmt.Errorf("failed to prepare database path: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS devices (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL DEFAULT 'webcam',
			active BOOLEAN NOT NULL DEFAULT 1,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			device_id TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'active',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS photos (
			id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL DEFAULT '',
			file_path TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_photos_session_id ON photos(session_id)`,
		`CREATE INDEX IF NOT EXISTS idx_photos_created_at ON photos(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_device_id ON sessions(device_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON sessions(created_at)`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute statement: %w\n%s", err, stmt)
		}
	}

	return nil
}

func ensureDBDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if dir == "." || dir == "" {
		return nil
	}

	return os.MkdirAll(dir, 0755)
}

// parseTime is a shared helper to parse RFC3339Nano timestamps stored in SQLite.
// Returns zero time on parse failure so callers don't need to handle the error.
func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339Nano, s)
	return t
}

// countByDate is a shared helper for COUNT queries filtered by a calendar day.
func countByDate(db *sql.DB, table string, date time.Time) (int, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, time.UTC).Format(time.RFC3339)

	var count int
	err := db.QueryRow(
		`SELECT COUNT(*) FROM `+table+` WHERE created_at >= ? AND created_at <= ?`,
		startOfDay, endOfDay,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count %s by date: %w", table, err)
	}
	return count, nil
}
