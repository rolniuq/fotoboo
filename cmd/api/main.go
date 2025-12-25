package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fotoboo/fotoboo/internal/handler"
	"github.com/fotoboo/fotoboo/internal/repository"
	"github.com/fotoboo/fotoboo/internal/usecase"
)

func main() {
	storagePath := getEnv("STORAGE_PATH", "./data/photos")
	port := getEnv("PORT", "8080")

	photoRepo := repository.NewFilePhotoRepository(storagePath)
	photoUseCase := usecase.NewPhotoUseCase(photoRepo)
	photoHandler := handler.NewPhotoHandler(photoUseCase)

	mux := http.NewServeMux()

	mux.HandleFunc("/photos", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		photoHandler.UploadPhoto(w, r)
	})

	mux.HandleFunc("/photos/", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		photoHandler.GetPhoto(w, r)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("FotoBoo API server starting on port %s", port)
	log.Printf("Storage path: %s", storagePath)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return strings.TrimSpace(value)
	}
	return defaultValue
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
