import { ref } from 'vue'

export function useCountdown(duration = 3) {
  const count = ref(0)
  const isRunning = ref(false)

  function start() {
    return new Promise((resolve) => {
      count.value = duration
      isRunning.value = true

      const interval = setInterval(() => {
        count.value--
        if (count.value <= 0) {
          clearInterval(interval)
          isRunning.value = false
          resolve()
        }
      }, 1000)
    })
  }

  return {
    count,
    isRunning,
    start,
  }
}
