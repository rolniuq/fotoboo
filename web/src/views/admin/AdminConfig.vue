<template>
  <div class="admin-config">
    <h1>Configuration</h1>

    <div v-if="loaded" class="card">
      <form @submit.prevent="handleSave">
        <div class="form-group">
          <label>Event Name</label>
          <input v-model="config.event_name" type="text" />
        </div>

        <div class="form-group">
          <label>Countdown Duration (seconds)</label>
          <input v-model.number="config.countdown_duration" type="number" min="1" max="10" />
        </div>

        <div class="form-group">
          <label>Max Upload Size (MB)</label>
          <input v-model.number="config.max_upload_size_mb" type="number" min="1" max="50" />
        </div>

        <div class="form-group">
          <label>Photo Retention (days)</label>
          <input v-model.number="config.photo_retention_days" type="number" min="1" max="365" />
        </div>

        <div class="form-group">
          <label>Available Frames</label>
          <div class="checkbox-group">
            <label v-for="f in allFrames" :key="f" class="checkbox-label">
              <input type="checkbox" :value="f" v-model="config.available_frames" />
              {{ f }}
            </label>
          </div>
        </div>

        <div class="form-group">
          <label>Available Filters</label>
          <div class="checkbox-group">
            <label v-for="f in allFilters" :key="f" class="checkbox-label">
              <input type="checkbox" :value="f" v-model="config.available_filters" />
              {{ f }}
            </label>
          </div>
        </div>

        <div class="form-actions">
          <button type="submit" class="btn btn-primary btn-sm" :disabled="saving">
            {{ saving ? 'Saving...' : 'Save Configuration' }}
          </button>
          <button type="button" class="btn btn-secondary btn-sm" @click="loadConfig">
            Reset
          </button>
        </div>

        <p v-if="message" :class="['message', messageType]">{{ message }}</p>
      </form>
    </div>

    <div v-else class="card">
      <p class="empty-text">Loading configuration...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getAdminConfig, updateAdminConfig } from '../../composables/useApi'

const allFrames = ['none', 'simple', 'event', 'party']
const allFilters = ['none', 'grayscale', 'vintage', 'brightness', 'contrast']

const config = reactive({
  event_name: '',
  countdown_duration: 3,
  available_frames: [],
  available_filters: [],
  max_upload_size_mb: 10,
  photo_retention_days: 30,
})

const loaded = ref(false)
const saving = ref(false)
const message = ref('')
const messageType = ref('')

onMounted(() => loadConfig())

async function loadConfig() {
  try {
    const data = await getAdminConfig()
    Object.assign(config, data)
    loaded.value = true
    message.value = ''
  } catch (err) {
    message.value = 'Failed to load config: ' + err.message
    messageType.value = 'error'
  }
}

async function handleSave() {
  saving.value = true
  message.value = ''

  try {
    await updateAdminConfig({ ...config })
    message.value = 'Configuration saved successfully!'
    messageType.value = 'success'
  } catch (err) {
    message.value = 'Failed to save: ' + err.message
    messageType.value = 'error'
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.admin-config h1 {
  margin-bottom: 1.5rem;
  color: var(--color-dark);
}

.card {
  background: white;
  padding: 2rem;
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-sm);
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group > label {
  display: block;
  font-weight: 600;
  color: var(--color-dark);
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.form-group input[type='text'],
.form-group input[type='number'] {
  width: 100%;
  max-width: 400px;
  padding: 0.6rem 0.8rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 0.9rem;
}

.form-group input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.checkbox-group {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.9rem;
  color: var(--color-dark);
  cursor: pointer;
}

.checkbox-label input {
  cursor: pointer;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
}

.message {
  margin-top: 1rem;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  font-size: 0.9rem;
}

.message.success {
  background: #e8f5e9;
  color: #2e7d32;
}

.message.error {
  background: #ffebee;
  color: #c62828;
}

.empty-text {
  color: var(--color-muted);
  text-align: center;
  padding: 2rem;
}
</style>
