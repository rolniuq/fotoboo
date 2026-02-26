# Deployment Guide

How to deploy and configure FotoBoo for different environments.

---

## Deployment Options

### Option -1: Oracle Always Free VM (No-Cost Production)

Use an Oracle Cloud Always Free VM for persistent, always-on hosting.

#### Why choose this

- No monthly hosting cost
- No idle sleep behavior
- Persistent storage for SQLite + photos
- Full control of networking and SSL

#### Included in this repository

- `docker-compose.yml` with `prod` profile (`app` + `caddy`)
- `.env.example` for production config
- `deploy/oracle/Caddyfile` for automatic HTTPS
- `deploy/oracle/oracle-first-boot.sh` for Docker setup
- `deploy/oracle/deploy-nipio.sh` for instant HTTPS without DNS setup
- `deploy/oracle/README.md` full step-by-step guide

Quick start guide: see [Oracle Always Free Deployment](../deploy/oracle/README.md).

### Option 0: Render Free (No-Cost Cloud, Fastest)

Use Render to deploy directly from GitHub with no server setup. This is the fastest free path for a live demo URL.

#### What You Need

- A Render account (free)
- This repository on GitHub
- `Dockerfile` and `render.yaml` in the project root

#### One-Time GitHub Setup

Push these files to your repository:

- `Dockerfile`
- `.dockerignore`
- `render.yaml`

#### Deploy Steps

1. Open Render dashboard: https://dashboard.render.com
2. Click **New +** → **Blueprint**
3. Connect your GitHub repo
4. Render detects `render.yaml` automatically
5. Click **Apply** to create and deploy service

Render creates one web service named `fotoboo` on the free plan with:

- Docker runtime
- Health check: `/health`
- Automatic deploy on every push
- Environment variables configured from `render.yaml`

#### Important Free-Plan Notes

- The service can spin down when idle and take time to wake up.
- `/tmp` storage is ephemeral. Uploaded photos and SQLite data may reset on redeploy/restart.
- This is ideal for demos and staging. For persistent storage, use Oracle Always Free VM or a paid persistent disk.

#### Environment Variables Used on Render

| Variable | Value |
|----------|-------|
| `PORT` | `10000` |
| `WEB_DIR` | `/app/web` |
| `STORAGE_PATH` | `/tmp/fotoboo/photos` |
| `DB_PATH` | `/tmp/fotoboo/fotoboo.db` |

---

### Option 1: Direct Binary (Simplest)

Build and run the compiled Go binary directly on the host machine.

```bash
# Build
go build -o bin/fotoboo-api ./cmd/api

# Run
PORT=8080 STORAGE_PATH=/var/data/fotoboo/photos ./bin/fotoboo-api
```

**Pros:** Simple, no container overhead, fast startup
**Cons:** Must install Go to build, manual process management

### Option 2: Systemd Service (Linux)

Create a systemd unit file for automatic startup and management:

```ini
# /etc/systemd/system/fotoboo.service
[Unit]
Description=FotoBoo Photo Booth API
After=network.target

[Service]
Type=simple
User=fotoboo
WorkingDirectory=/opt/fotoboo
ExecStart=/opt/fotoboo/bin/fotoboo-api
Restart=always
RestartSec=5
Environment=PORT=8080
Environment=STORAGE_PATH=/var/data/fotoboo/photos
Environment=WEB_DIR=/opt/fotoboo/web

[Install]
WantedBy=multi-user.target
```

```bash
# Install the service
sudo cp fotoboo.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable fotoboo
sudo systemctl start fotoboo

# Check status
sudo systemctl status fotoboo

# View logs
sudo journalctl -u fotoboo -f
```

### Option 3: Docker (Recommended for Production)

> **Note:** This repository already includes a production-ready `Dockerfile`.

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o fotoboo-api ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/fotoboo-api .
COPY --from=builder /app/web ./web
RUN mkdir -p /app/data/photos
EXPOSE 8080
CMD ["./fotoboo-api"]
```

```bash
# Build image
docker build -t fotoboo:latest .

# Run container
docker run -d \
  --name fotoboo \
  -p 8080:8080 \
  -v fotoboo-data:/app/data/photos \
  -e PORT=8080 \
  -e STORAGE_PATH=/app/data/photos \
  fotoboo:latest
```

---

## Configuration Reference

### Environment Variables

#### Core Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | HTTP server listen port | `8080` | No |
| `DB_PATH` | SQLite database path | `./data/fotoboo.db` | No |
| `WEB_DIR` | Directory containing frontend files | `./web` | No |
| `BASE_URL` | Base URL for QR codes and links | `http://localhost:8080` | No |

#### Storage Configuration

FotoBoo supports two storage backends: **local file system** and **MinIO/S3**.

##### Local Storage (Default)

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `STORAGE_PATH` | Directory for photo files | `./data/photos` | No |
| `USE_MINIO` | Set to `false` for local storage | `false` | No |

##### MinIO/S3 Storage

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `USE_MINIO` | Enable MinIO storage | `false` | Yes (for MinIO) |
| `MINIO_ENDPOINT` | MinIO server endpoint | `localhost:9000` | Yes (for MinIO) |
| `MINIO_ACCESS_KEY` | MinIO access key | `minioadmin` | Yes (for MinIO) |
| `MINIO_SECRET_KEY` | MinIO secret key | `minioadmin` | Yes (for MinIO) |
| `MINIO_BUCKET` | Bucket name for photos | `fotoboo` | Yes (for MinIO) |
| `MINIO_USE_SSL` | Use HTTPS for MinIO | `false` | No |

See [MinIO Setup Guide](./minio-setup.md) for detailed configuration instructions.

