<template>
  <TooltipProvider :delay-duration="0">
    <ResizablePanelGroup id="resize-panel-group-1"
                         direction="horizontal"
                         class="h-full items-stretch">
      <Separator orientation="vertical" />
      <ResizablePanel id="resize-panel-2"
                      :default-size="defaultLayout[1]"
                      :min-size="20">
        <Tabs default-value="all">
          <div class="flex items-center px-4 py-2">
            <h1 class="text-xl font-bold">
              Reports
            </h1>
            <TabsList class="ml-auto">
              <TabsTrigger value="all"
                           class="text-zinc-600 dark:text-zinc-200">
                All
              </TabsTrigger>
              <TabsTrigger value="unread"
                           class="text-zinc-600 dark:text-zinc-200">
                Archived
              </TabsTrigger>
            </TabsList>
          </div>
          <Separator />
          <HelperHint hint-key="reports.main">
            View and filter available reports.
          </HelperHint>
          <div class="bg-background/95 p-4 backdrop-blur supports-[backdrop-filter]:bg-background/60">
            <div class="relative">
              <Search class="absolute left-2 top-2.5 size-4 text-muted-foreground" />
              <Input v-model="searchValue"
                     placeholder="filter..."
                     class="pl-8" />
            </div>
          </div>
          <TabsContent value="all"
                       class="m-0">
            <ReportsList v-model:selected-report="selectedReport"
                         :items="filteredReportList" />
          </TabsContent>
          <TabsContent value="unread"
                       class="m-0">
            <ReportsList v-model:selected-report="selectedReport"
                         :items="filteredReportList" />
          </TabsContent>
        </Tabs>
      </ResizablePanel>
      <ResizableHandle id="resiz-handle-2"
                       with-handle />
      <ResizablePanel id="resize-panel-3"
                      :default-size="defaultLayout[2]">
        <ReportsDisplay :report="selectedReportData" />
      </ResizablePanel>
    </ResizablePanelGroup>
  </TooltipProvider>
</template>

<script lang="ts" setup>
import { computed, ref, onMounted } from 'vue'

import { Search } from 'lucide-vue-next'
import { refDebounced } from '@vueuse/core'
import { Separator } from '@/components/ui/separator'
import { Input } from '@/components/ui/input'
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs'
import { TooltipProvider } from '@/components/ui/tooltip'
import ReportsList from './ReportsList.vue'
import ReportsDisplay from './ReportsDisplay.vue'
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from '@/components/ui/resizable'
import HelperHint from '@/components/HelperHint.vue'

import { useReportStore } from '@/stores/report'
import type { Report } from '@/stores/report'

const reportStore = useReportStore()
const reports = computed(() => reportStore.reports)

const selectedReport = ref<number | undefined>(reports.value.length > 0 ? reports.value[0].id : undefined)
const searchValue = ref('')
const debouncedSearch = refDebounced(searchValue, 250)

const defaultLayout = ref([20, 30, 70])
// const navCollapsedSize = ref(2)

const filteredReportList = computed(() => {
  let output: Report[] = []
  const searchValue = debouncedSearch.value?.trim()
  if (!searchValue) {
    output = reports.value
  }
  else {
    output = reports.value.filter((item) => {
      return item.domain.includes(debouncedSearch.value)
    })
  }

  return output
})


const selectedReportData = computed(() => reports.value.find(item => item.id === selectedReport.value))

onMounted(() => {
  reportStore.getReports()
  if (reports.value.length > 0) {
    selectedReport.value = reports.value[0].id
  }
})
</script>
