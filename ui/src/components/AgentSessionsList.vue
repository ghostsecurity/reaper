<template>
  <ScrollArea class="flex h-screen">
    <div class="flex flex-1 flex-col gap-2 py-4 pt-0">
      <TransitionGroup name="list"
                       appear>
        <button v-for="item of items"
                class="flex flex-col items-start gap-2 rounded-md border bg-background/95 p-3 text-left text-sm transition-all hover:bg-accent/50"
                :key="item.id"
                :class="selectedSession === item.id && 'bg-muted'"
                @click="selectedSession = item.id">
          <div class="flex w-full flex-col gap-1">
            <div class="flex w-full items-center">
              <div class="flex w-full items-center gap-2">
                <BrainCircuit class="size-4" />
                <div class="flex w-full justify-between">
                  <div class="text-xs font-semibold">{{ item.description }}</div>
                  <div v-if="item.messages && item.messages.length > 0"
                       class="text-right text-xs font-semibold">{{ item.messages.length }} message{{
                        item.messages.length === 1 ? '' : 's' }}</div>
                </div>
              </div>
            </div>

            <div class="flex items-center justify-between gap-2">
              <div class="text-xs font-medium text-muted-foreground">
                reaper agent session {{ item.id }}
              </div>
              <div class="ml-auto text-xs"
                   :class="selectedSession === item.id ? 'text-foreground' : 'text-muted-foreground'">
                {{ formatDistanceToNow(new Date(item.created_at ?? ''), { addSuffix: true }) }}
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
import { BrainCircuit } from 'lucide-vue-next'
import type { AgentSession } from '@/stores/agent'

interface AgentSessionListProps {
  items: AgentSession[]
}

defineProps<AgentSessionListProps>()
const selectedSession = defineModel<number>('selectedSession', { required: false })
</script>
