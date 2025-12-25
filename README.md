# FotoBoo

A photo booth API service built with Go using Clean Architecture principles. FotoBoo enables photo capture, storage, and retrieval for event photo booths, retail stores, and SaaS applications.

## Features

- Photo upload via REST API
- Photo retrieval by ID
- Local filesystem storage with metadata persistence
- CORS support for browser-based clients
- Health check endpoint

## Project Structure

```
fotoboo/
├── cmd/
│   └── api/           # Application entrypoint
├── internal/
│   ├── domain/        # Business entities and interfaces
│   ├── usecase/       # Business logic
│   ├── handler/       # HTTP handlers
│   └── repository/    # Data persistence
├── data/
│   └── photos/        # Photo storage directory
└── bin/               # Compiled binaries
```

## Requirements

- Go 1.21+

## Getting Started

### Build

```bash
go build -o bin/fotoboo-api ./cmd/api
```

### Run

```bash
./bin/fotoboo-api
```

The server starts on port `8080` by default.

### Configuration

Environment variables:

| Variable       | Description                | Default          |
|----------------|----------------------------|------------------|
| `PORT`         | Server port                | `8080`           |
| `STORAGE_PATH` | Photo storage directory    | `./data/photos`  |

## API Endpoints

### Upload Photo

```
POST /photos
Content-Type: application/octet-stream

Body: <raw image data>
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-01-01T12:00:00Z"
}
```

### Get Photo

```
GET /photos/{id}
```

Returns the photo as `image/jpeg`.

### Health Check

```
GET /health
```

Response:
```json
{"status":"ok"}
```

## Development

```bash
# Install dependencies
go mod tidy

# Run directly
go run ./cmd/api

# Build
go build -o bin/fotoboo-api ./cmd/api
```

## License

MIT
