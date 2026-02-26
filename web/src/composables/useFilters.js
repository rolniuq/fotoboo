import { ref, computed } from 'vue'

const FILTER_PRESETS = {
  none: '',
  grayscale: 'grayscale(100%)',
  vintage: 'sepia(50%) contrast(1.2) brightness(0.9)',
  brightness: 'brightness(1.3)',
  contrast: 'contrast(1.5)',
}

const FRAMES = ['none', 'simple', 'event', 'party']

export function useFilters() {
  const currentFilter = ref('none')
  const currentFrame = ref('none')
  const brightness = ref(100)
  const contrast = ref(100)

  const filterString = computed(() => {
    const base = `brightness(${brightness.value / 100}) contrast(${contrast.value / 100})`
    const preset = FILTER_PRESETS[currentFilter.value] || ''
    return preset ? `${base} ${preset}` : base
  })

  function applyToCanvas(canvas, imageSrc) {
    return new Promise((resolve) => {
      const ctx = canvas.getContext('2d')
      const img = new Image()

      img.onload = () => {
        canvas.width = img.width
        canvas.height = img.height
        ctx.filter = filterString.value
        ctx.drawImage(img, 0, 0)
        resolve()
      }

      img.src = imageSrc
    })
  }

  function getCanvasBlob(canvas, quality = 0.92) {
    return new Promise((resolve) => {
      canvas.toBlob((blob) => resolve(blob), 'image/jpeg', quality)
    })
  }

  function reset() {
    currentFilter.value = 'none'
    currentFrame.value = 'none'
    brightness.value = 100
    contrast.value = 100
  }

  return {
    currentFilter,
    currentFrame,
    brightness,
    contrast,
    filterString,
    filters: Object.keys(FILTER_PRESETS),
    frames: FRAMES,
    applyToCanvas,
    getCanvasBlob,
    reset,
  }
}
