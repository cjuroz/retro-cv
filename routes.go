package main

import (
	"net/http"
)

func (a *App) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready\n"))
	})

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		a.renderPage(w, r, "home", "Home")
	})

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		a.renderPage(w, r, "about", "About")
	})

	mux.HandleFunc("GET /experience", func(w http.ResponseWriter, r *http.Request) {
		a.renderPage(w, r, "experience", "Experience")
	})

	mux.Handle("GET /static/", http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("static")),
	))

	mux.HandleFunc("GET /my-resume", a.handleResumeMarkdown)
	mux.HandleFunc("GET /my-resume-pdf", a.handleResumePDF)

	return mux
}
