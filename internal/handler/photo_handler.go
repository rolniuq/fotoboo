package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

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
	CreatedAt string `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *PhotoHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10MB max

	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}

	photo, err := h.useCase.UploadPhoto(data)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPhoto) {
			h.writeError(w, http.StatusBadRequest, "invalid photo data")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "failed to save photo")
		return
	}

	response := UploadResponse{
		ID:        photo.ID,
		CreatedAt: photo.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *PhotoHandler) GetPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/photos/")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "photo id is required")
		return
	}

	photo, data, err := h.useCase.GetPhotoData(id)
	if err != nil {
		if errors.Is(err, domain.ErrPhotoNotFound) {
			h.writeError(w, http.StatusNotFound, "photo not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "failed to retrieve photo")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "inline; filename=\""+photo.ID+".jpg\"")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *PhotoHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
