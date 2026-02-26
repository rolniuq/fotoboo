package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fotoboo/fotoboo/internal/usecase"
	qrcode "github.com/skip2/go-qrcode"
)

type QRHandler struct {
	photoUseCase *usecase.PhotoUseCase
	baseURL      string
}

func NewQRHandler(photoUseCase *usecase.PhotoUseCase, baseURL string) *QRHandler {
	return &QRHandler{
		photoUseCase: photoUseCase,
		baseURL:      baseURL,
	}
}

// GetPhotoQR is called from PhotoHandler as a sub-route
func (h *PhotoHandler) GetPhotoQR(w http.ResponseWriter, r *http.Request, id string) {
	// Verify photo exists
	_, err := h.useCase.GetPhoto(id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "photo not found")
		return
	}

	// Build the photo URL
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	forwardedProto := r.Header.Get("X-Forwarded-Proto")
	if forwardedProto != "" {
		scheme = forwardedProto
	}

	host := r.Host
	photoURL := fmt.Sprintf("%s://%s/photos/%s", scheme, host, id)

	// Generate QR code PNG
	png, err := qrcode.Encode(photoURL, qrcode.Medium, 256)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to generate QR code")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	w.Write(png)
}

// GenerateQR is a standalone handler for arbitrary URL QR codes
func GenerateQR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"method not allowed"}`))
		return
	}

	text := r.URL.Query().Get("text")
	if text == "" {
		// Only extract from path if it's a sub-path like /qr/something
		path := strings.TrimPrefix(r.URL.Path, "/qr")
		path = strings.TrimPrefix(path, "/") // Remove leading slash
		if path != "" && path != "qr" {
			text = path
		}
	}

	if text == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"text parameter is required"}`))
		return
	}

	png, err := qrcode.Encode(text, qrcode.Medium, 256)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to generate QR code"}`))
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	w.Write(png)
}
