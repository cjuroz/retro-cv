# Retro CV

[![retro-cv](https://github.com/cjuroz/retro-cv/actions/workflows/ci.yaml/badge.svg?branch=master)](https://github.com/cjuroz/retro-cv/actions/workflows/ci.yaml)
[![Go 1.27.1](https://img.shields.io/badge/Go-1.27.1-00ADD8?logo=go&logoColor=white)](https://go.dev/dl/)
[![GHCR image](https://img.shields.io/badge/GHCR-ghcr.io%2Fcjuroz%2Fretro--cv-2496ED?logo=docker&logoColor=white)](https://github.com/cjuroz/retro-cv/pkgs/container/retro-cv)

Retro CV is a retro-styled resume website built and served by Go, with CSS and JavaScript for the browser experience. Go's `net/http` server renders HTML templates using profile data from JSON; there is no separate frontend framework or application server. The site runs as a Docker service and includes a GitHub Actions CI/CD workflow intended to publish container images and deploy to a remote EC2 host over SSH.

## Stack and requirements

- Go 1.27.1 or later for local development
- Go standard library for the HTTP server, routing, and HTML templates
- HTML templates, CSS, and vanilla JavaScript for the site UI
- Docker and the Compose plugin for containerized development and deployment
- Kubernetes and `kubectl` are optional for using the included manifests

## Run locally

Run from the repository root so the server can find its templates, static assets, and profile data:

```sh
go run .
```

Open <http://localhost:8080>. Set `PORT` to use another port, for example `PORT=3000 go run .` (PowerShell: `$env:PORT = "3000"`).

## Pages and endpoints

| Path | Description |
| --- | --- |
| `/` | Home page |
| `/about` | About page |
| `/experience` | Experience page |
| `/my-resume` | Download the profile as Markdown |
| `/static/` | CSS, JavaScript, and other static assets |
| `/health` | Health check; returns `ok` |
| `/ready` | Readiness check; returns `ready` |

## Resume content

Edit [`data/profile.json`](data/profile.json) to update the displayed resume. The profile includes contact and social links, headline and summary, skills, areas of expertise, education, languages, and work experience. Experience entries include role, company, period, location, summary, tasks, and technologies. The profile must have a non-empty `name`.

The running server checks the profile file for changes and reloads valid JSON automatically, so content edits do not require a restart. Keep the JSON valid; invalid updates are logged and the last valid profile remains in use.

## Tests

Run the Go test suite from the repository root:

```sh
go test ./...
```

## Docker Compose

Build and start the service:

```sh
docker compose up --build -d
```

Open <http://localhost:8080>. Compose binds to `127.0.0.1` on the host, mounts `data/` into the container for profile updates, restarts the service automatically, and checks `/health`.

```sh
docker compose down
```

To build and run the image directly:

```sh
docker build -t retro-cv:local .
docker run --rm -p 8080:8080 retro-cv:local
```

The Dockerfile builds a static Go binary in a multi-stage build and runs it in a small Alpine image as a non-root user.

## GitHub Actions CI/CD

The workflow in [`.github/workflows/ci.yaml`](.github/workflows/ci.yaml) runs on pushes and pull requests targeting `master`. Both event types run Go tests with coverage, `go vet`, and the race detector. On a push to `master`, successful tests are followed by:

1. A multi-architecture Docker build for `linux/amd64` and `linux/arm64`.
2. Publishing to GitHub Container Registry (`ghcr.io/<owner>/<repository>`) with `latest`, the value in `VERSION`, and a full commit-SHA tag.
3. SSH deployment to the configured EC2 host. The script updates the checkout in `~/retro-cv`, pulls the commit-tagged image, starts the Compose service, checks `/health`, and attempts to restore the previous image if deployment or health validation fails.

Configure these GitHub Actions repository secrets for the SSH step: `EC2_HOST`, `EC2_USER`, `EC2_SSH_KEY`, and `EC2_HOST_FINGERPRINT`. The EC2 host must have Git access to the repository, Docker with the Compose plugin, and permission to pull the GHCR package (or the package must be public).

**Deployment configuration note:** the current workflow expects the remote Compose configuration to use the published GHCR image and the supplied `IMAGE_TAG`. The checked-in [`compose.yaml`](compose.yaml) currently builds `retro-cv:local` instead and does not use `IMAGE_TAG`, so the remote `docker compose pull` will not pull the image published by this workflow until that Compose configuration is aligned.

## Kubernetes

The manifests in `k8s/` define a two-replica Deployment and an internal ClusterIP Service. They currently use the local image name `retro-cv:local`; build or load that image into the cluster before applying them:

```sh
docker build -t retro-cv:local .
kubectl apply -f k8s/
kubectl port-forward service/retro-cv 8080:80
```

Open <http://localhost:8080>. The Deployment includes liveness and readiness probes, runs as a non-root user, and sets CPU and memory requests and limits.

## Project layout

```text
.
├── .github/workflows/ci.yaml # Test, publish, and remote deployment workflow
├── data/profile.json         # Resume profile content
├── internal/profile/         # Profile model and model tests
├── k8s/                      # Kubernetes Deployment and Service
├── static/css/               # Site styles
├── static/js/                # Browser-side interactions
├── templates/                # Shared layout and page templates
├── app.go                    # Application and page data types
├── handlers.go               # Page rendering and Markdown resume handler
├── main.go                   # Startup, template setup, and graceful shutdown
├── middleware.go             # HTTP middleware
├── profile_store.go          # Profile loading, snapshots, and live reload
├── resume.go                 # Markdown resume generation
├── routes.go                 # HTTP routes
├── compose.yaml              # Local Docker Compose service
├── Dockerfile                # Multi-stage container build
└── *_test.go                 # Go tests
```

The server loads `data/profile.json` and templates relative to its working directory. Run the local binary from the repository root; the container image uses `/app` as its working directory.