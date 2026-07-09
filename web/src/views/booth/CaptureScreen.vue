<template>
  <div class="capture-screen">
    <div class="camera-container">
      <video ref="videoRef" autoplay playsinline class="camera-video"></video>

      <!-- Face Overlay Canvas (transparent, on top of video) -->
      <canvas ref="overlayCanvasRef" class="overlay-canvas"></canvas>

      <!-- Countdown Overlay -->
      <Transition name="pulse">
        <div v-if="countdown.isRunning.value" class="countdown-overlay">
          {{ countdown.count.value }}
        </div>
      </Transition>

      <!-- Flash -->
      <div :class="['flash-overlay', { active: showFlash }]"></div>
    </div>

    <!-- Face Effect Selector (TikTok-style pill tray) -->
    <div class="face-effect-tray">
      <button
        v-for="key in faceFilters.overlayKeys.value"
        :key="key"
        :class="['effect-pill', { active: faceFilters.currentOverlay.value === key }]"
        @click="selectEffect(key)"
        :disabled="countdown.isRunning.value"
      >
        {{ faceFilters.overlays[key].label }}
      </button>
    </div>

    <div class="controls">
      <button
        class="btn btn-capture"
        :disabled="!camera.isReady.value || countdown.isRunning.value"
        @click="handleCapture"
      >
        Capture
      </button>
      <button class="btn btn-secondary" @click="$router.push({ name: 'welcome' })">
        Back
      </button>
    </div>

    <p v-if="camera.error.value" class="error-text">{{ camera.error.value }}</p>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCamera } from '../../composables/useCamera'
import { useCountdown } from '../../composables/useCountdown'
import { useFaceFilters } from '../../composables/useFaceFilters'
import { useBoothStore } from '../../stores/booth'

const router = useRouter()
const booth = useBoothStore()
const camera = useCamera()
const countdown = useCountdown(3)
const faceFilters = useFaceFilters()
const showFlash = ref(false)

// Bind refs
const videoRef = camera.videoRef
const overlayCanvasRef = ref(null)

onMounted(async () => {
  // Restore previously selected overlay (if coming from Preview "Retake")
  const prevOverlay = booth.overlayKey
  if (prevOverlay && prevOverlay !== 'none') {
    await faceFilters.selectOverlay(prevOverlay)
  }

  await camera.start()

  // Start real-time face overlay after camera is ready
  if (overlayCanvasRef.value) {
    faceFilters.startVideoOverlay(videoRef.value, overlayCanvasRef.value)
  }
})

onUnmounted(() => {
  faceFilters.stopVideoOverlay()
  camera.stop()
})

async function selectEffect(key) {
  await faceFilters.selectOverlay(key)
  booth.setOverlayKey(key)
}

async function handleCapture() {
  await countdown.start()

  // Capture raw frame (without overlay) for possible re-render
  const rawDataUrl = faceFilters.captureRawVideoFrame(videoRef.value)

  // Capture final frame with overlay baked in
  const finalDataUrl = faceFilters.captureVideoFrame(videoRef.value)

  if (finalDataUrl && rawDataUrl) {
    booth.setRawCapturedPhoto(rawDataUrl)
    booth.setCapturedPhoto(finalDataUrl)
  } else if (finalDataUrl) {
    // Fallback: just use the overlay frame
    booth.setCapturedPhoto(finalDataUrl)
  }

  // Flash
  showFlash.value = true

  setTimeout(() => {
    showFlash.value = false
    router.push({ name: 'preview' })
  }, 300)
}
</script>

<style scoped>
.capture-screen {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  animation: fadeIn 0.3s ease;
}

.camera-container {
  position: relative;
  width: 640px;
  max-width: 90vw;
  height: 480px;
  max-height: 55vh;
  background: #000;
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-lg);
}

.camera-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* Transparent canvas overlay on top of video for face effects */
.overlay-canvas {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 5;
  pointer-events: none;
}

.countdown-overlay {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 8rem;
  font-weight: bold;
  color: var(--color-white);
  text-shadow: 4px 4px 8px rgba(0, 0, 0, 0.8);
  z-index: 20;
  animation: countdown-pulse 1s ease-in-out;
}

.flash-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: white;
  opacity: 0;
  z-index: 15;
  pointer-events: none;
}

.flash-overlay.active {
  animation: flash-effect 0.3s ease-out;
}

/* Face effect pill tray */
.face-effect-tray {
  display: flex;
  gap: 0.4rem;
  flex-wrap: wrap;
  justify-content: center;
  padding: 0.6rem 1rem;
  margin-top: 0.5rem;
  max-width: 90vw;
  background: rgba(0, 0, 0, 0.5);
  border-radius: var(--radius-lg);
  backdrop-filter: blur(8px);
}

.effect-pill {
  padding: 0.35rem 0.85rem;
  border: 2px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
  border-radius: var(--radius-lg);
  cursor: pointer;
  font-weight: 500;
  font-size: 0.8rem;
  transition: all 0.2s;
  white-space: nowrap;
}

.effect-pill.active {
  background: var(--color-primary);
  border-color: var(--color-primary);
  color: white;
}

.effect-pill:hover:not(:disabled) {
  border-color: rgba(255, 255, 255, 0.6);
  transform: translateY(-1px);
}

.effect-pill:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-capture {
  background: var(--color-accent);
  color: var(--color-white);
  font-size: 1.3rem;
  padding: 1.2rem 2.5rem;
}

.btn-capture:hover:not(:disabled) {
  background: #ff5252;
  transform: scale(1.05);
}

.error-text {
  color: #ff8a80;
  margin-top: 1rem;
  font-size: 0.9rem;
}

.pulse-enter-active { animation: countdown-pulse 1s ease-in-out; }
.pulse-leave-active { opacity: 0; transition: opacity 0.1s; }
</style>
