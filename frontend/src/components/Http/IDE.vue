<script lang="ts" setup>
import { ArrowsPointingInIcon, ArrowsPointingOutIcon, ChevronDownIcon, XMarkIcon } from '@heroicons/vue/20/solid'
import { computed, PropType, ref, watch } from 'vue'
import {
  Listbox,
  ListboxButton,
  ListboxOption,
  ListboxOptions,
  Menu,
  MenuButton,
  MenuItem,
  MenuItems,
} from '@headlessui/vue'
import { Headers } from '../../lib/http'
import CodeEditor from './CodeEditor.vue'
import KeyValEditor from '../KeyValEditor.vue'
import Client from '../../lib/api/Client'
import { HttpRequest, HttpResponse, KeyValue } from '../../lib/api/packaging'

const props = defineProps({
  request: { type: Object as PropType<HttpRequest>, required: true },
  actions: {
    type: Object as PropType<Map<string, string>>,
    required: false,
    default: new Map<string, string>([['send', 'Send']]),
  },
  readonly: { type: Boolean, default: true },
  fullscreen: { type: Boolean, default: false },
  showButtons: { type: Boolean, default: true },
  client: { type: Object as PropType<Client>, required: true },
})

const emit = defineEmits(['action', 'close', 'fullscreen', 'request-update', 'create-workflow-from-request'])
const defaultAction = ref(props.actions.keys().next().value)
const extraActions = ref([...props.actions.keys()].filter(key => key !== defaultAction.value))

watch(
  () => props.actions,
  () => {
    defaultAction.value = props.actions.keys().next().value
    extraActions.value = [...props.actions.keys()].filter(key => key !== defaultAction.value)
  },
)

const selectedMethod = ref(props.request.method)

watch(
  () => props.request,
  () => {
    selectedMethod.value = props.request.method
  },
)

const methods = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS', 'TRACE']

const requestTabs = ref([
  {
    id: 'headers',
    name: 'Headers',
    current: true,
  },
  {
    id: 'query',
    name: 'Query',
  },
  {
    id: 'body',
    name: 'Body',
  },
])

const responseTabs = ref([
  {
    id: 'headers',
    name: 'Headers',
  },
  {
    id: 'body',
    name: 'Body',
    current: true,
  },
])

const requestTab = computed(() => requestTabs.value.find(tab => tab.current)?.id)
const responseTab = computed(() => responseTabs.value.find(tab => tab.current)?.id)

function switchRequestTab(id: string) {
  requestTabs.value = requestTabs.value.map(tab => {
    const updatedTab = tab
    updatedTab.current = updatedTab.id === id
    return updatedTab
  })
}

function switchResponseTab(id: string) {
  responseTabs.value = responseTabs.value.map(tab => {
    const updatedTab = tab
    updatedTab.current = updatedTab.id === id
    return updatedTab
  })
}

function selectRequestTab(e: Event) {
  switchRequestTab((e.target as HTMLSelectElement).value)
}

function selectResponseTab(e: Event) {
  switchResponseTab((e.target as HTMLSelectElement).value)
}

function toggleFullscreen(on: boolean) {
  emit('fullscreen', on)
}

function buildQueryString(params: KeyValue[]): string {
  let qs = ''
  for (let i = 0; i < params.length; i += 1) {
    const param = params[i]
    if (i > 0) {
      qs += '&'
    }
    qs += param.key
    if (param.value !== '') {
      qs += `=${encodeURIComponent(param.value)}`
    }
  }
  return qs
}

function updateURL(u: string) {
  const clone = cloneRequest(props.request) as HttpRequest
  try {
    const url = new URL(u)
    clone.url = u
    clone.host = url.host
    clone.path = url.pathname
    clone.scheme = url.protocol
    clone.query = [] as KeyValue[]
    url.searchParams.forEach((value: string, key: string) => {
      clone.query.push({ key, value })
    })
    clone.query_string = buildQueryString(clone.query)
  } catch (e) {
    clone.url = u
  }
  emit('request-update', clone)
}

function cloneRequest(input: HttpRequest | null): HttpRequest | null {
  if (input === null || input === undefined) {
    return null
  }
  return {
    id: input.id,
    local_id: input.local_id,
    url: input.url,
    method: input.method,
    scheme: input.scheme,
    host: input.host,
    path: input.path,
    body: input.body,
    query_string: input.query_string,
    tags: cloneTags(input.tags),
    query: cloneKeyValues(input.query),
    headers: cloneKeyValues(input.headers),
    response: cloneResponse(input.response),
  }
}

