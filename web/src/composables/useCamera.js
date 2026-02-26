import { ref, onUnmounted } from 'vue'

export function useCamera() {
  const videoRef = ref(null)
  const stream = ref(null)
  const isReady = ref(false)
  const error = ref(null)

  async function start() {
    try {
      error.value = null
      stream.value = await navigator.mediaDevices.getUserMedia({
        video: {
          width: { ideal: 1280 },
          height: { ideal: 720 },
          facingMode: 'user',
        },
        audio: false,
      })

      if (videoRef.value) {
        videoRef.value.srcObject = stream.value
        await new Promise((resolve) => {
          videoRef.value.onloadedmetadata = resolve
        })
        isReady.value = true
      }
    } catch (err) {
      error.value = err.message || 'Unable to access camera'
      console.error('Camera error:', err)
    }
  }

  function stop() {
    if (stream.value) {
      stream.value.getTracks().forEach((track) => track.stop())
      stream.value = null
      isReady.value = false
    }
  }

  function captureFrame() {
    if (!videoRef.value || !isReady.value) return null

    const video = videoRef.value
    const canvas = document.createElement('canvas')
    canvas.width = video.videoWidth
    canvas.height = video.videoHeight

    const ctx = canvas.getContext('2d')
    ctx.drawImage(video, 0, 0)

    return canvas.toDataURL('image/jpeg', 0.9)
  }

  onUnmounted(() => {
    stop()
  })

  return {
    videoRef,
    stream,
    isReady,
    error,
    start,
    stop,
    captureFrame,
  }
}
