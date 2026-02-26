# FotoBoo QA Test Execution Report

- Date: 2026-02-26
- Scope: `test-cases/test-cases.md`

## Execution Summary

- Total cases: 33
- Passed: 31
- Failed: 2

## Command-Based Test Runs

### 1) Backend unit tests

- Command: `go test ./...`
- Result: PASS

### 2) Frontend build smoke test

- Command: `npm run build` (in `web/`)
- Result: PASS

### 3) API functional suite

- Environment:
  - `PORT=18080`
  - temporary `DB_PATH`
  - temporary `STORAGE_PATH`
- Result: PARTIAL PASS (24/26 API cases passed)

## Failed Cases

| ID | Case | Expected | Actual | Log |
|---|---|---|---|---|
| API-018 | `DELETE /devices/nonexistent` | `404` | `204` | `error-logs/2026-02-26-api-018-delete-device-nonexistent.md` |
| API-024 | `GET /qr` (without `text`) | `400` | `200` (PNG returned) | `error-logs/2026-02-26-api-024-qr-missing-text.md` |

## Notes

- The two failures indicate backend behavior gaps/inconsistencies, not test harness failures.
- Detailed repro, impact, and code-level observations are documented in `error-logs/`.
