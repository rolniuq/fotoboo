<template>
  <div class="welcome-screen">
    <h1 class="title">FotoBoo</h1>
    <p class="subtitle">Your Photo Booth Experience</p>
    
    <div class="mode-buttons">
      <button class="btn btn-primary btn-start" @click="handleStart('single')">
        Single Photo
      </button>
      <button class="btn btn-secondary btn-start" @click="handleStart('collage')">
        Photo Collage
      </button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useBoothStore } from '../../stores/booth'
import { useFilters } from '../../composables/useFilters'

const router = useRouter()
const booth = useBoothStore()

const { reset: resetFilters } = useFilters()

async function handleStart(mode) {
  booth.reset()
  resetFilters()
  await booth.createSession()
  
  if (mode === 'collage') {
    router.push({ name: 'collage' })
  } else {
    router.push({ name: 'capture' })
  }
}
</script>

<style scoped>
.welcome-screen {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  color: var(--color-white);
  animation: fadeIn 0.5s ease;
}

.title {
  font-size: 4rem;
  margin-bottom: 1rem;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.subtitle {
  font-size: 1.5rem;
  margin-bottom: 2rem;
  opacity: 0.9;
}

.mode-buttons {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  justify-content: center;
}

.btn-start {
  font-size: 1.3rem;
  padding: 1.2rem 3rem;
}

@media (max-width: 768px) {
  .title {
    font-size: 3rem;
  }

  .mode-buttons {
    flex-direction: column;
  }
}
</style>
