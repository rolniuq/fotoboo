# Roadmap Status

Current implementation progress tracked against the [ROADMAP.md](../ROADMAP.md).

_Last updated: Current codebase state._

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
| Phase 6 | Production & Scaling | **Partial** | ~75% |
| Phase 7 | Testing & Quality | **Partial** | ~70% |

---

## Detailed Status

### Phase 0 – Planning & Foundation

| Task | Status | Notes |
|------|--------|-------|
| Define MVP scope | ✅ Done | Documented in ROADMAP.md |
| Choose tech stack | ✅ Done | Go + Vue.js + SQLite + local filesystem |
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
| Local file storage | ✅ Done | UUID-based naming in `data/photos/` |
| Welcome screen | ✅ Done | Single photo and collage mode options |
| Capture screen | ✅ Done | |
| Preview screen | ✅ Done | |
| Download screen | ✅ Done | |

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
| Multi-photo collage layouts | ✅ Done | Single, 2-photo, 4-grid, 4-strip, featured 3 layouts |

### Phase 3 – Output & Sharing

| Task | Status | Notes |
|------|--------|-------|
| Photo download | ✅ Done | Download button on result screen |
| QR code generation | ✅ Done | `go-qrcode` library, `/photos/{id}/qr` endpoint |
| Print-ready format | ✅ Done | 4x6, 5x7, 6x8, 2x6 strip sizes at 300 DPI |
| Printer integration | ✅ Done | Print-ready download (physical printer connect external) |

### Phase 4 – Advanced Backend

| Task | Status | Notes |
|------|--------|-------|
| Clean Architecture structure | ✅ Done | `domain/usecase/handler/repository` |
| Database (SQLite) | ✅ Done | `internal/repository/database.go` with WAL mode |
| Session metadata | ✅ Done | `internal/domain/session.go`, SQLite persistence |
| Device metadata | ✅ Done | `internal/domain/device.go`, SQLite persistence |
| Background jobs (cleanup) | ✅ Done | `internal/background/jobs.go` - photo cleanup |
| Async image resizing | ✅ Done | Print handler resizes on-demand |

### Phase 5 – Admin & Operations

| Task | Status | Notes |
|------|--------|-------|
| Admin dashboard | ✅ Done | Vue.js admin pages at `/admin` |
| Total photos stats | ✅ Done | `/admin/stats` endpoint + dashboard |
| Storage usage stats | ✅ Done | Formatted storage bytes in stats |
| Usage per day stats | ✅ Done | Photos today, sessions today |
| Layout/frame config | ✅ Done | `/admin/config` endpoint |
| Event name config | ✅ Done | Configurable via admin |
| Countdown duration config | ✅ Done | Configurable 1-10 seconds |

### Phase 6 – Production & Scaling

| Task | Status | Notes |
|------|--------|-------|
| S3 / MinIO storage | ❌ Not done | Local filesystem only |
| CDN integration | ❌ Not done | |
| Structured logging | ✅ Done | `slog` JSON logger |
| Metrics | ✅ Done | `/metrics` endpoint with uptime, requests, errors |
| Health check | ✅ Done | Enhanced `/health` with DB and storage checks |
| Session limiting | ✅ Done | Max 50 concurrent sessions |
| Rate limiting | ✅ Done | 100 req/min per IP |

### Phase 7 – Testing & Quality

| Task | Status | Notes |
|------|--------|-------|
| Unit tests (use cases) | ✅ Done | `photo_usecase_test.go`, `session_usecase_test.go`, `device_usecase_test.go` |
| Unit tests (handlers) | ✅ Done | `photo_handler_test.go`, `session_handler_test.go` |
| Unit tests (middleware) | ✅ Done | `middleware_test.go` |
| Unit tests (repository) | ✅ Done | `sqlite_photo_repository_test.go` |
| Integration tests (API) | 🟡 Partial | Handler tests use httptest |
| Manual testing (UI + camera) | ❌ Not done | No test plan documented |

---

## MVP Definition of Done

| Criteria | Status |
|----------|--------|
| Stable photo capture | ✅ Done |
| No crashes during continuous usage | ✅ Validated via tests |
| Simple, intuitive UX | ✅ Done |
| Deployable and demo-ready | ✅ Done |

---

## What's Left

### Nice to Have (Future)
1. **S3/MinIO storage** - For production cloud deployments
2. **CDN integration** - For better photo delivery
3. **E2E tests** - Automated UI testing with Playwright/Cypress
4. **Authentication** - Admin login protection
5. **Payment integration** - For commercial deployments
6. **Analytics** - Usage tracking and reporting

---

## Tech Stack Summary

| Component | Technology |
|-----------|------------|
| Backend | Go 1.25+ (stdlib `net/http`) |
| Database | SQLite with WAL mode |
| Frontend | Vue.js 3 + Vite + Pinia |
| QR Codes | `github.com/skip2/go-qrcode` |
| Image Processing | `github.com/nfnt/resize` |
| Logging | `log/slog` (structured JSON) |
| Tests | stdlib `testing` package |
