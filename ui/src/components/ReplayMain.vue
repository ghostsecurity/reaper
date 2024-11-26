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
              Replay
            </h1>
            <TabsList class="ml-auto">
              <TabsTrigger value="all"
                           class="text-zinc-600 dark:text-zinc-200">
                All
              </TabsTrigger>
              <TabsTrigger value="unread"
                           class="text-zinc-600 dark:text-zinc-200">
                APIs
              </TabsTrigger>
            </TabsList>
          </div>
          <Separator />
          <HelperHint hint-key="replay.main">
            Select a request to manipulate and replay.
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
            <ReplayList v-model:selected-request="selectedRequest"
                        :items="filteredRequestList" />
          </TabsContent>
          <TabsContent value="unread"
                       class="m-0">
            <ReplayList v-model:selected-request="selectedRequest"
                        :items="unreadRequestList" />
          </TabsContent>
        </Tabs>
      </ResizablePanel>
      <ResizableHandle id="resiz-handle-2"
                       with-handle />
      <ResizablePanel id="resize-panel-3"
                      :default-size="defaultLayout[2]">
        <ReplayDisplay :request="selectedRequestData" />
      </ResizablePanel>
    </ResizablePanelGroup>
  </TooltipProvider>
</template>

<script lang="ts" setup>
import { computed, ref, onMounted } from 'vue'

import { Search } from 'lucide-vue-next'
import { refDebounced } from '@vueuse/core'
import ReplayList from './ReplayList.vue'
import ReplayDisplay from './ReplayDisplay.vue'
import { Separator } from '@/components/ui/separator'
import { Input } from '@/components/ui/input'
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs'
import { TooltipProvider } from '@/components/ui/tooltip'
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from '@/components/ui/resizable'
import HelperHint from '@/components/HelperHint.vue'

import { useRequestStore } from '@/stores/request'
import { type ReaperRequest } from '@/stores/request'

const requestStore = useRequestStore()
const requests = computed(() => requestStore.requests)

const selectedRequest = ref<number | undefined>(requests.value.length > 0 ? requests.value[0].id : undefined)
const searchValue = ref('')
const debouncedSearch = refDebounced(searchValue, 250)

const defaultLayout = ref([20, 30, 70])
// const navCollapsedSize = ref(2)

const filteredRequestList = computed(() => {
  let output: ReaperRequest[] = []
  const searchValue = debouncedSearch.value?.trim()
  if (!searchValue) {
    output = requests.value
  }
  else {
    output = requests.value.filter((item) => {
      return item.url.includes(debouncedSearch.value)
        || item.method.includes(debouncedSearch.value)
        || item.response.status.toString().includes(debouncedSearch.value)
        || item.body.includes(debouncedSearch.value)
    })
  }

  return output
})

const unreadRequestList = computed(() => requests.value.filter(item => item.response.content_type === 'application/json'))

const selectedRequestData = computed(() => requests.value.find(item => item.id === selectedRequest.value))

onMounted(() => {
  requestStore.getRequests()
})
</script>
