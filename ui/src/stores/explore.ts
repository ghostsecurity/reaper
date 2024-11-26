import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useAxiosClient } from '@/stores/axios'

const axios = useAxiosClient()

export type Host = {
  id?: number
  name: string
  endpoints: Endpoint[]
}

export type Endpoint = {
  id?: number
  method?: string
  path: string
  status: number
}

type Proxy = {
  enabled: boolean
}

export const useExploreStore = defineStore('explore', () => {
  const hosts = ref<Host[]>([])
  const proxy = ref<Proxy>({ enabled: false })

  const addHost = (host: Host) => {
    console.info("adding host", host.name)
    hosts.value.push(host)
  }

  const addEndpoint = (hostname: string, endpoint: Endpoint) => {
    const host = hosts.value.find(h => h.name === hostname)
    if (host) {
      if (!host.endpoints) {
        console.info("adding first endpoint to host", host.name, endpoint.path)
        host.endpoints = [endpoint]
      } else {
        console.info("adding endpoint to host", host.name, endpoint.path)
        // TODO: improve precision of dupe checking
        if (!host.endpoints.find(e => (e.path === endpoint.path && e.method === endpoint.method))) {
          host.endpoints.push(endpoint)
        }
      }
    } else {
      console.warn(`Host ${hostname} not found. Creating new host.`)
      addHost({ name: hostname, endpoints: [endpoint] })
    }
  }

  const numEndpoints = computed(() => {
    return hosts.value.reduce((total, host) => total + host.endpoints.length, 0)
  })

  const proxyStart = () => {
    axios.post('/api/proxy/start')
      .then(() => {
        proxy.value.enabled = true
      })
      .catch(() => {
        proxy.value.enabled = false
      })
  }

  const proxyStop = () => {
    axios.post('/api/proxy/stop')
      .then(() => {
        proxy.value.enabled = false
      })
      .catch(() => {
        proxy.value.enabled = false
      })
  }

  const proxyStatus = () => {
    axios.get('/api/proxy/status')
      .then(() => {
        proxy.value.enabled = true
      })
      .catch(() => {
        proxy.value.enabled = false
      })
  }

  return {
    hosts,
    addHost,
    addEndpoint,
    numEndpoints,
    proxy,
    proxyStart,
    proxyStop,
    proxyStatus,
  }
})
