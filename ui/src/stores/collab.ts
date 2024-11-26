import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAxiosClient } from '@/stores/axios'

const axios = useAxiosClient()

type Tunnel = {
  enabled: boolean
  url?: string
}

export const useCollabStore = defineStore('collab', () => {
  const tunnel = ref<Tunnel>({ enabled: false })

  const tunnelStart = () => {
    axios.post('/api/tunnel/start')
      .then(() => {
        tunnel.value.enabled = true
      })
      .catch(() => {
        tunnel.value.enabled = false
      })
  }

  const tunnelStop = () => {
    axios.post('/api/tunnel/stop')
      .then(() => {
        tunnel.value.enabled = false
      })
      .catch(() => {
        tunnel.value.enabled = false
      })
      .finally(() => {
        tunnel.value.url = ""
      })
  }

  const tunnelStatus = () => {
    axios.get('/api/tunnel/status')
      .then((res) => {
        tunnel.value.enabled = true
        tunnel.value.url = res.data.url
      })
      .catch(() => {
        tunnel.value.enabled = false
      })
  }

  return {
    tunnel,
    tunnelStart,
    tunnelStop,
    tunnelStatus,
  }
})

