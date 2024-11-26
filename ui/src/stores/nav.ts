import { defineStore } from 'pinia'
import { useAxiosClient } from '@/stores/axios'

const axios = useAxiosClient()

type NavigationRecord = {
  to: string
  from: string
}

export const useNavStore = defineStore('nav', () => {
  const recordNavigation = (to: string, from: string) => {
    const record: NavigationRecord = { to, from }
    axios
      .post('/api/navigation', record)
      .catch(() => {
        // ignore errors
      })
  }

  return {
    recordNavigation,
  }
})
