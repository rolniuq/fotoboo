<template>
  <div class="animal-screen">
    <!-- Capture Mode -->
    <template v-if="mode === 'capture'">
      <div class="camera-container">
        <video ref="videoRef" autoplay playsinline class="camera-video"></video>

        <Transition name="pulse">
          <div v-if="countdown.isRunning.value" class="countdown-overlay">
            {{ countdown.count.value }}
          </div>
        </Transition>

        <div :class="['flash-overlay', { active: showFlash }]"></div>
      </div>

      <p class="mode-label">Choose your animal, then capture!</p>

      <div class="animal-tray">
        <button
          v-for="key in animalTransform.animalKeys.value"
          :key="key"
          :class="['animal-pill', { active: animalTransform.currentAnimal.value === key }]"
          @click="animalTransform.selectAnimal(key)"
          :disabled="countdown.isRunning.value"
        >
          {{ animalTransform.animals[key].label }}
        </button>
      </div>

      <div class="controls">
        <button
          class="btn btn-capture"
          :disabled="!camera.isReady.value || countdown.isRunning.value"
          @click="handleCapture"
        >
          Transform Me!
        </button>
        <button class="btn btn-secondary" @click="$router.push({ name: 'welcome' })">
          Back
        </button>
      </div>

      <p v-if="camera.error.value" class="error-text">{{ camera.error.value }}</p>
    </template>

    <!-- Preview Mode -->
    <template v-else>
      <div class="preview-container">
        <canvas ref="resultCanvasRef" class="result-canvas"></canvas>
        <p v-if="animalTransform.isProcessing.value" class="processing-label">
          Transforming...
        </p>
      </div>

      <p class="mode-label">Try a different animal or save your photo!</p>

      <div class="animal-tray">
        <button
          v-for="key in animalTransform.animalKeys.value"
          :key="key"
          :class="['animal-pill', { active: animalTransform.currentAnimal.value === key }]"
          @click="switchAnimal(key)"
          :disabled="animalTransform.isProcessing.value"
        >
          {{ animalTransform.animals[key].label }}
        </button>
      </div>

      <div class="controls">
        <button class="btn btn-secondary" @click="handleRetake">
          Retake
        </button>
        <button
          class="btn btn-primary"
          :disabled="isSaving || animalTransform.isProcessing.value"
          @click="handleSave"
        >
          {{ isSaving ? 'Saving...' : 'Save Photo' }}
        </button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCamera } from '../../composables/useCamera'
import { useCountdown } from '../../composables/useCountdown'
import { useAnimalTransform } from '../../composables/useAnimalTransform'
import { useBoothStore } from '../../stores/booth'
import { useFilters } from '../../composables/useFilters'

const router = useRouter()
const booth = useBoothStore()
const camera = useCamera()
const countdown = useCountdown(3)
const animalTransform = useAnimalTransform()
const filters = useFilters()

const mode = ref('capture')
const showFlash = ref(false)
const isSaving = ref(false)
const capturedDataUrl = ref(null)

const videoRef = camera.videoRef
const resultCanvasRef = ref(null)

onMounted(async () => {
  await camera.start()
})

onUnmounted(() => {
  camera.stop()
  animalTransform.reset()
})

async function handleCapture() {
  await countdown.start()

  const dataUrl = camera.captureFrame()
  if (!dataUrl) return

  capturedDataUrl.value = dataUrl

  showFlash.value = true
  setTimeout(() => {
    showFlash.value = false
    mode.value = 'preview'

    // Detect face + transform after switching mode (next tick for canvas render)
    setTimeout(() => runTransform(), 50)
  }, 300)
}

async function runTransform() {
  if (!resultCanvasRef.value || !capturedDataUrl.value) return

  // Detect face once for consistent placement across animal switches
  if (!animalTransform.faceBounds.value) {
    await animalTransform.detectFace(capturedDataUrl.value)
  }

  await animalTransform.transform(
    resultCanvasRef.value,
    capturedDataUrl.value,
    animalTransform.currentAnimal.value,
  )
}

async function switchAnimal(key) {
  animalTransform.selectAnimal(key)
  await runTransform()
}

function handleRetake() {
  mode.value = 'capture'
  capturedDataUrl.value = null
  animalTransform.reset()
}

async function handleSave() {
  if (!resultCanvasRef.value) return

  isSaving.value = true
  try {
    const blob = await filters.getCanvasBlob(resultCanvasRef.value)
    // Set captured photos so PreviewScreen / ResultScreen can reference them
    const dataUrl = resultCanvasRef.value.toDataURL('image/jpeg', 0.9)
    booth.setCapturedPhoto(dataUrl)
    booth.setRawCapturedPhoto(dataUrl)
    await booth.savePhoto(blob)
    router.push({ name: 'result' })
  } catch (err) {
    console.error('Failed to save animal photo:', err)
  } finally {
    isSaving.value = false
  }
}
</script>

<style scoped>
.animal-screen {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  animation: fadeIn 0.3s ease;
}

.mode-label {
  color: var(--color-white);
  font-size: 1rem;
  margin: 0.5rem 0;
  opacity: 0.9;
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

.preview-container {
  position: relative;
  width: 640px;
  max-width: 90vw;
  height: 480px;
  max-height: 55vh;
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-lg);
  background: #000;
}

.result-canvas {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.processing-label {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: var(--color-white);
  font-size: 1.3rem;
  font-weight: 600;
  text-shadow: 2px 2px 6px rgba(0, 0, 0, 0.7);
  background: rgba(0, 0, 0, 0.5);
  padding: 0.5rem 1.5rem;
  border-radius: var(--radius-lg);
}

.animal-tray {
  display: flex;
  gap: 0.4rem;
  flex-wrap: wrap;
  justify-content: center;
  padding: 0.6rem 1rem;
  max-width: 90vw;
  background: rgba(0, 0, 0, 0.5);
  border-radius: var(--radius-lg);
  backdrop-filter: blur(8px);
}

.animal-pill {
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

.animal-pill.active {
  background: var(--color-primary);
  border-color: var(--color-primary);
  color: white;
}

.animal-pill:hover:not(:disabled) {
  border-color: rgba(255, 255, 255, 0.6);
  transform: translateY(-1px);
}

.animal-pill:disabled {
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
