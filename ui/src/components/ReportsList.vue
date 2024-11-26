<template>
  <ScrollArea class="flex h-screen">
    <div class="flex flex-1 flex-col gap-2 p-4 pt-0">
      <TransitionGroup name="list"
                       appear>
        <button v-for="item of items"
                class="flex flex-col items-start gap-2 rounded-md border bg-background/95 p-3 text-left text-sm transition-all hover:bg-accent/50"
                :key="item.id"
                :class="selectedReport === item.id && 'bg-muted'"
                @click="selectedReport = item.id">
          <div class="flex w-full flex-col gap-1">
            <div class="flex items-center">
              <div class="flex items-center gap-2">
                <Badge variant="outline">
                  <div class="text-2xs font-semibold">
                    {{ item.id }}
                  </div>
                </Badge>
                <span class="text-xs font-semibold">{{ item.domain }}</span>
              </div>
            </div>

            <div class="flex items-center justify-between gap-2">
              <div class="text-xs font-medium text-muted-foreground">
                {{ item.domain }}
              </div>
              <div class="ml-auto text-xs"
                   :class="selectedReport === item.id ? 'text-foreground' : 'text-muted-foreground'">
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
import type { Report } from '@/stores/report'

interface ReportListProps {
  items: Report[]
}

defineProps<ReportListProps>()
const selectedReport = defineModel<number>('selectedReport', { required: false })
</script>