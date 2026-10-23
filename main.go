package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cjuroz/retro-cv/internal/profile"
)

var (
	version = "dev"
	commit  = "unknown"
)

type PageData struct {
	Title   string
	Profile profile.Profile
	Version string
	Commit  string
}

type App struct {
	pages        map[string]*template.Template
	profileStore *profileStore
}

func uniqueTechnologies(experiences []profile.Experience) []string {
	seen := make(map[string]struct{})
	technologies := make([]string, 0)

	for _, experience := range experiences {
		for _, technology := range experience.Technologies {
			if _, exists := seen[technology]; exists {
				continue
			}
			seen[technology] = struct{}{}
			technologies = append(technologies, technology)
		}
	}

	return technologies
}

func loadProfile(path string) (profile.Profile, error) {
	var profile profile.Profile

	data, err := os.ReadFile(path)
	if err != nil {
		return profile, err
	}

	if err := json.Unmarshal(data, &profile); err != nil {
		return profile, err
	}

	if profile.Name == "" {
		return profile, fmt.Errorf("profile name is required")
	}

	return profile, nil
}

func newApp() (*App, error) {
	pages := make(map[string]*template.Template)

	for _, name := range []string{"home", "about", "experience"} {
		tmpl, err := template.New(name).Funcs(template.FuncMap{
			"uniqueTechnologies": uniqueTechnologies,
		}).ParseFiles(
			"templates/layout.html",
			"templates/"+name+".html",
		)
		if err != nil {
			return nil, err
		}

		pages[name] = tmpl
	}

	store := &profileStore{}

	if err := store.Reload("data/profile.json"); err != nil {
		return nil, fmt.Errorf("load initial profile: %w", err)
	}

	return &App{
		pages:        pages,
		profileStore: store,
	}, nil
}

func (a *App) renderPage(
	w http.ResponseWriter,
	name string,
	title string,
) {
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
		a.renderPage(w, "home", "Home")
	})

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		a.renderPage(w, "about", "About")
	})

	mux.HandleFunc("GET /experience", func(w http.ResponseWriter, r *http.Request) {
		a.renderPage(w, "experience", "Experience")
	})

	mux.Handle("GET /static/", http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("static")),
	))

	mux.HandleFunc("GET /my-resume", func(w http.ResponseWriter, r *http.Request) {
		resume := generateResumeMarkdown(a.profileStore.Snapshot())

		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		w.Header().Set(
			"Content-Disposition",
			`attachment; filename="`+strings.ToLower(strings.ReplaceAll(a.profileStore.Snapshot().Name, " ", "-"))+`-resume.md"`)
		w.Header().Set("Cache-Control", "no-cache")

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(resume)); err != nil {
			log.Printf("write resume response: %v", err)
		}
	})

	return mux
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; style-src 'self'; img-src 'self' data:; script-src 'self'; object-src 'none'; base-uri 'self'; frame-ancestors 'none'",
		)

		next.ServeHTTP(w, r)
	})
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		log.Printf(
			"method=%s path=%s duration=%s",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

func main() {
	app, err := newApp()
	if err != nil {
		log.Printf("initialize application: %v", err)
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Application lifecycle: stop on SIGINT or SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	watchCtx, stopWatcher := context.WithCancel(ctx)

	watcherDone, err := watchProfile(watchCtx, "data/profile.json", app.profileStore)
	if err != nil {
		stopWatcher()
		log.Printf("start profile watcher: %v", err)
		return
	}

	// Stop and join the watcher before main exits.
	defer func() {
		stopWatcher()
		<-watcherDone
	}()

	handler := securityHeaders(requestLogger(app.routes()))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("server listening on :%s", port)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server failed: %v", err)
		}

	case <-ctx.Done():
		log.Println("shutdown signal received")
	}

	// Allow active requests up to 10 seconds to finish.
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)

		if closeErr := server.Close(); closeErr != nil {
			log.Printf("forcing server close: %v", closeErr)
		}
	}

	log.Println("server stopped")
}
