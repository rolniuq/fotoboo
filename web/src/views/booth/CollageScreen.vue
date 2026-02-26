<template>
  <div class="collage-screen">
    <!-- Layout Selection -->
    <div class="layout-selector">
      <h3>Choose Layout</h3>
      <div class="layout-options">
        <button
          v-for="key in collage.layoutKeys"
          :key="key"
          :class="['layout-btn', { active: collage.currentLayout.value === key }]"
          @click="collage.setLayout(key)"
        >
          <div :class="['layout-preview', `layout-${key}`]"></div>
          <span>{{ collage.layouts[key].name }}</span>
        </button>
      </div>
    </div>

    <!-- Camera / Preview Area -->
    <div class="capture-area">
      <div class="camera-container" v-if="!collage.isComplete.value">
        <video ref="videoRef" autoplay playsinline class="camera-video"></video>
        
        <!-- Countdown Overlay -->
        <Transition name="pulse">
          <div v-if="countdown.isRunning.value" class="countdown-overlay">
            {{ countdown.count.value }}
          </div>
        </Transition>

        <!-- Flash -->
        <div :class="['flash-overlay', { active: showFlash }]"></div>
      </div>

      <div class="collage-preview" v-else>
        <canvas ref="collageCanvasRef" class="collage-canvas"></canvas>
      </div>
    </div>

    <!-- Progress Indicator -->
    <div class="progress-info">
      <span>Photos: {{ collage.capturedCount.value }} / {{ collage.requiredPhotos.value }}</span>
      <div class="thumbnail-strip">
        <div
          v-for="(photo, idx) in collage.photos.value"
          :key="idx"
          class="thumb"
        >
          <img :src="photo" :alt="`Photo ${idx + 1}`" />
          <button class="thumb-remove" @click="removeAndRetake(idx)">×</button>
        </div>
        <div
          v-for="n in (collage.requiredPhotos.value - collage.capturedCount.value)"
          :key="`empty-${n}`"
          class="thumb thumb-empty"
        >
          {{ collage.capturedCount.value + n }}
        </div>
      </div>
    </div>

    <!-- Controls -->
    <div class="controls">
      <button
        v-if="!collage.isComplete.value"
        class="btn btn-capture"
        :disabled="!camera.isReady.value || countdown.isRunning.value"
        @click="handleCapture"
      >
        Capture {{ collage.capturedCount.value + 1 }} of {{ collage.requiredPhotos.value }}
      </button>

      <template v-else>
        <button class="btn btn-secondary" @click="handleRetake">
          Retake All
        </button>
        <button class="btn btn-primary" :disabled="isSaving" @click="handleSave">
          {{ isSaving ? 'Saving...' : 'Save Collage' }}
        </button>
      </template>

      <button class="btn btn-secondary btn-sm" @click="$router.push({ name: 'welcome' })">
        Back
      </button>
    </div>

    <p v-if="camera.error.value" class="error-text">{{ camera.error.value }}</p>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useCamera } from '../../composables/useCamera'
import { useCountdown } from '../../composables/useCountdown'
import { useCollage } from '../../composables/useCollage'
import { useBoothStore } from '../../stores/booth'

const router = useRouter()
const booth = useBoothStore()
const camera = useCamera()
const countdown = useCountdown(3)
const collage = useCollage()

const showFlash = ref(false)
const isSaving = ref(false)
const videoRef = camera.videoRef
const collageCanvasRef = ref(null)

onMounted(() => {
  camera.start()
})

onUnmounted(() => {
  camera.stop()
})

// Re-render collage preview when complete
watch(
  () => collage.isComplete.value,
  async (complete) => {
    if (complete && collageCanvasRef.value) {
      await collage.renderCollage(collageCanvasRef.value)
    }
  }
)

async function handleCapture() {
  await countdown.start()

  // Flash
  showFlash.value = true

  const dataUrl = camera.captureFrame()
  if (dataUrl) {
    collage.addPhoto(dataUrl)
  }

  setTimeout(() => {
    showFlash.value = false
  }, 300)
}

