# Copilot instructions for peer-link-server

This file gives targeted guidance for Copilot (and similar assistants) to operate effectively in this repository.

---

## Build, run, test, and lint (repo-specific)

- Build the server binary:
  - go build -o bin/peer-link-server ./cmd/server
- Run locally (uses config/config.yaml by default):
  - go run ./cmd/server -config ./config/config.yaml
  - or set environment variables (prefix APP_) and run: APP_SERVER_PORT=8080 APP_DATABASE_DSN="..." go run ./cmd/server
- Docker Compose (provides Postgres + LiveKit used in local/dev):
  - cd deployments && docker compose up -d
  - docker compose logs -f app
- Tests:
  - This repository currently has no Go test files. General commands if/when tests are added:
    - Run all tests: go test ./...
    - Run package tests: go test ./internal/service -v
    - Run a single test by name: go test ./internal/service -run TestCreateUser -v
- Lint/format (no repo-provided tooling currently):
  - Format: go fmt ./...
  - Vet: go vet ./...
  - If you add golangci-lint, run: golangci-lint run

---

## High-level architecture (big picture)

- Language & framework: Go + Gin HTTP framework.
- Entrypoint: cmd/server/main.go — loads configuration, initializes logger and DB, wires repositories/services/handlers, and starts HTTP server.
- Layers (clear separation):
  - config/         - Viper-based configuration (supports -config file and APP_* env overrides).
  - internal/       - core application packages:
    - database      - GORM + Postgres initialization and AutoMigrate.
    - repository    - data access interfaces and GORM implementations (UserRepository, RoomRepository).
    - service       - business logic interfaces and implementations (UserService, LiveKitService).
    - handler       - HTTP handlers that map requests to services.
    - router        - Gin router assembly, middleware wiring, and API routes (prefixed /api/v1).
    - middleware    - auth (JWT), logging, recovery, CORS, request ID.
  - pkg/            - shared utilities (errors, response envelope, logger helpers).
- External integrations:
  - LiveKit (real-time/RTC service) via github.com/livekit/server-sdk-go — room creation, token issuing, participant management.
  - Postgres via GORM (database DSN configured via config or APP_DATABASE_DSN).
- Dependency injection: manual wiring in cmd/server (simple constructor-based DI, not automated).

---

## Key conventions and patterns (repo-specific)

- Config / env
  - Viper reads config/config.yaml by default. Environment variables override config using the APP_ prefix. Example: APP_SERVER_PORT, APP_JWT_SECRET, APP_LIVEKIT_API_KEY.
  - Use -config to point to an alternate YAML file when running locally.

- Layered interfaces
  - Services and repositories expose Go interfaces (service.UserService, repository.UserRepository). Handlers depend on interfaces; concrete implementations are constructed in cmd/server.
  - Follow the existing pattern when adding new entities: model -> repository -> service -> handler -> router.

- Error & response handling
  - pkg/errors defines application error types (AppError). pkg/response.Envelope is the canonical API response format. Use response.Success / Error / Fail / Created consistently in handlers.

- Routing & middleware
  - Routes are grouped under /api/v1. Add new resources inside internal/router/router.go.
  - Middleware order is deliberate (request ID -> logger -> recovery -> CORS). JWT middleware (middleware.Auth) reads Authorization Bearer tokens and injects user_id into the Gin context.

- Database
  - database.New auto-runs GORM AutoMigrate for model.User and model.Room. For production, prefer an explicit migration tool and avoid AutoMigrate on startup.

- LiveKit usage
  - LiveKit configuration lives under livekit in config. LiveKitService wraps the SDK and the DB-backed room model. Tokens are built with LiveKit auth access tokens.

- Internationalization and timezones
  - Config examples use TimeZone settings in DSN; be mindful of timestamps and DB timezone behavior if adding time-sensitive features.

---

## Where to look for common tasks

- Add new HTTP endpoints: internal/handler + internal/service + internal/repository + add route in internal/router/router.go
- Add DB models: internal/model, update database.AutoMigrate list and repository implementations
- Change logging: pkg/logger

---

## Existing docs integrated

This file pulls key runtime details from README.md, config/config.yaml, and deployments/docker-compose.yml (database and livekit local setup).

---

If this file already exists, update it rather than replacing; prefer minimal edits that preserve hand-written notes.

---

Summary

- Provides concrete build/run/test commands, describes the layered Go architecture, lists repo-specific conventions (config, DI, error/response pattern, routing/middleware, LiveKit integration), and points to where to make common changes.

If you'd like adjustments (more examples for running in containers, CI commands, or adding lint/test automation), say which area to expand.