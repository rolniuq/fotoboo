package domain

import (
	"sync"
)

// Config holds booth configuration that can be changed at runtime
type Config struct {
	EventName         string   `json:"event_name"`
	CountdownDuration int      `json:"countdown_duration"` // seconds
	AvailableFrames   []string `json:"available_frames"`
	AvailableFilters  []string `json:"available_filters"`
	MaxUploadSizeMB   int      `json:"max_upload_size_mb"`
	PhotoRetentionDay int      `json:"photo_retention_days"`
}

func DefaultConfig() *Config {
	return &Config{
		EventName:         "FotoBoo Event",
		CountdownDuration: 3,
		AvailableFrames:   []string{"none", "simple", "event", "party"},
		AvailableFilters:  []string{"none", "grayscale", "vintage", "brightness", "contrast"},
		MaxUploadSizeMB:   10,
		PhotoRetentionDay: 30,
	}
}

// ConfigStore provides thread-safe access to booth configuration
type ConfigStore struct {
	mu     sync.RWMutex
	config *Config
}

func NewConfigStore() *ConfigStore {
	return &ConfigStore{
		config: DefaultConfig(),
	}
}

func (cs *ConfigStore) Get() Config {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	// Return a copy
	c := *cs.config
	c.AvailableFrames = make([]string, len(cs.config.AvailableFrames))
	copy(c.AvailableFrames, cs.config.AvailableFrames)
	c.AvailableFilters = make([]string, len(cs.config.AvailableFilters))
	copy(c.AvailableFilters, cs.config.AvailableFilters)
	return c
}

func (cs *ConfigStore) Update(cfg Config) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.config = &cfg
}


