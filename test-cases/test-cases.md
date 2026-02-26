# FotoBoo QA Test Cases

## Scope

- Backend API functional testing
- Backend unit test execution
- Frontend build smoke test
- Basic error handling verification

## Environment

- OS: macOS (darwin)
- Go: project-managed version
- Node.js/npm: required for frontend build check
- Test server config:
  - `PORT=18080`
  - `DB_PATH=<temp>/fotoboo.db`
  - `STORAGE_PATH=<temp>/photos`

## Automated Unit Test Cases

| ID | Area | Test Case | Expected Result |
|---|---|---|---|
| UT-001 | Use case | Run all Go unit tests (`go test ./...`) | All tests pass |
| UT-002 | Use case | Photo use case tests | Pass |
| UT-003 | Use case | Session use case tests | Pass |
| UT-004 | Use case | Device use case tests | Pass |
| UT-005 | Handler | Photo handler tests | Pass |
| UT-006 | Handler | Session handler tests | Pass |

## API Functional Test Cases

| ID | Endpoint | Scenario | Expected Status |
|---|---|---|---|
| API-001 | `GET /health` | Health check responds | `200` |
| API-002 | `POST /photos` | Upload valid JPEG bytes | `201` |
| API-003 | `GET /photos/{id}` | Retrieve uploaded photo | `200` |
| API-004 | `GET /photos` | List photos | `200` |
| API-005 | `GET /photos/{id}/qr` | Generate QR for existing photo | `200` |
| API-006 | `GET /photos/{id}/print?size=4x6` | Generate print-ready photo | `200` |
| API-007 | `DELETE /photos/{id}` | Delete existing photo | `204` |
| API-008 | `GET /photos/{id}` | Get deleted photo | `404` |
| API-009 | `POST /sessions` | Start session | `201` |
| API-010 | `GET /sessions` | List sessions | `200` |
| API-011 | `GET /sessions/{id}` | Get existing session | `200` |
| API-012 | `POST /sessions/{id}/complete` | Complete existing session | `200` |
| API-013 | `GET /sessions/{id}/photos` | Get session photos | `200` |
| API-014 | `POST /devices` | Register device | `201` |
| API-015 | `GET /devices/{id}` | Get registered device | `200` |
| API-016 | `PUT /devices/{id}` | Update existing device | `200` |
| API-017 | `DELETE /devices/{id}` | Delete existing device | `204` |
| API-018 | `DELETE /devices/nonexistent` | Delete unknown device | `404` |
| API-019 | `GET /admin/stats` | Fetch dashboard stats | `200` |
| API-020 | `GET /admin/config` | Fetch config | `200` |
| API-021 | `PUT /admin/config` | Update config with invalid countdown `0` | `400` |
| API-022 | `PUT /admin/config` | Update config with valid payload | `200` |
| API-023 | `GET /print-sizes` | List print sizes | `200` |
| API-024 | `GET /qr` | Missing text param | `400` |
| API-025 | `GET /qr?text=hello` | Generate QR from text | `200` |
| API-026 | `GET /metrics` | Read metrics endpoint | `200` |

## Frontend Smoke Test Cases

| ID | Area | Test Case | Expected Result |
|---|---|---|---|
| FE-001 | Build | Run `npm run build` in `web/` | Build succeeds |

## Notes

- API-018 is intentionally included to validate 404 semantics for deleting unknown resources.
- If any case fails, details must be written under `error-logs/`.
