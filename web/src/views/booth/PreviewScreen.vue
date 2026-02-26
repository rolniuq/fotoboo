<template>
  <div class="preview-screen">
    <div class="preview-container">
      <canvas ref="canvasRef" class="preview-canvas"></canvas>
      <div :class="frameOverlayClass"></div>
    </div>

    <!-- Filters -->
    <div class="control-card">
      <h3>Filters</h3>
      <div class="pill-group">
        <button
          v-for="f in filters.filters"
          :key="f"
          :class="['pill', { active: filters.currentFilter.value === f }]"
          @click="filters.currentFilter.value = f; renderPreview()"
        >
          {{ f === 'none' ? 'None' : f.charAt(0).toUpperCase() + f.slice(1) }}
        </button>
      </div>

      <div class="slider-group">
        <label>
          Brightness: {{ filters.brightness.value }}%
          <input
            type="range"
            min="50"
            max="150"
            v-model.number="filters.brightness.value"
            @input="renderPreview()"
          >
        </label>
        <label>
          Contrast: {{ filters.contrast.value }}%
          <input
            type="range"
            min="50"
            max="150"
            v-model.number="filters.contrast.value"
            @input="renderPreview()"
          >
        </label>
      </div>
    </div>

    <!-- Frames -->
    <div class="control-card">
      <h3>Frames</h3>
      <div class="pill-group">
        <button
          v-for="f in filters.frames"
          :key="f"
          :class="['pill', { active: filters.currentFrame.value === f }]"
          @click="filters.currentFrame.value = f"
        >
          {{ f === 'none' ? 'None' : f.charAt(0).toUpperCase() + f.slice(1) }}
        </button>
      </div>
    </div>

    <div class="controls">
      <button class="btn btn-secondary" @click="$router.push({ name: 'capture' })">
        Retake
      </button>
      <button
        class="btn btn-primary"
        :disabled="booth.isSaving"
        @click="handleSave"
      >
        {{ booth.isSaving ? 'Saving...' : 'Save Photo' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useBoothStore } from '../../stores/booth'
import { useFilters } from '../../composables/useFilters'

const router = useRouter()
const booth = useBoothStore()
const filters = useFilters()
const canvasRef = ref(null)

const frameOverlayClass = computed(() => {
  const frame = filters.currentFrame.value
  return ['frame-overlay', frame !== 'none' ? `frame-${frame}` : '']
})

onMounted(() => {
  if (!booth.capturedPhoto) {
    router.push({ name: 'capture' })
    return
  }
  renderPreview()
})

async function renderPreview() {
  if (!canvasRef.value || !booth.capturedPhoto) return
  await filters.applyToCanvas(canvasRef.value, booth.capturedPhoto)
}

async function handleSave() {
  if (!canvasRef.value) return

  // Ensure latest filters applied
  await renderPreview()

  const blob = await filters.getCanvasBlob(canvasRef.value)
  await booth.savePhoto(blob)
  router.push({ name: 'result' })
}
</script>

<style scoped>
.preview-screen {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 0.75rem;
  padding: 1rem;
  overflow-y: auto;
  animation: fadeIn 0.3s ease;
}

.preview-container {
  position: relative;
  width: 640px;
  max-width: 90vw;
  height: 400px;
  max-height: 45vh;
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-lg);
}

.preview-canvas {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.frame-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 10;
}

.frame-simple {
  border: 15px solid #fff;
  box-shadow: inset 0 0 20px rgba(0, 0, 0, 0.1);
}

.frame-event {
  border: 20px solid var(--color-gold);
}

.frame-event::before {
  content: 'Special Event';
  position: absolute;
  bottom: 10px;
  left: 0;
  right: 0;
  text-align: center;
  color: var(--color-dark);
  font-size: 1.2rem;
  font-weight: bold;
  background: rgba(255, 215, 0, 0.8);
  padding: 4px;
}

.frame-party {
  border: 20px solid var(--color-accent);
  border-radius: var(--radius-md);
}

.frame-party::before {
  content: 'PARTY TIME';
  position: absolute;
  top: 10px;
  left: 0;
  right: 0;
  text-align: center;
  color: var(--color-white);
  font-size: 1.4rem;
  font-weight: bold;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.7);
  background: rgba(255, 23, 68, 0.7);
  padding: 4px;
}

/* Control cards */
.control-card {
  background: rgba(255, 255, 255, 0.95);
  padding: 1rem 1.5rem;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-md);
  width: 640px;
  max-width: 90vw;
}

.control-card h3 {
  margin-bottom: 0.75rem;
  color: var(--color-dark);
  text-align: center;
  font-size: 0.95rem;
}

.pill-group {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  justify-content: center;
  margin-bottom: 0.75rem;
}

.pill {
  padding: 0.4rem 1rem;
  border: 2px solid var(--color-primary);
  background: white;
  color: var(--color-primary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  font-weight: 500;
  font-size: 0.85rem;
  transition: var(--transition);
}

.pill.active {
  background: var(--color-primary);
  color: white;
}

.pill:hover {
  transform: translateY(-1px);
}

.slider-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.slider-group label {
  font-weight: 500;
  color: var(--color-dark);
  font-size: 0.85rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.slider-group input[type='range'] {
  width: 100%;
}

@media (max-width: 768px) {
  .pill-group {
    flex-direction: column;
    align-items: center;
  }

  .pill {
    width: 180px;
  }
}
</style>
