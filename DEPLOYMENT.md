# GitHub Deployment Guide for FotoBoo

This guide explains how to deploy your FotoBoo project on GitHub with Docker Compose using Render.

## Overview

Your project is configured to:
- Use **Docker Compose** locally with the `docker-compose.yml` file
- Deploy automatically to **Render** (free tier) when you push to the `master` branch
- Use a **GitHub Actions workflow** to automate the deployment process

## Setup Instructions

### Step 1: Push Code to GitHub

First, ensure your latest code is committed and pushed:

```bash
git add .
git commit -m "Add GitHub Actions deployment workflow"
git push origin master
```

### Step 2: Create Render Account

1. Go to [https://render.com](https://render.com)
2. Sign up with your GitHub account
3. Authorize Render to access your GitHub repositories

### Step 3: Create a Service on Render

1. Log in to Render dashboard
2. Click **New +** → **Blueprint**
3. Select your `fotoboo` repository
4. Render will automatically detect `render.yaml` and configure the service
5. Click **Apply** to deploy
6. Wait for the build and deployment to complete (usually 2-5 minutes)

### Step 4: Configure GitHub Secrets (Optional)

If you want the GitHub Actions workflow to trigger Render deployments automatically:

1. Go to your GitHub repository → **Settings** → **Secrets and variables** → **Actions**
2. Create two new secrets:
   - `RENDER_SERVICE_ID`: Found in your Render service URL or dashboard
   - `RENDER_API_TOKEN`: Generate from Render → **Account** → **API Tokens**

## Deployment Architecture

```
┌─────────────────┐
│  Local Machine  │
│                 │
│ docker-compose  │ (Development)
│      up         │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  GitHub Repo    │
│                 │
│  Push to master │ (Trigger)
└────────┬────────┘
         │
         ▼
┌──────────────────────────┐
│  GitHub Actions          │
│                          │
│  deploy.yml workflow     │ (CI/CD)
│  (optional automation)   │
└────────┬─────────────────┘
         │
         ▼
┌──────────────────────────┐
│  Render Platform         │
│                          │
│  - Auto-builds Docker    │ (Production)
│  - Runs docker-compose   │
│  - Serves at public URL  │
└──────────────────────────┘
```

## Environment Variables

Render will use the `render.yaml` configuration:

| Variable | Value | Purpose |
|----------|-------|---------|
| `PORT` | 10000 | Server port on Render |
| `WEB_DIR` | /app/web | Frontend static files |
| `STORAGE_PATH` | /tmp/fotoboo/photos | Photo storage |
| `DB_PATH` | /tmp/fotoboo/fotoboo.db | SQLite database |

## Running Locally with Docker Compose

To test your setup locally before deploying:

```bash
# Basic setup (local storage)
docker compose up -d --build

# With MinIO (S3-compatible storage)
USE_MINIO=true docker compose --profile minio up -d --build
```

Access the app at: `http://localhost:8080`
MinIO console: `http://localhost:9001` (if enabled)

## Monitoring Your Deployment

1. **Render Dashboard**: https://dashboard.render.com
   - View build logs
   - Monitor service health
   - View traffic metrics

2. **Health Check Endpoint**:
   ```bash
   curl https://your-fotoboo-url.onrender.com/health
   ```

## Automatic Deployment Flow

Every time you push to the `master` branch:

1. GitHub Actions workflow runs
2. Triggers Render deployment
3. Render pulls latest code
4. Builds Docker image
5. Starts `docker-compose` services
6. Application is live in ~2-5 minutes

## Important Notes

- **Free Tier Limitations**: Render free instances may sleep when idle
- **Storage**: `/tmp` storage is not persistent, data resets on redeploy
- **Better for Production**: Use Oracle Always Free VM for persistent storage (see `deploy/oracle/README.md`)

## Troubleshooting

### Build Fails
- Check Render build logs in dashboard
- Ensure `Dockerfile` and `docker-compose.yml` are valid
- Verify all environment variables are set

### App Won't Start
- Check Render logs: `Service logs` in dashboard
- Verify port binding (should be 10000)
- Check database path exists

### GitHub Actions Not Triggering
- Verify workflow file in `.github/workflows/deploy.yml`
- Check branch is `master`
- Confirm `RENDER_SERVICE_ID` and `RENDER_API_TOKEN` secrets are set

## Next Steps

1. ✅ Push code to GitHub
2. ✅ Create Render account and service
3. ✅ Test local deployment with docker-compose
4. ✅ Monitor your live deployment
5. (Optional) Configure GitHub Actions secrets for automation

Good luck with your deployment! 🚀