function removeAndRetake(index) {
  collage.removePhoto(index)
  // Camera should already be running
}

function handleRetake() {
  collage.clearPhotos()
}

async function handleSave() {
  if (!collageCanvasRef.value) return

  isSaving.value = true
  try {
    await collage.renderCollage(collageCanvasRef.value)
    const blob = await collage.getCanvasBlob(collageCanvasRef.value)
    await booth.savePhoto(blob)
    router.push({ name: 'result' })
  } catch (err) {
    console.error('Failed to save collage:', err)
  } finally {
    isSaving.value = false
  }
}
</script>

<style scoped>
.collage-screen {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1rem;
  height: 100%;
  overflow-y: auto;
  animation: fadeIn 0.3s ease;
}

.layout-selector {
  margin-bottom: 1rem;
  text-align: center;
}

.layout-selector h3 {
  margin-bottom: 0.5rem;
  color: var(--color-white);
  font-size: 1rem;
}

.layout-options {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  justify-content: center;
}

.layout-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.3rem;
  padding: 0.5rem;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  cursor: pointer;
  transition: all 0.2s;
}

.layout-btn.active {
  border-color: var(--color-primary);
  background: rgba(102, 126, 234, 0.3);
}

.layout-btn:hover {
  border-color: rgba(255, 255, 255, 0.5);
}

.layout-btn span {
  font-size: 0.7rem;
}

.layout-preview {
  width: 50px;
  height: 40px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
  position: relative;
}

/* Layout preview patterns */
.layout-single { background: rgba(255, 255, 255, 0.4); }

.layout-horizontal2 {
  background: linear-gradient(to right, rgba(255,255,255,0.4) 49%, #333 49%, #333 51%, rgba(255,255,255,0.4) 51%);
}

.layout-vertical2 {
  background: linear-gradient(to bottom, rgba(255,255,255,0.4) 49%, #333 49%, #333 51%, rgba(255,255,255,0.4) 51%);
}

.layout-grid4 {
  background:
    linear-gradient(to right, transparent 49%, #333 49%, #333 51%, transparent 51%),
    linear-gradient(to bottom, transparent 49%, #333 49%, #333 51%, transparent 51%),
    rgba(255,255,255,0.4);
}

.layout-strip4 {
  background:
    repeating-linear-gradient(
      to bottom,
      rgba(255,255,255,0.4) 0%,
      rgba(255,255,255,0.4) 24%,
      #333 24%,
      #333 26%
    );
}

.layout-featured3 {
  background: linear-gradient(to right, rgba(255,255,255,0.4) 59%, #333 59%, #333 61%, rgba(255,255,255,0.3) 61%);
}

.capture-area {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 1rem;
}

.camera-container {
  position: relative;
  width: 640px;
  max-width: 90vw;
  height: 480px;
  max-height: 50vh;
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

.collage-preview {
  width: 640px;
  max-width: 90vw;
  aspect-ratio: 4 / 3;
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-lg);
}

.collage-canvas {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.countdown-overlay {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 6rem;
  font-weight: bold;
  color: var(--color-white);
  text-shadow: 4px 4px 8px rgba(0, 0, 0, 0.8);
  z-index: 20;
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

.progress-info {
  text-align: center;
  margin-bottom: 1rem;
  color: var(--color-white);
}

.thumbnail-strip {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
  margin-top: 0.5rem;
}

.thumb {
  width: 80px;
  height: 60px;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumb-remove {
  position: absolute;
  top: 2px;
  right: 2px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  border: none;
  cursor: pointer;
  font-size: 14px;
  line-height: 1;
}

.thumb-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  color: rgba(255, 255, 255, 0.5);
}

.btn-capture {
  background: var(--color-accent);
  color: var(--color-white);
  font-size: 1.2rem;
  padding: 1rem 2rem;
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

@keyframes flash-effect {
  0% { opacity: 0; }
  50% { opacity: 0.8; }
  100% { opacity: 0; }
}
</style>
