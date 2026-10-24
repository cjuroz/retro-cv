package main

import (
	"html/template"

	"github.com/cjuroz/retro-cv/internal/profile"
)

type App struct {
	pages        map[string]*template.Template
	profileStore *profileStore
}

type PageData struct {
	Title   string
	SEO     SEOData
	Profile profile.Profile
	Version string
	Commit  string
}

type SEOData struct {
	Description string
	Canonical   string
	OGTitle     string
	OGImage     string
}