### Storage Requirements

#### Local File System

Estimate storage based on expected usage:

| Photos/day | Avg size | Daily storage | Monthly storage |
|------------|----------|---------------|-----------------|
| 100 | ~500 KB | ~50 MB | ~1.5 GB |
| 500 | ~500 KB | ~250 MB | ~7.5 GB |
| 1,000 | ~500 KB | ~500 MB | ~15 GB |

#### MinIO/S3

When using MinIO or S3:
- Photos are stored as objects in the cloud bucket
- No local disk space needed for photos (only database)
- Consider S3 storage costs: ~$0.023/GB/month for standard storage
- Use MinIO for self-hosted object storage with unlimited capacity
- See [MinIO Setup Guide](./minio-setup.md) for deployment options

---

## Network Configuration

### Reverse Proxy (Nginx)

For production, place FotoBoo behind a reverse proxy:

```nginx
# /etc/nginx/sites-available/fotoboo
server {
    listen 80;
    server_name fotoboo.example.com;

    client_max_body_size 10M;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### HTTPS (Required for Camera Access)

Browsers require HTTPS for `getUserMedia()` on non-localhost domains. Use Let's Encrypt or your own certificate:

```nginx
server {
    listen 443 ssl;
    server_name fotoboo.example.com;

    ssl_certificate /etc/letsencrypt/live/fotoboo.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/fotoboo.example.com/privkey.pem;

    client_max_body_size 10M;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

> **Important:** Camera access (`getUserMedia`) only works on:
> - `localhost` / `127.0.0.1` (any protocol)
> - HTTPS domains
> - It will **not work** on plain HTTP with a non-localhost domain

---

## Health Monitoring

### Health Check Endpoint

```bash
curl -f http://localhost:8080/health || echo "Service is down"
```

### Docker Healthcheck

```dockerfile
HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
  CMD wget -q --spider http://localhost:8080/health || exit 1
```

### Uptime Monitoring

Point your monitoring service (UptimeRobot, Pingdom, etc.) at:

```
GET https://fotoboo.example.com/health
Expected: 200 OK, body contains {"status":"ok"}
```

---

## Backup Strategy

### What to Back Up

#### Local Storage

| Data | Location | Priority |
|------|----------|----------|
| Photo files | `$STORAGE_PATH/*.jpg` | **High** — irreplaceable |
| Database | `$DB_PATH` (fotoboo.db) | **High** — metadata and relationships |
| Configuration | Environment variables | **Low** — easily recreated |

#### MinIO/S3 Storage

| Data | Location | Priority |
|------|----------|----------|
| Photo objects | MinIO bucket or S3 | **High** — use MinIO mirroring or S3 versioning |
| Database | `$DB_PATH` (fotoboo.db) | **High** — metadata and relationships |
| Configuration | Environment variables | **Low** — easily recreated |

### Backup Commands

#### Local Storage

```bash
# Full backup
tar -czf fotoboo-backup-$(date +%Y%m%d).tar.gz \
  /var/data/fotoboo/photos/ \
  /var/data/fotoboo/fotoboo.db

# Incremental backup (files modified in last 24h)
find /var/data/fotoboo/photos/ -mtime -1 -type f | \
  tar -czf fotoboo-incremental-$(date +%Y%m%d).tar.gz -T -

# Restore from backup
tar -xzf fotoboo-backup-20250115.tar.gz -C /
```

#### MinIO Storage

```bash
# Backup MinIO bucket
mc mirror myminio/fotoboo ./backup/fotoboo

# Backup database
cp /var/data/fotoboo/fotoboo.db ./backup/

# Restore MinIO bucket
mc mirror ./backup/fotoboo myminio/fotoboo

# Restore database
cp ./backup/fotoboo.db /var/data/fotoboo/
```

See [MinIO Setup Guide](./minio-setup.md) for detailed backup instructions.

---

## Event Deployment Checklist

For deploying FotoBoo at a physical event:

- [ ] **Hardware:** Computer/tablet with webcam, stable power supply
- [ ] **Network:** Reliable WiFi or ethernet (for multi-device setups)
- [ ] **Storage:** Sufficient disk space for expected photo count
- [ ] **HTTPS:** Required if accessing from non-localhost (use self-signed cert for local network)
- [ ] **Browser:** Chrome/Edge in kiosk mode (`--kiosk` flag)
- [ ] **Camera permissions:** Pre-grant camera access in browser settings
- [ ] **Auto-start:** Configure service to start on boot (systemd/launchd)
- [ ] **Monitoring:** Set up health check alerts
- [ ] **Backup:** Schedule periodic backups during the event
- [ ] **Test run:** Capture and download at least 10 photos before the event

### Kiosk Mode (Chrome)

```bash
# Linux
google-chrome --kiosk --disable-infobars http://localhost:8080

# macOS
open -a "Google Chrome" --args --kiosk http://localhost:8080

# Windows
chrome.exe --kiosk http://localhost:8080
```

---

## Troubleshooting

### Server won't start

```bash
# Check if port is in use
lsof -i :8080

# Check storage directory permissions
ls -la ./data/photos/

# Run with verbose output
go run ./cmd/api 2>&1 | tee fotoboo.log
```

### Photos not persisting after restart

- Verify `STORAGE_PATH` points to a persistent volume (not tmpfs)
- In Docker, ensure the volume is mounted: `-v fotoboo-data:/app/data/photos`

### Camera not working in browser

- Must be HTTPS or localhost
- Check browser camera permissions
- Verify no other application is using the webcam
- Try a different browser

### High memory usage

- The in-memory photo index grows with photo count
- For very large datasets (100k+ photos), consider migrating to a database (Phase 4)
