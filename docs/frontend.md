# Frontend Architecture

Documentation for the FotoBoo web-based frontend — a vanilla JavaScript single-page application (SPA) with no build step.

---

## Overview

The frontend is a **browser-based photo booth interface** that:

- Captures photos via the device webcam (WebRTC)
- Applies filters and frames using Canvas 2D API
- Uploads photos to the backend API
- Provides download functionality

**Key design decisions:**
- No JavaScript framework (vanilla JS)
- No build tools (no webpack, vite, etc.)
- Single HTML file with 4 screen states
- One JavaScript class manages all application logic

---

## File Structure

```
web/
├── index.html              # Single-page app (4 screens in one file)
└── static/
    ├── css/
    │   └── style.css       # All styles + animations (366 lines)
    └── js/
        └── app.js          # FotoBooApp class (309 lines)
```

---

## Screen Flow

The UI consists of 4 screens, shown/hidden via CSS classes:

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Welcome    │───>│   Capture   │───>│   Preview   │───>│   Result    │
│              │    │             │    │             │    │             │
│  "Start      │    │  Live cam   │    │  Filters    │    │  Final      │
│   Session"   │    │  Countdown  │    │  Frames     │    │  Download   │
│              │    │  Capture    │    │  Sliders    │    │  New photo  │
└─────────────┘    └──────┬──────┘    │  Retake     │    └──────┬──────┘
                          ▲           │  Save       │           │
                          │           └──────┬──────┘           │
                          │                  │                  │
                          └──── (Retake) ────┘                  │
                          ▲                                     │
                          └──────── (Take Another) ─────────────┘
```

### Screen Visibility

Each screen is a `<div>` with an ID. Only one screen is visible at a time, controlled by the `.active` CSS class:

```html
<div id="welcome-screen" class="screen active">...</div>
<div id="capture-screen" class="screen">...</div>
<div id="preview-screen" class="screen">...</div>
<div id="result-screen" class="screen">...</div>
```

```javascript
showScreen(screenId) {
    document.querySelectorAll('.screen').forEach(s => s.classList.remove('active'));
    document.getElementById(screenId).classList.add('active');
}
```

---

## FotoBooApp Class

All frontend logic lives in a single class instantiated on page load:

```javascript
document.addEventListener('DOMContentLoaded', () => {
    new FotoBooApp();
});
```

### Key Properties

| Property | Type | Description |
|----------|------|-------------|
| `video` | HTMLVideoElement | Live camera feed element |
| `canvas` | HTMLCanvasElement | Used for photo capture and processing |
| `capturedImage` | string (dataURL) | The captured photo as base64 JPEG |
| `currentPhotoId` | string (UUID) | ID returned from server after upload |
| `stream` | MediaStream | Active camera stream reference |

### Key Methods

| Method | Description |
|--------|-------------|
| `startCamera()` | Opens webcam via `getUserMedia` at 1280x720 |
| `stopCamera()` | Stops all media stream tracks |
| `capturePhoto()` | Runs 3-2-1 countdown, captures canvas frame |
| `applyFilters()` | Renders photo to canvas with current filter settings |
| `getFilterString()` | Builds CSS filter string from current selections |
| `savePhoto()` | Converts dataURL → blob → POST to `/photos` |
| `downloadPhoto()` | Creates download link for `/photos/{id}` |
| `showScreen(id)` | Navigates to a screen by toggling `.active` class |

---

## Camera Integration

### WebRTC Setup

```javascript
const stream = await navigator.mediaDevices.getUserMedia({
    video: {
        width: { ideal: 1280 },
        height: { ideal: 720 },
        facingMode: 'user'   // Front-facing camera
    },
    audio: false
});
this.video.srcObject = stream;
```

### Photo Capture

1. Countdown animation (3 → 2 → 1)
2. Flash effect (white overlay)
3. Draw video frame to canvas
4. Export as JPEG data URL (`canvas.toDataURL('image/jpeg', 0.9)`)

---

## Image Processing

### Filters

All filters use the Canvas 2D `context.filter` property with CSS filter functions:

| Filter | CSS Filter Applied |
|--------|-------------------|
| None | (none) |
| Grayscale | `grayscale(100%)` |
| Vintage | `sepia(50%) contrast(120%) brightness(90%)` |
| Bright | `brightness(130%)` |
| Contrast | `contrast(150%)` |

### Adjustments

| Control | Range | Default | CSS Filter |
|---------|-------|---------|------------|
| Brightness | 50% – 150% | 100% | `brightness(value%)` |
| Contrast | 50% – 150% | 100% | `contrast(value%)` |

### Filter Composition

Filters are composed into a single CSS filter string:

```javascript
getFilterString() {
    let filter = `brightness(${this.brightness}%) contrast(${this.contrast}%)`;
    // Append preset filter if selected
    if (this.currentFilter === 'grayscale') filter += ' grayscale(100%)';
    // ...
    return filter;
}
```

### Frames

Frames are applied via CSS classes on the preview container:

| Frame | CSS Class | Visual |
|-------|-----------|--------|
| None | (none) | No border |
| Simple | `.frame-simple` | White border |
| Event | `.frame-event` | Gold border + "Special Event" text |
| Party | `.frame-party` | Red border + 🎉 emoji |

---

## API Communication

### Photo Upload

```javascript
async savePhoto() {
    // Convert data URL to blob
    const response = await fetch(this.capturedImage);
    const blob = await response.blob();

    // Upload to backend
    const result = await fetch('/photos', {
        method: 'POST',
        body: blob
    });
    const data = await result.json();
    this.currentPhotoId = data.id;
}
```

### Photo Download

```javascript
downloadPhoto() {
    const link = document.createElement('a');
    link.href = `/photos/${this.currentPhotoId}`;
    link.download = `fotoboo-${this.currentPhotoId}.jpg`;
    link.click();
}
```

---

## Animations & Effects

### Countdown

- Numbers (3, 2, 1) displayed with CSS scale animation
- Each number shows for 1 second with a pulse effect

### Flash

- Full-screen white overlay
- Opacity transitions from 0 → 1 → 0
- Simulates camera flash

### Screen Transitions

```css
.screen {
    opacity: 0;
    visibility: hidden;
    transition: opacity 0.3s ease, visibility 0.3s ease;
}

.screen.active {
    opacity: 1;
    visibility: visible;
}
```

---

## Styling

### Design System

| Element | Value |
|---------|-------|
| Background | Purple gradient (`#667eea` → `#764ba2`) |
| Camera container | 640×480 max, rounded corners, drop shadow |
| Buttons | White with purple text, rounded |
| Font | System default |

### Responsive Design

The layout adapts for mobile screens (≤768px):
- Controls stack vertically
- Camera container fills available width
- Touch-friendly button sizes

---

## Planned Improvements (from Roadmap)

| Feature | Phase | Description |
|---------|-------|-------------|
| QR code display | Phase 3 | Show QR code linking to photo URL (currently placeholder text) |
| More layouts | Phase 2 | Additional photo frame designs |
| Multi-photo collage | Phase 2 | Capture 2-4 photos in a single layout |
| Print button | Phase 3 | Trigger physical printing from the UI |
| Event branding | Phase 5 | Configurable logo/text overlay per event |
