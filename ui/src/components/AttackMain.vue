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
              Dynamic Testing
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
          <HelperHint hint-key="test.main">
            Select an endpoint to launch a dynamic test against it.
          </HelperHint>
          <div class="bg-background/95 p-4 backdrop-blur supports-[backdrop-filter]:bg-background/60">
            <div class="relative">
              <Search class="absolute left-2 top-2.5 size-4 text-muted-foreground" />
              <Input v-model="searchValue"
                     spellcheck="false"
                     placeholder="filter..."
                     class="pl-8" />
            </div>
          </div>
          <TabsContent value="all"
                       class="m-0">
            <AttackList v-model:selected-endpoint="selectedEndpoint"
                        :items="filteredEndpointList" />
          </TabsContent>
          <TabsContent value="unread"
                       class="m-0">
            <AttackList v-model:selected-endpoint="selectedEndpoint"
                        :items="filteredEndpointList" />
          </TabsContent>
        </Tabs>
      </ResizablePanel>
      <ResizableHandle id="resiz-handle-2"
                       with-handle />
      <ResizablePanel id="resize-panel-3"
                      :default-size="defaultLayout[2]">
        <AttackDisplay :endpoint="selectedEndpointData" />
      </ResizablePanel>
    </ResizablePanelGroup>
  </TooltipProvider>
</template>

<script lang="ts" setup>
import { computed, ref, onMounted } from 'vue'

import { Search } from 'lucide-vue-next'
import { refDebounced } from '@vueuse/core'
import AttackList from './AttackList.vue'
import AttackDisplay from './AttackDisplay.vue'
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

import { useEndpointStore } from '@/stores/endpoint'
import type { Endpoint } from '@/stores/endpoint'

const endpointStore = useEndpointStore()
const endpoints = computed(() => endpointStore.endpoints)

const selectedEndpoint = ref<number | undefined>(endpoints.value.length > 0 ? endpoints.value[0].id : undefined)
const searchValue = ref('')
const debouncedSearch = refDebounced(searchValue, 250)

const defaultLayout = ref([20, 30, 70])
// const navCollapsedSize = ref(2)

const filteredEndpointList = computed(() => {
  let output: Endpoint[] = []
  const searchValue = debouncedSearch.value?.trim()
  if (!searchValue) {
    output = endpoints.value
  }
  else {
    output = endpoints.value.filter((item) => {
      return item.hostname.includes(debouncedSearch.value)
        || item.path.includes(debouncedSearch.value)
    })
  }

  return output
})


const selectedEndpointData = computed(() => endpoints.value.find(item => item.id === selectedEndpoint.value))

onMounted(() => {
  endpointStore.getEndpoints()
  if (endpoints.value.length > 0) {
    selectedEndpoint.value = endpoints.value[0].id
  }
})
</script>
