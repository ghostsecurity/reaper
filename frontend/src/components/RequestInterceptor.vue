<script lang="ts" setup>
import {PropType, ref, watch} from 'vue'
import {HandRaisedIcon} from '@heroicons/vue/20/solid'
import IDE from './Http/IDE.vue'
import Client from "../lib/api/Client";
import {HttpRequest} from "../lib/api/packaging";

const props = defineProps(
    {
      emptyTitle: {type: String, required: false, default: 'Nothing intercepted yet.'},
      emptyMessage: {
        type: String,
        required: false,
        default: 'Configure your interception rules and start intercepting requests.',
      },
      emptyIcon: {type: Object, required: false, default: HandRaisedIcon},
      request: {type: Object as PropType<HttpRequest | null>, required: false, default: null},
      previous: {type: Object as PropType<HttpRequest | null>, required: false, default: null},
      count: {type: Number, required: false, default: 0},
      client: {type: Object as PropType<Client>, required: true},
    },
)

const req = ref<HttpRequest | null>(props.request)
const previous = ref<HttpRequest | null>(props.previous)

watch(
    () => props.request,
    () => {
      req.value = props.request
    },
)

watch(
    () => props.previous,
    () => {
      previous.value = props.previous
    },
)

const writeActions = new Map<string, string>([
  ['send', 'Send'],
  ['drop', 'Drop'],
])

const readActions = new Map<string, string>([
  ['close', 'Next'],
])

const emit = defineEmits(['send', 'drop', 'close'])

function forwardRequest() {
  emit('send', req.value)
  previous.value = req.value
}

function dropRequest() {
  emit('drop', req.value)
  req.value = null
}

function update(r: HttpRequest | null) {
  req.value = r
}

function action(a: string) {
  switch (a) {
    case 'send':
      forwardRequest()
      break
    case 'drop':
      dropRequest()
      break
    case 'close':
      emit('close')
      break
    default:
  }
}

function closePrevious() {
  emit('close')
}
</script>

<template>
  <div class="h-full w-full">
    <IDE v-if="!!previous" :client="client" :request="previous" :readonly="true"
         :actions="readActions"
         @action="action"
         @close="closePrevious"/>
    <IDE v-else-if="!!req" :client="client" :request="req" :readonly="false"
         :actions="writeActions"
         @action="action"
         @request-update="update($event)"
         @close="dropRequest"/>
    <div v-else class="pl-8 pt-8 text-center text-frost-3">
      <component :is="emptyIcon" class="mx-auto h-12 w-12"/>
      <h3 class="mt-2 text-sm font-medium">{{ emptyTitle }}</h3>
      <p class="mt-1 text-sm">{{ emptyMessage }}</p>
    </div>
  </div>
</template>
