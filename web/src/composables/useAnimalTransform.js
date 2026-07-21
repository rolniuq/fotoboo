import { ref, computed } from 'vue'

const ANIMALS = {
  cat: { label: 'Cat' },
  dog: { label: 'Dog' },
  fox: { label: 'Fox' },
  bear: { label: 'Bear' },
  rabbit: { label: 'Rabbit' },
  panda: { label: 'Panda' },
  lion: { label: 'Lion' },
}

export function useAnimalTransform() {
  const currentAnimal = ref('cat')
  const isProcessing = ref(false)
  const faceBounds = ref(null)

  const imageCache = new Map()

  const animalKeys = computed(() => Object.keys(ANIMALS))

  async function loadAnimalImage(animalKey) {
    if (imageCache.has(animalKey)) return imageCache.get(animalKey)
    const img = new Image()
    img.src = `/overlays/animal-${animalKey}.svg`
    await new Promise((resolve, reject) => {
      img.onload = resolve
      img.onerror = reject
    })
    imageCache.set(animalKey, img)
    return img
  }

  function getProportionalBounds(imgWidth, imgHeight) {
    const faceWidth = imgWidth * 0.32
    const faceHeight = imgHeight * 0.45
    return {
      x: (imgWidth - faceWidth) / 2,
      y: imgHeight * 0.12,
      width: faceWidth,
      height: faceHeight,
    }
  }

  async function detectFace(imageSrc) {
    try {
      const img = await new Promise((resolve, reject) => {
        const i = new Image()
        i.onload = () => resolve(i)
        i.onerror = reject
        i.src = imageSrc
      })

      if (window.FaceDetector) {
        try {
          const detector = new FaceDetector({ maxDetectedFaces: 1 })
          const faces = await detector.detect(img)
          if (faces.length > 0) {
            const bb = faces[0].boundingBox
            faceBounds.value = { x: bb.x, y: bb.y, width: bb.width, height: bb.height }
            return faceBounds.value
          }
        } catch {
          // fall through
        }
      }

      faceBounds.value = getProportionalBounds(img.width, img.height)
      return faceBounds.value
    } catch {
      return null
    }
  }

  function getAnimalPlacement(animalKey, bounds) {
    const { x, y, width, height } = bounds

    switch (animalKey) {
      case 'cat':
        return { x: x - width * 0.3, y: y - height * 0.5, w: width * 1.6, h: height * 1.4 }
      case 'dog':
        return { x: x - width * 0.25, y: y - height * 0.35, w: width * 1.5, h: height * 1.5 }
      case 'fox':
        return { x: x - width * 0.3, y: y - height * 0.5, w: width * 1.6, h: height * 1.4 }
      case 'bear':
        return { x: x - width * 0.35, y: y - height * 0.35, w: width * 1.7, h: height * 1.4 }
      case 'rabbit':
        return { x: x - width * 0.3, y: y - height * 0.8, w: width * 1.6, h: height * 1.6 }
      case 'panda':
        return { x: x - width * 0.35, y: y - height * 0.35, w: width * 1.7, h: height * 1.35 }
      case 'lion':
        return { x: x - width * 0.6, y: y - height * 0.7, w: width * 2.2, h: height * 2.2 }
      default:
        return { x: x - width * 0.3, y: y - height * 0.4, w: width * 1.6, h: height * 1.4 }
    }
  }

  function applyCartoonEffect(ctx, width, height) {
    const imageData = ctx.getImageData(0, 0, width, height)
    const data = imageData.data

    // Color quantization — reduce to fewer levels per channel
    const levels = 5
    const step = 255 / levels
    for (let i = 0; i < data.length; i += 4) {
      data[i]     = Math.round(data[i] / step) * step
      data[i + 1] = Math.round(data[i + 1] / step) * step
      data[i + 2] = Math.round(data[i + 2] / step) * step
    }

    // Edge detection via simple luminance gradient
    const gray = new Uint8ClampedArray(width * height)
    for (let i = 0; i < data.length; i += 4) {
      gray[i / 4] = data[i] * 0.299 + data[i + 1] * 0.587 + data[i + 2] * 0.114
    }

    const edges = new Float32Array(width * height)
    for (let y = 1; y < height - 1; y++) {
      for (let x = 1; x < width - 1; x++) {
        const idx = y * width + x
        const gx =
          -gray[(y - 1) * width + (x - 1)] + gray[(y - 1) * width + (x + 1)]
          - 2 * gray[y * width + (x - 1)] + 2 * gray[y * width + (x + 1)]
          - gray[(y + 1) * width + (x - 1)] + gray[(y + 1) * width + (x + 1)]
        const gy =
          -gray[(y - 1) * width + (x - 1)] - 2 * gray[(y - 1) * width + x] - gray[(y - 1) * width + (x + 1)]
          + gray[(y + 1) * width + (x - 1)] + 2 * gray[(y + 1) * width + x] + gray[(y + 1) * width + (x + 1)]
        edges[idx] = Math.sqrt(gx * gx + gy * gy)
      }
    }

    // Darken pixels along strong edges (creates cartoon-like outlines)
    for (let i = 0; i < data.length; i += 4) {
      const edge = edges[i / 4]
      if (edge > 30) {
        const factor = Math.min(1, (edge - 30) / 100) * 0.45
        data[i]     *= (1 - factor)
        data[i + 1] *= (1 - factor)
        data[i + 2] *= (1 - factor)
      }
    }

    ctx.putImageData(imageData, 0, 0)
  }

  async function transform(canvas, sourceDataUrl, animalKey) {
    isProcessing.value = true
    try {
      const img = await new Promise((resolve, reject) => {
        const i = new Image()
        i.onload = () => resolve(i)
        i.onerror = reject
        i.src = sourceDataUrl
      })

      canvas.width = img.width
      canvas.height = img.height
      const ctx = canvas.getContext('2d')

      ctx.drawImage(img, 0, 0)

      applyCartoonEffect(ctx, canvas.width, canvas.height)

      const animalImg = await loadAnimalImage(animalKey)

      let bounds = faceBounds.value
      if (!bounds) {
        bounds = getProportionalBounds(img.width, img.height)
      }
      const placement = getAnimalPlacement(animalKey, bounds)

      ctx.drawImage(animalImg, placement.x, placement.y, placement.w, placement.h)
    } finally {
      isProcessing.value = false
    }
  }

  function selectAnimal(key) {
    currentAnimal.value = key
  }

  function reset() {
    currentAnimal.value = 'cat'
    isProcessing.value = false
    faceBounds.value = null
  }

  return {
    currentAnimal,
    isProcessing,
    faceBounds,
    animalKeys,
    animals: ANIMALS,
    loadAnimalImage,
    detectFace,
    selectAnimal,
    transform,
    reset,
  }
}
