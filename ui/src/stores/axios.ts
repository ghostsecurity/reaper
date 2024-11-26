/**
 * Axios http client store
 */
import axios from 'axios'
import { AxiosInstance } from 'axios'

// axios setup
const BASE_URL = import.meta.env.PROD ? '/' : import.meta.env.VITE_BASE_URL
const token = localStorage.getItem('reaper.token')
const axiosInstance = axios.create({
  baseURL: BASE_URL,
  headers: {
    'accept': 'application/json',
    'x-reaper-token': token,
  }
})

export function useAxiosClient(): AxiosInstance {
  return axiosInstance
}
