<template>
  <TooltipProvider :delay-duration="0">
    <ResizablePanelGroup id="resize-panel-group-1"
                         direction="horizontal"
                         class="h-full items-stretch">
      <Separator orientation="vertical" />
      <ResizablePanel id="resize-panel-2"
                      :default-size="defaultLayout[1]"
                      :min-size="20">
        <div class="flex items-center justify-between px-4 py-3">
          <h1 class="text-xl font-bold">
            Explore
          </h1>
          <div class="flex items-center space-x-2">
            <Label for="proxy-status"
                   class="text-xs text-muted-foreground">Proxy {{ proxyStatus ? 'on' : 'off' }}</Label>
            <Switch id="proxy-status"
                    :checked="proxyStatus"
                    @click.prevent="handleProxyStatusToggle" />
          </div>
        </div>

        <Separator />
        <HelperHint hint-key="explore.main">
          Enable the proxy to explore targets.
        </HelperHint>
        <div class="bg-background/95 p-4 backdrop-blur supports-[backdrop-filter]:bg-background/60">
          <form>
            <div class="relative">
              <Search class="absolute left-2 top-2.5 size-4 text-muted-foreground" />
              <Input v-model="searchValue"
                     placeholder="filter..."
                     class="pl-8" />
            </div>
          </form>
        </div>
        <ExploreList :hosts="filteredHostList"
                     v-model:searchValue="searchValue" />
      </ResizablePanel>
      <ResizableHandle id="resize-handle-2"
                       with-handle />
      <ResizablePanel id="resize-panel-3"
                      :default-size="defaultLayout[2]">
        <ExploreDisplay :endpoint="undefined" />
      </ResizablePanel>
    </ResizablePanelGroup>
  </TooltipProvider>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue'
import { refDebounced } from '@vueuse/core'
import ExploreList from './ExploreList.vue'
import ExploreDisplay from './ExploreDisplay.vue'
import { Separator } from '@/components/ui/separator'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { TooltipProvider } from '@/components/ui/tooltip'
import { Search } from 'lucide-vue-next'
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from '@/components/ui/resizable'
import HelperHint from './HelperHint.vue'

import { useExploreStore } from '@/stores/explore'
import { type Host } from '@/stores/explore'
const searchValue = ref('')
const debouncedSearch = refDebounced(searchValue, 250)
const defaultLayout = ref([20, 30, 70])
const exploreStore = useExploreStore()
const hosts = computed(() => exploreStore.hosts)
/**
 * TODO: fix this; we don't have and endpoint ID, only the current endpoints stored 
 * in the store. If we want to navigate to a specific endpoint, we need to find it
 * by host index (ExploreList.vue `i`) and endpoint index (ExploreList.vue `j`).
 */
// const selectedEndpoint = defineModel<number>('selectedEndpoint', { required: false })
// const selectedEndpointData = computed(() => exploreStore.endpoints.find(item => item.id === selectedEndpoint.value))
const proxyStatus = computed(() => exploreStore.proxy.enabled)

const handleProxyStatusToggle = () => {
  if (exploreStore.proxy.enabled) {
    exploreStore.proxyStop()
  } else {
    exploreStore.proxyStart()
  }
}

const filteredHostList = computed(() => {
  let output: Host[] = []
  const searchValue = debouncedSearch.value?.trim()
  if (!searchValue) {
    output = hosts.value
  }
  else {
    output = hosts.value.filter((item) => {
      return item.name.includes(debouncedSearch.value)
        || item.endpoints.some(e => e.path.includes(debouncedSearch.value))
    })
  }

  return output
})

onMounted(() => {
  exploreStore.proxyStatus()
})
</script>
