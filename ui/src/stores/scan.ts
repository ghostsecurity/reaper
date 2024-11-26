/**
 * Scan management store
 */
import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAxiosClient } from '@/stores/axios'

const axios = useAxiosClient()

export type Domain = {
  id?: number
  name: string
  status?: 'pending' | 'scanning' | 'probing' | 'success' | 'error' | 'completed'
  auto_scan: boolean
  host_count?: number
  last_scanned_at?: Date
  selected?: boolean
}

type Host = {
  id?: number
  domain_id: number
  project_id?: number
  name: string
  status: string
  status_code?: number
  source: string
  scheme?: string
  content_type?: string
  cdn_name?: string
  cdn_type?: string
  webserver?: string
  tech?: string
  created_at: Date
  updated_at: Date
}

export const useScanStore = defineStore('scan', () => {
  const domains = ref<Domain[]>([])
  const hosts = ref<Host[]>([])
  const errors = ref<string>('')
  const selectedDomain = ref<Domain | null>(null)

  const addDomain = (domain: Domain) => {
    domains.value.push(domain)
  }

  /**
   * Create a domain
   * @param domain - The domain to create
   * @returns Promise<Domain> - The created domain
   */
  const createDomain = (domain: Domain) => {
    return new Promise((resolve, reject) => {
      axios
        .post('/api/scan/domains', domain)
        .then((response) => {
          addDomain(response.data)
          errors.value = '' // clear any previous errors
          resolve(response.data)
        })
        .catch((error) => {
          errors.value = error?.response?.data?.error ?? error.message
          reject(error)
        })
    })
  }

  /**
   * Get all domains
   */
  const getDomains = () => {
    axios
      .get('/api/scan/domains')
      .then((response) => {
        domains.value = response.data
      })
      .catch((error) => {
        console.error(error)
        errors.value = error.message
      })
  }

  /**
   * Delete a domain
   * @param domain - The domain to delete
   */
  const deleteDomain = (domain: Domain) => {
    axios
      .delete(`/api/scan/domains/${domain.id}`)
      .then(() => {
        domains.value = domains.value.filter((d) => d.id !== domain.id)
      })
      .catch((error) => {
        console.error("[scan.ts]", error)
        errors.value = error.message
      })
      .finally(() => {
        if (domains.value.length > 0) {
          const newestDomain = domains.value[domains.value.length - 1]
          selectDomain(newestDomain)
        }
      })
  }

  /**
   * Get hosts for the domain
   * @param domain - The domain to get hosts for
   */
  const getHosts = (domain: Domain) => {
    axios
      .get(`/api/scan/domains/${domain.id}/hosts`)
      .then((response) => {
        hosts.value = response.data
      })
      .catch((error) => {
        console.error("[scan.ts]", error)
        errors.value = error.message
      })
  }

  /**
   * Select active domain to show in Display component
   * @param domain - The domain to select
   */
  const selectDomain = (domain: Domain) => {
    selectedDomain.value = domain
  }

  /**
   * Sync the domain with the latest data from websocket broadcast message
   * @param domain - The domain to sync
   */
  const syncDomain = (domain: Domain) => {
    const index = domains.value.findIndex((d) => d.id === domain.id)
    if (index !== -1) {
      domains.value[index] = { ...domains.value[index], ...domain }
    }
  }

  return {
    addDomain,
    getDomains,
    getHosts,
    createDomain,
    deleteDomain,
    selectDomain,
    syncDomain,
    domains,
    hosts,
    selectedDomain,
    errors,
  }
})
