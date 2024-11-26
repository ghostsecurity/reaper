/**
 * Reports store
 */
import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAxiosClient } from '@/stores/axios'
const axios = useAxiosClient()

export type Report = {
  id: number
  domain: string
  markdown: string
  status: string
  created_at: string
}

export const useReportStore = defineStore('report', () => {
  const reports = ref<Report[]>([])
  const selectedReport = ref<Report | null>(null)

  /**
   * Get all reports
   */
  const getReports = () => {
    axios
      .get('/api/reports')
      .then((response) => {
        reports.value = response.data
      })
      .catch((error) => {
        console.error("[report.ts]", error)
      })
  }

  /**
   * Delete a report
   * @param report - The report to delete
   */
  const deleteReport = (report: Report) => {
    axios
      .delete(`/api/reports/${report.id}`)
      .then(() => {
        reports.value = reports.value.filter((r) => r.id !== report.id)
      })
      .catch((error) => {
        console.error("[report.ts]", error)
      })
  }

  /**
   * Select the active report
   * @param report - The report to select
   */
  const selectReport = (report: Report) => {
    selectedReport.value = report
  }

  return {
    getReports,
    deleteReport,
    selectReport,
    reports,
    selectedReport,
  }
})
