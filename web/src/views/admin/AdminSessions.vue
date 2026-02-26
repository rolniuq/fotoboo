<template>
  <div class="admin-sessions">
    <div class="page-header">
      <h1>Sessions</h1>
      <span class="count">{{ sessions.length }} total</span>
    </div>

    <div class="card">
      <table v-if="sessions.length" class="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Status</th>
            <th>Device</th>
            <th>Created</th>
            <th>Updated</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in sessions" :key="s.id">
            <td class="mono">{{ s.id.slice(0, 12) }}...</td>
            <td>
              <span :class="['badge', `badge-${s.status}`]">{{ s.status }}</span>
            </td>
            <td>{{ s.device_id ? s.device_id.slice(0, 8) + '...' : '—' }}</td>
            <td>{{ formatDate(s.created_at) }}</td>
            <td>{{ formatDate(s.updated_at) }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else class="empty-text">No sessions recorded yet</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listSessions } from '../../composables/useApi'

const sessions = ref([])

onMounted(async () => {
  try {
    sessions.value = (await listSessions()) || []
  } catch (err) {
    console.error('Failed to load sessions:', err)
  }
})

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

.card {
  background: white;
  padding: 1.5rem;
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-sm);
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

.badge-active {
  background: #e8f5e9;
  color: #2e7d32;
}

.badge-completed {
  background: #e3f2fd;
  color: #1565c0;
}

.empty-text {
  color: var(--color-muted);
  text-align: center;
  padding: 2rem;
}
</style>
