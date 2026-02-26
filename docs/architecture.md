# Architecture

This document describes the system architecture of FotoBoo, a photo booth platform built with Go and vanilla JavaScript following Clean Architecture principles.

---

## High-Level Overview

```
┌─────────────────────────────────┐
│   Client (Browser / Touch UI)   │
│   - Camera capture (WebRTC)     │
│   - Image processing (Canvas)   │
│   - UI navigation               │
└──────────────┬──────────────────┘
               │ HTTP (REST)
               ▼
┌─────────────────────────────────┐
│        Go Backend API           │
│   - Photo upload & retrieval    │
│   - Static file serving         │
│   - Health check                │
└──────────────┬──────────────────┘
               │
               ▼
┌─────────────────────────────────┐
│       Storage Layer             │
│   - Local filesystem (photos)   │
│   - metadata.json (index)       │
└─────────────────────────────────┘
```

---

## Clean Architecture

FotoBoo follows **Clean Architecture** (also known as Hexagonal Architecture / Ports & Adapters). The key principle is that **dependencies always point inward** — outer layers know about inner layers, but inner layers have no knowledge of outer layers.

### Layer Diagram

```
┌──────────────────────────────────────────────────┐
│                  cmd/api/main.go                  │  ← Composition Root
│              (wiring & configuration)             │
├──────────────────────────────────────────────────┤
│                                                  │
│  ┌────────────────┐      ┌────────────────────┐  │
│  │    handler/     │      │    repository/     │  │  ← Adapters (outer)
│  │  (HTTP in)      │      │  (persistence out) │  │
│  └───────┬────────┘      └────────┬───────────┘  │
│          │                        │              │
│          ▼                        ▼              │
│  ┌────────────────────────────────────────────┐  │
│  │              usecase/                       │  │  ← Application Logic
│  │         (business operations)               │  │
│  └────────────────────┬───────────────────────┘  │
│                       │                          │
│                       ▼                          │
│  ┌────────────────────────────────────────────┐  │
│  │              domain/                        │  │  ← Core Domain
│  │   (entities, interfaces, errors)            │  │
│  └────────────────────────────────────────────┘  │
│                                                  │
└──────────────────────────────────────────────────┘
```

### Layer Responsibilities

| Layer | Package | Responsibility |
|-------|---------|----------------|
| **Domain** | `internal/domain/` | Entities (`Photo`), repository interfaces (`PhotoRepository`), domain errors |
| **Use Case** | `internal/usecase/` | Business logic — orchestrates domain entities and repositories |
| **Handler** | `internal/handler/` | HTTP transport — parses requests, calls use cases, formats responses |
| **Repository** | `internal/repository/` | Data persistence — implements domain interfaces using filesystem |
| **Composition Root** | `cmd/api/` | Wires all layers together via dependency injection |

### Dependency Rule

```
handler → usecase → domain ← repository
```

- `handler` depends on `usecase` (calls business methods)
- `usecase` depends on `domain` (uses entities and repository interfaces)
- `repository` depends on `domain` (implements repository interfaces)
- `domain` depends on **nothing** — it is the innermost layer

---

## Domain Model

### Current Entities

#### Photo

```go
type Photo struct {
    ID        string    `json:"id"`
    FilePath  string    `json:"file_path"`
    CreatedAt time.Time `json:"created_at"`
}
```

### Planned Entities (from Roadmap)

| Entity | Purpose | Phase |
|--------|---------|-------|
| `Session` | Groups photos from one photo-taking flow | Phase 4 |
| `Device` | Tracks camera/booth devices | Phase 4 |

### Repository Interfaces (Ports)

Repository interfaces are defined in the **domain layer**, making them ports that can be implemented by any adapter:

```go
type PhotoRepository interface {
    Save(photo *Photo, data []byte) error
    FindByID(id string) (*Photo, error)
    GetFileData(photo *Photo) ([]byte, error)
}
```

---

## Data Flow

### Photo Upload Flow

```
Browser                    Handler              UseCase              Repository
   │                         │                     │                     │
   │  POST /photos (bytes)   │                     │                     │
   │────────────────────────>│                     │                     │
   │                         │  UploadPhoto(data)  │                     │
   │                         │────────────────────>│                     │
   │                         │                     │  NewPhoto("")       │
   │                         │                     │  repo.Save(photo,   │
   │                         │                     │          data)      │
   │                         │                     │────────────────────>│
   │                         │                     │                     │ Write {id}.jpg
   │                         │                     │                     │ Update metadata.json
   │                         │                     │                     │ Update in-memory map
   │                         │                     │<────────────────────│
   │                         │<────────────────────│                     │
   │  201 {id, created_at}   │                     │                     │
   │<────────────────────────│                     │                     │
```

### Photo Retrieval Flow

```
Browser                    Handler              UseCase              Repository
   │                         │                     │                     │
   │  GET /photos/{id}       │                     │                     │
   │────────────────────────>│                     │                     │
   │                         │  GetPhotoData(id)   │                     │
   │                         │────────────────────>│                     │
   │                         │                     │  FindByID(id)       │
   │                         │                     │────────────────────>│
   │                         │                     │<────────────────────│
   │                         │                     │  GetFileData(photo) │
   │                         │                     │────────────────────>│
   │                         │                     │                     │ Read {id}.jpg
   │                         │                     │<────────────────────│
   │                         │<────────────────────│                     │
   │  200 image/jpeg         │                     │                     │
   │<────────────────────────│                     │                     │
```

---

## Storage Architecture

### Current (MVP)

- **Photos:** Stored as `{uuid}.jpg` files in `STORAGE_PATH` directory
- **Metadata:** JSON array in `metadata.json` file alongside photos
- **Runtime index:** In-memory `map[string]*Photo` protected by `sync.RWMutex`

```
data/photos/
├── metadata.json
├── 550e8400-e29b-41d4-a716-446655440000.jpg
├── 6ba7b810-9dad-11d1-80b4-00c04fd430c8.jpg
└── ...
```

### Planned (Production)

| Phase | Storage | Database |
|-------|---------|----------|
| MVP (current) | Local filesystem | metadata.json (flat file) |
| Phase 4 | Local filesystem | SQLite |
| Phase 6 | S3 / MinIO + CDN | PostgreSQL |

---

## Thread Safety

The `FilePhotoRepository` uses `sync.RWMutex` to ensure thread-safe access:

- **Read operations** (`FindByID`, `GetFileData`) acquire a read lock (`RLock`)
- **Write operations** (`Save`) acquire an exclusive write lock (`Lock`)
- The in-memory map and metadata.json are always kept in sync

---

## Frontend Architecture

The frontend is a **single-page application (SPA)** with no build step:

```
web/
├── index.html          # 4 screens, all in one HTML file
└── static/
    ├── css/style.css   # Styling + animations
    └── js/app.js       # FotoBooApp class (all logic)
```

### Screen Flow

```
Welcome → Capture → Preview → Result
            ↑          │
            └──────────┘
              (Retake)
```

See [frontend.md](./frontend.md) for detailed frontend documentation.

---

## Future Architecture Considerations

Based on the roadmap, the architecture will evolve:

1. **Database layer** (Phase 4) — Replace `metadata.json` with SQLite/PostgreSQL. The repository interface stays the same; only the implementation changes.
2. **Background jobs** (Phase 4) — Add async processing for image resizing and cleanup.
3. **Cloud storage** (Phase 6) — Implement S3-compatible `PhotoRepository` adapter.
4. **Observability** (Phase 6) — Add structured logging, metrics, and health checks.
5. **Admin API** (Phase 5) — New handlers and use cases for dashboard data.
