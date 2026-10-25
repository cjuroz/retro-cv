package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/cjuroz/retro-cv/internal/profile"
)

func generateResumeMarkdown(p profile.Profile) (string, error) {
	tmpl, err := template.ParseFiles("templates/resume.md.tmpl")
	if err != nil {
		return "", fmt.Errorf("parse resume template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p); err != nil {
		return "", fmt.Errorf("render resume template: %w", err)
	}

	return buf.String(), nil
}
