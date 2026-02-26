<template>
  <div class="admin-photos">
    <div class="page-header">
      <h1>Photos</h1>
      <span class="count">{{ photos.length }} total</span>
    </div>

    <div v-if="photos.length" class="photo-grid">
      <div v-for="photo in photos" :key="photo.id" class="photo-card">
        <div class="photo-image">
          <img :src="`/photos/${photo.id}`" :alt="photo.id" />
        </div>
        <div class="photo-info">
          <div class="photo-id mono">{{ photo.id.slice(0, 12) }}...</div>
          <div class="photo-date">{{ formatDate(photo.created_at) }}</div>
          <div v-if="photo.session_id" class="photo-session">
            Session: {{ photo.session_id.slice(0, 8) }}...
          </div>
        </div>
        <div class="photo-actions">
          <a :href="`/photos/${photo.id}`" download class="btn btn-sm btn-primary">
            Download
          </a>
          <button class="btn btn-sm btn-danger" @click="handleDelete(photo.id)">
            Delete
          </button>
        </div>
      </div>
    </div>

    <p v-else class="empty-text">No photos uploaded yet</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listPhotos, deletePhoto } from '../../composables/useApi'

const photos = ref([])

onMounted(async () => {
  await loadPhotos()
})

async function loadPhotos() {
  try {
    photos.value = (await listPhotos()) || []
  } catch (err) {
    console.error('Failed to load photos:', err)
  }
}

async function handleDelete(id) {
  if (!confirm('Delete this photo? This cannot be undone.')) return
  try {
    await deletePhoto(id)
    photos.value = photos.value.filter((p) => p.id !== id)
  } catch (err) {
    alert('Failed to delete photo: ' + err.message)
  }
}

function formatDate(dateStr) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleString()
}
</script>

<style scoped>
.page-header {
  display: flex;
  align-items: baseline;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.page-header h1 {
  color: var(--color-dark);
}

.count {
  color: var(--color-muted);
  font-size: 0.9rem;
}

.photo-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
}

.photo-card {
  background: white;
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.photo-image {
  aspect-ratio: 4 / 3;
  background: #f0f0f0;
}

.photo-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.photo-info {
  padding: 0.75rem 1rem;
}

.photo-id {
  font-size: 0.8rem;
  color: var(--color-muted);
}

.photo-date {
  font-size: 0.85rem;
  color: var(--color-dark);
  margin-top: 0.25rem;
}

.photo-session {
  font-size: 0.75rem;
  color: var(--color-muted);
  margin-top: 0.25rem;
}

.photo-actions {
  padding: 0 1rem 0.75rem;
  display: flex;
  gap: 0.5rem;
}

.mono {
  font-family: 'SF Mono', 'Fira Code', monospace;
}

.empty-text {
  color: var(--color-muted);
  text-align: center;
  padding: 3rem;
  background: white;
  border-radius: var(--radius-sm);
}
</style>
