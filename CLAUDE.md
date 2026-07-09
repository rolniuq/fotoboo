# FotoBoo — Claude Code Skill File

## Project Overview

FotoBoo is a **photo booth system** built with Go (Clean Architecture) and Vue 3.  
It enables photo capture, preview, filtering, framing, collage creation, and export for event photo booths, retail stores, and SaaS applications.

**Module:** `github.com/fotoboo/fotoboo`  
**Go Version:** 1.25.3  
**Frontend:** Vue 3 + Vite + Pinia + Vue Router  

---

## Architecture

Clean Architecture / Hexagonal Architecture with strict dependency inversion:

```
domain (entities + interfaces)
   ↑
usecase (business logic)
   ↑
handler (HTTP transport) + repository (persistence)
```

**Dependency flow is always inward** — outer layers depend on inner layers, never the reverse.

### Directory Structure

```
fotoboo/
├── cmd/api/              # Application entrypoint (main.go)
├── internal/
│   ├── domain/           # Entities, repository interfaces (ports)
│   ├── usecase/          # Business logic / application services
│   ├── handler/          # HTTP handlers (inbound adapters)
│   ├── repository/       # SQLite persistence (outbound adapters)
│   ├── middleware/        # HTTP middleware (logging, metrics, rate limiter)
│   └── background/       # Background jobs (photo cleanup)
├── pkg/storage/          # Storage abstraction (local, MinIO/S3)
├── web/                  # Vue 3 SPA
│   ├── src/              # Components, stores, composables
│   ├── public/           # Static assets (favicon, manifest, OG image)
│   └── index.html
├── data/photos/          # Photo storage (runtime, gitignored)
├── docs/                 # Documentation
├── deploy/               # Deployment scripts
└── test-cases/           # QA test plans and reports
```

---

## Coding Conventions

### Go Backend

- **No external HTTP framework** — use stdlib `net/http` only
- **No ORM** — repositories handle SQLite directly
- **Manual dependency injection** in `cmd/api/main.go` (no DI container)
- **Domain errors** use sentinel values (`var ErrPhotoNotFound = errors.New(...)`) in `internal/domain/errors.go`
- **Entity constructors** follow `NewXxx(...)` pattern (e.g., `domain.NewPhoto(sessionID, filePath)`)
- **Repository interfaces** are defined in the `domain` package (ports), not in `repository`
- **File naming:** snake_case (e.g., `photo_handler.go`, `photo_usecase.go`)
- **Package naming:** single lowercase word matching directory name
- **Logging:** `log/slog` with structured JSON output
- **UUID generation:** `github.com/google/uuid` for entity IDs
- **Config:** environment variables read in `main.go` with sensible defaults
- **CORS:** use `middleware.WithCORS()` when registering routes in `main.go`
- **JSON responses:** use shared `writeJSON(w, status, data)` from the `handler` package

### Frontend

- **Vue 3** with Composition API (`<script setup>`)
- **Pinia** for state management (`useBoothStore`)
- **Vue Router** for screen navigation (booth + admin)
- **Camera API** — `navigator.mediaDevices.getUserMedia` at 1280x720
- **Image processing** — Canvas 2D API with CSS filters
- **Vite** as build tool with API proxy for development

### General

- **Keep it simple** — prioritize simplicity and reliability over features
- **No over-engineering** — add abstractions only when needed
- **Domain-first** — start with domain entities and interfaces, then build outward
- **Test-first** when fixing bugs: reproduce with a test, then fix

---

## Key Commands

```bash
# Build the API server
go build -o bin/fotoboo-api ./cmd/api

# Run directly (development)
go run ./cmd/api

# Run all tests
go test ./... -count=1

# Run tests with verbose output
go test ./... -v -count=1

# Tidy dependencies
go mod tidy

# Frontend dev
cd web && npm install && npm run dev

# Frontend build
cd web && npm run build
```

### Environment Variables

| Variable | Description | Default |
|---|---|---|
| `PORT` | Server listen port | `8080` |
| `STORAGE_PATH` | Photo storage directory | `./data/photos` |
| `DB_PATH` | SQLite database path | `./data/fotoboo.db` |
| `WEB_DIR` | Frontend files path | `./web` |
| `USE_MINIO` | Enable MinIO/S3 storage | `false` |

