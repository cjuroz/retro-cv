package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"text/template"
	"unicode"

	"github.com/cjuroz/retro-cv/internal/profile"
	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/pdf"
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
	latinFont, err := os.ReadFile("fonts/NotoSans-Regular.ttf")
	if err != nil {
		return nil, fmt.Errorf("read Noto Sans font: %w", err)
	}
	cjkFont, err := os.ReadFile("fonts/NotoSansTC-Regular.ttf")
	if err != nil {
		return nil, fmt.Errorf("read Noto Sans TC font: %w", err)
	}

	doc := pdftemplate.New(
		pdftemplate.WithPageSize(document.A4),
		pdftemplate.WithMargins(document.UniformEdges(document.Mm(17))),
		pdftemplate.WithFont("NotoSans", latinFont),
		pdftemplate.WithFont("NotoSansTC", cjkFont),
		pdftemplate.WithDefaultFont("NotoSans", 10),
	)

	page := doc.AddPage()
	addTextRow := func(text string, options ...pdftemplate.TextOption) {
		page.AutoRow(func(r *pdftemplate.RowBuilder) {
			r.Col(12, func(c *pdftemplate.ColBuilder) {
				if containsHan(text) {
					options = append(options, pdftemplate.FontFamily("NotoSansTC"))
				}
				c.Text(text, options...)
				c.Spacer(document.Mm(1.5))
			})
		})
	}
	addHeadingRow := func(text string, options ...pdftemplate.TextOption) {
		addTextRow(text, options...)
		page.AutoRow(func(r *pdftemplate.RowBuilder) {
			r.Col(12, func(c *pdftemplate.ColBuilder) {
				c.Spacer(document.Mm(2))
				c.Line()
				c.Spacer(document.Mm(4))
			})
		})
	}
	addRichTextRow := func(build func(*pdftemplate.RichTextBuilder), options ...pdftemplate.TextOption) {
		page.AutoRow(func(r *pdftemplate.RowBuilder) {
			r.Col(12, func(c *pdftemplate.ColBuilder) {
				c.RichText(build, options...)
				c.Spacer(document.Mm(1.5))
			})
		})
	}
	addSpan := func(rt *pdftemplate.RichTextBuilder, text string, options ...pdftemplate.TextOption) {
		if containsHan(text) {
			options = append(options, pdftemplate.FontFamily("NotoSansTC"))
		}
		rt.Span(text, options...)
	}
	addListRow := func(items []string) {
		for _, item := range items {
			item := item
			addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
				addSpan(rt, "•    "+item)
			}, pdftemplate.FontSize(10), pdftemplate.TextIndent(document.Mm(6)))
		}
	}
	addSpacer := func(height float64) {
		page.AutoRow(func(r *pdftemplate.RowBuilder) {
			r.Col(12, func(c *pdftemplate.ColBuilder) {
				c.Spacer(document.Mm(height))
			})
		})
	}

	addHeadingRow(p.Name+" "+p.NameZh, pdftemplate.FontSize(25), pdftemplate.Bold())
	addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
		addSpan(rt, "ONLINE RESUME (GO RETRO-WEBSITE) -> ", pdftemplate.Bold())
		addSpan(rt, p.ResumeURL, pdftemplate.TextColor(pdf.Blue))
	}, pdftemplate.FontSize(9))
	addSpacer(4)

	addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
		addSpan(rt, p.Address+" | ")
		addSpan(rt, p.Email, pdftemplate.TextColor(pdf.Blue))
	}, pdftemplate.FontSize(9))
	addTextRow("Work Status: "+p.WorkStatus, pdftemplate.FontSize(9))
	addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
		addSpan(rt, "LinkedIn: ")
		addSpan(rt, p.Socials.LinkedIn, pdftemplate.TextColor(pdf.Blue))
		addSpan(rt, " | GitHub: ")
		addSpan(rt, p.Socials.GitHub, pdftemplate.TextColor(pdf.Blue))
	}, pdftemplate.FontSize(9))
	addSpacer(6)

	addHeadingRow("PROFESSIONAL SUMMARY", pdftemplate.FontSize(16), pdftemplate.Bold())
	addTextRow(p.Headline, pdftemplate.FontSize(12), pdftemplate.Bold())
	addSpacer(2)
	addTextRow(p.Summary, pdftemplate.FontSize(10))
	addSpacer(6)

	addHeadingRow("CORE EXPERTISE", pdftemplate.FontSize(16), pdftemplate.Bold())
	addListRow(p.Expertises)
	addSpacer(5)

	addHeadingRow("TECHNICAL SKILLS", pdftemplate.FontSize(16), pdftemplate.Bold())
	addListRow(p.Skills)
	addSpacer(5)

	addHeadingRow("WORK EXPERIENCE", pdftemplate.FontSize(16), pdftemplate.Bold())

	for _, exp := range p.Experience {
		addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
			addSpan(rt, exp.Role+" @ "+exp.Company+" | "+exp.Location+" "+exp.Period, pdftemplate.Bold())
			if exp.Website != "" {
				addSpan(rt, " | ")
				addSpan(rt, exp.Website, pdftemplate.Bold(), pdftemplate.TextColor(pdf.Blue))
			}
		}, pdftemplate.FontSize(12))
		addSpacer(4)
		addTextRow(exp.Summary, pdftemplate.FontSize(10))
		addSpacer(4)
		addListRow(exp.Tasks)
		addSpacer(4)
		if len(exp.Technologies) > 0 {
			addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
				addSpan(rt, "Technologies: ", pdftemplate.Bold())
				addSpan(rt, joinTechnologies(exp.Technologies))
			}, pdftemplate.FontSize(9))
		}
		addSpacer(6)
	}

	addHeadingRow("EDUCATION & LANGUAGES", pdftemplate.FontSize(16), pdftemplate.Bold())
	for _, education := range p.Education {
		addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
			addSpan(rt, education.Degree, pdftemplate.Bold())
			addSpan(rt, " "+education.Period)
		}, pdftemplate.FontSize(10))
		addTextRow(education.Institution, pdftemplate.FontSize(10))
		addSpacer(5)
	}
	if len(p.Languages) > 0 {
		addRichTextRow(func(rt *pdftemplate.RichTextBuilder) {
			for i, language := range p.Languages {
				if i > 0 {
					addSpan(rt, " | ")
				}
				addSpan(rt, language.Language, pdftemplate.Bold())
				addSpan(rt, ": "+language.Level)
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

func containsHan(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}
