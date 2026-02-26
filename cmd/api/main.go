package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/fotoboo/fotoboo/internal/background"
	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/fotoboo/fotoboo/internal/handler"
	"github.com/fotoboo/fotoboo/internal/middleware"
	"github.com/fotoboo/fotoboo/internal/repository"
	"github.com/fotoboo/fotoboo/internal/usecase"
)

func main() {
	startTime := time.Now()

	storagePath := getEnv("STORAGE_PATH", "./data/photos")
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "./data/fotoboo.db")
	webDir := getEnv("WEB_DIR", "./web")
	baseURL := getEnv("BASE_URL", "http://localhost:"+port)

	if err := ensureDataPaths(storagePath, dbPath); err != nil {
		log.Fatalf("Failed to prepare data paths: %v", err)
	}

	// MinIO configuration
	useMinIO := getEnv("USE_MINIO", "false") == "true"
	minioEndpoint := getEnv("MINIO_ENDPOINT", "localhost:9000")
	minioAccessKey := getEnv("MINIO_ACCESS_KEY", "minioadmin")
	minioSecretKey := getEnv("MINIO_SECRET_KEY", "minioadmin")
	minioBucket := getEnv("MINIO_BUCKET", "fotoboo")
	minioUseSSL := getEnv("MINIO_USE_SSL", "false") == "true"

	// Structured logger
	logger := middleware.NewStructuredLogger()

	// Initialize database
	db, err := repository.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Config store
	configStore := domain.NewConfigStore()

	// Repositories - choose between local storage or MinIO
	var photoRepo domain.PhotoRepository
	if useMinIO {
		minioConfig := domain.NewMinioConfig(minioEndpoint, minioAccessKey, minioSecretKey, minioBucket, minioUseSSL)
		photoRepo, err = repository.NewMinioPhotoRepository(db, minioConfig)
		if err != nil {
			log.Fatalf("Failed to initialize MinIO photo repository: %v", err)
		}
		logger.Info("Using MinIO for photo storage", "endpoint", minioEndpoint, "bucket", minioBucket)
	} else {
		photoRepo = repository.NewSQLitePhotoRepository(db, storagePath)
		logger.Info("Using local file storage for photos", "path", storagePath)
	}

	sessionRepo := repository.NewSQLiteSessionRepository(db)
	deviceRepo := repository.NewSQLiteDeviceRepository(db)

	// Use cases
	photoUseCase := usecase.NewPhotoUseCase(photoRepo)
	sessionUseCase := usecase.NewSessionUseCase(sessionRepo, photoRepo)
	deviceUseCase := usecase.NewDeviceUseCase(deviceRepo)
	adminUseCase := usecase.NewAdminUseCase(photoRepo, sessionRepo, deviceRepo, configStore)

	// Handlers
	photoHandler := handler.NewPhotoHandler(photoUseCase)
	sessionHandler := handler.NewSessionHandler(sessionUseCase)
	deviceHandler := handler.NewDeviceHandler(deviceUseCase)
	printHandler := handler.NewPrintHandler(photoUseCase)
	adminHandler := handler.NewAdminHandler(adminUseCase)
	_ = handler.NewQRHandler(photoUseCase, baseURL)

	// Middleware
	metrics := middleware.NewMetrics()
	rateLimiter := middleware.NewRateLimiter(100, time.Minute) // 100 req/min per IP
	sessionLimiter := middleware.NewSessionLimiter(50)         // max 50 concurrent sessions

	// Background jobs
	jobRunner := background.NewJobRunner(photoRepo, storagePath)
	jobRunner.Start()

	mux := http.NewServeMux()

	// Photo endpoints
	mux.HandleFunc("/photos", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		switch r.Method {
		case http.MethodPost:
			photoHandler.UploadPhoto(w, r)
		case http.MethodGet:
			photoHandler.ListPhotos(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/photos/", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/photos/")
		parts := strings.SplitN(path, "/", 2)
		id := parts[0]

		if len(parts) > 1 && parts[1] == "print" {
			printHandler.HandlePrint(w, r, id)
			return
		}

		if r.Method == http.MethodDelete {
			photoHandler.DeletePhoto(w, r)
			return
		}

		photoHandler.GetPhoto(w, r)
	})

	// Session endpoints
	mux.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Session limiting on POST (creating new sessions)
		if r.Method == http.MethodPost {
			if !sessionLimiter.Acquire() {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"error":"maximum concurrent sessions reached"}`))
				return
			}
			// Note: Release is called when session completes
		}

		sessionHandler.HandleSessions(w, r)
	})

	mux.HandleFunc("/sessions/", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Release session slot when completing
		path := strings.TrimPrefix(r.URL.Path, "/sessions/")
		parts := strings.SplitN(path, "/", 2)
		if len(parts) > 1 && parts[1] == "complete" && r.Method == http.MethodPost {
			sessionLimiter.Release()
		}

		sessionHandler.HandleSession(w, r)
	})

	// Device endpoints
	mux.HandleFunc("/devices", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		deviceHandler.HandleDevices(w, r)
	})

	mux.HandleFunc("/devices/", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		deviceHandler.HandleDevice(w, r)
	})

	// Admin endpoints
	mux.HandleFunc("/admin/stats", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		adminHandler.HandleStats(w, r)
	})

	mux.HandleFunc("/admin/config", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		adminHandler.HandleConfig(w, r)
	})

	// Print sizes
	mux.HandleFunc("/print-sizes", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler.HandlePrintSizes(w, r)
	})

	// QR code
	mux.HandleFunc("/qr", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler.GenerateQR(w, r)
	})

	// Metrics endpoint
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		metrics.HandleMetrics(w, r)
	})

	// Enhanced health check
	mux.HandleFunc("/health", middleware.EnhancedHealthCheck(db, storagePath, startTime))

	// Serve Vue.js SPA
	distDir := webDir + "/dist"
	if _, err := os.Stat(distDir); err == nil {
		webDir = distDir
	}
	mux.Handle("/", spaHandler(http.Dir(webDir)))

	// Apply middleware chain: rate limiter → metrics → logger → mux
	var chain http.Handler = mux
	chain = middleware.Logger(logger)(chain)
	chain = metrics.Middleware(chain)
	chain = rateLimiter.Middleware(chain)

	if useMinIO {
		logger.Info("FotoBoo API server starting",
			"port", port,
			"db_path", dbPath,
			"web_dir", webDir,
			"storage_type", "minio",
			"minio_endpoint", minioEndpoint,
			"minio_bucket", minioBucket,
		)
	} else {
		logger.Info("FotoBoo API server starting",
			"port", port,
			"db_path", dbPath,
			"web_dir", webDir,
			"storage_type", "local",
			"storage_path", storagePath,
		)
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		logger.Info("Shutting down...")
		jobRunner.Stop()
		db.Close()
		os.Exit(0)
	}()

	if err := http.ListenAndServe(":"+port, chain); err != nil {
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
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func ensureDataPaths(storagePath, dbPath string) error {
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return err
	}

	dbDir := filepath.Dir(dbPath)
	if dbDir == "." || dbDir == "" {
		return nil
	}

	return os.MkdirAll(dbDir, 0755)
}

// spaHandler serves static files and falls back to index.html for SPA routes
func spaHandler(fs http.FileSystem) http.Handler {
	fileServer := http.FileServer(fs)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		f, err := fs.Open(path)
		if err != nil {
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}
		f.Close()

		fileServer.ServeHTTP(w, r)
	})
}
