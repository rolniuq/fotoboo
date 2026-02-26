# FotoBoo - Claude Code Skill File

## Project Overview

FotoBoo is a **photo booth system** built with Go (backend) and vanilla JavaScript (frontend). It enables photo capture, preview, filtering, framing, and export for event photo booths, retail stores, and SaaS applications.

**Module:** `github.com/fotoboo/fotoboo`
**Go Version:** 1.25.3
**External Dependencies:** `github.com/google/uuid`

---

## Architecture

This project follows **Clean Architecture / Hexagonal Architecture** with strict dependency rules:

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
│   ├── domain/           # Entities, value objects, repository interfaces (ports)
│   ├── usecase/          # Business logic / application services
│   ├── handler/          # HTTP handlers (adapters - inbound)
│   └── repository/       # Data persistence (adapters - outbound)
├── pkg/                  # Shared/public packages (future)
├── web/                  # Frontend SPA
│   ├── index.html        # Single-page app with 4 screens
│   └── static/
│       ├── css/style.css
│       └── js/app.js     # FotoBooApp class
├── data/photos/          # Photo storage (runtime, gitignored)
├── bin/                  # Compiled binaries
└── docs/                 # Documentation
```

---

## Coding Conventions

### Go Backend

- **No external HTTP framework** — use stdlib `net/http` only
- **No ORM** — repositories handle persistence directly
- **Manual dependency injection** in `cmd/api/main.go` (no DI container)
- **Domain errors** use sentinel values (`var ErrPhotoNotFound = errors.New(...)`) in `internal/domain/errors.go`
- **Entity constructors** follow `NewXxx(...)` pattern (e.g., `domain.NewPhoto(filePath)`)
- **Repository interfaces** are defined in the `domain` package (ports), not in `repository`
- **File naming:** snake_case (e.g., `photo_handler.go`, `photo_usecase.go`)
- **Package naming:** single lowercase word matching directory name
- **Logging:** stdlib `log` package (no structured logging yet)
- **UUID generation:** `github.com/google/uuid` for entity IDs
- **Config:** environment variables read in `main.go` with sensible defaults

### Frontend

- **Vanilla JavaScript** — no frameworks, no build step
- **Single class pattern** — `FotoBooApp` in `app.js` manages all state and behavior
- **Screen-based navigation** — show/hide screens via `.active` CSS class
- **Camera API** — `navigator.mediaDevices.getUserMedia` at 1280x720
- **Image processing** — Canvas 2D API with CSS filters
- **No external JS libraries** — keep it dependency-free

### General

- **Keep it simple** — prioritize simplicity and reliability over features
- **No over-engineering** — add abstractions only when needed
- **Domain-first** — start with domain entities and interfaces, then build outward

---

## Key Commands

```bash
# Build the API server
go build -o bin/fotoboo-api ./cmd/api

# Run directly (development)
go run ./cmd/api

# Install/tidy dependencies
go mod tidy

# Run the compiled binary
./bin/fotoboo-api
```

### Environment Variables

| Variable       | Description             | Default         |
|----------------|-------------------------|-----------------|
| `PORT`         | Server listen port      | `8080`          |
| `STORAGE_PATH` | Photo storage directory | `./data/photos` |
| `WEB_DIR`      | Frontend files path     | `./web`         |

---

## API Endpoints

| Method | Path           | Description          |
|--------|----------------|----------------------|
| POST   | `/photos`      | Upload photo (raw body, max 10MB) |
| GET    | `/photos/{id}` | Get photo as JPEG    |
| GET    | `/health`      | Health check         |
| GET    | `/`            | Serve frontend SPA   |

---

## Current State (vs Roadmap)

### Implemented
- **Phase 1 (MVP Core):** Camera capture, backend API (upload + retrieve), local storage, basic UI flow (welcome → capture → preview → result)
- **Phase 2 (UX):** Filters (grayscale, vintage, brightness, contrast), frames (simple, event, party), countdown, flash animation

### Not Yet Implemented
- **Phase 3:** QR code generation (placeholder only), print integration
- **Phase 4:** Database (SQLite/PostgreSQL), background jobs, advanced architecture
- **Phase 5:** Admin dashboard, configuration management
- **Phase 6:** Cloud storage (S3/MinIO), observability, rate limiting
- **Phase 7:** Testing (no tests exist yet)

---

## When Adding New Features

1. **Start with the domain** — define entities and repository interfaces in `internal/domain/`
2. **Implement business logic** — create use case methods in `internal/usecase/`
3. **Add HTTP handlers** — create handler methods in `internal/handler/`
4. **Implement persistence** — create/update repository in `internal/repository/`
5. **Wire it up** — inject dependencies in `cmd/api/main.go`
6. **Update frontend** — add screens/controls in `web/`

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
2. Register route in `cmd/api/main.go` with CORS wrapper
3. CORS pattern: set `Access-Control-Allow-Origin: *` and `Access-Control-Allow-Methods`

---

## Testing Guidelines

- **No tests exist yet** — this is a priority item from the roadmap
- When adding tests:
  - Unit tests for use cases (mock repository interfaces)
  - Integration tests for handlers (use `httptest`)
  - Test files go next to source: `photo_usecase_test.go`
  - Use stdlib `testing` package — no external test frameworks unless necessary

---

## Data & Storage

- Photos stored as `{uuid}.jpg` in `STORAGE_PATH` directory
- Metadata persisted in `metadata.json` (JSON array) alongside photos
- In-memory map (`map[string]*Photo`) provides fast lookups at runtime
- Thread safety via `sync.RWMutex` in repository

---

## Common Pitfalls

- `UploadPhoto` passes `""` as filePath to `NewPhoto` — the repository sets the actual path after saving
- Photo handler returns raw JPEG bytes with hardcoded `Content-Type: image/jpeg`
- Max upload size is 10MB (`http.MaxBytesReader`)
- Frontend sends photo as raw blob body (not multipart form)
- Route matching uses `strings.TrimPrefix` — no path parameter library
