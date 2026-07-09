# FotoBoo Project — Development Roadmap

This document describes the development roadmap for the FotoBoo project, from MVP to production-ready.

---

## Overall Goals

### User-facing
- Capture photos via browser camera
- Preview, filter, frame, and enhance results
- Create multi-photo collages
- Export via download, QR code, or print-ready format

### System
- Run reliably for long event sessions
- Clean, extensible architecture
- Scalable storage (local → S3/MinIO)
- Easy deployment (Docker, Render, Oracle Always Free)

---

## Phases

### ✅ Phase 0 — Planning & Foundation
- [x] Define MVP scope
- [x] Choose tech stack (Go + Vue 3 + SQLite)
- [x] Set up Clean Architecture project structure
- [x] Define entity model (Photo, Session, Device, Config)

### ✅ Phase 1 — MVP Core Features
- [x] Camera integration (`getUserMedia` at 1280x720)
- [x] Photo capture with countdown + flash animation
- [x] Preview screen
- [x] Upload to backend (Go API, raw body)
- [x] Retrieve photo by ID
- [x] Local filesystem storage
- [x] Session management
- [x] Welcome → Capture → Preview → Result flow

### ✅ Phase 2 — UX & Photo Enhancement
- [x] CSS filters: grayscale, vintage, brightness, contrast
- [x] Frame overlays: simple, event, party
- [x] Brightness/contrast sliders
- [x] Multi-photo collage layouts (6 layouts)
- [x] Animated countdown (3-2-1)
- [x] Flash animation effect

### ✅ Phase 3 — Output & Sharing
- [x] Photo download
- [x] QR code generation per photo
- [x] Print-ready export (4x6, 5x7, 6x8, 2x6 at 300 DPI)
- [x] Print size catalog

### ✅ Phase 4 — Advanced Backend
- [x] Clean Architecture with dependency inversion
- [x] SQLite persistence with WAL mode
- [x] Session + Device CRUD
- [x] Thread-safe runtime configuration store
- [x] Background job: cleanup old photos
- [x] Dual storage backends (local + MinIO/S3 via Strategy Pattern)

### ✅ Phase 5 — Admin & Operations
- [x] Admin dashboard (Vue.js)
- [x] Stats: total photos, sessions, devices, storage
- [x] Daily usage stats (photos/sessions today)
- [x] Configuration management (event name, countdown, frames, filters)
- [x] Upload size and retention configuration

### ✅ Phase 6 — Production & Scaling
- [x] MinIO / S3 object storage support
- [x] Structured JSON logging (`log/slog`)
- [x] Metrics endpoint (requests, errors, uptime, rate)
- [x] Enhanced health check (DB + storage)
- [x] Rate limiting (per-IP token bucket)
- [x] Session limiting (max concurrent)
- [x] CORS middleware
- [x] Docker multi-stage build
- [x] Docker Compose (local, prod, minio profiles)
- [x] Render Blueprint deployment
- [x] Oracle Always Free deployment guide
- [x] GitHub Actions CI/CD

### ✅ Phase 7 — Testing & Quality
- [x] Domain entity tests (Photo, Session, Device, Config, Errors)
- [x] Use case tests (Photo, Session, Device, Admin)
- [x] HTTP handler tests (Photo, Session, Device, Admin, Print, QR)
- [x] Middleware tests (Metrics, RateLimiter, SessionLimiter, Logger)
- [x] Repository tests (all SQLite repositories)
- [x] Storage tests (LocalStorage — full filesystem coverage)
- [x] Code quality: `go vet` clean, effective go, no dead code

---

## Future Work

| Priority | Feature | Notes |
|----------|---------|-------|
| Medium | E2E / integration tests | Playwright or Cypress for full UI flow |
| Low | Admin authentication | Login protection for admin endpoints |
| Low | Usage analytics | Event-level reporting |
| Low | CDN integration | Optimized photo delivery |

---

## Architecture

```
[Browser Camera] → [Vue 3 SPA] → [Go API] → [SQLite + Local/MinIO Storage]
                                      ↑
                              Clean Architecture layers:
                              domain → usecase → handler / repository
```

---

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.25 (stdlib `net/http`) |
| Frontend | Vue 3 + Vite + Pinia + Vue Router |
| Database | SQLite (WAL mode) |
| Storage | Local filesystem or MinIO/S3 |
| Image | `nfnt/resize` (Lanczos3) |
| QR Code | `skip2/go-qrcode` |
| Testing | `testing` + `testify` |
| CI/CD | GitHub Actions + Render |
