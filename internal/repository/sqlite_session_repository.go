package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fotoboo/fotoboo/internal/domain"
)

type SQLiteSessionRepository struct {
	db *sql.DB
}

func NewSQLiteSessionRepository(db *sql.DB) *SQLiteSessionRepository {
	return &SQLiteSessionRepository{db: db}
}

func (r *SQLiteSessionRepository) Save(session *domain.Session) error {
	_, err := r.db.Exec(
		`INSERT INTO sessions (id, device_id, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		session.ID, session.DeviceID, string(session.Status),
		session.CreatedAt.UTC().Format(time.RFC3339Nano),
		session.UpdatedAt.UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}

func (r *SQLiteSessionRepository) FindByID(id string) (*domain.Session, error) {
	row := r.db.QueryRow(
		`SELECT id, device_id, status, created_at, updated_at FROM sessions WHERE id = ?`, id,
	)
	return r.scanSession(row)
}

func (r *SQLiteSessionRepository) Update(session *domain.Session) error {
	session.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		`UPDATE sessions SET device_id = ?, status = ?, updated_at = ? WHERE id = ?`,
		session.DeviceID, string(session.Status),
		session.UpdatedAt.UTC().Format(time.RFC3339Nano),
		session.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	return nil
}

func (r *SQLiteSessionRepository) ListAll() ([]*domain.Session, error) {
	rows, err := r.db.Query(
		`SELECT id, device_id, status, created_at, updated_at FROM sessions ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	defer rows.Close()

	return r.scanSessions(rows)
}

func (r *SQLiteSessionRepository) CountAll() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM sessions`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count sessions: %w", err)
	}
	return count, nil
}

func (r *SQLiteSessionRepository) CountByDate(date time.Time) (int, error) {
	return countByDate(r.db, "sessions", date)
}

func (r *SQLiteSessionRepository) scanSession(row *sql.Row) (*domain.Session, error) {
	var session domain.Session
	var status, createdAt, updatedAt string

	err := row.Scan(&session.ID, &session.DeviceID, &status, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrSessionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan session: %w", err)
	}

	session.Status = domain.SessionStatus(status)
	session.CreatedAt = parseTime(createdAt)
	session.UpdatedAt = parseTime(updatedAt)

	return &session, nil
}

func (r *SQLiteSessionRepository) scanSessions(rows *sql.Rows) ([]*domain.Session, error) {
	var sessions []*domain.Session
	for rows.Next() {
		var session domain.Session
		var status, createdAt, updatedAt string

		err := rows.Scan(&session.ID, &session.DeviceID, &status, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session row: %w", err)
		}

		session.Status = domain.SessionStatus(status)
		session.CreatedAt = parseTime(createdAt)
		session.UpdatedAt = parseTime(updatedAt)
		sessions = append(sessions, &session)
	}

	return sessions, nil
}
