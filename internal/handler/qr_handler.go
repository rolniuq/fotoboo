package handler

import (
	"fmt"
	"net/http"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// GetPhotoQR serves a QR code for a specific photo
func (h *PhotoHandler) GetPhotoQR(w http.ResponseWriter, r *http.Request, id string) {
	// Verify photo exists
	_, err := h.useCase.GetPhoto(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "photo not found"})
		return
	}

	// Build the photo URL
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if fwd := r.Header.Get("X-Forwarded-Proto"); fwd != "" {
		scheme = fwd
	}

	photoURL := fmt.Sprintf("%s://%s/photos/%s", scheme, r.Host, id)

	// Generate QR code PNG
	png, err := qrcode.Encode(photoURL, qrcode.Medium, 256)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to generate QR code"})
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
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	text := r.URL.Query().Get("text")
	if text == "" {
		path := strings.TrimPrefix(r.URL.Path, "/qr")
		path = strings.TrimPrefix(path, "/")
		if path != "" && path != "qr" {
			text = path
		}
	}

	if text == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "text parameter is required"})
		return
	}

	png, err := qrcode.Encode(text, qrcode.Medium, 256)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to generate QR code"})
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	w.Write(png)
}
