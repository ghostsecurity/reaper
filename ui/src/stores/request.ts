/**
 * Requests store
 */
import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAxiosClient } from '@/stores/axios'
const axios = useAxiosClient()

export type ReaperResponse = {
  content_type: string
  content_length: number
  status: number
  status_text: string
  status_code: number
  body: string
}

export type ReaperRequest = {
  id: number
  headers: string
  host: string
  method: string
  url: string
  body: string
  response: ReaperResponse
  created_at: string
}

export const useRequestStore = defineStore('request', () => {
  const requests = ref<ReaperRequest[]>([])

  const getRequests = () => {
    axios
      .get('/api/requests')
      .then((response) => {
        requests.value = response.data
      })
      .catch((error) => {
        console.error("[request.ts]", error)
      })
  }

  return {
    requests,
    getRequests,
  }
})
