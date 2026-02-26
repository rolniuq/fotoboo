package handler

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"net/http"
	"strconv"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/internal/usecase"
	"github.com/nfnt/resize"
)

// PrintHandler handles print-ready photo generation
type PrintHandler struct {
	photoUseCase *usecase.PhotoUseCase
}

func NewPrintHandler(photoUseCase *usecase.PhotoUseCase) *PrintHandler {
	return &PrintHandler{photoUseCase: photoUseCase}
}

// PrintSize represents a print-ready photo size
type PrintSize struct {
	Name     string
	WidthPx  uint // width in pixels at target DPI
	HeightPx uint
	DPI      int
}

var printSizes = map[string]PrintSize{
	"4x6": {
		Name:     "4x6 inches",
		WidthPx:  1200, // 4 * 300 DPI
		HeightPx: 1800, // 6 * 300 DPI
		DPI:      300,
	},
	"5x7": {
		Name:     "5x7 inches",
		WidthPx:  1500,
		HeightPx: 2100,
		DPI:      300,
	},
	"6x8": {
		Name:     "6x8 inches",
		WidthPx:  1800,
		HeightPx: 2400,
		DPI:      300,
	},
	"2x6": {
		Name:     "2x6 inches (strip)",
		WidthPx:  600,
		HeightPx: 1800,
		DPI:      300,
	},
}

// HandlePrint serves a print-ready version of the photo
// GET /photos/{id}/print?size=4x6
func (h *PrintHandler) HandlePrint(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	// Get the photo
	_, data, err := h.photoUseCase.GetPhotoData(id)
	if err != nil {
		if errors.Is(err, domain.ErrPhotoNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "photo not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to retrieve photo"})
		return
	}

	// Get requested size
	sizeParam := r.URL.Query().Get("size")
	if sizeParam == "" {
		sizeParam = "4x6"
	}

	printSize, ok := printSizes[sizeParam]
	if !ok {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid size. supported: 4x6, 5x7, 6x8, 2x6"})
		return
	}

	// Decode the image
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to decode image"})
		return
	}

	// Resize to print dimensions while maintaining aspect ratio
	resized := resize.Thumbnail(printSize.WidthPx, printSize.HeightPx, img, resize.Lanczos3)

	// Encode as high-quality JPEG
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 95})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to encode print image"})
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "attachment; filename=\"fotoboo-"+id+"-"+sizeParam+".jpg\"")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// HandlePrintSizes lists available print sizes
// GET /print-sizes
func HandlePrintSizes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	type SizeInfo struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		Width  uint   `json:"width_px"`
		Height uint   `json:"height_px"`
		DPI    int    `json:"dpi"`
	}

	sizes := make([]SizeInfo, 0, len(printSizes))
	for key, ps := range printSizes {
		sizes = append(sizes, SizeInfo{
			Key:    key,
			Name:   ps.Name,
			Width:  ps.WidthPx,
			Height: ps.HeightPx,
			DPI:    ps.DPI,
		})
	}

	writeJSON(w, http.StatusOK, sizes)
}
