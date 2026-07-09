# FotoBoo

A production-ready photo booth system built with Go (Clean Architecture) and Vue 3.  
FotoBoo enables photo capture, filtering, framing, collage creation, and export for event photo booths, retail stores, and SaaS applications.

---

## Features

### Core
- Photo capture via browser camera (1280x720)
- Live preview with CSS filters (grayscale, vintage, brightness, contrast)
- Frame overlays (simple, event, party)
- Photo collage layouts (single, horizontal, vertical, grid, strip, featured)
- Countdown timer + flash animation
- Print-ready export (4x6, 5x7, 6x8, 2x6 strip at 300 DPI)
- QR code generation for instant photo sharing

### Backend
- RESTful API for photo, session, and device management
- **Dual storage backends:** Local filesystem or MinIO/S3 object storage
- SQLite database for metadata persistence
- Sessions with device tracking
- Admin dashboard with usage statistics
- Configurable booth settings (countdown, frames, filters, retention)
- Background cleanup of expired photos
- Rate limiting, metrics, health checks
- CORS support for browser clients

### Deployment
- Docker Compose (local, production, MinIO profiles)
- Render Blueprint (free tier)
- Oracle Always Free VM with automatic HTTPS
- GitHub Actions CI/CD

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

## Project Structure

```
fotoboo/
├── cmd/api/              # Application entrypoint
├── internal/
│   ├── domain/           # Entities, value objects, repository interfaces
│   ├── usecase/          # Business logic / application services
│   ├── handler/          # HTTP handlers (inbound adapters)
│   ├── repository/       # Data persistence (outbound adapters)
│   ├── middleware/        # HTTP middleware (logging, metrics, rate limiter)
│   └── background/       # Background jobs (photo cleanup)
├── pkg/storage/          # Storage abstraction (local, MinIO)
├── web/                  # Vue 3 frontend SPA
│   ├── src/              # Vue components, stores, composables
│   ├── public/           # Static assets (favicon, manifest, OG image)
│   └── index.html
├── deploy/               # Deployment scripts
├── docs/                 # Documentation
└── test-cases/           # QA test plans and reports
```

---

## Quick Start

### Requirements

- Go 1.25+
- Node.js 20+ (for frontend build)

### Backend

```bash
# Build the API server
go build -o bin/fotoboo-api ./cmd/api

# Run directly
go run ./cmd/api

# Run with options
STORAGE_PATH=./data/photos PORT=8080 go run ./cmd/api
```

### Frontend (Development)

```bash
cd web
npm install
npm run dev          # Vite dev server on :5173, proxies API to :8080
```

### Production Frontend Build

```bash
cd web
npm install && npm run build   # Outputs to web/dist/
```

---

## Environment Variables

| Variable | Description | Default |
|---|---|---|
| `PORT` | Server listen port | `8080` |
| `STORAGE_PATH` | Photo storage directory | `./data/photos` |
| `DB_PATH` | SQLite database path | `./data/fotoboo.db` |
| `WEB_DIR` | Frontend files path | `./web` |
| `BASE_URL` | Public base URL | `http://localhost:8080` |
| `USE_MINIO` | Enable MinIO/S3 storage | `false` |
| `MINIO_ENDPOINT` | MinIO server endpoint | `localhost:9000` |
| `MINIO_ACCESS_KEY` | MinIO access key | `minioadmin` |
| `MINIO_SECRET_KEY` | MinIO secret key | `minioadmin` |
| `MINIO_BUCKET` | MinIO bucket name | `fotoboo` |

---

## API Overview

| Method | Path | Description |
|--------|------|-------------|
| POST | `/photos` | Upload a photo (raw body, max 10MB) |
| GET | `/photos` | List all photos |
| GET | `/photos/{id}` | Get photo as JPEG |
| GET | `/photos/{id}/qr` | Get QR code for photo |
| GET | `/photos/{id}/print` | Get print-ready JPEG |
| DELETE | `/photos/{id}` | Delete a photo |
| POST | `/sessions` | Start a new session |
| GET | `/sessions` | List all sessions |
| GET | `/sessions/{id}` | Get session details |
| POST | `/sessions/{id}/complete` | Complete a session |
| GET | `/sessions/{id}/photos` | Get photos in session |
| POST | `/devices` | Register a device |
| GET | `/devices` | List all devices |
| GET | `/devices/{id}` | Get device details |
| PUT | `/devices/{id}` | Update device |
| DELETE | `/devices/{id}` | Delete device |
| GET | `/print-sizes` | List available print sizes |
| GET | `/qr` | Generate QR code for arbitrary text |
| GET | `/admin/stats` | Dashboard statistics |
| GET | `/admin/config` | Get booth configuration |
| PUT | `/admin/config` | Update booth configuration |
| GET | `/metrics` | Server metrics |
| GET | `/health` | Health check |

Full API reference: [docs/api.md](./docs/api.md)

---

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v -count=1

# All 96+ tests pass covering:
# - Domain entities and errors
# - All use cases (photo, session, device, admin)
# - All HTTP handlers (photo, session, device, admin, print, QR)
# - All SQLite repositories (photo, session, device)
# - All middleware (metrics, rate limiter, logger, session limiter)
# - Local storage backend
```

---

## Docker Deployment

```bash
# Local storage mode
docker compose up -d --build

# MinIO mode
USE_MINIO=true docker compose --profile minio up -d --build

# Production mode (with Caddy reverse proxy)
docker compose --profile prod up -d --build
```

---

## Deployment Options

- **Render (Free):** Uses `render.yaml` Blueprint — push to GitHub, deploy in 1 click
- **Railway / Fly.io:** See [PUBLIC_DEPLOYMENT.md](./PUBLIC_DEPLOYMENT.md)
- **Oracle Always Free:** Full guide at [deploy/oracle/README.md](./deploy/oracle/README.md)
- **Quick HTTPS launch:**
  ```bash
  bash deploy/oracle/deploy-nipio.sh
  ```

---

## Documentation

| Document | Description |
|---|---|
| [docs/api.md](./docs/api.md) | Full REST API reference |
| [docs/architecture.md](./docs/architecture.md) | Architecture and data flow |
| [docs/frontend.md](./docs/frontend.md) | Frontend architecture guide |
| [docs/deployment.md](./docs/deployment.md) | Comprehensive deployment guide |
| [docs/getting-started.md](./docs/getting-started.md) | Local development setup |
| [docs/minio-setup.md](./docs/minio-setup.md) | MinIO/S3 storage configuration |
| [docs/roadmap-status.md](./docs/roadmap-status.md) | Current progress against roadmap |
| [docs/storage-refactoring.md](./docs/storage-refactoring.md) | Storage backend architecture |

---

## UI Screens

| Screen | Description |
|--------|-------------|
| Welcome | Start single photo or collage mode |
| Capture | Live camera with countdown + flash |
| Preview | Apply filters, frames, brightness/contrast |
| Collage | Multi-photo layouts with progressive capture |
| Result | Final photo with download, QR, print controls |
| Admin | Dashboard, photo/session/device management, config |

---

## License

MIT
