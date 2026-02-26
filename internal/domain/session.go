package domain

import (
	"time"

	"github.com/google/uuid"
)

type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
)

type Session struct {
	ID        string        `json:"id"`
	DeviceID  string        `json:"device_id"`
	Status    SessionStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func NewSession(deviceID string) *Session {
	now := time.Now()
	return &Session{
		ID:        uuid.New().String(),
		DeviceID:  deviceID,
		Status:    SessionStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (s *Session) Complete() {
	s.Status = SessionStatusCompleted
	s.UpdatedAt = time.Now()
}

type SessionRepository interface {
	Save(session *Session) error
	FindByID(id string) (*Session, error)
	Update(session *Session) error
	ListAll() ([]*Session, error)
	CountAll() (int, error)
	CountByDate(date time.Time) (int, error)
}
