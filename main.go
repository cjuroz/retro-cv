package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/cjuroz/retro-cv/internal/cv"
)

type PageData struct {
	Title   string
	Profile cv.Profile
}

type App struct {
	pages   map[string]*template.Template
	profile cv.Profile
}

func newApp() (*App, error) {
	pages := make(map[string]*template.Template)

	for _, name := range []string{"home", "about", "experience"} {
		tmpl, err := template.ParseFiles(
			"templates/layout.html",
			"templates/"+name+".html",
		)
		if err != nil {
			return nil, err
		}

		pages[name] = tmpl
	}

	profile := cv.Profile{
		Name:     "Cristian Uroz",
		Headline: "Staff Engineer",
		Summary:  "I'm an engineer interested in building reliable systems, solving complex problems, and helping teams deliver better software.",
		Experience: []cv.Experience{
			{
				Company:      "Company One",
				Role:         "Senior Software Engineer",
				Period:       "2020 - 2024",
				Description:  "Replace this example with a real achievement from your career.",
				Technologies: []string{"Go", "Linux", "AWS"},
			},
			{
				Company:      "Company Two",
				Role:         "Software Engineer",
				Period:       "2016 - 2020",
				Description:  "Replace this example with a real project or responsibility.",
				Technologies: []string{"Networking", "Docker"},
			},
		},
	}

	return &App{
		pages:   pages,
		profile: profile,
	}, nil
}

func (a *App) renderPage(
	w http.ResponseWriter,
	name string,
	title string,
) {
	data := PageData{
		Title:   title,
		Profile: a.profile,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := a.pages[name].ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		log.Printf("render %s page: %v", name, err)
	}
}

func (a *App) routes() *http.ServeMux {
	mux := http.NewServeMux()

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

	return mux
}

func main() {
	app, err := newApp()
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: app.routes(),
	}

	log.Println("Server running on http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
