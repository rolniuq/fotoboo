<template>
  <div class="capture-screen">
    <div class="camera-container">
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
import { useBoothStore } from '../../stores/booth'

const router = useRouter()
const booth = useBoothStore()
const camera = useCamera()
const countdown = useCountdown(3)
const showFlash = ref(false)

// Bind the video ref from composable
const videoRef = camera.videoRef

onMounted(() => {
  camera.start()
})

onUnmounted(() => {
  camera.stop()
})

async function handleCapture() {
  await countdown.start()

  // Flash
  showFlash.value = true

  const dataUrl = camera.captureFrame()
  if (dataUrl) {
    booth.setCapturedPhoto(dataUrl)
  }

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
  max-height: 60vh;
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
