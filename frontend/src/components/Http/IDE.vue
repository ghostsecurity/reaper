<script lang="ts" setup>
import { PropType, ref, watch } from 'vue'
import { HttpRequest } from '../../lib/Http.js';
import Request from "./Request.vue";
import Response from "./Response.vue";

const props = defineProps({
  request: { type: Object as PropType<HttpRequest>, required: true },
})

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
      <div class="border-b border:snow-storm-3 dark:border-polar-night-4">
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
