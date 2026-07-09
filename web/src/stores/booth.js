import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as api from '../composables/useApi'

export const useBoothStore = defineStore('booth', () => {
  // State
  const sessionId = ref(null)
  const photoId = ref(null)
  const capturedPhoto = ref(null) // base64 data URL (with face overlay baked in, if any)
  const rawCapturedPhoto = ref(null) // base64 data URL (without face overlay, for re-rendering)
  const overlayKey = ref('none')    // Currently selected face overlay
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

  function setRawCapturedPhoto(dataUrl) {
    rawCapturedPhoto.value = dataUrl
  }

  function setOverlayKey(key) {
    overlayKey.value = key
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
    rawCapturedPhoto.value = null
    overlayKey.value = 'none'
    isSaving.value = false
    error.value = null
  }

  return {
    sessionId,
    photoId,
    capturedPhoto,
    rawCapturedPhoto,
    overlayKey,
    isSaving,
    error,
    createSession,
    setCapturedPhoto,
    setRawCapturedPhoto,
    setOverlayKey,
    savePhoto,
    reset,
  }
})
