<template>
  <div id="wrapper"
       class="flex h-screen flex-col">
    <div id="messages"
         ref="messagesContainer"
         class="flex-1 overflow-y-auto">
      <div class="flex flex-col justify-end">
        <div class="space-y-4 p-4">
          <div v-for="(message, index) in messages"
               :key="index"
               class="flex w-max max-w-[75%] flex-col gap-2 rounded-lg px-3 py-2 text-sm"
               :class="message.author_role !== 'agent' ? 'ml-auto bg-primary text-primary-foreground' : 'bg-muted'">
            {{ message.content }}
          </div>
        </div>
      </div>
    </div>
    <div id="input"
         class="mb-8 h-36 flex-none border-t bg-background p-4">
      <form class="flex w-full items-center space-x-2"
            @submit.prevent="handleUserMessageSubmit">
        <Input v-model="input"
               placeholder="Type a message..."
               class="flex-1" />
        <Button class="flex items-center justify-center p-2.5"
                type="submit"
                :disabled="inputLength === 0">
          <PaperPlaneIcon class="size-4" />
          <span class="sr-only">Send</span>
        </Button>
      </form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch, nextTick, onMounted } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { PaperPlaneIcon } from '@radix-icons/vue'
import type { AgentMessage } from '@/stores/agent'
import { useAgentStore } from '@/stores/agent'

interface AgentSessionMessagesProps {
  messages: AgentMessage[]
  selectedSessionId: number
}

const props = defineProps<AgentSessionMessagesProps>()
const agentStore = useAgentStore()
const input = ref('')
const inputLength = computed(() => input.value.trim().length)
const messages = computed(() => props.messages)
const messagesContainer = ref<HTMLElement | null>(null)

const scrollToBottom = () => {
  console.log("tyring to scroool")
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

const handleUserMessageSubmit = () => {
  agentStore.sendUserMessage(props.selectedSessionId, input.value)
  input.value = ''
}

// when message list changes, reset input
watch(messages, () => {
  input.value = ''
})

// scroll to bottom when messages change
watch(messages, scrollToBottom, { deep: true })

onMounted(() => {
  scrollToBottom()
})
</script>
