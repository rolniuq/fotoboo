package background

import (
	"log"
	"sync"
	"time"

	"github.com/fotoboo/fotoboo/internal/domain"
)

// JobRunner manages background tasks
type JobRunner struct {
	photoRepo domain.PhotoRepository
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

func NewJobRunner(photoRepo domain.PhotoRepository) *JobRunner {
	return &JobRunner{
		photoRepo: photoRepo,
		stopCh:    make(chan struct{}),
	}
}

// Start begins all background jobs
func (jr *JobRunner) Start() {
	log.Println("[background] Starting background jobs")

	// Cleanup old photos (run every hour, delete photos older than 30 days)
	jr.wg.Add(1)
	go jr.runPeriodic("cleanup", 1*time.Hour, jr.cleanupOldPhotos)
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
			deleted++
		}
	}

	if deleted > 0 {
		log.Printf("[background] cleanup: deleted %d old photos", deleted)
	}
}
