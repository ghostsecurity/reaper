/**
 * Replay store
 */
import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAxiosClient } from '@/stores/axios'
const axios = useAxiosClient()
import { type ReaperRequest } from '@/stores/request'

export type ReplayPayload = {
  request: ReaperRequest
  method: string
  url: string
  headers: string
  body: string
}

export type ReplayResponse = {
  headers: string
  body: string
}

export const useReplayStore = defineStore('replay', () => {
  const replay = ref<ReaperRequest | undefined>(undefined)

  /**
   * Replay a request
   * TODO: support other methods
   * @param payload - Replay payload
   * @returns Promise<ReaperRequest> - The response to the replayed request
   */
  const replayRequest = (payload: ReplayPayload) => {
    const newPayload: ReplayPayload = {
      request: payload.request,
      method: payload.method,
      url: payload.request.url,
      // new payload with modified headers and body
      headers: payload.headers,
      body: payload.body,
    }
    return new Promise<ReplayResponse>((resolve, reject) => {
      axios
        .post('/api/replay', newPayload)
        .then((response) => {
          replay.value = response.data
          resolve({
            headers: response.data.headers,
            body: response.data.body,
          })
        })
        .catch((error) => {
          reject(error)
        })
    })
  }

  /**
   * Base64 decode the response body
   * @param payload - The payload to decode
   * @returns The decoded payload
   */
  const decodePayload = (payload: string) => {
    try {
      return atob(payload)
    }
    catch (error) {
      console.error("[replay.ts] decodePayload", error)
      return payload
    }
  }

  return {
    replay,
    replayRequest,
    decodePayload,
  }
})
