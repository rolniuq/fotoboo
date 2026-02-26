# Getting Started

This guide walks you through setting up and running FotoBoo locally.

---

## Prerequisites

- **Go 1.21+** (project uses Go 1.25.3)
- A modern web browser with webcam access (Chrome, Firefox, Edge, Safari)
- (Optional) A webcam or built-in camera for photo capture
- (Optional) Docker for running MinIO (if using object storage)

---

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/fotoboo/fotoboo.git
cd fotoboo
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run the server

```bash
go run ./cmd/api
```

You should see:

```
Starting FotoBoo API server on :8080
Storage path: ./data/photos
Serving web UI from: ./web
```

### 4. Open in browser

Navigate to [http://localhost:8080](http://localhost:8080)

---

## Building from Source

Compile a production binary:

```bash
go build -o bin/fotoboo-api ./cmd/api
```

Run the compiled binary:

```bash
./bin/fotoboo-api
```

---

## Configuration

FotoBoo is configured via environment variables.

### Basic Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `DB_PATH` | SQLite database path | `./data/fotoboo.db` |
| `WEB_DIR` | Directory for frontend files | `./web` |

### Storage Options

#### Option 1: Local File Storage (Default)

| Variable | Description | Default |
|----------|-------------|---------|
| `STORAGE_PATH` | Directory for photo storage | `./data/photos` |

#### Option 2: MinIO/S3 Object Storage

| Variable | Description | Default |
|----------|-------------|---------|
| `USE_MINIO` | Enable MinIO storage | `false` |
| `MINIO_ENDPOINT` | MinIO server endpoint | `localhost:9000` |
| `MINIO_ACCESS_KEY` | MinIO access key | `minioadmin` |
| `MINIO_SECRET_KEY` | MinIO secret key | `minioadmin` |
| `MINIO_BUCKET` | Bucket name | `fotoboo` |
| `MINIO_USE_SSL` | Use HTTPS | `false` |

For detailed MinIO setup, see [MinIO Setup Guide](./minio-setup.md).

### Examples

#### Local Storage

```bash
# Run on a custom port
PORT=3000 go run ./cmd/api

# Use a custom storage directory
STORAGE_PATH=/tmp/fotoboo-photos go run ./cmd/api

# Combine multiple settings
PORT=3000 STORAGE_PATH=/mnt/data/photos go run ./cmd/api
```

#### MinIO Storage

```bash
# Start MinIO server (Docker)
docker-compose -f docker-compose.minio.yml up -d

# Run FotoBoo with MinIO
USE_MINIO=true \
MINIO_ENDPOINT=localhost:9000 \
MINIO_ACCESS_KEY=minioadmin \
MINIO_SECRET_KEY=minioadmin \
MINIO_BUCKET=fotoboo \
go run ./cmd/api
```

---

## Project Structure

```
fotoboo/
├── cmd/api/              # Application entrypoint
│   └── main.go           # Server setup & dependency injection
├── internal/
│   ├── domain/           # Core entities & interfaces
│   │   ├── photo.go      # Photo entity + PhotoRepository interface
│   │   └── errors.go     # Domain error definitions
│   ├── usecase/          # Business logic
│   │   └── photo_usecase.go
│   ├── handler/          # HTTP handlers
│   │   └── photo_handler.go
│   └── repository/       # Persistence implementations
│       └── photo_repository.go
├── web/                  # Frontend SPA
│   ├── index.html
│   └── static/
│       ├── css/style.css
│       └── js/app.js
├── data/photos/          # Photo storage (created at runtime)
├── bin/                  # Compiled binaries
├── docs/                 # Documentation
├── ROADMAP.md            # Development roadmap
├── CLAUDE.md             # Claude Code skill file
├── go.mod
└── go.sum
```

---

## Using the Photo Booth

### Step 1: Welcome Screen

Click **"Start Photo Session"** to begin.

### Step 2: Capture

- Your webcam feed appears live on screen
- Click the **capture button** to take a photo
- A 3-2-1 countdown starts before capture
- A flash animation confirms the shot

### Step 3: Preview & Edit

After capture, you can:

- **Apply filters:** None, Grayscale, Vintage, Bright, Contrast
- **Choose a frame:** None, Simple (white border), Event (gold + text), Party (red + emoji)
- **Adjust brightness/contrast** using sliders
- Click **"Retake"** to go back and try again
- Click **"Save & Continue"** to finalize

### Step 4: Download

- View your final photo
- Click **"Download Photo"** to save it to your device
- Click **"Take Another Photo"** to start over

---

## Testing the API Manually

### Upload a photo

```bash
curl -X POST http://localhost:8080/photos \
  --data-binary @test-photo.jpg \
  -H "Content-Type: application/octet-stream"
```

Response:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-01-15T14:30:00Z"
}
```

### Retrieve a photo

```bash
curl http://localhost:8080/photos/550e8400-e29b-41d4-a716-446655440000 \
  --output downloaded.jpg
```

### Health check

```bash
curl http://localhost:8080/health
```

---

## Development Workflow

```bash
# Run in development mode (auto-restart not included)
go run ./cmd/api

# Build and test binary
go build -o bin/fotoboo-api ./cmd/api && ./bin/fotoboo-api

# Tidy dependencies
go mod tidy
```

### File Watching (Optional)

For automatic restarts during development, you can use tools like [air](https://github.com/air-verse/air):

```bash
go install github.com/air-verse/air@latest
air
```

---

## Troubleshooting

### Camera not working

- Ensure your browser has camera permissions for `localhost`
- Check that no other application is using the webcam
- Try a different browser (Chrome recommended)
- HTTPS is required for camera access on non-localhost domains

### Port already in use

```bash
# Find what's using the port
lsof -i :8080

# Use a different port
PORT=3001 go run ./cmd/api
```

### Photos not saving

- Check that the `STORAGE_PATH` directory exists and is writable
- The server creates it automatically, but parent directories must exist
- Check server logs for error messages

---

## Next Steps

- Read the [Architecture Guide](./architecture.md) to understand the codebase structure
- Check the [API Reference](./api.md) for detailed endpoint documentation
- See the [Roadmap](../ROADMAP.md) for planned features
