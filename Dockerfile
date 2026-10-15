# Stage 1: Build the Go application.
FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/retro-cv \
    .

# Stage 2: Run the application.
FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app

COPY --from=builder --chown=app:app /out/retro-cv ./retro-cv
COPY --from=builder --chown=app:app /src/templates ./templates
COPY --from=builder --chown=app:app /src/static ./static
COPY --from=builder --chown=app:app /src/data ./data

ENV PORT=8080

USER app

EXPOSE 8080

ENTRYPOINT ["./retro-cv"]