# urlshortner
# Go URL Shortener with SQLite and Prometheus

A simple URL shortener written in Go, using SQLite for storage and Prometheus for metrics.

## Features
- Shorten URLs
- Redirect using short codes
- Prometheus metrics exposed at `/metrics`
- SQLite persistence
- Unit tests
- GitHub Actions CI

## Run Locally
```bash
go mod tidy
go run main.go

API Endpoints
POST /shorten — Shorten a URL
Request: {"original_url":"https://example.com"}

GET /r/{shortCode} — Redirect to original URL

GET /metrics — Prometheus metrics



