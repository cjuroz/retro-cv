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
	Profile profile.Profile
	Version string
	Commit  string
}