function cloneTags(tags: Array<string>): Array<string> {
  const clone: Array<string> = []
  for (let i = 0; i < tags.length; i += 1) {
    clone.push(tags[i])
  }
  return clone
}

function cloneKeyValues(kv: KeyValue[]): KeyValue[] {
  const clone: KeyValue[] = []
  for (let i = 0; i < kv.length; i += 1) {
    clone.push({ key: kv[i].key, value: kv[i].value })
  }
  return clone
}

function cloneResponse(input: HttpResponse | null): HttpResponse | null {
  if (input === null || input === undefined) {
    return null
  }
  return {
    id: input.id,
    local_id: input.local_id,
    body: input.body,
    status_code: input.status_code,
    request: cloneRequest(input.request),
    headers: cloneKeyValues(input.headers),
    tags: cloneTags(input.tags),
    body_size: input.body_size,
    cookies: cloneKeyValues(input.cookies),
  }
}

function buildURL(c: HttpRequest): string {
  let url = c.scheme
  if (!url) {
    url = 'https'
  }
  if (!url.endsWith(':')) {
    url += ':'
  }
  url += `//${c.host}${c.path}`
  if (c.query_string !== '') {
    url += `?${c.query_string}`
  }
  return url
}

function updateQuery(params: KeyValue[]) {
  const clone = cloneRequest(props.request) as HttpRequest
  clone.query = params
  clone.query_string = buildQueryString(params)
  clone.url = buildURL(clone)
  emit('request-update', clone)
}

function updateMethod(method: string) {
  const clone = cloneRequest(props.request) as HttpRequest
  clone.method = method
  emit('request-update', clone)
}

function updateHeaders(headers: KeyValue[]) {
  const clone = cloneRequest(props.request) as HttpRequest
  clone.headers = headers
  emit('request-update', clone)
}

function updateBody(b: string) {
  const clone = cloneRequest(props.request) as HttpRequest
  clone.body = b
  emit('request-update', clone)
}

function updateResponseHeaders(headers: KeyValue[]) {
  const clone = cloneRequest(props.request) as HttpRequest
  if (clone.response === null) {
    return
  }
  clone.response.headers = headers
  emit('request-update', clone)
}

function updateResponseBody(b: string) {
  const clone = cloneRequest(props.request) as HttpRequest
  if (clone.response === null) {
    return
  }
  clone.response.body = b
  emit('request-update', clone)
}

function getMime(): string {
  const headers = props.request?.headers || []
  const mime = getContentType(headers, '')
  if (mime !== '') {
    return mime
  }
  if (props.request?.body === '') {
    return 'text/plain'
  }
  const body = props.request?.body.trim() || ''
  if (body.startsWith('{') || body.startsWith('[')) {
    return 'application/json'
  }
  if (body.startsWith('<')) {
    return 'application/xml'
  }
  return 'text/plain'
}

function getResponseMime(): string {
  const headers = props.request?.response?.headers || []
  const mime = getContentType(headers, '')
  if (mime !== '') {
    return mime
  }
  if (props.request?.response?.body === '') {
    return 'text/plain'
  }
  const body = props.request?.response?.body.trim() || ''
  if (body.startsWith('{') || body.startsWith('[')) {
    return 'application/json'
  }
  if (body.startsWith('<')) {
    return 'text/html'
  }
  return 'text/plain'
}

function getContentType(headers: KeyValue[], def: string): string {
  for (let i = 0; i < headers.length; i += 1) {
    const header = headers[i]
    if (header.key.toLowerCase() === 'content-type') {
      return header.value
    }
  }
  return def
}
</script>

