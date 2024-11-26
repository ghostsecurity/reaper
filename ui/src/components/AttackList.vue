<template>
  <ScrollArea class="flex h-screen">
    <div class="flex flex-1 flex-col gap-2 p-4 pt-0">
      <TransitionGroup name="list"
                       appear>
        <button v-for="item of items"
                class="flex flex-col items-start gap-2 rounded-md border bg-background/95 p-3 text-left text-sm transition-all hover:bg-accent/50"
                :key="item.id"
                :class="selectedEndpoint === item.id && 'bg-muted'"
                @click="selectedEndpoint = item.id">
          <div class="flex w-full flex-col gap-1">
            <div class="flex items-center">
              <div class="flex items-center gap-2">
                <Badge variant="outline"
                       :class="badgeColorFromMethod(item.method)">
                  <div class="text-2xs font-semibold">
                    {{ item.method }}
                  </div>
                </Badge>
                <span class="text-xs font-semibold">{{ item.path }}</span>
              </div>
            </div>

            <div class="flex items-center justify-between gap-2">
              <div class="text-xs font-medium text-muted-foreground">
                {{ item.hostname }}
              </div>
              <div class="ml-auto text-xs"
                   :class="selectedEndpoint === item.id ? 'text-foreground' : 'text-muted-foreground'">
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
import { Badge } from '@/components/ui/badge'
import type { Endpoint } from '@/stores/endpoint'

interface EndpointListProps {
  items: Endpoint[]
}

defineProps<EndpointListProps>()
const selectedEndpoint = defineModel<number>('selectedEndpoint', { required: false })

function badgeColorFromMethod(method: string) {
  if (method === 'GET')
    return 'bg-green-50 border-green-600/20 text-green-700'

  if (method === 'POST')
    return 'bg-yellow-50 border-yellow-600/20 text-yellow-700'

  return 'bg-red-50 border-red-600/20 text-red-700'
}
</script>