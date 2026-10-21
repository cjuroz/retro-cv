# Retro CV

A retro-styled, server-rendered resume website built with Go's standard library. Profile content is loaded from JSON and rendered into the home, about, and experience pages.

## Requirements

- Go 1.27.1 or later
- Docker with the Compose plugin (for containerized development)
- Kubernetes and `kubectl` (optional, for cluster deployment)

## Run locally

Run the server from the project root so it can find the templates, static assets, and profile data:

```sh
go run .
```

The site listens on <http://localhost:8080> by default. Set `PORT` to use another port:

```sh
PORT=3000 go run .
```

On PowerShell, set the environment variable with `$env:PORT = "3000"` before running `go run .`.

## Pages and endpoints

| Path | Description |
| --- | --- |
| `/` | Home page |
| `/about` | About page |
| `/experience` | Professional experience |
| `/static/` | CSS, JavaScript, and other static assets |
| `/health` | Liveness check; returns `ok` |
| `/ready` | Readiness check; returns `ready` |

## Update the resume

Edit [`data/profile.json`](data/profile.json) to update the profile displayed by the site. Its top-level fields are `name`, `headline`, `summary`, and `experience`. Each experience entry has `company`, `role`, `period`, `description`, and `technologies` fields. The JSON profile must include a non-empty `name`.

The HTML layout and pages are in `templates/`; styles and browser scripts are in `static/`.

## Tests

Run the Go test suite from the project root:

```sh
go test ./...
```

## Docker Compose

Build and start the app in the background:

```sh
docker compose up --build -d
```

Open <http://localhost:8080>. Compose binds the service to `127.0.0.1` on the host and checks `/health`. Stop it with:

```sh
docker compose down
```

To build and run the image without Compose:

```sh
docker build -t retro-cv:local .
docker run --rm -p 8080:8080 retro-cv:local
```

## Kubernetes

The manifests in `k8s/` define a two-replica Deployment and an internal ClusterIP Service. Build or load the image `retro-cv:local` into the cluster first, then apply the manifests:

```sh
docker build -t retro-cv:local .
kubectl apply -f k8s/
kubectl port-forward service/retro-cv 8080:80
```

Open <http://localhost:8080>. The Deployment configures liveness and readiness probes, runs as a non-root user, and sets CPU and memory requests and limits.

## Project layout

```text
.
├── data/             # Resume profile JSON
├── internal/profile/ # Profile data types
├── k8s/              # Kubernetes Deployment and Service
├── static/           # Stylesheets and browser scripts
├── templates/        # Server-rendered HTML pages
├── main.go           # HTTP server, routes, and template rendering
└── main_test.go      # HTTP route and profile tests
```

The server reads `data/profile.json` and the templates from the working directory at startup. Run the binary from the project root, or from the `/app` working directory used by the provided container image.