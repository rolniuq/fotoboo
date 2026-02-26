# 🚀 FotoBoo Public Deployment Guide

Your code is now on GitHub with automatic deployment configured. Follow these steps to get your project **live and public** in the next 5 minutes.

## 🎯 Quick Start (5 minutes)

### Option 1: Deploy on Render (Recommended - FREE)

1. **Go to Render Dashboard**
   - Visit: https://render.com
   - Sign up with GitHub account

2. **Deploy Your Repository**
   - Click **New +** → **Blueprint**
   - Select `fotoboo` repository from your GitHub account
   - Render auto-detects `render.yaml` configuration
   - Click **Apply** to deploy

3. **Wait for Deployment**
   - Build takes ~2-5 minutes
   - You'll get a public URL like: `https://fotoboo-xxx.onrender.com`
   - Your app is **LIVE** 🎉

### Option 2: Deploy on Railway.app (Alternative - FREE)

1. **Go to Railway**
   - Visit: https://railway.app
   - Sign up with GitHub

2. **Create New Project**
   - Connect your `fotoboo` repository
   - Railway auto-detects Docker configuration
   - Deploy starts automatically

3. **Get Public URL**
   - Your app will be live in 2-3 minutes
   - URL format: `https://fotoboo-xxxxx.up.railway.app`

### Option 3: Deploy on Fly.io (FREE tier)

1. **Install Fly CLI**
   ```bash
   curl -L https://fly.io/install.sh | sh
   ```

2. **Launch Your App**
   ```bash
   cd /Users/quynh.vo/Workplace/fotoboo
   flyctl launch
   ```
   - Follow prompts
   - Choose a region
   - Deploy

3. **Get Public URL**
   ```bash
   flyctl open
   ```

---

## 📊 Comparison Table

| Platform | Cost | Startup Time | Always On | Best For |
|----------|------|--------------|-----------|----------|
| **Render** | Free | 2-5 min | ⚠️ Sleeps idle | Simple projects |
| **Railway** | Free tier | 2-3 min | ⚠️ Limited | Development |
| **Fly.io** | Free tier | 1-2 min | ✅ Yes | Always-on apps |
| **Oracle Always Free** | Free | 5-10 min | ✅ Yes | Production |

---

## 🌍 What People Can Do

Once deployed, users can:

1. **Upload Photos**
   ```bash
   curl -X POST https://your-app.onrender.com/photos \
     --data-binary @photo.jpg
   ```

2. **Retrieve Photos**
   ```
   https://your-app.onrender.com/photos/{id}
   ```

3. **Health Check**
   ```bash
   curl https://your-app.onrender.com/health
   ```

---

## 🔄 Automatic Updates

Every time you push to GitHub:

```bash
git add .
git commit -m "Your changes"
git push origin master
```

Your deployment automatically updates! ✨

---

## 📝 Environment Configuration

Your app uses:
- **Port**: 10000 (Render)
- **Database**: SQLite in `/tmp/fotoboo/fotoboo.db`
- **Storage**: Local filesystem in `/tmp/fotoboo/photos`
- **API Base**: `https://your-app.onrender.com`

### To Add MinIO (S3-compatible storage):

1. Update `render.yaml`:
```yaml
services:
  - type: web
    name: fotoboo
    runtime: docker
    envVars:
      - key: USE_MINIO
        value: "true"
      - key: MINIO_ENDPOINT
        value: "your-minio-endpoint.com"
      - key: MINIO_ACCESS_KEY
        value: ${{ secrets.MINIO_KEY }}
```

2. Add secrets in platform dashboard

---

## 🛡️ Security Checklist

Before sharing with public:

- [ ] No sensitive credentials in code (check `.env` files)
- [ ] CORS is configured for your domain
- [ ] Rate limiting is enabled
- [ ] Database backups are configured
- [ ] HTTPS is enabled (all platforms provide this)

---

## 🐛 Troubleshooting

### App Won't Start
1. Check platform logs (Dashboard → Service Logs)
2. Verify Dockerfile builds locally:
   ```bash
   docker build -t fotoboo .
   docker run -p 8080:8080 fotoboo
   ```

### Build Fails
1. Ensure `go.mod` and `go.sum` are committed
2. Check Docker configuration
3. Verify all environment variables are set

### Photos Not Persisting
- This is expected on Render free tier (uses `/tmp`)
- For production, use Oracle Always Free or Fly.io persistent volumes

---

## 📈 Scaling Up (When You Get Users)

If you need persistent storage and better performance:

### Move to Production Setup

**Option A: Oracle Always Free VM** (Recommended)
- Always-on, truly free
- Persistent storage
- See: `deploy/oracle/README.md`

**Option B: Paid Tier on Render/Railway/Fly**
- Better performance
- Guaranteed uptime
- More storage

**Option C: AWS/Azure/GCP**
- Enterprise features
- But requires payment

---

## 🚀 Next Steps

**Right Now:**
1. ✅ Code is on GitHub
2. ✅ Workflow is configured
3. Pick one platform above and deploy

**After Deployment:**
1. Test your endpoints
2. Share the public URL with users
3. Monitor logs and metrics
4. Gather feedback

**When Ready for Production:**
1. Set up domain name (optional)
2. Configure custom email alerts
3. Set up database backups
4. Monitor API usage

---

## 📞 Support Resources

- **Render Docs**: https://render.com/docs
- **Railway Docs**: https://docs.railway.app
- **Fly.io Docs**: https://fly.io/docs
- **Docker Docs**: https://docs.docker.com
- **Go Docs**: https://golang.org/doc

---

## ✨ Summary

Your FotoBoo project is ready for the public! Choose your deployment platform from the options above and follow the steps. Within 5 minutes, your app will be live and people can start using it.

**Good luck! 🎉**
