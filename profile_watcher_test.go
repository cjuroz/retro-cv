package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTestProfile(t *testing.T, path, headline string) {
	t.Helper()

	data := []byte(`{
		"name": "Test User",
		"headline": "` + headline + `",
		"summary": "Test summary"
	}`)

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write profile: %v", err)
	}
}

func waitForHeadline(t *testing.T, store *profileStore, want string) {
	t.Helper()

	deadline := time.Now().Add(4 * time.Second)

	for time.Now().Before(deadline) {
		if got := store.Snapshot().Headline; got == want {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatalf(
		"timed out waiting for headline %q; got %q",
		want,
		store.Snapshot().Headline,
	)
}

func TestProfileStoreReload(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profile.json")
	store := &profileStore{}

	writeTestProfile(t, path, "Original headline")

	if err := store.Reload(path); err != nil {
		t.Fatalf("Reload() error = %v", err)
	}

	if got := store.Snapshot().Headline; got != "Original headline" {
		t.Fatalf("Headline = %q, want %q", got, "Original headline")
	}
}

func TestProfileStoreReloadInvalidJSONPreservesCurrentProfile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profile.json")
	store := &profileStore{}

	writeTestProfile(t, path, "Original headline")

	if err := store.Reload(path); err != nil {
		t.Fatalf("initial Reload() error = %v", err)
	}

	if err := os.WriteFile(path, []byte(`{"name":`), 0644); err != nil {
		t.Fatalf("write invalid JSON: %v", err)
	}

	if err := store.Reload(path); err == nil {
		t.Fatal("Reload() error = nil, want an error for invalid JSON")
	}

	if got := store.Snapshot().Headline; got != "Original headline" {
		t.Fatalf("Headline after failed reload = %q, want %q", got, "Original headline")
	}
}

func TestWatchProfileReloadsChangedFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profile.json")
	store := &profileStore{}

	writeTestProfile(t, path, "Original headline")

	if err := store.Reload(path); err != nil {
		t.Fatalf("initial Reload() error = %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done, err := watchProfile(ctx, path, store)
	if err != nil {
		cancel()
		t.Fatalf("watchProfile() error = %v", err)
	}

	writeTestProfile(t, path, "Updated headline")

	waitForHeadline(t, store, "Updated headline")

	t.Cleanup(func() {
		cancel()
		<-done
	})
}
