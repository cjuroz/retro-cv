# Stage 1: Build the Go application.
FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

ARG APP_VERSION=dev
ARG APP_COMMIT=unknown

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${APP_VERSION} -X main.commit=${APP_COMMIT}" \
    -o /out/retro-cv \
    .

# Stage 2: Run the application.
FROM alpine:3.22

RUN apk add --no-cache ca-certificates curl \
    && addgroup -S -g 1001 app \
    && adduser -S -u 1001 -G app app

WORKDIR /app

COPY --from=builder --chown=app:app /out/retro-cv ./retro-cv
COPY --from=builder --chown=app:app /src/templates ./templates
COPY --from=builder --chown=app:app /src/static ./static
COPY --from=builder --chown=app:app /src/data ./data

ENV PORT=8080

USER 1001:1001

EXPOSE 8080

ENTRYPOINT ["./retro-cv"]
