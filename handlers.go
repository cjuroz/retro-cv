package main

import (
	"bytes"
	"log"
	"net/http"
	"strings"
)

func (a *App) renderPage(w http.ResponseWriter, r *http.Request, name string, title string) {
	tmpl, ok := a.pages[name]
	if !ok {
		http.NotFound(w, nil)
		return
	}

	p := a.profileStore.Snapshot()

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Fallback check for proxies/load balancers (optional but good practice)
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	seo := SEOData{
		Description: "Cloud and software engineering leader with 12+ years in AWS, SRE, DevOps and distributed systems.",
		Canonical:   scheme + "://" + r.Host,
		OGTitle:     title + " | " + p.Name,
		OGImage:     scheme + "://" + r.Host + "/static/img/og-image.png",
	}

	switch name {
	case "about":
		seo.Description = "Engineering leadership, cloud architecture, production reliability and backend engineering experience."
	case "experience":
		seo.Description = "Professional experience spanning cloud infrastructure, enterprise security SaaS, fintech and software architecture."
	}

	data := PageData{
		Title:   title,
		SEO:     seo,
		Profile: p,
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

func (a *App) handleResumeMarkdown(w http.ResponseWriter, r *http.Request) {
	p := a.profileStore.Snapshot()

	md, err := generateResumeMarkdown(p)
	if err != nil {
		log.Printf("generate resume markdown: %v", err)
		http.Error(w, "Could not generate resume markdown", http.StatusInternalServerError)
		return
	}

	filename := strings.ToLower(strings.ReplaceAll(p.Name, " ", "-"))

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`-resume.md"`)
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(md)); err != nil {
		log.Printf("write resume markdown response: %v", err)
	}
}

func (a *App) handleResumePDF(w http.ResponseWriter, r *http.Request) {
	p := a.profileStore.Snapshot()

	pdf, err := a.cachedResumePDF(p)
	if err != nil {
		log.Printf("generate resume PDF: %v", err)
		http.Error(w, "Could not generate resume PDF", http.StatusInternalServerError)
		return
	}

	filename := strings.ToLower(strings.ReplaceAll(p.Name, " ", "-"))

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`-resume.pdf"`)
	w.Header().Set("Cache-Control", "no-store")

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(pdf); err != nil {
		log.Printf("write resume PDF response: %v", err)
	}
}
