package usecase

import (
	"github.com/fotoboo/fotoboo/internal/domain"
)

type SessionUseCase struct {
	sessionRepo domain.SessionRepository
	photoRepo   domain.PhotoRepository
}

func NewSessionUseCase(sessionRepo domain.SessionRepository, photoRepo domain.PhotoRepository) *SessionUseCase {
	return &SessionUseCase{
		sessionRepo: sessionRepo,
		photoRepo:   photoRepo,
	}
}

func (uc *SessionUseCase) StartSession(deviceID string) (*domain.Session, error) {
	session := domain.NewSession(deviceID)

	if err := uc.sessionRepo.Save(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (uc *SessionUseCase) GetSession(id string) (*domain.Session, error) {
	return uc.sessionRepo.FindByID(id)
}

func (uc *SessionUseCase) CompleteSession(id string) (*domain.Session, error) {
	session, err := uc.sessionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	session.Complete()

	if err := uc.sessionRepo.Update(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (uc *SessionUseCase) GetSessionPhotos(sessionID string) ([]*domain.Photo, error) {
	// Verify session exists
	if _, err := uc.sessionRepo.FindByID(sessionID); err != nil {
		return nil, err
	}

	return uc.photoRepo.FindBySessionID(sessionID)
}

func (uc *SessionUseCase) ListSessions() ([]*domain.Session, error) {
	return uc.sessionRepo.ListAll()
}

func (uc *SessionUseCase) CountSessions() (int, error) {
	return uc.sessionRepo.CountAll()
}
