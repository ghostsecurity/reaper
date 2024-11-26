/**
 * Endpoints store
 */
import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAxiosClient } from '@/stores/axios'
const axios = useAxiosClient()

export type Endpoint = {
  id?: number
  method: string
  hostname: string
  path: string
  params: string
  created_at: string
}


export type AttackResult = {
  id: number
  endpoint_id: number
  hostname: string
  ip_address: string
  port: string
  scheme: string
  template_author: string
  template_name: string
  template_severity: string
  url: string
  request: string
  response: string
  created_at: string
}

export const useEndpointStore = defineStore('endpoint', () => {
  const endpoints = ref<Endpoint[]>([])
  const results = ref<AttackResult[]>([])
  const errors = ref<string[]>([])
  const attackRunning = ref(false)
  const attackComplete = ref(false)
  const selectedParams = ref<Set<string>>(new Set())

  const getEndpoints = () => {
    axios
      .get('/api/endpoints')
      .then((response) => {
        endpoints.value = response.data
      })
      .catch((error) => {
        console.error("[request.ts]", error)
      })
  }

  const startAttack = (endpoint: Endpoint) => {
    attackRunning.value = true
    attackComplete.value = false

    const payload = {
      endpoint_id: endpoint.id,
      params: Array.from(selectedParams.value),
    }
    axios
      .post(`/api/attack`, payload)
      .then((response) => {
        console.log("attack started", response)
      })
      .catch((error) => {
        console.error("[request.ts]", error)
      })
  }

  const addResult = (result: AttackResult) => {
    results.value.push(result)
  }

  const attackCompleted = () => {
    attackComplete.value = true
    attackRunning.value = false
  }

  const clearResults = (endpoint: Endpoint) => {
    axios
      .delete(`/api/attack/${endpoint.id}/results`)
      .then(() => {
        results.value = []
        attackComplete.value = false
        attackRunning.value = false
      })
  }

  /**
   * Empty the local results array for navigation followers
   */
  const emptyResults = () => {
    results.value = []
  }

  const isParamSelected = (param: string): boolean => {
    return selectedParams.value.has(param)
  }

  const toggleParam = (param: string, checked: boolean) => {
    if (checked) {
      selectedParams.value.add(param)
    } else {
      selectedParams.value.delete(param)
    }
  }

  return {
    addResult,
    attackComplete,
    attackCompleted,
    attackRunning,
    clearResults,
    emptyResults,
    endpoints,
    errors,
    results,
    getEndpoints,
    startAttack,
    isParamSelected,
    toggleParam,
  }
})
