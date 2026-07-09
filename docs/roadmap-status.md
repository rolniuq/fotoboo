# Roadmap Status

Current implementation progress tracked against the [ROADMAP.md](../ROADMAP.md).

---

## Summary

| Phase | Name | Status | Progress |
|-------|------|--------|----------|
| Phase 0 | Planning & Foundation | **Done** | 100% |
| Phase 1 | MVP Core Features | **Done** | 100% |
| Phase 2 | UX & Photo Enhancement | **Done** | 100% |
| Phase 3 | Output & Sharing | **Done** | 100% |
| Phase 4 | Advanced Backend | **Done** | 100% |
| Phase 5 | Admin & Operations | **Done** | 100% |
| Phase 6 | Production & Scaling | **Done** | 100% |
| Phase 7 | Testing & Quality | **Done** | 100% |

---

## Detailed Status

### Phase 0 – Planning & Foundation

| Task | Status | Notes |
|------|--------|-------|
| Define MVP scope | ✅ Done | Documented in ROADMAP.md |
| Choose tech stack | ✅ Done | Go + Vue.js + SQLite + local/MinIO storage |
| Set up project structure | ✅ Done | Clean Architecture directories |

### Phase 1 – MVP Core Features

| Task | Status | Notes |
|------|--------|-------|
| Connect to webcam | ✅ Done | `getUserMedia` at 1280x720 |
| Capture photo | ✅ Done | Canvas capture + countdown |
| Show preview | ✅ Done | Preview screen with captured image |
| Retake photo | ✅ Done | Navigates back to capture screen |
| `POST /photos` – upload | ✅ Done | Raw bytes, 10MB limit, returns UUID |
| `GET /photos/{id}` – retrieve | ✅ Done | Returns JPEG with proper content type |
| Photo domain entity | ✅ Done | `internal/domain/photo.go` |
| Session domain entity | ✅ Done | `internal/domain/session.go` |
| Device domain entity | ✅ Done | `internal/domain/device.go` |
| Local file storage | ✅ Done | UUID-based naming in `data/photos/` |
| MinIO/S3 storage | ✅ Done | Strategy pattern via `pkg/storage` interface |
| Welcome screen | ✅ Done | Single photo and collage mode options |
| Capture screen | ✅ Done | |
| Preview screen | ✅ Done | |
| Result screen | ✅ Done | Download, QR, print controls |

### Phase 2 – UX & Photo Enhancement

| Task | Status | Notes |
|------|--------|-------|
| Fixed photo frames (2-4 layouts) | ✅ Done | Simple, Event, Party frames |
| Overlay logo / event name | ✅ Done | Event frame has "Special Event" text |
| Grayscale filter | ✅ Done | Canvas CSS filter |
| Vintage filter | ✅ Done | Sepia + contrast + brightness |
| Brightness adjustment | ✅ Done | Slider (50-150%) |
| Contrast adjustment | ✅ Done | Slider (50-150%) |
| 3-2-1 countdown | ✅ Done | Animated countdown before capture |
| Flash animation | ✅ Done | White overlay effect |
| Multi-photo collage layouts | ✅ Done | 6 layouts: single, horizontal2, vertical2, grid4, strip4, featured3 |

### Phase 3 – Output & Sharing

| Task | Status | Notes |
|------|--------|-------|
| Photo download | ✅ Done | Download button on result screen |
| QR code generation | ✅ Done | `go-qrcode` library, `/photos/{id}/qr` endpoint |
| Print-ready format | ✅ Done | 4x6, 5x7, 6x8, 2x6 strip sizes at 300 DPI |
| Print size catalog | ✅ Done | `GET /print-sizes` endpoint |

### Phase 4 – Advanced Backend

| Task | Status | Notes |
|------|--------|-------|
| Clean Architecture structure | ✅ Done | `domain → usecase → handler + repository` |
| Database (SQLite) | ✅ Done | WAL mode, busy timeout, auto-migration |
| Session metadata | ✅ Done | SQLite persistence with status tracking |
| Device metadata | ✅ Done | SQLite persistence with CRUD |
| Background jobs (cleanup) | ✅ Done | Hourly cleanup of photos older than 30 days |
| Config store | ✅ Done | Thread-safe runtime configuration |
| Session limiting | ✅ Done | Max 50 concurrent sessions |

### Phase 5 – Admin & Operations

| Task | Status | Notes |
|------|--------|-------|
| Admin dashboard | ✅ Done | Vue.js admin pages at `/admin` |
| Total photos stats | ✅ Done | `/admin/stats` endpoint |
| Storage usage stats | ✅ Done | Formatted bytes in stats |
| Usage per day stats | ✅ Done | Photos today, sessions today |
| Layout/frame config | ✅ Done | `/admin/config` endpoint |
| Event name config | ✅ Done | Configurable via admin |
| Countdown duration config | ✅ Done | Configurable 1-10 seconds |

### Phase 6 – Production & Scaling

| Task | Status | Notes |
|------|--------|-------|
| MinIO / S3 object storage | ✅ Done | `pkg/storage/minio.go` with Strategy Pattern |
| Structured logging | ✅ Done | `slog` JSON logger |
| Metrics | ✅ Done | `/metrics` with uptime, requests, errors, rate |
| Health check | ✅ Done | Enhanced `/health` with DB + storage checks |
| Rate limiting | ✅ Done | Per-IP token bucket, 100 req/min |
| CORS middleware | ✅ Done | Extracted reusable `WithCORS()` wrapper |
| Docker deployment | ✅ Done | Dockerfile + docker-compose (3 profiles) |
| Render deployment | ✅ Done | `render.yaml` Blueprint |
| Oracle deployment | ✅ Done | Full guide + scripts in `deploy/oracle/` |
| GitHub Actions CI/CD | ✅ Done | Build, test, docker push, deploy workflows |

### Phase 7 – Testing & Quality

| Task | Status | Notes |
|------|--------|-------|
| Unit tests (domain entities) | ✅ Done | Photo, Session, Device, Config, Errors |
| Unit tests (use cases) | ✅ Done | Photo, Session, Device, Admin (37 tests) |
| Unit tests (handlers) | ✅ Done | Photo, Session, Device, Admin, Print, QR (48 tests) |
| Unit tests (middleware) | ✅ Done | Metrics, RateLimiter, SessionLimiter, Logger |
| Unit tests (repositories) | ✅ Done | Photo, Session, Device — all SQLite repos (25 tests) |
| Unit tests (storage) | ✅ Done | LocalStorage — full filesystem coverage (14 tests) |
| Code quality | ✅ Done | `go vet` clean, effective go |
| Dead code removal | ✅ Done | Completed refactoring pass |

**Total: 96+ individual test cases across all layers.**

---

## What's Left (Future)

1. **E2E / integration tests** — Playwright or Cypress for full UI flow
2. **Authentication** — Admin login protection
3. **Analytics** — Usage tracking and reporting
4. **Cloud CDN** — For optimized photo delivery

---

## Tech Stack Summary

| Component | Technology |
|-----------|------------|
| Backend | Go 1.25 (stdlib `net/http`, Clean Architecture) |
| Database | SQLite with WAL mode |
| Storage | Local filesystem or MinIO/S3 (Strategy Pattern) |
| Frontend | Vue 3 + Vite + Pinia + Vue Router |
| QR Codes | `github.com/skip2/go-qrcode` |
| Image Processing | `github.com/nfnt/resize` (Lanczos3) |
| Logging | `log/slog` (structured JSON) |
| Testing | `testing` package + `testify` suite/assert |
| CI/CD | GitHub Actions + Render |
| Deployment | Docker Compose (multi-profile) |
