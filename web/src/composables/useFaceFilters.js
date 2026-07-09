import { ref, computed, shallowRef } from 'vue'

/**
 * Overlay definitions for face effects.
 * Each overlay has:
 *   - label: Human-readable name
 *   - svg: Path to SVG asset in public/overlays/
 */
const OVERLAYS = {
  none: { label: 'None' },
  'dog-ears': { label: 'Dog Ears', svg: '/overlays/dog-ears.svg' },
  'cat-ears': { label: 'Cat Ears', svg: '/overlays/cat-ears.svg' },
  'fox-ears': { label: 'Fox Ears', svg: '/overlays/fox-ears.svg' },
  glasses: { label: 'Glasses', svg: '/overlays/glasses.svg' },
  sunglasses: { label: 'Sunglasses', svg: '/overlays/sunglasses.svg' },
  mustache: { label: 'Mustache', svg: '/overlays/mustache.svg' },
  crown: { label: 'Crown', svg: '/overlays/crown.svg' },
  'clown-nose': { label: 'Clown Nose', svg: '/overlays/clown-nose.svg' },
  halo: { label: 'Halo', svg: '/overlays/halo.svg' },
}

/**
 * Composable for face overlay effects with REAL-TIME video support (TikTok-style).
 *
 * Two modes:
 *   1. STATIC MODE: detectFace() + applyOverlay() for images (PreviewScreen)
 *   2. LIVE MODE: startVideoOverlay() for real-time video (CaptureScreen)
 *
 * Uses the native FaceDetector API when available, falls back
 * to proportional placement.
 */
