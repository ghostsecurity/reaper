<template>
  <div class="flex h-full flex-col">
    <div class="flex items-center p-2">
      <div class="flex items-center gap-2">
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost"
                    size="icon"
                    :disabled="!session"
                    @click="handleDeleteSession">
              <Trash2 class="size-4" />
              <span class="sr-only">Delete session</span>
            </Button>
          </TooltipTrigger>
          <TooltipContent>Delete session</TooltipContent>
        </Tooltip>
      </div>
    </div>
    <Separator />
    <div v-if="session"
         class="flex flex-1 flex-col">
      <div class="flex items-start p-4">
        <div class="flex items-start gap-4 text-sm">
          <div class="grid gap-1">
            <div class="font-semibold">
              {{ session.description }}
              <Badge v-if="session.messages && session.messages.length > 0"
                     variant="outline"
                     class="text-muted-foreground">
                {{ session.messages.length }}
              </Badge>
            </div>
          </div>
        </div>
        <div v-if="session.created_at"
             class="ml-auto text-xs text-muted-foreground">
          {{ format(new Date(session.created_at), "PPpp") }}
        </div>
      </div>
      <Separator />
      <div class="min-h-0 flex-1">
        <AgentSessionMessages id="session-messages"
                              :messages="session.messages ?? []"
                              :selected-session-id="session.id ?? 0" />
      </div>
    </div>
    <div v-else
         class="p-8 text-center text-muted-foreground">
      No session selected
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Trash2 } from 'lucide-vue-next'
import { format } from 'date-fns'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Badge } from '@/components/ui/badge'
import AgentSessionMessages from './AgentSessionMessages.vue'

import type { AgentSession } from '@/stores/agent'
import { useAgentStore } from '@/stores/agent'

interface AgentSessionDisplayProps {
  session: AgentSession | undefined
}

const agentStore = useAgentStore()
const props = defineProps<AgentSessionDisplayProps>()

const handleDeleteSession = () => {
  if (props.session) {
    agentStore.deleteSession(props.session)
  }
}
</script>
