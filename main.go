package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cjuroz/retro-cv/internal/profile"
)

var (
	version = "dev"
	commit  = "unknown"
)

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
