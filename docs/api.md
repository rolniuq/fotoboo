# API Reference

Complete API documentation for the FotoBoo backend service.

---

## Base URL

```
http://localhost:8080
```

Configure via `PORT` environment variable.

---

## CORS

All API endpoints include CORS headers:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type
```

`OPTIONS` preflight requests are handled automatically with a `200 OK` response.

---

## Endpoints

### Upload Photo

Upload a photo to the server.

```
POST /photos
```

**Request:**

| Header | Value |
|--------|-------|
| `Content-Type` | `application/octet-stream` |

- **Body:** Raw image data (JPEG bytes)
- **Max size:** 10 MB

**Success Response (201 Created):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-01-15T14:30:00.123456Z"
}
```

**Error Responses:**

| Status | Condition | Body |
|--------|-----------|------|
| `400 Bad Request` | Empty body / invalid data | `{"error": "invalid photo data"}` |
| `405 Method Not Allowed` | Non-POST method | `{"error": "method not allowed"}` |
| `500 Internal Server Error` | Storage failure | `{"error": "internal server error"}` |

**Example:**

```bash
# Upload a photo file
curl -X POST http://localhost:8080/photos \
  --data-binary @photo.jpg \
  -H "Content-Type: application/octet-stream"

# Upload from webcam capture (JavaScript)
# Frontend sends: fetch('/photos', { method: 'POST', body: blob })
```

---

### Get Photo

Retrieve a photo by its ID.

```
GET /photos/{id}
```

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Photo identifier returned from upload |

**Success Response (200 OK):**

| Header | Value |
|--------|-------|
| `Content-Type` | `image/jpeg` |

- **Body:** Raw JPEG image data

**Error Responses:**

| Status | Condition | Body |
|--------|-----------|------|
| `404 Not Found` | Photo ID doesn't exist | `{"error": "photo not found"}` |
| `405 Method Not Allowed` | Non-GET method | `{"error": "method not allowed"}` |
| `500 Internal Server Error` | File read failure | `{"error": "internal server error"}` |

**Example:**

```bash
# Download a photo
curl http://localhost:8080/photos/550e8400-e29b-41d4-a716-446655440000 \
  --output photo.jpg

# View in browser
# Navigate to: http://localhost:8080/photos/550e8400-e29b-41d4-a716-446655440000
```

---

### Health Check

Check if the server is running and responsive.

```
GET /health
```

**Success Response (200 OK):**

```json
{
  "status": "ok"
}
```

**Example:**

```bash
curl http://localhost:8080/health
```

---

### Serve Frontend

Serves the static frontend SPA files.

```
GET /
GET /static/*
```

The server serves all files from the `web/` directory (configurable via `WEB_DIR` env var).

---

## Response Formats

### Success Responses

All JSON responses use `Content-Type: application/json` except for photo retrieval which returns `image/jpeg`.

### Error Responses

All errors follow a consistent format:

```json
{
  "error": "human-readable error message"
}
```

---

## Rate Limits & Constraints

| Constraint | Value |
|------------|-------|
| Max upload size | 10 MB |
| Supported image format | JPEG (hardcoded) |
| Concurrent connections | Unlimited (no rate limiting yet) |

---

## Planned Endpoints (from Roadmap)

These endpoints are not yet implemented but are planned for future phases:

| Method | Path | Description | Phase |
|--------|------|-------------|-------|
| GET | `/photos/{id}/qr` | Generate QR code for photo URL | Phase 3 |
| POST | `/photos/{id}/print` | Send photo to printer | Phase 3 |
| GET | `/admin/stats` | Dashboard statistics | Phase 5 |
| GET | `/admin/config` | Get current configuration | Phase 5 |
| PUT | `/admin/config` | Update configuration | Phase 5 |
