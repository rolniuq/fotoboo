package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/internal/usecase"
)

type AdminHandler struct {
	adminUseCase *usecase.AdminUseCase
}

func NewAdminHandler(adminUseCase *usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{adminUseCase: adminUseCase}
}

func (h *AdminHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	stats, err := h.adminUseCase.GetStats()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to get stats"})
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

func (h *AdminHandler) HandleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getConfig(w, r)
	case http.MethodPut:
		h.updateConfig(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
	}
}

func (h *AdminHandler) getConfig(w http.ResponseWriter, r *http.Request) {
	cfg := h.adminUseCase.GetConfig()
	writeJSON(w, http.StatusOK, cfg)
}

func (h *AdminHandler) updateConfig(w http.ResponseWriter, r *http.Request) {
	var cfg domain.Config
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid config data"})
		return
	}

	// Validate
	if cfg.CountdownDuration < 1 || cfg.CountdownDuration > 10 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "countdown_duration must be 1-10"})
		return
	}
	if cfg.MaxUploadSizeMB < 1 || cfg.MaxUploadSizeMB > 50 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "max_upload_size_mb must be 1-50"})
		return
	}
	if cfg.PhotoRetentionDay < 1 || cfg.PhotoRetentionDay > 365 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "photo_retention_days must be 1-365"})
		return
	}

	h.adminUseCase.UpdateConfig(cfg)

	writeJSON(w, http.StatusOK, cfg)
}
