<template>
  <div class="dashboard">
    <h1>Dashboard</h1>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-value">{{ stats.total_photos }}</div>
        <div class="stat-label">Total Photos</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ stats.total_sessions }}</div>
        <div class="stat-label">Total Sessions</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ stats.total_devices }}</div>
        <div class="stat-label">Devices</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ stats.photos_today }}</div>
        <div class="stat-label">Photos Today</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ stats.sessions_today }}</div>
        <div class="stat-label">Sessions Today</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ stats.storage_formatted || '0 B' }}</div>
        <div class="stat-label">Storage Used</div>
      </div>
    </div>

    <!-- Metrics -->
    <div class="card">
      <h2>Server Metrics</h2>
      <div v-if="metrics" class="metrics-grid">
        <div class="metric"><strong>Uptime:</strong> {{ metrics.uptime }}</div>
        <div class="metric"><strong>Total Requests:</strong> {{ metrics.total_requests }}</div>
        <div class="metric"><strong>Active Requests:</strong> {{ metrics.active_requests }}</div>
        <div class="metric"><strong>Error Rate:</strong> {{ metrics.error_rate }}</div>
      </div>
      <p v-else class="empty-text">Loading metrics...</p>
    </div>

    <!-- Recent Photos -->
    <div class="card">
      <h2>Recent Photos</h2>
      <div v-if="recentPhotos.length" class="photo-grid">
        <div v-for="photo in recentPhotos" :key="photo.id" class="photo-thumb">
          <img :src="`/photos/${photo.id}`" :alt="photo.id" />
          <div class="photo-meta">
            <span class="photo-date">{{ formatDate(photo.created_at) }}</span>
          </div>
        </div>
      </div>
      <p v-else class="empty-text">No photos yet</p>
    </div>

    <!-- Recent Sessions -->
    <div class="card">
      <h2>Recent Sessions</h2>
      <table v-if="recentSessions.length" class="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Status</th>
            <th>Device</th>
            <th>Created</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in recentSessions" :key="s.id">
            <td class="mono">{{ s.id.slice(0, 8) }}...</td>
            <td><span :class="['badge', `badge-${s.status}`]">{{ s.status }}</span></td>
            <td>{{ s.device_id ? s.device_id.slice(0, 8) + '...' : '—' }}</td>
            <td>{{ formatDate(s.created_at) }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else class="empty-text">No sessions yet</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getAdminStats, getMetrics, listPhotos, listSessions } from '../../composables/useApi'

const stats = reactive({
  total_photos: 0,
  total_sessions: 0,
  total_devices: 0,
  photos_today: 0,
  sessions_today: 0,
  storage_bytes: 0,
  storage_formatted: '0 B',
})

const metrics = ref(null)
const recentPhotos = ref([])
const recentSessions = ref([])

onMounted(async () => {
  try {
    const [statsData, metricsData, photos, sessions] = await Promise.all([
      getAdminStats(),
      getMetrics(),
      listPhotos(),
      listSessions(),
    ])

    Object.assign(stats, statsData)
    metrics.value = metricsData
    recentPhotos.value = (photos || []).slice(0, 8)
    recentSessions.value = (sessions || []).slice(0, 10)
  } catch (err) {
    console.error('Failed to load dashboard data:', err)
  }
})

function formatDate(dateStr) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleString()
}
</script>

<style scoped>
.dashboard h1 {
  margin-bottom: 1.5rem;
  color: var(--color-dark);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(170px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  padding: 1.5rem;
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-sm);
  text-align: center;
}

.stat-value {
  font-size: 2.2rem;
  font-weight: 700;
  color: var(--color-primary);
}

.stat-label {
  font-size: 0.85rem;
  color: var(--color-muted);
  margin-top: 0.25rem;
}

.card {
  background: white;
  padding: 1.5rem;
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-sm);
  margin-bottom: 1.5rem;
}

.card h2 {
  margin-bottom: 1rem;
  color: var(--color-dark);
  font-size: 1.2rem;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 0.75rem;
}

.metric {
  padding: 0.5rem 0;
  font-size: 0.9rem;
  color: var(--color-dark);
}

.photo-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 0.75rem;
}

.photo-thumb {
  border-radius: 8px;
  overflow: hidden;
  background: #f0f0f0;
  aspect-ratio: 4 / 3;
  position: relative;
}

.photo-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.photo-meta {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  padding: 4px 8px;
  font-size: 0.7rem;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid #eee;
  font-size: 0.9rem;
}

.data-table th {
  font-weight: 600;
  color: var(--color-muted);
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.mono {
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 0.8rem;
}

.badge {
  padding: 0.2rem 0.6rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.badge-active { background: #e8f5e9; color: #2e7d32; }
.badge-completed { background: #e3f2fd; color: #1565c0; }

.empty-text {
  color: var(--color-muted);
  text-align: center;
  padding: 2rem;
}
</style>
