<template>
  <ScrollArea class="flex h-screen">
    <div class="flex flex-1 flex-col gap-2 p-4 pt-0">
      <TransitionGroup name="list"
                       appear>
        <button v-for="item of items"
                class="flex flex-col items-start gap-2 rounded-md border bg-background/95 p-3 text-left text-sm transition-all hover:bg-accent/50"
                :key="item.id"
                :class="selectedRequest === item.id && 'bg-muted'"
                @click="selectedRequest = item.id">
          <div class="flex w-full flex-col gap-1">
            <div class="flex items-center">
              <div class="flex min-w-0 flex-1 items-center gap-2">
                <RequestMethod :code="item.response.status_code">
                  {{ item.method }}
                </RequestMethod>
                <div class="truncate text-xs font-semibold">{{ pathFromURI(item.url) }}</div>
              </div>
            </div>

            <div class="flex items-center justify-between gap-2">
              <div class="text-xs font-medium text-muted-foreground">
                {{ hostFromURI(item.url) }}
              </div>
              <div class="ml-auto text-xs"
                   :class="selectedRequest === item.id ? 'text-foreground' : 'text-muted-foreground'">
                {{ formatDistanceToNow(new Date(item.created_at), { addSuffix: true }) }}
              </div>
            </div>
          </div>
        </button>
      </TransitionGroup>
    </div>
  </ScrollArea>
</template>

<script lang="ts" setup>
import { formatDistanceToNow } from 'date-fns'
import { ScrollArea } from '@/components/ui/scroll-area'
import RequestMethod from '@/components/shared/RequestMethod.vue'

import { type ReaperRequest } from '@/stores/request'

interface ReplayListProps {
  items: ReaperRequest[]
}

defineProps<ReplayListProps>()
const selectedRequest = defineModel<number>('selectedRequest', { required: false })

const pathFromURI = (uri: string) => {
  const url = new URL(uri)
  return url.pathname
}

const hostFromURI = (uri: string) => {
  const url = new URL(uri)
  return url.host
}
</script>