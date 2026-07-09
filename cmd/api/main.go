package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
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
	"github.com/fotoboo/fotoboo/pkg/storage"
)

func main() {
	startTime := time.Now()

	storagePath := getEnv("STORAGE_PATH", "./data/photos")
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "./data/fotoboo.db")
	webDir := getEnv("WEB_DIR", "./web")
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

	// Initialize storage backend
	var (
		storageBackend storage.Storage
		storageType    string
	)
	if useMinIO {
		storageType = "minio"
		minioConfig := &storage.MinioConfig{
			Endpoint:        minioEndpoint,
			AccessKeyID:     minioAccessKey,
			SecretAccessKey: minioSecretKey,
			BucketName:      minioBucket,
			UseSSL:          minioUseSSL,
			Region:          "us-east-1",
		}
		var err error
		storageBackend, err = storage.NewMinioStorage(minioConfig)
		if err != nil {
			log.Fatalf("Failed to initialize MinIO storage: %v", err)
		}
	} else {
		storageType = "local"
		var err error
		storageBackend, err = storage.NewLocalStorage(storagePath)
		if err != nil {
			log.Fatalf("Failed to initialize local storage: %v", err)
		}
	}
	defer storageBackend.Close()

	// Repositories
	photoRepo := repository.NewPhotoRepository(db, storageBackend)
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

	// Middleware
	metrics := middleware.NewMetrics()
	rateLimiter := middleware.NewRateLimiter(100, time.Minute) // 100 req/min per IP
	sessionLimiter := middleware.NewSessionLimiter(50)         // max 50 concurrent sessions

	// Background jobs
	jobRunner := background.NewJobRunner(photoRepo)
	jobRunner.Start()

	mux := http.NewServeMux()

	// Photo endpoints
	mux.HandleFunc("/photos", middleware.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			photoHandler.UploadPhoto(w, r)
		case http.MethodGet:
			photoHandler.ListPhotos(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/photos/", middleware.WithCORS(func(w http.ResponseWriter, r *http.Request) {
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
	}))

	// Session endpoints
	mux.HandleFunc("/sessions", middleware.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		// Session limiting on POST (creating new sessions)
		if r.Method == http.MethodPost {
			if !sessionLimiter.Acquire() {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"error":"maximum concurrent sessions reached"}`))
				return
			}
		}

		sessionHandler.HandleSessions(w, r)
	}))

	mux.HandleFunc("/sessions/", middleware.WithCORS(func(w http.ResponseWriter, r *http.Request) {
		// Release session slot when completing
		path := strings.TrimPrefix(r.URL.Path, "/sessions/")
		parts := strings.SplitN(path, "/", 2)
		if len(parts) > 1 && parts[1] == "complete" && r.Method == http.MethodPost {
			sessionLimiter.Release()
		}

		sessionHandler.HandleSession(w, r)
	}))

	// Device endpoints
	mux.HandleFunc("/devices", middleware.WithCORS(deviceHandler.HandleDevices))
	mux.HandleFunc("/devices/", middleware.WithCORS(deviceHandler.HandleDevice))

	// Admin endpoints
	mux.HandleFunc("/admin/stats", middleware.WithCORS(adminHandler.HandleStats))
	mux.HandleFunc("/admin/config", middleware.WithCORS(adminHandler.HandleConfig))

	// Print sizes & QR code
	mux.HandleFunc("/print-sizes", middleware.WithCORS(handler.HandlePrintSizes))
	mux.HandleFunc("/qr", middleware.WithCORS(handler.GenerateQR))

	// Metrics endpoint (CORS but no OPTIONS handling needed)
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
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

	logger.Info("startup",
		"port", port,
		"db_path", dbPath,
		"web_dir", webDir,
		"storage_type", storageType,
		"storage_path", storagePath,
		"minio_endpoint", minioEndpoint,
		"minio_bucket", minioBucket,
	)

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
