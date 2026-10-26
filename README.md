# Retro CV

[![retro-cv](https://github.com/cjuroz/retro-cv/actions/workflows/ci.yaml/badge.svg?branch=master)](https://github.com/cjuroz/retro-cv/actions/workflows/ci.yaml)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go 1.27.1](https://img.shields.io/badge/Go-1.27.1-00ADD8?logo=go&logoColor=white)](https://go.dev/dl/)
[![Docker Compose](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)](https://docs.docker.com/compose/)
[![GHCR image](https://img.shields.io/badge/GHCR-ghcr.io%2Fcjuroz%2Fretro--cv-2496ED?logo=docker&logoColor=white)](https://github.com/cjuroz/retro-cv/pkgs/container/retro-cv)

Retro CV is a lightweight Go-based personal resume and portfolio site with a retro terminal aesthetic. It serves a small set of HTML pages, provides Markdown and PDF resume downloads, and reloads profile content from JSON at runtime so updates can be made without restarting the app.

## Features

- Retro-styled resume website with terminal-inspired UI across multiple pages
- Go HTTP server with HTML templates and static assets
- Live profile reload from JSON data files
- Downloadable Markdown and PDF resume exports
- Health and readiness endpoints for Docker and Kubernetes checks
- Container-ready build with Docker Compose support

## Tech stack

- Go 1.27.1+
- gpdf for PDF generation
- Standard library HTTP server and templates
- HTML, CSS, and vanilla JavaScript
- Docker and Docker Compose
- Optional Kubernetes manifests

## Prerequisites

- Go 1.27.1 or newer
- Docker and Docker Compose plugin
- Optional: kubectl for local cluster testing

## Quick start

From the repository root:

```bash
go run .
```

Then open:

- http://localhost:8080

If you want to run on a different port:

```bash
PORT=3000 go run .
```

PowerShell:

```powershell
$env:PORT = "3000"
go run .
```

## Configuration

The resume content is stored in [data/profile.json](data/profile.json). Update the JSON there to change:

- Name, headline, and summary
- Contact details and social links
- Skills, expertise, education, and languages
- Work experience and technologies

The application watches this file and reloads valid JSON automatically. Invalid JSON is ignored and the last known-good profile remains active.

## Routes

| Route | Purpose |
| --- | --- |
| / | Home page |
| /about | About page |
| /experience | Experience page |
| /resume.md | Download resume as Markdown |
| /resume.pdf | Download resume as PDF |
| /health | Health endpoint |
| /ready | Readiness endpoint |
| /static/ | Static CSS/JS assets |

## Testing

Run the test suite from the repository root:

```bash
go test ./...
```

## Docker

Build and start locally:

```bash
docker compose up --build -d
```

Stop it with:

```bash
docker compose down
```

Run the image directly:

```bash
docker build -t retro-cv:local .
docker run --rm -p 8080:8080 retro-cv:local
```

## Kubernetes

The manifests in [k8s](k8s) provide a Deployment and Service. To apply them locally:

```bash
docker build -t retro-cv:local .
kubectl apply -f k8s/
kubectl port-forward service/retro-cv 8080:80
```

Then open http://localhost:8080.

## Project structure

```text
.
├── .github/
│   └── workflows/
│       └── ci.yaml             # CI/CD workflow for tests and image publishing
├── data/
│   └── profile.json            # Resume content used by the app
├── fonts/
│   ├── NotoSans-Regular.ttf    # Default PDF font
│   └── NotoSansTC-Regular.ttf  # CJK PDF font
├── internal/
│   └── profile/                # Profile model and validation logic
├── k8s/
│   ├── deployment.yaml
│   └── service.yaml            # Kubernetes deployment manifests
├── static/
│   ├── css/
│   ├── img/
│   └── js/
├── templates/
│   ├── about.html
│   ├── experience.html
│   ├── home.html
│   ├── layout.html
│   └── resume.md.tmpl          # Markdown resume template
├── app.go                      # Core app and page model definitions
├── compose.yaml                # Docker Compose configuration
├── Dockerfile                  # Multi-stage container build
├── go.mod                      # Go module definition
├── handlers.go                 # Page rendering handlers
├── LICENSE                     # Project license
├── main.go                     # Server startup and lifecycle
├── middleware.go               # HTTP middleware
├── profile_store.go            # Profile loader and live reload logic
├── resume.go                   # Markdown and PDF resume generation
├── routes.go                   # Route registration
├── main_test.go
├── profile_store_test.go
├── resume_test.go
├── README.md
└── retro-cv.code-workspace
```

## Notes

- The app expects to run from the repository root so it can locate templates, static assets, and the profile data file.
- The container image uses a non-root user and a minimal runtime image.
- The included CI workflow publishes container images and can deploy via SSH to a remote host.

