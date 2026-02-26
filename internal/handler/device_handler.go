package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/internal/usecase"
)

type DeviceHandler struct {
	deviceUseCase *usecase.DeviceUseCase
}

func NewDeviceHandler(deviceUseCase *usecase.DeviceUseCase) *DeviceHandler {
	return &DeviceHandler{
		deviceUseCase: deviceUseCase,
	}
}

type RegisterDeviceRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type DeviceResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateDeviceRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Active bool   `json:"active"`
}

func deviceToResponse(d *domain.Device) DeviceResponse {
	return DeviceResponse{
		ID:        d.ID,
		Name:      d.Name,
		Type:      d.Type,
		Active:    d.Active,
		CreatedAt: d.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: d.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *DeviceHandler) HandleDevices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.registerDevice(w, r)
	case http.MethodGet:
		h.listDevices(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
	}
}

func (h *DeviceHandler) HandleDevice(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/devices/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "device id is required"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getDevice(w, r, id)
	case http.MethodPut:
		h.updateDevice(w, r, id)
	case http.MethodDelete:
		h.deleteDevice(w, r, id)
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
	}
}

func (h *DeviceHandler) registerDevice(w http.ResponseWriter, r *http.Request) {
	var req RegisterDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Type == "" {
		req.Type = "webcam"
	}

	device, err := h.deviceUseCase.RegisterDevice(req.Name, req.Type)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidDevice) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "device name is required"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to register device"})
		return
	}

	writeJSON(w, http.StatusCreated, deviceToResponse(device))
}

func (h *DeviceHandler) getDevice(w http.ResponseWriter, r *http.Request, id string) {
	device, err := h.deviceUseCase.GetDevice(id)
	if err != nil {
		if errors.Is(err, domain.ErrDeviceNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "device not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to get device"})
		return
	}

	writeJSON(w, http.StatusOK, deviceToResponse(device))
}

func (h *DeviceHandler) updateDevice(w http.ResponseWriter, r *http.Request, id string) {
	var req UpdateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	device, err := h.deviceUseCase.UpdateDevice(id, req.Name, req.Type, req.Active)
	if err != nil {
		if errors.Is(err, domain.ErrDeviceNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "device not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to update device"})
		return
	}

	writeJSON(w, http.StatusOK, deviceToResponse(device))
}

func (h *DeviceHandler) deleteDevice(w http.ResponseWriter, r *http.Request, id string) {
	err := h.deviceUseCase.DeleteDevice(id)
	if err != nil {
		if errors.Is(err, domain.ErrDeviceNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "device not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to delete device"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DeviceHandler) listDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.deviceUseCase.ListDevices()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to list devices"})
		return
	}

	responses := make([]DeviceResponse, len(devices))
	for i, d := range devices {
		responses[i] = deviceToResponse(d)
	}

	writeJSON(w, http.StatusOK, responses)
}
