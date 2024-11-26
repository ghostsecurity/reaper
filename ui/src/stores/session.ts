/**
 * Session management store
 */
import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAxiosClient } from '@/stores/axios'
import { useRouter } from 'vue-router'
const axios = useAxiosClient()

type User = {
  username: string
  invite_code: string
  token?: string
  role?: string
}

export const useSessionStore = defineStore('session', () => {
  const router = useRouter()
  const currentUser = ref<User | null>(null)
  const loggedIn = ref(false)
  const loading = ref(true)
  const errors = ref<string>('')

  const signOut = () => {
    localStorage.removeItem('reaper.token')
    loggedIn.value = false
    window.location.href = '/'
  }

  /**
   * Get auth status of the user
   */
  const status = () => {
    axios
      .get('/status')
      .then((response) => {
        // if valid status, set loggedIn to true and set currentUser to the user object
        loggedIn.value = true
        currentUser.value = response.data.user
      })
      .catch(() => {
        // console.error(error)
      })
      .finally(() => {
        loading.value = false
      })
  }

  const navigationFollow = (to: string) => {
    // only viewers should follow
    if (currentUser.value?.role === 'viewer') {
      router.push(to)
    }
  }

  const register = (user: User) => {
    axios
      .post('/register', {
        username: user.username,
        invite_code: user.invite_code,
      })
      .then((response) => {
        errors.value = ''
        currentUser.value = response.data
        console.log('[session]', currentUser.value)
        localStorage.setItem('reaper.token', response.data.user.token)
        window.location.href = '/'
      })
      .catch((error) => {
        errors.value = error?.response?.data?.error ?? error.message
        console.error(error)
      })
  }

  return {
    currentUser,
    errors,
    loading,
    loggedIn,
    navigationFollow,
    signOut,
    status,
    register,
  }
})
