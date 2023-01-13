<script lang="ts" setup>
import { XMarkIcon } from '@heroicons/vue/20/solid'
import { PropType, ref, watch } from 'vue'
import { HttpRequest } from '../../lib/Http'
import Request from './HttpRequest.vue'
import Response from './HttpResponse.vue'

const props = defineProps({
  request: { type: Object as PropType<HttpRequest>, required: true },
})

const emit = defineEmits(['close'])

const currentTab = ref('request')

watch(() => props.request, () => {
  if (props.request.Response === null) {
    currentTab.value = 'request'
  }
})

function selectTab(tab: string) {
  currentTab.value = tab
}
</script>

<template>
  <div style="height: calc(100% - 4rem)">
    <div>
      <div class="border-b border:snow-storm-3 dark:border-polar-night-4 flex">
        <div class="flex-1">
          <nav class="-mb-px flex space-x-8" aria-label="Tabs">
            <a @click="selectTab('request')" :class="['request' == currentTab ?
                         'border-frost text-frost' :
                         'border-transparent text-polar-night-4 hover:text-frost-4 hover:border-frost-4',
            'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm cursor-pointer']"
              :aria-current="'request' == currentTab ? 'page' : undefined">
              Request
            </a>
            <a v-if="request.Response" @click="selectTab('response')" :class="['response' == currentTab ?
            'border-frost text-frost' :
            'border-transparent text-polar-night-4 hover:text-frost-4 hover:border-frost-4',
            'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm cursor-pointer']"
              :aria-current="'response' == currentTab ? 'page' : undefined">
              Response
            </a>
          </nav>
        </div>
        <div class="flex-0 cursor-pointer pt-3 text-gray-400 hover:text-snow-storm-1" @click="emit('close')">
          <XMarkIcon class="w-6 h-6" />
        </div>
      </div>
    </div>
    <div class="pt-4 h-full">
      <div :class="{ 'hidden': currentTab != 'request', 'h-full': true }">
        <Request :request="request" :readonly="true" />
      </div>
      <div :class="{ 'hidden': currentTab != 'response', 'h-full': true }">
        <Response v-if="request.Response" :response="request.Response" :readonly="true" />
      </div>
    </div>
  </div>
</template>
