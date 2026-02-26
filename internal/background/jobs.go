package background

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fotoboo/fotoboo/internal/domain"
	"github.com/nfnt/resize"
)

// JobRunner manages background tasks
type JobRunner struct {
	photoRepo   domain.PhotoRepository
	storagePath string
	stopCh      chan struct{}
	wg          sync.WaitGroup
}

func NewJobRunner(photoRepo domain.PhotoRepository, storagePath string) *JobRunner {
	return &JobRunner{
		photoRepo:   photoRepo,
		storagePath: storagePath,
		stopCh:      make(chan struct{}),
	}
}

// Start begins all background jobs
func (jr *JobRunner) Start() {
	log.Println("[background] Starting background jobs")

	// Cleanup old photos (run every hour, delete photos older than 30 days)
	jr.wg.Add(1)
	go jr.runPeriodic("cleanup", 1*time.Hour, jr.cleanupOldPhotos)

	// Generate thumbnails for photos that don't have one
	jr.wg.Add(1)
	go jr.runPeriodic("thumbnails", 5*time.Minute, jr.generateMissingThumbnails)
}

// Stop gracefully stops all background jobs
func (jr *JobRunner) Stop() {
	log.Println("[background] Stopping background jobs...")
	close(jr.stopCh)
	jr.wg.Wait()
	log.Println("[background] All background jobs stopped")
}

func (jr *JobRunner) runPeriodic(name string, interval time.Duration, fn func()) {
	defer jr.wg.Done()

	// Run immediately on start
	fn()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-jr.stopCh:
			log.Printf("[background] Job '%s' stopped", name)
			return
		case <-ticker.C:
			fn()
		}
	}
}

// cleanupOldPhotos removes photos older than 30 days
func (jr *JobRunner) cleanupOldPhotos() {
	maxAge := 30 * 24 * time.Hour
	cutoff := time.Now().Add(-maxAge)

	photos, err := jr.photoRepo.ListAll()
	if err != nil {
		log.Printf("[background] cleanup: failed to list photos: %v", err)
		return
	}

	deleted := 0
	for _, photo := range photos {
		if photo.CreatedAt.Before(cutoff) {
			if err := jr.photoRepo.Delete(photo.ID); err != nil {
				log.Printf("[background] cleanup: failed to delete photo %s: %v", photo.ID, err)
				continue
			}
			// Also remove thumbnail
			thumbPath := filepath.Join(jr.storagePath, "thumbnails", photo.ID+"_thumb.jpg")
			os.Remove(thumbPath)
			deleted++
		}
	}

	if deleted > 0 {
		log.Printf("[background] cleanup: deleted %d old photos", deleted)
	}
}

// generateMissingThumbnails creates thumbnails for photos that don't have one
func (jr *JobRunner) generateMissingThumbnails() {
	thumbDir := filepath.Join(jr.storagePath, "thumbnails")
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		log.Printf("[background] thumbnails: failed to create directory: %v", err)
		return
	}

	photos, err := jr.photoRepo.ListAll()
	if err != nil {
		log.Printf("[background] thumbnails: failed to list photos: %v", err)
		return
	}

	generated := 0
	for _, photo := range photos {
		thumbPath := filepath.Join(thumbDir, photo.ID+"_thumb.jpg")

		// Skip if thumbnail already exists
		if _, err := os.Stat(thumbPath); err == nil {
			continue
		}

		// Read original photo
		data, err := jr.photoRepo.GetFileData(photo)
		if err != nil {
			log.Printf("[background] thumbnails: failed to read photo %s: %v", photo.ID, err)
			continue
		}

		// Generate thumbnail
		thumbData, err := generateThumbnail(data, 320, 240)
		if err != nil {
			log.Printf("[background] thumbnails: failed to generate thumbnail for %s: %v", photo.ID, err)
			continue
		}

		// Save thumbnail
		if err := os.WriteFile(thumbPath, thumbData, 0644); err != nil {
			log.Printf("[background] thumbnails: failed to save thumbnail for %s: %v", photo.ID, err)
			continue
		}

		generated++
	}

	if generated > 0 {
		log.Printf("[background] thumbnails: generated %d thumbnails", generated)
	}
}

// generateThumbnail creates a thumbnail from image data
func generateThumbnail(data []byte, maxWidth, maxHeight uint) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	thumb := resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, thumb, &jpeg.Options{Quality: 80})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
