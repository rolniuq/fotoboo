import { ref, computed } from 'vue'

// Collage layout definitions
const LAYOUTS = {
  single: {
    name: 'Single',
    photos: 1,
    positions: [{ x: 0, y: 0, w: 1, h: 1 }],
  },
  horizontal2: {
    name: '2 Side by Side',
    photos: 2,
    positions: [
      { x: 0, y: 0, w: 0.5, h: 1 },
      { x: 0.5, y: 0, w: 0.5, h: 1 },
    ],
  },
  vertical2: {
    name: '2 Stacked',
    photos: 2,
    positions: [
      { x: 0, y: 0, w: 1, h: 0.5 },
      { x: 0, y: 0.5, w: 1, h: 0.5 },
    ],
  },
  grid4: {
    name: '4 Grid',
    photos: 4,
    positions: [
      { x: 0, y: 0, w: 0.5, h: 0.5 },
      { x: 0.5, y: 0, w: 0.5, h: 0.5 },
      { x: 0, y: 0.5, w: 0.5, h: 0.5 },
      { x: 0.5, y: 0.5, w: 0.5, h: 0.5 },
    ],
  },
  strip4: {
    name: '4 Strip',
    photos: 4,
    positions: [
      { x: 0, y: 0, w: 1, h: 0.25 },
      { x: 0, y: 0.25, w: 1, h: 0.25 },
      { x: 0, y: 0.5, w: 1, h: 0.25 },
      { x: 0, y: 0.75, w: 1, h: 0.25 },
    ],
  },
  featured3: {
    name: '1 Large + 2 Small',
    photos: 3,
    positions: [
      { x: 0, y: 0, w: 0.6, h: 1 },
      { x: 0.6, y: 0, w: 0.4, h: 0.5 },
      { x: 0.6, y: 0.5, w: 0.4, h: 0.5 },
    ],
  },
}

export function useCollage() {
  const currentLayout = ref('single')
  const photos = ref([]) // Array of data URLs
  const canvasWidth = 1280
  const canvasHeight = 960

  const layout = computed(() => LAYOUTS[currentLayout.value])
  const requiredPhotos = computed(() => layout.value.photos)
  const capturedCount = computed(() => photos.value.length)
  const isComplete = computed(() => photos.value.length >= requiredPhotos.value)

  function setLayout(layoutKey) {
    if (LAYOUTS[layoutKey]) {
      currentLayout.value = layoutKey
      // Clear photos if layout requires different count
      if (photos.value.length > LAYOUTS[layoutKey].photos) {
        photos.value = photos.value.slice(0, LAYOUTS[layoutKey].photos)
      }
    }
  }

  function addPhoto(dataUrl) {
    if (photos.value.length < requiredPhotos.value) {
      photos.value.push(dataUrl)
    }
  }

  function removePhoto(index) {
    if (index >= 0 && index < photos.value.length) {
      photos.value.splice(index, 1)
    }
  }

  function clearPhotos() {
    photos.value = []
  }

  // Render collage to canvas
  function renderCollage(canvas, filterString = '') {
    return new Promise((resolve) => {
      const ctx = canvas.getContext('2d')
      canvas.width = canvasWidth
      canvas.height = canvasHeight

      // Fill background
      ctx.fillStyle = '#ffffff'
      ctx.fillRect(0, 0, canvasWidth, canvasHeight)

      const positions = layout.value.positions
      let loaded = 0

      if (photos.value.length === 0) {
        resolve()
        return
      }

      photos.value.forEach((photoSrc, index) => {
        if (index >= positions.length) return

        const pos = positions[index]
        const img = new Image()

        img.onload = () => {
          const x = pos.x * canvasWidth
          const y = pos.y * canvasHeight
          const w = pos.w * canvasWidth
          const h = pos.h * canvasHeight

          // Calculate aspect-ratio-aware drawing
          const imgAspect = img.width / img.height
          const slotAspect = w / h

          let sx = 0, sy = 0, sw = img.width, sh = img.height

          if (imgAspect > slotAspect) {
            // Image wider than slot - crop sides
            sw = img.height * slotAspect
            sx = (img.width - sw) / 2
          } else {
            // Image taller than slot - crop top/bottom
            sh = img.width / slotAspect
            sy = (img.height - sh) / 2
          }

          // Apply filter if provided
          if (filterString) {
            ctx.filter = filterString
          }

          ctx.drawImage(img, sx, sy, sw, sh, x, y, w, h)
          ctx.filter = 'none'

          // Add subtle border between photos
          ctx.strokeStyle = '#ffffff'
          ctx.lineWidth = 4
          ctx.strokeRect(x, y, w, h)

          loaded++
          if (loaded === Math.min(photos.value.length, positions.length)) {
            resolve()
          }
        }

        img.onerror = () => {
          loaded++
          if (loaded === Math.min(photos.value.length, positions.length)) {
            resolve()
          }
        }

        img.src = photoSrc
      })
    })
  }

  function getCanvasBlob(canvas, quality = 0.92) {
    return new Promise((resolve) => {
      canvas.toBlob((blob) => resolve(blob), 'image/jpeg', quality)
    })
  }

  function reset() {
    currentLayout.value = 'single'
    photos.value = []
  }

  return {
    currentLayout,
    photos,
    layout,
    layouts: LAYOUTS,
    layoutKeys: Object.keys(LAYOUTS),
    requiredPhotos,
    capturedCount,
    isComplete,
    setLayout,
    addPhoto,
    removePhoto,
    clearPhotos,
    renderCollage,
    getCanvasBlob,
    reset,
  }
}