---

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/photos` | Upload photo (raw body, max 10MB) |
| GET | `/photos` | List photos |
| GET | `/photos/{id}` | Get photo as JPEG |
| GET | `/photos/{id}/qr` | QR code for photo |
| GET | `/photos/{id}/print?size=4x6` | Print-ready JPEG |
| DELETE | `/photos/{id}` | Delete photo |
| POST | `/sessions` | Start session |
| GET | `/sessions` | List sessions |
| GET | `/sessions/{id}` | Get session |
| POST | `/sessions/{id}/complete` | Complete session |
| GET | `/sessions/{id}/photos` | Session photos |
| POST | `/devices` | Register device |
| GET | `/devices` | List devices |
| GET | `/devices/{id}` | Get device |
| PUT | `/devices/{id}` | Update device |
| DELETE | `/devices/{id}` | Delete device |
| GET | `/print-sizes` | List print sizes |
| GET | `/qr?text=...` | Generate QR code |
| GET | `/admin/stats` | Dashboard stats |
| GET | `/admin/config` | Get config |
| PUT | `/admin/config` | Update config |
| GET | `/metrics` | Server metrics |
| GET | `/health` | Health check |
| GET | `/` | Serve frontend SPA |

---

## Current State

All 7 phases (0–7) are **100% complete**. See [ROADMAP.md](./ROADMAP.md) for details.

Key accomplishments:
- Full Clean Architecture implementation (domain → usecase → handler + repository)
- Vue 3 SPA with booth flow (capture, filters, frames, collages) and admin panel
- SQLite persistence with WAL mode
- Dual storage backends (local filesystem + MinIO/S3 via Strategy Pattern)
- Structured logging, metrics, rate limiting, CORS middleware
- 150+ automated tests across all layers
- Docker, Render, Oracle Always Free deployment paths

---

## When Adding New Features

1. **Start with the domain** — define entities and repository interfaces in `internal/domain/`
2. **Implement business logic** — create use case methods in `internal/usecase/`
3. **Write use case tests first** — mock repository interfaces
4. **Add HTTP handlers** — create handler methods in `internal/handler/`
5. **Write handler tests** — use `httptest` with mock repositories
6. **Implement persistence** — create/update repository in `internal/repository/`
7. **Wire it up** — inject dependencies in `cmd/api/main.go`, use `middleware.WithCORS()`
8. **Update frontend** — add screens/controls in `web/`

### Adding a New Domain Entity

```go
// internal/domain/session.go
type Session struct {
    ID        string    `json:"id"`
    CreatedAt time.Time `json:"created_at"`
}

type SessionRepository interface {
    Save(session *Session) error
    FindByID(id string) (*Session, error)
}
```

### Adding a New API Endpoint

1. Add handler method in `internal/handler/`
2. Register route in `cmd/api/main.go` with `middleware.WithCORS()`
3. Use shared `writeJSON()` for JSON responses
4. Add tests in a `*_test.go` file in the same package

---

## Testing

All tests are written using stdlib `testing` + `testify` suite/assert.

**Pattern:** Each handler or use case test file defines its own mocks (same `package xxx_test`).

```go
// Example: handler_test package shares mocks across test files
// MockPhotoRepository, MockSessionRepository, MockDeviceRepository
// are defined once in their respective test files and shared within the package
```

Run tests:
```bash
go test ./... -count=1
```

---

## Data & Storage

- Photos stored as `photos/{uuid}.jpg` in the storage backend
- Metadata persisted in SQLite tables (photos, sessions, devices)
- Storage backend selected at startup (local or MinIO)
- Background job cleans up photos older than 30 days

---

## Common Pitfalls

- `UploadPhoto` passes `""` as filePath to `NewPhoto` — the repository sets the actual path after saving
- Photo handler returns raw JPEG bytes with `Content-Type: image/jpeg`
- Max upload size is 10MB (`http.MaxBytesReader`)
- Frontend sends photo as raw blob body (not multipart form)
- Route matching uses `strings.TrimPrefix` — no path parameter library
- Use `middleware.WithCORS()` not raw `enableCORS()` for new route handlers
- Tests in `handler_test` package share mock types — don't redefine them
