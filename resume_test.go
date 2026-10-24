package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResumeDownload(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/my-resume", nil)
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
