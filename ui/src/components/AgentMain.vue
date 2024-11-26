<template>
  <TooltipProvider :delay-duration="0">
    <ResizablePanelGroup id="resize-panel-group-1"
                         direction="horizontal"
                         class="h-full items-stretch">
      <Separator orientation="vertical" />
      <ResizablePanel id="resize-panel-2"
                      :default-size="defaultLayout[1]"
                      :min-size="20">
        <div class="flex items-center justify-between px-4 py-2.5">
          <h1 class="text-xl font-bold">
            Sessions
          </h1>
          <Dialog v-model:open="isModalOpen">
            <DialogTrigger as-child>
              <Button variant="outline"
                      size="sm">
                <CircleFadingPlus class="mr-2 w-4" />New Session
              </Button>
            </DialogTrigger>
            <DialogContent class="sm:max-w-[425px]">
              <DialogHeader>
                <DialogTitle>New session</DialogTitle>
                <DialogDescription>
                  Start with a descriptive name of your AI driven workflow.
                </DialogDescription>
              </DialogHeader>
              <div class="py-0">
                <Input id="name"
                       v-model="sessionName"
                       placeholder="example.com dynamic security test"
                       @keyup.enter.prevent="handleCreateSession"
                       class="w-full" />
                <div class="h-4 text-xs font-medium text-destructive">{{ errors }}</div>
              </div>
              <DialogFooter>
                <Button type="submit"
                        @click.prevent="handleCreateSession"
                        :disabled="!sessionName">
                  Create session
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
        <Separator />
        <HelperHint hint-key="agent.session.list">
          Start an AI session to experiance how context-aware AI can help you drive intelligent and efficient workflows.
        </HelperHint>
        <div class="bg-background/95 p-4 backdrop-blur supports-[backdrop-filter]:bg-background/60">
          <div class="relative">
            <Search class="absolute left-2 top-2.5 size-4 text-muted-foreground" />
            <Input v-model="searchValue"
                   placeholder="filter..."
                   class="pl-8" />
          </div>
        </div>
        <div class="px-4">
          <AgentSessionsList v-model:selected-session="selectedSession"
                             :items="filteredSessionList" />
        </div>
      </ResizablePanel>
      <ResizableHandle id="resiz-handle-2"
                       with-handle />
      <ResizablePanel id="resize-panel-3"
                      :default-size="defaultLayout[2]">
        <AgentSessionDisplay :session="selectedSessionData" />
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
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { TooltipProvider } from '@/components/ui/tooltip'
import { CircleFadingPlus } from 'lucide-vue-next'
import AgentSessionsList from './AgentSessionsList.vue'
import AgentSessionDisplay from './AgentSessionDisplay.vue'
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from '@/components/ui/resizable'
import HelperHint from '@/components/HelperHint.vue'

import { useAgentStore } from '@/stores/agent'
import type { AgentSession } from '@/stores/agent'

const agentStore = useAgentStore()
const sessions = computed(() => agentStore.sessions)
const errors = computed(() => agentStore.errors)

const selectedSession = ref<number | undefined>(sessions.value.length > 0 ? sessions.value[0].id : undefined)
const searchValue = ref('')
const debouncedSearch = refDebounced(searchValue, 250)

const isModalOpen = ref(false)
const isSubmitting = ref(false)
const sessionName = ref('')

const defaultLayout = ref([20, 30, 70])
// const navCollapsedSize = ref(2)

const handleCreateSession = () => {
  if (!sessionName.value) return

  isSubmitting.value = true
  agentStore.createSession({
    description: sessionName.value,
  })
    .then(() => {
      // close the modal and reset the form
      isModalOpen.value = false
      sessionName.value = ''
    })
    .finally(() => {
      isSubmitting.value = false
    })
}

const filteredSessionList = computed(() => {
  let output: AgentSession[] = []
  const searchValue = debouncedSearch.value?.trim()
  if (!searchValue) {
    output = sessions.value
  }
  else {
    output = sessions.value.filter((item) => {
      return item.description.includes(debouncedSearch.value)
    })
  }

  return output
})


const selectedSessionData = computed(() => sessions.value.find(item => item.id === selectedSession.value))

onMounted(() => {
  agentStore.getSessions()
  if (sessions.value.length > 0) {
    selectedSession.value = sessions.value[0].id
  }
})
</script>
