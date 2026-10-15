package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoutes(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}

	handler := securityHeaders(app.routes())
	server := httptest.NewServer(handler)
	defer server.Close()

	tests := []struct {
		path       string
		wantStatus int
		wantBody   string
	}{
		{"/", http.StatusOK, "CRISTIAN"},
		{"/about", http.StatusOK, "WHO I AM"},
		{"/experience", http.StatusOK, "PROFESSIONAL HISTORY"},
		{"/missing", http.StatusNotFound, "404 page not found"},
		{"/health", http.StatusOK, "ok"},
		{"/ready", http.StatusOK, "ready"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			resp, err := http.Get(server.URL + tt.path)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if !strings.Contains(string(body), tt.wantBody) {
				t.Errorf("response does not contain %q", tt.wantBody)
			}

			if resp.Header.Get("X-Content-Type-Options") != "nosniff" {
				t.Error("missing X-Content-Type-Options header")
			}
		})
	}
}

func TestLoadProfile(t *testing.T) {
	profile, err := loadProfile("data/profile.json")
	if err != nil {
		t.Fatal(err)
	}

	if profile.Name == "" {
		t.Fatal("profile name should not be empty")
	}

	if len(profile.Experience) == 0 {
		t.Fatal("expected at least one example experience")
	}
}
