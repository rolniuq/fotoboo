<template>
  <div class="admin-devices">
    <div class="page-header">
      <h1>Devices</h1>
      <button class="btn btn-primary btn-sm" @click="showForm = !showForm">
        {{ showForm ? 'Cancel' : '+ Add Device' }}
      </button>
    </div>

    <!-- Add Device Form -->
    <Transition name="slide">
      <div v-if="showForm" class="card form-card">
        <h2>Register New Device</h2>
        <form @submit.prevent="handleAdd">
          <div class="form-row">
            <label>
              Name
              <input v-model="newDevice.name" type="text" placeholder="e.g. Booth Camera 1" required />
            </label>
            <label>
              Type
              <select v-model="newDevice.type">
                <option value="webcam">Webcam</option>
                <option value="dslr">DSLR</option>
                <option value="phone">Phone</option>
              </select>
            </label>
            <button type="submit" class="btn btn-primary btn-sm">Register</button>
          </div>
        </form>
      </div>
    </Transition>

    <!-- Devices List -->
    <div class="card">
      <table v-if="devices.length" class="data-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Status</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="d in devices" :key="d.id">
            <td>{{ d.name }}</td>
            <td><span class="type-badge">{{ d.type }}</span></td>
            <td>
              <span :class="['badge', d.active ? 'badge-active' : 'badge-inactive']">
                {{ d.active ? 'Active' : 'Inactive' }}
              </span>
            </td>
            <td>{{ formatDate(d.created_at) }}</td>
            <td class="actions-cell">
              <button
                class="btn btn-sm"
                :class="d.active ? 'btn-secondary' : 'btn-success'"
                @click="toggleActive(d)"
              >
                {{ d.active ? 'Deactivate' : 'Activate' }}
              </button>
              <button class="btn btn-sm btn-danger" @click="handleDelete(d.id)">
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      <p v-else class="empty-text">No devices registered</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import {
  listDevices,
  registerDevice,
  updateDevice,
  deleteDevice,
} from '../../composables/useApi'

const devices = ref([])
const showForm = ref(false)
const newDevice = reactive({
  name: '',
  type: 'webcam',
})

onMounted(async () => {
  await loadDevices()
})

async function loadDevices() {
  try {
    devices.value = (await listDevices()) || []
  } catch (err) {
    console.error('Failed to load devices:', err)
  }
}

async function handleAdd() {
  try {
    await registerDevice(newDevice.name, newDevice.type)
    newDevice.name = ''
    newDevice.type = 'webcam'
    showForm.value = false
    await loadDevices()
  } catch (err) {
    alert('Failed to register device: ' + err.message)
  }
}

async function toggleActive(device) {
  try {
    await updateDevice(device.id, {
      name: device.name,
      type: device.type,
      active: !device.active,
    })
    await loadDevices()
  } catch (err) {
    alert('Failed to update device: ' + err.message)
  }
}

async function handleDelete(id) {
  if (!confirm('Delete this device?')) return
  try {
    await deleteDevice(id)
    devices.value = devices.value.filter((d) => d.id !== id)
  } catch (err) {
    alert('Failed to delete device: ' + err.message)
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
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
}

.page-header h1 {
  color: var(--color-dark);
}

.card {
  background: white;
  padding: 1.5rem;
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-sm);
  margin-bottom: 1rem;
}

.form-card h2 {
  margin-bottom: 1rem;
  font-size: 1.1rem;
  color: var(--color-dark);
}

.form-row {
  display: flex;
  gap: 1rem;
  align-items: flex-end;
  flex-wrap: wrap;
}

.form-row label {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--color-muted);
  flex: 1;
  min-width: 150px;
}

.form-row input,
.form-row select {
  padding: 0.6rem 0.8rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 0.9rem;
}

.form-row input:focus,
.form-row select:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
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

.actions-cell {
  display: flex;
  gap: 0.5rem;
}

.type-badge {
  background: #f0f0f0;
  padding: 0.15rem 0.5rem;
  border-radius: 8px;
  font-size: 0.8rem;
  color: var(--color-muted);
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

.badge-inactive {
  background: #ffebee;
  color: #c62828;
}

.empty-text {
  color: var(--color-muted);
  text-align: center;
  padding: 2rem;
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