<template>
  <div class="relative flex h-full max-h-full min-w-0 max-w-full flex-col overflow-hidden py-2">
    <div class="flex-0 flex items-stretch text-xs">
      <!-- Button for method (GET, POST etc.) -->
      <div class="flex-0 rounded-l-md border border-frost-1 bg-frost-1 py-2 hover:bg-frost-1/80">
        <span v-if="readonly" class="px-4">{{ selectedMethod }}</span>
        <Listbox v-else as="div" @update:model-value="updateMethod($event)" v-model="selectedMethod"
                 class="relative m-0 p-0 text-left align-middle">
          <div class="relative">
            <ListboxButton
                class="focus:ring-none relative w-full cursor-pointer pl-4 pr-10 text-left align-middle shadow-sm focus:border-none focus:outline-none">
              <span class="truncate">{{ selectedMethod }}</span>
              <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                <ChevronDownIcon class="h-5 w-5" aria-hidden="true"/>
              </span>
            </ListboxButton>

            <transition leave-active-class="transition ease-in duration-100" leave-from-class="opacity-100"
                        leave-to-class="opacity-0">
              <ListboxOptions style="z-index: 11"
                              class="max-h-70 absolute z-10 mt-1 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none dark:bg-polar-night-4 sm:text-sm">
                <ListboxOption as="template" v-for="method in methods" :key="method" :value="method"
                               v-slot="{ active, selected }">
                  <li :class="[
                    active ? 'bg-frost-1 text-white' : 'text-gray-900 dark:text-snow-storm-1',
                    'relative cursor-pointer select-none px-1 py-2',
                  ]">
                    <span :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']">
                      {{ method }}
                    </span>
                  </li>
                </ListboxOption>
              </ListboxOptions>
            </transition>
          </div>
        </Listbox>
      </div>

      <!-- Input for full URL -->
      <div class="flex-1 bg-snow-storm-1 align-middle dark:bg-polar-night-4">
        <input type="text" :readonly="readonly" :value="request.url"
               @input="updateURL(($event.target as HTMLInputElement).value)"
               @change="updateURL(($event.target as HTMLInputElement).value)"
               class="my-0 h-full w-full border-none bg-transparent py-0 text-sm outline-none ring-0 focus:border-transparent focus:outline-none focus:ring-0"/>
      </div>

      <!-- default action button -->
      <button v-if="defaultAction" type="button" :class="[
        extraActions.length > 0 ? 'rounded-r-none' : 'rounded-r-md',
        'flex-0 bg-frost-4 px-4 font-semibold text-snow-storm shadow-sm',
        'my-0 py-0 text-center align-middle hover:bg-frost-4/80 focus:outline-none',
      ]" @click="emit('action', defaultAction)">
        {{ actions.get(defaultAction) }}
      </button>

      <!-- Button/dropdown for Send/Export etc. -->
      <Menu v-if="extraActions.length > 0" as="div" class="flex-0 relative m-0 p-0 text-left align-middle">
        <div class="h-full rounded-r-md border border-frost-4 bg-frost-4 hover:dark:bg-frost-4/80">
          <MenuButton
              class="m-0 h-full w-full justify-center px-0 align-middle font-medium text-gray-700 focus:outline-none">
            <ChevronDownIcon class="h-5 w-5" aria-hidden="true"/>
          </MenuButton>
        </div>

        <transition enter-active-class="transition ease-out duration-100"
                    enter-from-class="transform opacity-0 scale-95"
                    enter-to-class="transform opacity-100 scale-100" leave-active-class="transition ease-in duration-75"
                    leave-from-class="transform opacity-100 scale-100" leave-to-class="transform opacity-0 scale-95">
          <MenuItems
              class="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-md bg-snow-storm-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none dark:bg-polar-night-4">
            <div class="py-1">
              <MenuItem v-slot="{ active }" v-for="action in extraActions" :key="action">
                <a @click="emit('action', action)" :class="[
                active ? 'bg-frost-1 text-white' : 'text-gray-800 dark:text-snow-storm-1',
                'block cursor-pointer px-4 py-2 text-sm',
              ]">
                  {{ actions.get(action) }}
                </a>
              </MenuItem>
            </div>
          </MenuItems>
        </transition>
      </Menu>

      <div v-if="!defaultAction"
           class="flex-0 w-2 rounded-r-md bg-snow-storm-1 align-middle dark:bg-polar-night-4"></div>

      <div v-if="showButtons" class="flex-0 -mr-1 flex pl-2 pt-1 text-right align-middle">
        <a class="cursor-pointer pt-1 text-gray-400 hover:text-polar-night-1 dark:hover:text-snow-storm-1"
           @click="toggleFullscreen(!fullscreen)">
          <ArrowsPointingOutIcon v-if="!fullscreen" class="h-5 w-5"/>
          <ArrowsPointingInIcon v-else class="h-5 w-5"/>
        </a>
        <a class="cursor-pointer text-gray-400 hover:text-polar-night-1 dark:hover:text-snow-storm-1"
           @click="emit('close')">
          <XMarkIcon class="h-7 w-7"/>
        </a>
      </div>
    </div>

    <!-- request -->
    <div :class="[request.response ? 'h-1/2' : 'h-full', 'overflow-y-hidden']">
      <div class="flex h-full flex-col">
        <!-- request tabs (headers, body, etc. )-->
        <div class="flex-0 min-h-16 h-16 max-h-16 px-2">
          <div class="sm:hidden">
            <label for="tabs" class="sr-only">Select a tab</label>
            <select @change="selectRequestTab" id="tabs" name="tabs"
                    class="block w-full rounded-md bg-polar-night-2 text-sm text-snow-storm-1 focus:border-frost focus:ring-frost">
              <option v-for="tab in requestTabs" :key="tab.id" :selected="tab.current" :value="tab.id">
                {{ tab.name }}
              </option>
            </select>
          </div>
          <div class="hidden sm:block">
            <div class="border-b dark:border-polar-night-4">
              <nav class="-mb-px flex space-x-8" aria-label="Tabs">
                <a v-for="tab in requestTabs" @click="switchRequestTab(tab.id)" :key="tab.id" :class="[
                  tab.current
                    ? 'border-frost text-frost'
                    : 'border-transparent text-gray-400 hover:border-gray-500 hover:text-gray-200',
                  'group cursor-pointer border-b-2 px-1 pb-1 pt-5 text-xs',
                ]" :aria-current="tab.current ? 'page' : undefined">
                  <span>{{ tab.name }}</span>
                </a>
              </nav>
            </div>
          </div>
        </div>

        <div class="w-full min-w-0 flex-1 overflow-y-auto">
          <KeyValEditor v-if="requestTab === 'headers'" :data="request.headers" :readonly="readonly"
                        :key-suggestions="Object.keys(Headers)" @publish="updateHeaders($event)"/>
          <KeyValEditor v-else-if="requestTab === 'query'" :data="request.query" :readonly="readonly"
                        @publish="updateQuery($event)"/>
          <CodeEditor v-else-if="requestTab === 'body'" :client="client" :code="request.body" :mime="getMime()"
                      :readonly="readonly"
                      @change="updateBody($event)"/>
        </div>
      </div>
    </div>

    <!-- response -->
    <div v-if="request.response" class="h-1/2 overflow-y-hidden pt-4">
      <div class="flex h-full flex-col">
        <h3 class="bg-snow-storm-1 p-2 text-left text-sm text-polar-night-4 dark:bg-polar-night-4 dark:text-snow-storm-1">
          Response ({{ request.response.status_code }})
        </h3>
        <div class="flex-0 min-h-16 h-16 max-h-16 px-2">
          <div class="sm:hidden">
            <label for="tabs" class="sr-only">Select a tab</label>
            <select @change="selectResponseTab" id="tabs" name="tabs"
                    class="block w-full rounded-md bg-polar-night-2 text-sm text-snow-storm-1 focus:border-frost focus:ring-frost">
              <option v-for="tab in responseTabs" :key="tab.id" :selected="tab.current" :value="tab.id">
                {{ tab.name }}
              </option>
            </select>
          </div>
          <div class="hidden sm:block">
            <div class="border-b dark:border-polar-night-4">
              <nav class="-mb-px flex space-x-8" aria-label="Tabs">
                <a v-for="tab in responseTabs" @click="switchResponseTab(tab.id)" :key="tab.id" :class="[
                  tab.current
                    ? 'border-frost text-frost'
                    : 'border-transparent text-gray-400 hover:border-gray-500 hover:text-gray-200',
                  'group cursor-pointer border-b-2 px-1 pb-1 pt-5 text-xs',
                ]" :aria-current="tab.current ? 'page' : undefined">
                  <span>{{ tab.name }}</span>
                </a>
              </nav>
            </div>
          </div>
        </div>

        <div class="w-full min-w-0 flex-1 overflow-y-auto">
          <KeyValEditor v-if="responseTab === 'headers'" :data="request.response.headers" :readonly="readonly"
                        :key-suggestions="Object.keys(Headers)" @publish="updateResponseHeaders($event)"/>
          <CodeEditor v-else-if="responseTab === 'body'" :client="client" :code="request.response.body"
                      :readonly="readonly"
                      :mime="getResponseMime()" @change="updateResponseBody($event)"/>
        </div>
      </div>
    </div>
  </div>
</template>
