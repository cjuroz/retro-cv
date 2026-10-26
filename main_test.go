package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cjuroz/retro-cv/internal/profile"
)

func TestUniqueTechnologies(t *testing.T) {
	experiences := []profile.Experience{
		{Technologies: []string{"Go", "AWS"}},
		{Technologies: []string{"AWS", "Docker"}},
	}

	got := uniqueTechnologies(experiences)
	want := []string{"Go", "AWS", "Docker"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("uniqueTechnologies() = %v, want %v", got, want)
	}
}

func TestRoutes(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatalf("newApp() error = %v", err)
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
		{"/about", http.StatusOK, "ABOUT"},
		{"/experience", http.StatusOK, "PROFESSIONAL HISTORY"},
		{"/missing", http.StatusNotFound, "404 page not found"},
		{"/health", http.StatusOK, "ok"},
		{"/ready", http.StatusOK, "ready"},
		{"/resume.md", http.StatusOK, "WORK EXPERIENCE"},
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
	got, err := loadProfile("data/profile.json")
	if err != nil {
		t.Fatalf("loadProfile() error = %v", err)
	}

	if strings.TrimSpace(got.Name) == "" {
		t.Fatal("profile name should not be empty")
	}

	if strings.TrimSpace(got.Headline) == "" {
		t.Error("profile headline should not be empty")
	}

	if len(got.Experience) == 0 {
		t.Fatal("expected at least one experience entry")
	}

	for i, exp := range got.Experience {
		if strings.TrimSpace(exp.Company) == "" {
			t.Errorf("experience[%d]: company should not be empty", i)
		}
		if strings.TrimSpace(exp.Summary) == "" {
			t.Errorf("experience[%d]: summary should not be empty", i)
		}
		if strings.TrimSpace(exp.Location) == "" {
			t.Errorf("experience[%d]: location should not be empty", i)
		}
		if strings.TrimSpace(exp.Role) == "" {
			t.Errorf("experience[%d]: role should not be empty", i)
		}
		if strings.TrimSpace(exp.Period) == "" {
			t.Errorf("experience[%d]: period should not be empty", i)
		}
		if len(exp.Tasks) == 0 {
			t.Errorf("experience[%d]: description should contain at least one bullet", i)
		}

		for j, item := range exp.Tasks {
			if strings.TrimSpace(item) == "" {
				t.Errorf("experience[%d]: description[%d] should not be empty", i, j)
			}
		}
	}
}
