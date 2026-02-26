import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as api from '../composables/useApi'

export const useBoothStore = defineStore('booth', () => {
  // State
  const sessionId = ref(null)
  const photoId = ref(null)
  const capturedPhoto = ref(null) // base64 data URL
  const isSaving = ref(false)
  const error = ref(null)

  // Actions
  async function createSession() {
    try {
      error.value = null
      const data = await api.startSession()
      sessionId.value = data.id
      return data
    } catch (err) {
      error.value = err.message
      // Continue without session — non-blocking
      sessionId.value = null
    }
  }

  function setCapturedPhoto(dataUrl) {
    capturedPhoto.value = dataUrl
  }

  async function savePhoto(blob) {
    try {
      error.value = null
      isSaving.value = true
      const data = await api.uploadPhoto(blob, sessionId.value)
      photoId.value = data.id

      // Complete session
      if (sessionId.value) {
        try {
          await api.completeSession(sessionId.value)
        } catch {
          // Non-critical
        }
      }

      return data
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      isSaving.value = false
    }
  }

  function reset() {
    sessionId.value = null
    photoId.value = null
    capturedPhoto.value = null
    isSaving.value = false
    error.value = null
  }

  return {
    sessionId,
    photoId,
    capturedPhoto,
    isSaving,
    error,
    createSession,
    setCapturedPhoto,
    savePhoto,
    reset,
  }
})
