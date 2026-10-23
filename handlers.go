package main

import (
	"bytes"
	"log"
	"net/http"
	"strings"
)

func (a *App) renderPage(w http.ResponseWriter, name string, title string) {
	tmpl, ok := a.pages[name]
	if !ok {
		http.NotFound(w, nil)
		return
	}

	data := PageData{
		Title:   title,
		Profile: a.profileStore.Snapshot(),
		Version: version,
		Commit:  commit,
	}

	// Render completely before sending the response.
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout.html", data); err != nil {
		log.Printf("render %s page: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("write %s response: %v", name, err)
	}
}

func (a *App) handleResume(w http.ResponseWriter, r *http.Request) {
	p := a.profileStore.Snapshot()
	resume := generateResumeMarkdown(p)

	filename := strings.ToLower(
		strings.ReplaceAll(p.Name, " ", "-"),
	)

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set(
		"Content-Disposition",
		`attachment; filename="`+filename+`-resume.md"`,
	)
	w.Header().Set("Cache-Control", "no-cache")

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(resume)); err != nil {
		log.Printf("write resume response: %v", err)
	}
}
