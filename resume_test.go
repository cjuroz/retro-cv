package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cjuroz/retro-cv/internal/profile"
)

func TestResumeDownload(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/resume.md", nil)
	rec := httptest.NewRecorder()

	app.routes().ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d",
			resp.StatusCode, http.StatusOK)
	}

	profile := app.profileStore.Snapshot()

	if got := resp.Header.Get("Content-Disposition"); got != `attachment; filename="`+strings.ToLower(strings.ReplaceAll(profile.Name, " ", "-"))+`-resume.md"` {
		t.Errorf("Content-Disposition = %q", got)
	}

	if got := resp.Header.Get("Content-Type"); got != "text/markdown; charset=utf-8" {
		t.Errorf("Content-Type = %q", got)
	}

	body := rec.Body.String()

	for _, want := range []string{
		"# " + profile.Name + " " + profile.NameZh,
		profile.Address,
		"[" + profile.Email + "](mailto:" + profile.Email + ")",
		profile.WorkStatus,
		"[LinkedIn](" + profile.Socials.LinkedIn + ")",
		"[GitHub](" + profile.Socials.GitHub + ")",
		profile.Headline,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("resume missing %q", want)
		}
	}
}

func TestRenderResumePDF(t *testing.T) {
	p := profile.Profile{
		Name:       "Example Name",
		ResumeURL:  "https://example.com/resume",
		Address:    "Taipei, Taiwan",
		Email:      "person@example.com",
		WorkStatus: "Open to opportunities",
		Headline:   "Engineering leader",
		Summary:    "Experienced engineering leader.",
		Socials: profile.Socials{
			LinkedIn: "https://example.com/linkedin",
			GitHub:   "https://example.com/github",
		},
		Expertises: []string{"Cloud architecture"},
		Skills:     []string{"Go", "AWS"},
		Experience: []profile.Experience{{
			Company:      "Example Corp",
			Summary:      "Built reliable systems.",
			Website:      "https://example.com",
			Location:     "Taipei",
			Role:         "Engineering Lead",
			Period:       "2020 - Present",
			Tasks:        []string{"Led a platform team."},
			Technologies: []string{"Go", "AWS"},
		}},
		Education: []profile.Education{{
			Institution: "Example University",
			Degree:      "BSc Computer Science",
			Period:      "2010 - 2014",
		}},
		Languages: []profile.Language{{Language: "English", Level: "Fluent"}},
	}

	pdf, err := generateResumePDF(p)
	if err != nil {
		t.Fatalf("renderResumePDF() error = %v", err)
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Fatal("renderResumePDF() did not return a PDF document")
	}
}
