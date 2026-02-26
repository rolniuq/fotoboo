const BASE = ''

async function request(url, options = {}) {
  const res = await fetch(`${BASE}${url}`, {
    headers: { 'Content-Type': 'application/json', ...options.headers },
    ...options,
  })

  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(body.error || `HTTP ${res.status}`)
  }

  if (res.status === 204) return null

  const contentType = res.headers.get('Content-Type') || ''
  if (contentType.includes('application/json')) return res.json()
  return res
}

// --- Photos ---
export function uploadPhoto(blob, sessionId = null) {
  let url = '/photos'
  if (sessionId) url += `?session_id=${sessionId}`
  return request(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/octet-stream' },
    body: blob,
  })
}

export function getPhotoUrl(id) {
  return `${BASE}/photos/${id}`
}

export function getPhotoQrUrl(id) {
  return `${BASE}/photos/${id}/qr`
}

export function getPhotoPrintUrl(id, size = '4x6') {
  return `${BASE}/photos/${id}/print?size=${size}`
}

export function listPhotos() {
  return request('/photos')
}

export function deletePhoto(id) {
  return request(`/photos/${id}`, { method: 'DELETE' })
}

// --- Sessions ---
export function startSession(deviceId = '') {
  return request('/sessions', {
    method: 'POST',
    body: JSON.stringify({ device_id: deviceId }),
  })
}

export function getSession(id) {
  return request(`/sessions/${id}`)
}

export function completeSession(id) {
  return request(`/sessions/${id}/complete`, { method: 'POST' })
}

export function listSessions() {
  return request('/sessions')
}

export function getSessionPhotos(id) {
  return request(`/sessions/${id}/photos`)
}

// --- Devices ---
export function registerDevice(name, type = 'webcam') {
  return request('/devices', {
    method: 'POST',
    body: JSON.stringify({ name, type }),
  })
}

export function getDevice(id) {
  return request(`/devices/${id}`)
}

export function updateDevice(id, data) {
  return request(`/devices/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export function deleteDevice(id) {
  return request(`/devices/${id}`, { method: 'DELETE' })
}

export function listDevices() {
  return request('/devices')
}

// --- Print Sizes ---
export function getPrintSizes() {
  return request('/print-sizes')
}

// --- Health ---
export function healthCheck() {
  return request('/health')
}

// --- Admin ---
export function getAdminStats() {
  return request('/admin/stats')
}

export function getAdminConfig() {
  return request('/admin/config')
}

export function updateAdminConfig(config) {
  return request('/admin/config', {
    method: 'PUT',
    body: JSON.stringify(config),
  })
}

// --- Metrics ---
export function getMetrics() {
  return request('/metrics')
}
