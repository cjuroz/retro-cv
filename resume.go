package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	textTmpl "text/template"
	"time"

	"github.com/cjuroz/retro-cv/internal/profile"
	"github.com/yuin/goldmark"
)

type resumePDFCache struct {
	mu    sync.Mutex
	key   [32]byte
	pdf   []byte
	valid bool
}

func (a *App) cachedResumePDF(markdown string) ([]byte, error) {
	key := sha256.Sum256([]byte(markdown))

	a.pdfCache.mu.Lock()
	defer a.pdfCache.mu.Unlock()

	// Return the cached PDF if the resume hasn't changed.
	if a.pdfCache.valid && a.pdfCache.key == key {
		return a.pdfCache.pdf, nil
	}

	// Generate a new PDF only when the content changes.
	pdf, err := generateResumePDF(markdown)
	if err != nil {
		return nil, err
	}

	a.pdfCache.key = key
	a.pdfCache.pdf = pdf
	a.pdfCache.valid = true

	return pdf, nil
}

func generateResumeMarkdown(p profile.Profile) (string, error) {
	tmpl, err := textTmpl.ParseFiles("templates/resume.md.tmpl")
	if err != nil {
		return "", fmt.Errorf("parse resume template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p); err != nil {
		return "", fmt.Errorf("render resume template: %w", err)
	}

	return buf.String(), nil
}

const resumeHTML = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
  @page {
    size: A4;
    margin: 16mm 17mm;
  }

  body {
    font-family: "Noto Sans", "Noto Sans CJK TC",
                 "Arial", sans-serif;
    font-size: 10pt;
    line-height: 1.45;
    color: #202820;
  }

  h1 {
    font-size: 25pt;
    margin-bottom: 4pt;
  }

  h2 {
    font-size: 16pt;
    border-bottom: 1px solid #aab5aa;
    padding-bottom: 4pt;
    margin-top: 20pt;
  }

  h3 {
    font-size: 12pt;
    margin-bottom: 4pt;
    break-after: avoid;
  }

  p {
    margin: 5pt 0;
  }

  ul {
    padding-left: 18pt;
  }

  li {
    margin-bottom: 3pt;
  }

  a {
    color: #245d3b;
    overflow-wrap: anywhere;
  }

  strong {
    color: #17271b;
  }

  @media print {
    h2, h3 {
      break-after: avoid;
    }

    a {
      text-decoration: none;
    }
  }
</style>
</head>
<body>
{{.Content}}
</body>
</html>`

var resumeHTMLTemplate = template.Must(
	template.New("resume").Parse(resumeHTML),
)

func markdownToHTML(markdown string) (string, error) {
	var rendered bytes.Buffer

	if err := goldmark.Convert([]byte(markdown), &rendered); err != nil {
		return "", fmt.Errorf("convert markdown to HTML: %w", err)
	}

	var page bytes.Buffer
	err := resumeHTMLTemplate.Execute(&page, struct{ Content template.HTML }{Content: template.HTML(rendered.String())})
	if err != nil {
		return "", fmt.Errorf("render resume HTML: %w", err)
	}

	return page.String(), nil
}

func generateResumePDF(markdown string) ([]byte, error) {
	htmlContent, err := markdownToHTML(markdown)
	if err != nil {
		return nil, err
	}

	dir, err := os.MkdirTemp("", "retro-cv-pdf-*")
	if err != nil {
		return nil, fmt.Errorf("create temporary directory: %w", err)
	}
	defer os.RemoveAll(dir)

	htmlPath := filepath.Join(dir, "resume.html")
	pdfPath := filepath.Join(dir, "resume.pdf")

	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0600); err != nil {
		return nil, fmt.Errorf("write resume HTML: %w", err)
	}

	chromium, err := findChromium()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		chromium,
		"--headless",
		"--no-sandbox",
		"--disable-gpu",
		"--no-pdf-header-footer",
		"--disable-dev-shm-usage",
		"--no-first-run",
		"--no-zygote",
		"--allow-file-access-from-files",
		"--print-to-pdf="+pdfPath,
		"file://"+htmlPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("Chromium timed out: %w", ctx.Err())
		}
		return nil, fmt.Errorf("Chromium PDF generation failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	pdf, err := os.ReadFile(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("read generated PDF: %w", err)
	}

	if len(pdf) == 0 {
		return nil, fmt.Errorf("Chromium generated an empty PDF")
	}

	return pdf, nil
}

func findChromium() (string, error) {
	for _, name := range []string{
		"chromium-browser",
		"chromium",
	} {
		path, err := exec.LookPath(name)
		if err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("Chromium executable not found")
}
