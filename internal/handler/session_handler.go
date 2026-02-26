package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/internal/usecase"
)

type SessionHandler struct {
	sessionUseCase *usecase.SessionUseCase
}

func NewSessionHandler(sessionUseCase *usecase.SessionUseCase) *SessionHandler {
	return &SessionHandler{
		sessionUseCase: sessionUseCase,
	}
}

type StartSessionRequest struct {
	DeviceID string `json:"device_id"`
}

type SessionResponse struct {
	ID        string `json:"id"`
	DeviceID  string `json:"device_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func sessionToResponse(s *domain.Session) SessionResponse {
	return SessionResponse{
		ID:        s.ID,
		DeviceID:  s.DeviceID,
		Status:    string(s.Status),
		CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *SessionHandler) HandleSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.startSession(w, r)
	case http.MethodGet:
		h.listSessions(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
	}
}

func (h *SessionHandler) HandleSession(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/sessions/")
	parts := strings.SplitN(path, "/", 2)
	id := parts[0]

	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "session id is required"})
		return
	}

	// Check for sub-resources: /sessions/{id}/complete or /sessions/{id}/photos
	if len(parts) > 1 {
		switch parts[1] {
		case "complete":
			h.completeSession(w, r, id)
			return
		case "photos":
			h.getSessionPhotos(w, r, id)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		h.getSession(w, r, id)
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
	}
}

func (h *SessionHandler) startSession(w http.ResponseWriter, r *http.Request) {
	var req StartSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Allow empty body — device_id is optional
		req = StartSessionRequest{}
	}

	session, err := h.sessionUseCase.StartSession(req.DeviceID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to start session"})
		return
	}

	writeJSON(w, http.StatusCreated, sessionToResponse(session))
}

func (h *SessionHandler) getSession(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	session, err := h.sessionUseCase.GetSession(id)
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "session not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to get session"})
		return
	}

	writeJSON(w, http.StatusOK, sessionToResponse(session))
}

func (h *SessionHandler) completeSession(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	session, err := h.sessionUseCase.CompleteSession(id)
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "session not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to complete session"})
		return
	}

	writeJSON(w, http.StatusOK, sessionToResponse(session))
}

func (h *SessionHandler) getSessionPhotos(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	photos, err := h.sessionUseCase.GetSessionPhotos(id)
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "session not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to get session photos"})
		return
	}

	writeJSON(w, http.StatusOK, photos)
}

func (h *SessionHandler) listSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.sessionUseCase.ListSessions()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to list sessions"})
		return
	}

	responses := make([]SessionResponse, len(sessions))
	for i, s := range sessions {
		responses[i] = sessionToResponse(s)
	}

	writeJSON(w, http.StatusOK, responses)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
