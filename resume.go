package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"text/template"

	"github.com/cjuroz/retro-cv/internal/profile"
	"github.com/gpdf-dev/gpdf/document"
	pdftemplate "github.com/gpdf-dev/gpdf/template"
)

type resumePDFCache struct {
	mu    sync.Mutex
	key   [32]byte
	pdf   []byte
	valid bool
}

func (a *App) cachedResumePDF(p profile.Profile) ([]byte, error) {
	profileJSON, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("marshal profile for PDF cache: %w", err)
	}

	key := sha256.Sum256(profileJSON)

	a.pdfCache.mu.Lock()
	defer a.pdfCache.mu.Unlock()

	if a.pdfCache.valid && a.pdfCache.key == key {
		return a.pdfCache.pdf, nil
	}

	pdfBytes, err := generateResumePDF(p)
	if err != nil {
		return nil, err
	}

	a.pdfCache.key = key
	a.pdfCache.pdf = pdfBytes
	a.pdfCache.valid = true

	return pdfBytes, nil
}

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

func generateResumePDF(p profile.Profile) ([]byte, error) {
	cjkFont, err := os.ReadFile("fonts/NotoSans-Regular.ttf")
	if err != nil {
		return nil, fmt.Errorf("read Noto Sans TC font: %w", err)
	}

	doc := pdftemplate.New(
		pdftemplate.WithPageSize(document.A4),
		pdftemplate.WithMargins(document.UniformEdges(document.Mm(17))),
		pdftemplate.WithFont("NotoSansTC", cjkFont),
		pdftemplate.WithDefaultFont("NotoSansTC", 10),
	)

	page := doc.AddPage()
	addTextRow := func(text string, options ...pdftemplate.TextOption) {
		page.AutoRow(func(r *pdftemplate.RowBuilder) {
			r.Col(12, func(c *pdftemplate.ColBuilder) {
				c.Text(text, options...)
			})
		})
	}
	addRichTextRow := func(build func(*pdftemplate.RichTextBuilder), options ...pdftemplate.TextOption) {
		page.AutoRow(func(r *pdftemplate.RowBuilder) {
			r.Col(12, func(c *pdftemplate.ColBuilder) {
				c.RichText(build, options...)
			})
		})
	}
	addListRow := func(items []string) {
		if len(items) == 0 {
			return
		}
		page.AutoRow(func(r *pdftemplate.RowBuilder) {
			r.Col(12, func(c *pdftemplate.ColBuilder) {
				c.List(items)
			})
		})
	}
	addSpacer := func(height float64) {
		page.AutoRow(func(r *pdftemplate.RowBuilder) {
			r.Col(12, func(c *pdftemplate.ColBuilder) {
				c.Spacer(document.Mm(height))
			})
		})
	}

	addTextRow(p.Name+" "+p.NameZh, pdftemplate.FontSize(25), pdftemplate.Bold())
	addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
		rt.Span("ONLINE RESUME (GO RETRO-WEBSITE) -> ", pdftemplate.Bold())
		rt.Span(p.ResumeURL)
	}, pdftemplate.FontSize(9))
	addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
		rt.Span(p.Address + " | " + p.Email)
	}, pdftemplate.FontSize(9))
	addTextRow("Work Status: "+p.WorkStatus, pdftemplate.FontSize(9))
	addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
		rt.Span("LinkedIn: " + p.Socials.LinkedIn + " | GitHub: " + p.Socials.GitHub)
	}, pdftemplate.FontSize(9))
	addSpacer(5)

	addTextRow("PROFESSIONAL SUMMARY", pdftemplate.FontSize(16), pdftemplate.Bold())
	addTextRow(p.Headline, pdftemplate.FontSize(12), pdftemplate.Bold())
	addTextRow(p.Summary, pdftemplate.FontSize(10))
	addSpacer(4)

	addTextRow("CORE EXPERTISE", pdftemplate.FontSize(16), pdftemplate.Bold())
	addListRow(p.Expertises)
	addSpacer(3)

	addTextRow("TECHNICAL SKILLS", pdftemplate.FontSize(16), pdftemplate.Bold())
	addListRow(p.Skills)
	addSpacer(3)

	addTextRow("WORK EXPERIENCE", pdftemplate.FontSize(16), pdftemplate.Bold())

	for _, exp := range p.Experience {
		heading := exp.Role + " @ " + exp.Company + " | " + exp.Location + " " + exp.Period
		if exp.Website != "" {
			heading += " | " + exp.Website
		}
		addTextRow(heading, pdftemplate.FontSize(12), pdftemplate.Bold())
		addTextRow(exp.Summary, pdftemplate.FontSize(10))
		addListRow(exp.Tasks)
		if len(exp.Technologies) > 0 {
			addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
				rt.Span("Technologies: ", pdftemplate.Bold())
				rt.Span(joinTechnologies(exp.Technologies))
			}, pdftemplate.FontSize(9))
		}
		addSpacer(4)
	}

	addTextRow("EDUCATION & LANGUAGES", pdftemplate.FontSize(16), pdftemplate.Bold())
	for _, education := range p.Education {
		addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
			rt.Span(education.Degree, pdftemplate.Bold())
			rt.Span(" " + education.Period)
		}, pdftemplate.FontSize(10))
		addTextRow(education.Institution, pdftemplate.FontSize(10))
		addSpacer(3)
	}
	if len(p.Languages) > 0 {
		addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
			for i, language := range p.Languages {
				if i > 0 {
					rt.Span(" | ")
				}
				rt.Span(language.Language, pdftemplate.Bold())
				rt.Span(": " + language.Level)
			}
		}, pdftemplate.FontSize(10))
	}

	pdfBytes, err := doc.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate resume PDF: %w", err)
	}

	return pdfBytes, nil
}

func joinTechnologies(items []string) string {
	result := ""
	for i, item := range items {
		if i > 0 {
			result += ", "
		}
		result += item
	}
	return result
}
