# TrackStack

A personal observability platform for tracking calories, macros, and life metrics with a "Lazy UX" design philosophy.

## Architecture

TrackStack implements the **"Schrödinger's Binary"** pattern - a single Go binary that can be deployed as either a Monolith or Microservices purely through configuration, without code changes.

### Project Structure

```
trackstack/
├── cmd/
│   └── monolith/          # Main entry point for monolithic deployment
│       └── main.go
├── internal/
│   ├── common/            # Shared utilities and infrastructure
│   │   ├── server/        # HTTP server with routing
│   │   └── db/            # Database layer (SQLite)
│   ├── calories/          # Calorie tracking domain module
│   └── ingest/            # AI-powered meal data ingestion
├── go.mod
├── go.sum
└── .golangci.yml          # Linter configuration
```

## Technology Stack

- **Language:** Go 1.22+
- **Web Framework:** Standard Library `net/http`
- **Templating:** Templ (Type-safe HTML)
- **Interactivity:** HTMX
- **Styling:** Tailwind CSS
- **Database:** SQLite 3 (CGO enabled)
- **ORM/Query:** SQLC (Type-safe SQL)
- **Infrastructure:** Terraform + Docker Compose + Caddy

## Getting Started

### Prerequisites

- Go 1.22+
- mise (for version management)

### Installation

```bash
# Install Go 1.22 via mise
mise use golang@1.22

# Build the project
go build -o bin/trackstack ./cmd/monolith

# Run the server
./bin/trackstack

# Server starts on http://localhost:8080
```

### Health Check

```bash
curl http://localhost:8080/health
# Returns: {"status":"ok"}
```

## Development

### Linting

```bash
golangci-lint run ./...
```

### Testing

```bash
go test ./...
```

## Architecture Principles

1. **Strict Modular Monolith:** Modules (`calories`, `ingest`) are isolated with clear boundaries
2. **Interface-based DI:** Modules communicate through interfaces, not concrete types
3. **No Cross-Module Imports:** `internal/calories` cannot import `internal/ingest`
4. **Co-located Tests:** Test files live next to the code they test

## Status

✅ **Story 1.1 Complete:** Project skeleton initialized with basic health check endpoint

## License

Private project - All rights reserved