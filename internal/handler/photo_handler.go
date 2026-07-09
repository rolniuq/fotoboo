package handler

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/internal/usecase"
)

type PhotoHandler struct {
	useCase *usecase.PhotoUseCase
}

func NewPhotoHandler(useCase *usecase.PhotoUseCase) *PhotoHandler {
	return &PhotoHandler{
		useCase: useCase,
	}
}

type UploadResponse struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id,omitempty"`
	CreatedAt string `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *PhotoHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10MB max

	data, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "failed to read request body"})
		return
	}

	// Get optional session_id from query param
	sessionID := r.URL.Query().Get("session_id")

	photo, err := h.useCase.UploadPhoto(sessionID, data)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPhoto) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid photo data"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to save photo"})
		return
	}

	response := UploadResponse{
		ID:        photo.ID,
		SessionID: photo.SessionID,
		CreatedAt: photo.CreatedAt.Format(time.RFC3339),
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *PhotoHandler) GetPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	// Parse path: /photos/{id} or /photos/{id}/qr
	path := strings.TrimPrefix(r.URL.Path, "/photos/")
	parts := strings.SplitN(path, "/", 2)
	id := parts[0]

	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "photo id is required"})
		return
	}

	// Check for sub-resources
	if len(parts) > 1 && parts[1] == "qr" {
		h.GetPhotoQR(w, r, id)
		return
	}

	photo, data, err := h.useCase.GetPhotoData(id)
	if err != nil {
		if errors.Is(err, domain.ErrPhotoNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "photo not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to retrieve photo"})
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "inline; filename=\""+photo.ID+".jpg\"")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *PhotoHandler) ListPhotos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	photos, err := h.useCase.ListPhotos()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to list photos"})
		return
	}

	writeJSON(w, http.StatusOK, photos)
}

func (h *PhotoHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/photos/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "photo id is required"})
		return
	}

	err := h.useCase.DeletePhoto(id)
	if err != nil {
		if errors.Is(err, domain.ErrPhotoNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "photo not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to delete photo"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
