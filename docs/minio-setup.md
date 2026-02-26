# MinIO Setup Guide for FotoBoo

This guide explains how to configure FotoBoo to use MinIO for photo storage instead of local file system.

## Quick Start with Docker Compose

1. Start MinIO server:
```bash
USE_MINIO=true docker compose --profile minio up -d
```

2. Access MinIO Console at http://localhost:9001
   - Username: `minioadmin`
   - Password: `minioadmin`

3. FotoBoo runs in the same compose stack with MinIO enabled:
```bash
docker compose ps
```

4. Access URLs:
   - FotoBoo: http://localhost:8080
   - MinIO Console: http://localhost:9001

## Environment Variables

### MinIO Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `USE_MINIO` | Enable MinIO storage (true/false) | `false` | Yes |
| `MINIO_ENDPOINT` | MinIO server endpoint | `localhost:9000` | Yes |
| `MINIO_ACCESS_KEY` | MinIO access key | `minioadmin` | Yes |
| `MINIO_SECRET_KEY` | MinIO secret key | `minioadmin` | Yes |
| `MINIO_BUCKET` | Bucket name for photos | `fotoboo` | Yes |
| `MINIO_USE_SSL` | Use HTTPS for MinIO (true/false) | `false` | No |

### Example .env file

```bash
# Storage configuration
USE_MINIO=true
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=fotoboo
MINIO_USE_SSL=false

# Database
DB_PATH=./data/fotoboo.db

# Server
PORT=8080
BASE_URL=http://localhost:8080
WEB_DIR=./web
```

## Production Setup

### Using AWS S3

MinIO SDK is compatible with AWS S3. Configure as follows:

```bash
export USE_MINIO=true
export MINIO_ENDPOINT=s3.amazonaws.com
export MINIO_ACCESS_KEY=your-aws-access-key
export MINIO_SECRET_KEY=your-aws-secret-key
export MINIO_BUCKET=your-bucket-name
export MINIO_USE_SSL=true
```

### Using MinIO in Production

1. Deploy MinIO server (see https://min.io/docs/minio/linux/operations/install-deploy-manage/deploy-minio-single-node-single-drive.html)

2. Create a dedicated access key:
```bash
mc alias set myminio http://your-minio-server:9000 minioadmin minioadmin
mc admin user add myminio fotoboo-app your-secure-password
mc admin policy attach myminio readwrite --user fotoboo-app
```

3. Configure FotoBoo:
```bash
export USE_MINIO=true
export MINIO_ENDPOINT=your-minio-server:9000
export MINIO_ACCESS_KEY=fotoboo-app
export MINIO_SECRET_KEY=your-secure-password
export MINIO_BUCKET=fotoboo
export MINIO_USE_SSL=true
```

## Storage Migration

### From Local Storage to MinIO

Currently, there's no automated migration tool. Manual migration steps:

1. Export all photo records from database:
```sql
SELECT id, file_path FROM photos;
```

2. Upload each photo to MinIO using the MinIO client:
```bash
mc cp ./data/photos/*.jpg myminio/fotoboo/photos/
```

3. Update database file_path column to use MinIO object paths:
```sql
UPDATE photos SET file_path = 'photos/' || id || '.jpg';
```

4. Enable MinIO in FotoBoo configuration

### From MinIO to Local Storage

1. Download all photos from MinIO:
```bash
mc cp --recursive myminio/fotoboo/photos/ ./data/photos/
```

2. Update database file_path to local paths
3. Disable MinIO in FotoBoo configuration

## Bucket Structure

Photos are stored with the following structure:
```
fotoboo/                    # Bucket name
└── photos/
    ├── {photo-id-1}.jpg
    ├── {photo-id-2}.jpg
    └── {photo-id-3}.jpg
```

## Troubleshooting

### Connection refused
- Verify MinIO is running: `docker ps`
- Check endpoint is correct (no http:// prefix)
- Ensure port 9000 is accessible

### Authentication failed
- Verify access key and secret key match MinIO configuration
- For AWS S3, ensure IAM user has S3 permissions

### Bucket not found
- The application automatically creates the bucket if it doesn't exist
- Ensure the access key has permission to create buckets

### SSL certificate errors
- For development, use `MINIO_USE_SSL=false`
- For production with self-signed certs, you may need to configure cert validation

## Performance Considerations

- MinIO provides better scalability than local file storage
- Photos are stored as objects, not files on disk
- Supports distributed storage across multiple servers
- Better for containerized/cloud deployments
- Network latency may be higher than local disk access

## Backup and Recovery

### Backup MinIO Data

```bash
# Backup entire bucket
mc mirror myminio/fotoboo ./backup/fotoboo

# Schedule daily backups
mc mirror --watch myminio/fotoboo ./backup/fotoboo
```

### Restore from Backup

```bash
mc mirror ./backup/fotoboo myminio/fotoboo
```

## Switching Between Storage Types

You can switch between local and MinIO storage at runtime by changing the `USE_MINIO` environment variable. However:

- Existing photos won't be automatically migrated
- The database still references the old storage paths
- You need to manually migrate data (see "Storage Migration" above)
