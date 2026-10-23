package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cjuroz/retro-cv/internal/profile"
)

type profileStore struct {
	mu      sync.RWMutex
	current profile.Profile
}

func (s *profileStore) Snapshot() profile.Profile {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current
}

func (s *profileStore) Reload(path string) error {
	next, err := loadProfile(path)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.current = next
	s.mu.Unlock()

	return nil
}

func watchProfile(ctx context.Context, path string, store *profileStore) (<-chan struct{}, error) {
	initial, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	done := make(chan struct{})

	go func() {
		defer close(done)

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		lastSeen := initial

		for {
			select {
			case <-ctx.Done():
				log.Println("profile watcher stopped")
				return

			case <-ticker.C:
				current, err := os.ReadFile(path)
				if err != nil {
					log.Printf("profile poll failed: %v", err)
					continue
				}

				if bytes.Equal(current, lastSeen) {
					continue
				}

				// Avoid repeatedly logging the same invalid JSON.
				lastSeen = current

				if err := store.Reload(path); err != nil {
					log.Printf("profile reload failed: %v", err)
					continue
				}

				log.Println("profile.json reloaded")
			}
		}
	}()

	return done, nil
}
