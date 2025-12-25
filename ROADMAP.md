# 📸 FotoBoo Project – Development Roadmap

This document describes a **step-by-step development roadmap** for the **FotoBoo project**, from MVP to production-ready. The goal is to build a system that is **stable, extensible, and practical for real-world deployment** (events, stores, or SaaS).

---

## 🎯 Overall Goals

Users should be able to:

- Capture photos (camera / webcam)
- Preview results
- Choose layouts / filters
- Export photos (download / QR / print)

The system should:

- Run reliably for long sessions
- Be easy to extend (themes, filters, payment, analytics)
- Have a clean backend architecture with clear domain boundaries

---

## 🧩 High-Level Architecture

```text
[Camera / Browser]
        ↓
[Frontend (Web / Touch UI)]
        ↓
[Backend API]
        ↓
[Storage] – [Database]
```

---

## 🗂️ Phase 0 – Planning & Foundation

### 0.1 Define MVP Scope

**Included (MVP):**

- Photo capture
- Photo preview
- Save photo
- Download photo

**Excluded (initially):**

- Payment
- Authentication
- Analytics

### 0.2 Suggested Tech Stack

**Frontend**

- Web-based (touch-friendly)
- HTML/CSS/JS or React
- Camera API (`getUserMedia`)

**Backend**

- Go (REST API)
- Clean Architecture / Hexagonal Architecture

**Storage**

- Local filesystem (MVP)
- Later: S3-compatible storage

**Database**

- SQLite (MVP)
- PostgreSQL (Production)

---

## 🚀 Phase 1 – MVP Core Features

### 1.1 Camera & Capture

- [ ] Connect to webcam
- [ ] Capture photo
- [ ] Show preview
- [ ] Retake photo

### 1.2 Backend API (Go)

**Core APIs**

- `POST /photos` – upload photo
- `GET /photos/{id}` – retrieve photo

**Core Domains**

- Photo
- Session (one photo-taking flow)

```text
Photo
 ├── id
 ├── file_path
 ├── created_at
```

### 1.3 Storage

- Store images locally (`/data/photos/...`)
- Use UUID-based file naming

### 1.4 Basic UI Flow

- Welcome screen
- Capture screen
- Preview screen
- Download screen

---

## 🎨 Phase 2 – UX & Photo Enhancement

### 2.1 Layouts & Frames

- [ ] Fixed photo frames (2–4 layouts)
- [ ] Overlay logo / event name

### 2.2 Filters & Image Processing

Basic filters:

- Grayscale
- Vintage
- Brightness / Contrast

Processing options:

- Frontend (Canvas)
- Backend (image processing libraries)

### 2.3 Countdown & Animations

- 3–2–1 countdown
- Flash animation on capture

---

## 🖨️ Phase 3 – Output & Sharing

### 3.1 Download & QR Code

- [ ] Photo download
- [ ] Generate QR code for photo URL

### 3.2 Print Integration

- Print-ready image format (size, DPI)
- Send jobs to printer

---

## 🧠 Phase 4 – Advanced Backend

### 4.1 Clean Architecture Structure

```text
/fotoboo
 ├── cmd/
 ├── internal/
 │   ├── domain/
 │   ├── usecase/
 │   ├── handler/
 │   ├── repository/
 ├── pkg/
```

### 4.2 Database Design

- Store metadata:
  - sessions
  - photos
  - devices

### 4.3 Background Jobs

- Cleanup old photos
- Async image resizing

---

## 🔐 Phase 5 – Admin & Operations

### 5.1 Admin Dashboard

- Total photos captured
- Storage usage
- Usage per day

### 5.2 Configuration Management

- Layout / frame selection
- Event name
- Countdown duration

---

## 📈 Phase 6 – Production & Scaling

### 6.1 Cloud Storage

- S3 / MinIO
- CDN integration

### 6.2 Observability

- Logging
- Metrics
- Health checks

### 6.3 Performance & Stability

- Limit concurrent sessions
- Rate limiting

---

## 🧪 Phase 7 – Testing & Quality

- Unit tests (use cases)
- Integration tests (API)
- Manual testing (UI + camera)

---

## 🗺️ Recommended Implementation Order

1. Camera + capture
2. Backend photo upload
3. Preview + download
4. Layouts / frames
5. QR code sharing
6. Printing
7. Admin dashboard

---

## ✅ MVP Definition of Done

- Stable photo capture
- No crashes during continuous usage
- Simple, intuitive UX
- Deployable and demo-ready

---

## 📌 Notes

- Prioritize **simplicity and reliability** in MVP
- Avoid over-engineering early
- Keep the domain clean to enable future scaling

---

📸 _Build for usability first, then scale with confidence._
