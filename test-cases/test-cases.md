# FotoBoo QA Test Cases

## Scope

- Automated Go unit tests across all layers
- Frontend build test
- Covers domain, usecase, handler, middleware, repository, and storage layers

## Automated Unit Test Cases

| ID | Area | Coverage | Test Count |
|----|------|----------|------------|
| UT-001 | Domain | Photo, Session, Device, Config, Errors | 17 |
| UT-002 | Use Cases | Photo, Session, Device, Admin use cases | 37 |
| UT-003 | Handlers | Photo, Session, Device, Admin, Print, QR | 48 |
| UT-004 | Middleware | Metrics, RateLimiter, SessionLimiter, Logger | 10 |
| UT-005 | Repositories | Photo (SQLite), Session (SQLite), Device (SQLite) | 25 |
| UT-006 | Storage | LocalStorage (filesystem) | 14 |
| **Total** | | | **151+** |

## API Endpoints Covered

| ID | Endpoint | Scenario | Status |
|----|----------|----------|--------|
| API-001 | `GET /health` | Health check | ✅ |
| API-002 | `POST /photos` | Upload valid data | ✅ |
| API-003 | `GET /photos/{id}` | Retrieve photo | ✅ |
| API-004 | `GET /photos` | List photos | ✅ |
| API-005 | `GET /photos/{id}/qr` | Generate QR | ✅ |
| API-006 | `GET /photos/{id}/print` | Print-ready photo | ✅ |
| API-007 | `DELETE /photos/{id}` | Delete photo | ✅ |
| API-008 | `GET /photos/nonexistent` | Not found | ✅ |
| API-009 | `POST /sessions` | Start session | ✅ |
| API-010 | `GET /sessions` | List sessions | ✅ |
| API-011 | `GET /sessions/{id}` | Get session | ✅ |
| API-012 | `POST /sessions/{id}/complete` | Complete session | ✅ |
| API-013 | `GET /sessions/{id}/photos` | Session photos | ✅ |
| API-014 | `POST /devices` | Register device | ✅ |
| API-015 | `GET /devices/{id}` | Get device | ✅ |
| API-016 | `PUT /devices/{id}` | Update device | ✅ |
| API-017 | `DELETE /devices/{id}` | Delete device | ✅ |
| API-018 | `DELETE /devices/nonexistent` | 404 on missing | ✅ |
| API-019 | `GET /admin/stats` | Dashboard stats | ✅ |
| API-020 | `GET /admin/config` | Get config | ✅ |
| API-021 | `PUT /admin/config` | Validate config | ✅ |
| API-022 | `PUT /admin/config` | Update config | ✅ |
| API-023 | `GET /print-sizes` | Print sizes | ✅ |
| API-024 | `GET /qr` | Missing text | ✅ |
| API-025 | `GET /qr?text=hello` | Generate QR | ✅ |
| API-026 | `GET /metrics` | Metrics | ✅ |

## Frontend

| ID | Test | Expected | Status |
|----|------|----------|--------|
| FE-001 | `npm run build` in `web/` | Build succeeds | ✅ |

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v -count=1

# Run specific package
go test ./internal/handler -v -count=1
```