export function useFaceFilters() {
  const currentOverlay = ref('none')
  const faceBounds = ref(null)
  const isDetecting = ref(false)
  const hasDetected = ref(false)
  const isOverlayRunning = ref(false)

  // Cache loaded overlay SVG Image objects
  const imageCache = new Map()

  // Animation + interval handles for live mode
  let animationFrameId = null
  let detectIntervalId = null

  // --- Computed ---

  const overlayKeys = computed(() => Object.keys(OVERLAYS))
  const currentOverlayLabel = computed(() => OVERLAYS[currentOverlay.value]?.label || 'None')
  const hasOverlay = computed(() => currentOverlay.value !== 'none')
  const isFaceDetectSupported = computed(() => typeof window.FaceDetector !== 'undefined')

  // --- Private helpers ---

  function loadImageAsPromise(src) {
    return new Promise((resolve, reject) => {
      const img = new Image()
      img.onload = () => resolve(img)
      img.onerror = reject
      img.src = src
    })
  }

  /**
   * Load (or retrieve from cache) an overlay SVG as an HTMLImageElement.
   */
  async function loadOverlayImage(overlayKey) {
    if (imageCache.has(overlayKey)) return imageCache.get(overlayKey)

    const meta = OVERLAYS[overlayKey]
    if (!meta || !meta.svg) return null

    try {
      const img = await loadImageAsPromise(meta.svg)
      imageCache.set(overlayKey, img)
      return img
    } catch {
      console.warn(`Failed to load overlay: ${overlayKey}`)
      return null
    }
  }

  /**
   * Proportional fallback: estimate face bounds assuming a centered selfie.
   */
  function getProportionalBounds(imgWidth, imgHeight) {
    const faceWidth = imgWidth * 0.32
    const faceHeight = imgHeight * 0.45
    return {
      x: (imgWidth - faceWidth) / 2,
      y: imgHeight * 0.12,
      width: faceWidth,
      height: faceHeight,
    }
  }

  /**
   * Calculate overlay placement rectangle from face bounds.
   * Each overlay type has specific positioning logic.
   */
  function getPlacement(overlayKey, bounds) {
    const { x, y, width, height } = bounds

    switch (overlayKey) {
      case 'dog-ears':
        return { x: x - width * 0.4, y: y - height * 0.4, w: width * 1.8, h: height * 1.2 }
      case 'cat-ears':
      case 'fox-ears':
        return { x: x - width * 0.4, y: y - height * 0.6, w: width * 1.8, h: height * 0.9 }
      case 'glasses':
        return { x: x - width * 0.05, y: y + height * 0.2, w: width * 1.1, h: height * 0.4 }
      case 'sunglasses':
        return { x: x - width * 0.05, y: y + height * 0.22, w: width * 1.1, h: height * 0.35 }
      case 'mustache':
        return { x: x + width * 0.1, y: y + height * 0.55, w: width * 0.8, h: height * 0.25 }
      case 'crown':
        return { x: x - width * 0.25, y: y - height * 0.55, w: width * 1.5, h: height * 0.9 }
      case 'clown-nose':
        return { x: x + width * 0.1, y: y + height * 0.42, w: width * 0.4, h: height * 0.22 }
      case 'halo':
        return { x: x - width * 0.3, y: y - height * 0.35, w: width * 1.6, h: height * 0.5 }
      default:
        return { x, y, w: width, h: height }
    }
  }

  /**
   * Synchronously draw an overlay on a canvas at the given placement.
   * The overlay image MUST be in the cache (call loadOverlayImage first).
   */
  function drawOverlaySync(ctx, overlayKey, placement) {
    const img = imageCache.get(overlayKey)
    if (!img) return
    ctx.drawImage(img, placement.x, placement.y, placement.w, placement.h)
  }

  // ================================================================
  //  STATIC MODE — for still images (PreviewScreen)
  // ================================================================

  /**
   * Detect face bounds from an image source.
   * Falls back to proportional estimation if FaceDetector unavailable.
   */
  async function detectFace(imageSrc) {
    if (hasDetected.value && faceBounds.value) return faceBounds.value

    isDetecting.value = true

    try {
      const img = await loadImageAsPromise(imageSrc)

      if (window.FaceDetector) {
        try {
          const detector = new FaceDetector({ maxDetectedFaces: 1 })
          const faces = await detector.detect(img)
          if (faces.length > 0) {
            const bb = faces[0].boundingBox
            faceBounds.value = { x: bb.x, y: bb.y, width: bb.width, height: bb.height }
          } else {
            faceBounds.value = getProportionalBounds(img.width, img.height)
          }
        } catch {
          faceBounds.value = getProportionalBounds(img.width, img.height)
        }
      } else {
        faceBounds.value = getProportionalBounds(img.width, img.height)
      }

      hasDetected.value = true
      return faceBounds.value
    } catch (err) {
      console.warn('Face detection error:', err)
      return null
    } finally {
      isDetecting.value = false
    }
  }

  /**
   * Draw a face overlay onto a canvas at the correct position (async load).
   */
  async function applyOverlay(canvas, overlayKey, bounds) {
    if (overlayKey === 'none') return
    const img = await loadOverlayImage(overlayKey)
    if (!img) return

    const ctx = canvas.getContext('2d')
    const placement = getPlacement(overlayKey, bounds)
    ctx.drawImage(img, placement.x, placement.y, placement.w, placement.h)
  }

  /**
   * Convenience: detect face then apply overlay to canvas.
   */
  async function detectAndApply(canvas, imageSrc) {
    if (currentOverlay.value === 'none') return
    if (!faceBounds.value) await detectFace(imageSrc)
    if (faceBounds.value) await applyOverlay(canvas, currentOverlay.value, faceBounds.value)
  }

  // ================================================================
  //  LIVE MODE — real-time video overlay (CaptureScreen, TikTok-style)
  // ================================================================

  /**
   * Select an overlay and preload its image (for live mode).
   * Returns a promise that resolves when the image is cached.
   */
  async function selectOverlay(key) {
    currentOverlay.value = key
    if (key !== 'none') {
      await loadOverlayImage(key)
    }
  }

  /**
   * Start real-time face detection + overlay rendering on a live video feed.
   *
   * Architecture:
   *   - Face detection runs every ~250ms via FaceDetector API
   *   - Overlay drawing runs every animation frame (~60fps) at last known position
   *   - Canvas is transparent; positioned on top of the <video> element
   *
   * @param {HTMLVideoElement} videoElement
   * @param {HTMLCanvasElement} canvasElement - transparent overlay canvas
   */
  function startVideoOverlay(videoElement, canvasElement) {
    if (isOverlayRunning.value) return
    if (!videoElement || !canvasElement) return

    isOverlayRunning.value = true
    const ctx = canvasElement.getContext('2d')

    // ---- Face detection loop (every ~250ms) ----
    async function detectLoop() {
      while (isOverlayRunning.value) {
        if (
          videoElement.readyState >= HTMLMediaElement.HAVE_CURRENT_DATA &&
          currentOverlay.value !== 'none'
        ) {
          try {
            // Native Shape Detection API
            if (window.FaceDetector) {
              const detector = new FaceDetector({ maxDetectedFaces: 1 })
              const faces = await detector.detect(videoElement)
              if (faces.length > 0) {
                const bb = faces[0].boundingBox
                faceBounds.value = { x: bb.x, y: bb.y, width: bb.width, height: bb.height }
              }
            }
            // If no FaceDetector, we keep the proportional fallback from last call
          } catch {
            // Silently continue — overlay renders at last known (or fallback) position
          }
        }
        await new Promise((r) => setTimeout(r, 250))
      }
    }
    detectLoop()

    // Initial proportional fallback (works even without FaceDetector)
    if (videoElement.videoWidth > 0) {
      faceBounds.value = getProportionalBounds(videoElement.videoWidth, videoElement.videoHeight)
    }

    // ---- Render loop (every animation frame, ~60fps) ----
    function renderLoop() {
      if (!isOverlayRunning.value) return

      // Match canvas to video dimensions
      const vw = videoElement.videoWidth
      const vh = videoElement.videoHeight
      if (canvasElement.width !== vw || canvasElement.height !== vh) {
        canvasElement.width = vw
        canvasElement.height = vh
      }

      // Clear previous overlay
      ctx.clearRect(0, 0, vw, vh)

      // Draw overlay at last known face position
      if (currentOverlay.value !== 'none' && faceBounds.value) {
        const placement = getPlacement(currentOverlay.value, faceBounds.value)
        drawOverlaySync(ctx, currentOverlay.value, placement)
      }

      animationFrameId = requestAnimationFrame(renderLoop)
    }
    renderLoop()
  }

  /**
   * Stop the real-time video overlay rendering + detection loops.
   */
  function stopVideoOverlay() {
    isOverlayRunning.value = false
    if (animationFrameId) {
      cancelAnimationFrame(animationFrameId)
      animationFrameId = null
    }
    // The detect loop uses its own async while-loop; it will exit on next tick
    // because isOverlayRunning is false.
  }

  /**
   * Capture the current video frame WITH the face overlay baked in.
   * Used when the user presses the shutter button.
   *
   * @param {HTMLVideoElement} videoElement
   * @returns {string|null} data URL (JPEG) with overlay, or null on failure
   */
  function captureVideoFrame(videoElement) {
    if (!videoElement || videoElement.readyState < HTMLMediaElement.HAVE_CURRENT_DATA) return null

    const w = videoElement.videoWidth
    const h = videoElement.videoHeight
    if (!w || !h) return null

    const canvas = document.createElement('canvas')
    canvas.width = w
    canvas.height = h
    const ctx = canvas.getContext('2d')

    // Draw current video frame
    ctx.drawImage(videoElement, 0, 0)

    // Draw overlay at last known position
    if (currentOverlay.value !== 'none' && faceBounds.value) {
      const placement = getPlacement(currentOverlay.value, faceBounds.value)
      drawOverlaySync(ctx, currentOverlay.value, placement)
    }

    return canvas.toDataURL('image/jpeg', 0.9)
  }

  /**
   * Capture the current video frame WITHOUT overlay (raw).
   * Used for storing the raw photo that can be re-rendered with different overlays.
   */
  function captureRawVideoFrame(videoElement) {
    if (!videoElement || videoElement.readyState < HTMLMediaElement.HAVE_CURRENT_DATA) return null

    const w = videoElement.videoWidth
    const h = videoElement.videoHeight
    if (!w || !h) return null

    const canvas = document.createElement('canvas')
    canvas.width = w
    canvas.height = h
    const ctx = canvas.getContext('2d')
    ctx.drawImage(videoElement, 0, 0)

    return canvas.toDataURL('image/jpeg', 0.9)
  }

  // ================================================================
  //  UNIVERSAL
  // ================================================================

  function reset() {
    stopVideoOverlay()
    currentOverlay.value = 'none'
    faceBounds.value = null
    isDetecting.value = false
    hasDetected.value = false
    // Keep imageCache for next use
  }

  return {
    // State
    currentOverlay,
    faceBounds,
    isDetecting,
    hasDetected,
    isOverlayRunning,

    // Computed
    overlayKeys,
    currentOverlayLabel,
    hasOverlay,
    isFaceDetectSupported,
    overlays: OVERLAYS,
    imageCache,

    // Static mode
    detectFace,
    applyOverlay,
    detectAndApply,

    // Live video mode
    selectOverlay,
    startVideoOverlay,
    stopVideoOverlay,
    captureVideoFrame,
    captureRawVideoFrame,

    // Universal
    reset,
  }
}
