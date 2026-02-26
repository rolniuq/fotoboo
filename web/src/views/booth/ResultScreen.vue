<template>
  <div class="result-screen">
    <div class="result-container">
      <img
        v-if="booth.photoId"
        :src="getPhotoUrl(booth.photoId)"
        alt="Your photo"
        class="final-photo"
      />

      <div v-if="booth.photoId" class="qr-card">
        <img
          :src="getPhotoQrUrl(booth.photoId)"
          alt="QR Code"
          class="qr-image"
          @error="qrError = true"
        />
        <p class="qr-label">
          {{ qrError ? photoFullUrl : 'Scan QR code to view photo' }}
        </p>
      </div>
    </div>

    <div class="controls">
      <button class="btn btn-primary" @click="handleDownload">
        Download
      </button>

      <div class="print-controls">
        <select v-model="printSize" class="print-select">
          <option value="4x6">4x6 inches</option>
          <option value="5x7">5x7 inches</option>
          <option value="6x8">6x8 inches</option>
          <option value="2x6">2x6 strip</option>
        </select>
        <button class="btn btn-success btn-sm" @click="handlePrint">
          Print-Ready
        </button>
      </div>

      <button class="btn btn-secondary" @click="handleNewPhoto">
        New Photo
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useBoothStore } from '../../stores/booth'
import { getPhotoUrl, getPhotoQrUrl, getPhotoPrintUrl } from '../../composables/useApi'

const router = useRouter()
const booth = useBoothStore()
const printSize = ref('4x6')
const qrError = ref(false)

const photoFullUrl = computed(() => {
  if (!booth.photoId) return ''
  return `${window.location.origin}/photos/${booth.photoId}`
})

onMounted(() => {
  if (!booth.photoId) {
    router.push({ name: 'welcome' })
  }
})

function handleDownload() {
  if (!booth.photoId) return
  const link = document.createElement('a')
  link.href = getPhotoUrl(booth.photoId)
  link.download = `fotoboo-${booth.photoId}.jpg`
  link.click()
}

function handlePrint() {
  if (!booth.photoId) return
  const link = document.createElement('a')
  link.href = getPhotoPrintUrl(booth.photoId, printSize.value)
  link.download = `fotoboo-${booth.photoId}-${printSize.value}.jpg`
  link.click()
}

function handleNewPhoto() {
  booth.reset()
  router.push({ name: 'welcome' })
}
</script>

<style scoped>
.result-screen {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: 1rem;
  animation: fadeIn 0.3s ease;
}

.result-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.final-photo {
  max-width: 640px;
  width: 90vw;
  max-height: 50vh;
  object-fit: contain;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
}

.qr-card {
  background: white;
  padding: 1rem;
  border-radius: var(--radius-sm);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  box-shadow: var(--shadow-sm);
}

.qr-image {
  width: 150px;
  height: 150px;
}

.qr-label {
  font-size: 0.85rem;
  color: var(--color-muted);
  text-align: center;
  max-width: 280px;
  word-break: break-all;
}

.print-controls {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.print-select {
  padding: 0.6rem 2rem 0.6rem 1rem;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: var(--radius-lg);
  background: rgba(255, 255, 255, 0.15);
  color: white;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  backdrop-filter: blur(10px);
  appearance: none;
  -webkit-appearance: none;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='white'%3e%3cpath d='M7 10l5 5 5-5z'/%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 0.5rem center;
  background-size: 1.5rem;
}

.print-select option {
  background: #333;
  color: white;
}

@media (max-width: 768px) {
  .print-controls {
    flex-direction: column;
    width: 100%;
    max-width: 300px;
  }

  .print-select {
    width: 100%;
  }

  .qr-image {
    width: 120px;
    height: 120px;
  }
}
</style>
