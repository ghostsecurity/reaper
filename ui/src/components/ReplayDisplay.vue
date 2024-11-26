<template>
  <div class="flex h-full flex-col">
    <div class="flex items-center p-2">
      <div class="flex items-center gap-2">
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost"
                    size="icon"
                    :disabled="!request">
              <Trash2 class="size-4" />
              <span class="sr-only">Move to trash</span>
            </Button>
          </TooltipTrigger>
          <TooltipContent>Move to trash</TooltipContent>
        </Tooltip>
      </div>
      <div class="ml-auto flex items-center gap-2">
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="outline"
                    @click="handleReplay"
                    :disabled="!request">
              {{ sendButtonText }}
              <ReplaceIcon class="ml-2 size-4 scale-x-[-1] transform" />
              <span class="sr-only">{{ sendButtonText }}</span>
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            {{ sendButtonText }}
          </TooltipContent>
        </Tooltip>
      </div>
    </div>
    <Separator />
    <div v-if="request"
         class="flex flex-1 flex-col">
      <div class="flex items-start p-4">
        <div class="mr-4 flex items-start gap-4 text-sm">
          <div class="flex flex-col gap-2">
            <RequestMethod :code="request.response.status_code">
              {{ request.method }}
            </RequestMethod>
            <div class="text-center text-xs font-medium">{{ request.response.status_code }}</div>
          </div>
          <div class="grid gap-1">
            <div class="mr-4 truncate font-semibold text-foreground/80">
              {{ request.url }}
            </div>
            <div class="space-y-1 text-xs">
              <div class="font-medium text-muted-foreground">
                {{ request.host }}
              </div>
              <div>
                {{ request.response.content_type }}
              </div>
            </div>
          </div>
        </div>
        <div v-if="request.created_at"
             class="ml-auto w-36 text-xs text-muted-foreground">
          {{ formatDistanceToNow(new Date(request.created_at), { addSuffix: true }) }}
        </div>
      </div>
      <Separator />
      <div class="h-screen space-y-2 overflow-y-auto whitespace-pre-wrap bg-muted/50 p-2 text-sm text-foreground/80">
        <div class="rounded-md bg-background shadow-sm">
          <div class="flex gap-2 rounded-t-md bg-muted p-2 text-xs font-semibold">
            Request Headers
            <LoaderCircle v-if="replayInProgress"
                          class="size-4 animate-spin" />
            <Asterisk v-if="requestHeadersIsModified"
                      class="size-4 text-primary" />
          </div>
          <textarea class="min-h-48 w-full resize-y whitespace-pre rounded-sm bg-background p-2 font-mono text-xs focus:outline-none"
                    name=""
                    spellcheck="false"
                    v-model="headersText"></textarea>
        </div>
        <div class="rounded-md bg-background shadow-sm">
          <div class="flex gap-2 rounded-t-md bg-muted p-2 text-xs font-semibold">
            Request Body
            <LoaderCircle v-if="replayInProgress"
                          class="size-4 animate-spin" />
            <Asterisk v-if="requestBodyIsModified"
                      class="size-4 text-primary" />
          </div>
          <textarea class="min-h-64 w-full resize-y whitespace-pre rounded-sm bg-background p-2 font-mono text-xs focus:outline-none"
                    name=""
                    spellcheck="false"
                    v-model="bodyText"></textarea>
        </div>
        <div class="rounded-md bg-muted shadow-sm">
          <div class="flex items-center justify-between rounded-t-md p-2 text-xs font-semibold">
            <div class="flex items-center">
              Response
              <LoaderCircle v-if="replayInProgress"
                            class="size-4 animate-spin" />
            </div>
            <div class="ml-4 flex cursor-pointer items-center text-xs"
                 @click="handleCopyResponse">
              <Copy class="size-4" />
              <span class="ml-1">
                Copy
              </span>
            </div>
          </div>
          <textarea class="min-h-64 w-full resize-y whitespace-pre rounded-sm bg-muted p-2 font-mono text-xs focus:outline-none"
                    name=""
                    disabled
                    spellcheck="false"
                    v-model="responseTextFormatted"></textarea>
        </div>
      </div>
    </div>
    <div v-else
         class="p-8 text-center text-sm text-muted-foreground">
      No requests selected
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { Asterisk, Copy, LoaderCircle, ReplaceIcon, Trash2 } from 'lucide-vue-next'
import { formatDistanceToNow } from 'date-fns'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import RequestMethod from '@/components/shared/RequestMethod.vue'
import { useReplayStore } from '@/stores/replay'
import { type ReplayPayload } from '@/stores/replay'
import { type ReaperRequest } from '@/stores/request'

interface ReplayDisplayProps {
  request: ReaperRequest | undefined
}
const props = defineProps<ReplayDisplayProps>()

const replayStore = useReplayStore()
const replayInProgress = ref(false)
const originalRequestHeadersText = ref('')
const originalRequestBodyText = ref('')
const headersText = ref('')
const bodyText = ref('')
const responseText = ref('')

const requestHeadersIsModified = computed(() => {
  return headersText.value !== originalRequestHeadersText.value
})

const requestBodyIsModified = computed(() => {
  return bodyText.value !== originalRequestBodyText.value
})

const sendButtonText = computed(() => {
  if (replayInProgress.value) {
    return 'Replaying...'
  }
  return requestHeadersIsModified.value || requestBodyIsModified.value ? 'Replay modified' : 'Replay original'
})

const handleReplay = () => {
  const request: ReplayPayload = {
    request: props.request!,
    method: props.request?.method || 'GET',
    url: props.request?.url || '',
    headers: headersText.value,
    body: bodyText.value,
  }
  replayInProgress.value = true
  replayStore.replayRequest(request)
    .then((response) => {
      responseText.value = replayStore.decodePayload(response.body)
      replayInProgress.value = false
    })
    .catch((error) => {
      console.error("[ReplayDisplay.vue] handleReplay", error)
      replayInProgress.value = false
    })
}

const handleCopyResponse = () => {
  navigator.clipboard.writeText(responseTextFormatted.value)
}

const responseTextFormatted = computed(() => {
  try {
    return JSON.stringify(JSON.parse(responseText.value), null, 2)
  }
  catch {
    return responseText.value
  }
})

watch(() => props.request, (newRequest) => {
  if (newRequest) {
    originalRequestHeadersText.value = newRequest.headers || ''
    originalRequestBodyText.value = newRequest.body || ''
    headersText.value = originalRequestHeadersText.value
    bodyText.value = originalRequestBodyText.value
    if (newRequest.response.body) {
      responseText.value = newRequest.response.body
    }
  }
}, { immediate: true })
</script>
