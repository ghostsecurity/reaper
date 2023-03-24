<script lang="ts" setup>
import {XMarkIcon, ChevronDownIcon, ArrowsPointingOutIcon, ArrowsPointingInIcon} from '@heroicons/vue/20/solid'
import {PropType, ref, computed, watch} from 'vue'
import {
  Menu,
  MenuButton,
  MenuItem,
  MenuItems,
  Listbox,
  ListboxButton,
  ListboxOption,
  ListboxOptions,
} from '@headlessui/vue'
import {HttpRequest, HttpResponse, Headers} from '../../lib/Http'
import CodeEditor from './CodeEditor.vue'
import KeyValEditor from '../KeyValEditor.vue'
import {KeyValue} from '../../lib/KeyValue'

const props = defineProps({
  request: {type: Object as PropType<HttpRequest>, required: true},
  actions: {
    type: Object as PropType<Map<string, string>>,
    required: false,
    default: new Map<string, string>([['send', 'Send']]),
  },
  readonly: {type: Boolean, default: true},
  fullscreen: {type: Boolean, default: false},
})

const emit = defineEmits(['action', 'close', 'fullscreen', 'request-update'])
const defaultAction = ref(props.actions.keys().next().value)
const extraActions = ref([...props.actions.keys()].filter(key => key !== defaultAction.value))

watch(
    () => props.actions,
    () => {
      defaultAction.value = props.actions.keys().next().value
      extraActions.value = [...props.actions.keys()].filter(key => key !== defaultAction.value)
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

const selectedMethod = ref(methods[0])

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
  for(let i =0; i < params.length; i++) {
    let param = params[i]
    if (i > 0) {
      qs += '&'
      }
    qs += param.Key
    if (param.Value !== '') {
      qs += '=' + encodeURIComponent(param.Value)
    }
  }
  return qs
}

function updateURL(u: string) {
  let clone = cloneRequest(props.request) as HttpRequest
  try{
  const url = new URL(u)
  clone.URL = u
  clone.Host = url.host
  clone.Path = url.pathname
  clone.Scheme = url.protocol
  clone.Method = clone.Method
  clone.Query = [] as KeyValue[]
  let params = url.searchParams
  for (let p of params) {
    clone.Query.push({Key: p[0], Value: p[1]})
  }
  clone.QueryString = buildQueryString(clone.Query)
  }catch(e) {
    clone.URL = u
  }
  emit('request-update', clone)
}

function cloneRequest(input: HttpRequest|null): HttpRequest|null {
  if (input === null || input === undefined) {
    return null
  }
  let clone: HttpRequest = {
    ID: input.ID,
    LocalID: input.LocalID,
    URL: input.URL,
    Method: input.Method,
    Scheme: input.Scheme,
    Host: input.Host,
    Path: input.Path,
    Body: input.Body,
    QueryString: input.QueryString,
    Tags: cloneTags(input.Tags),
    Query: cloneKeyValues(input.Query),
    Headers: cloneKeyValues(input.Headers),
    Response: cloneResponse(input.Response),
  }
  return clone
}

function cloneTags(tags: Array<string>): Array<string> {
  let clone: Array<string> = []
  for (let i = 0; i < tags.length; i++) {
    clone.push(tags[i])
  }
  return clone
}

function cloneKeyValues(kv: KeyValue[]): KeyValue[] {
  let clone: KeyValue[] = []
  for (let i = 0; i < kv.length; i++) {
    clone.push({Key: kv[i].Key, Value: kv[i].Value})
  }
  return clone
}

function cloneResponse(input: HttpResponse | null): HttpResponse|null {
  if (input === null || input === undefined) {
     return null
  }
  let clone: HttpResponse = {
    ID: input.ID,
    LocalID: input.LocalID,
    Body: input.Body,
    StatusCode: input.StatusCode,
    Request: cloneRequest(input.Request),
    Headers: cloneKeyValues(input.Headers),
    Tags: cloneTags(input.Tags),
    BodySize: input.BodySize,
  }
  return clone
}

function buildURL(c: HttpRequest): string {
  let url = c.Scheme
  if (!url) {
    url = 'https'
  }
  if (url.endsWith(':') === false) {
    url += ':'
  }
  url += '//' + c.Host + c.Path
  if (c.QueryString !== "") {
    url += '?' + c.QueryString
  }
  return url
}

function updateQuery(params: KeyValue[]) {
  let clone = cloneRequest(props.request) as HttpRequest
  clone.Query = params
  clone.QueryString = buildQueryString(params)
  clone.URL = buildURL(clone)
  emit('request-update', clone)
}

function updateMethod(method: string) {
  let clone = cloneRequest(props.request) as HttpRequest
  clone.Method = method
  emit('request-update', clone)
}

function updateHeaders(headers: KeyValue[]) {
  let clone = cloneRequest(props.request) as HttpRequest
  clone.Headers = headers
  emit('request-update', clone)
}

function updateBody(b: string) {
  let clone = cloneRequest(props.request) as HttpRequest
  clone.Body = b
  emit('request-update', clone)
}


function updateResponseHeaders(headers: KeyValue[]) {
  let clone = cloneRequest(props.request) as HttpRequest
  if(clone.Response === null) {
     return
  }
  clone.Response.Headers = headers
  emit('request-update', clone)
}

function updateResponseBody(b: string) {
  let clone = cloneRequest(props.request) as HttpRequest
  if(clone.Response === null) {
     return
  }
  clone.Response.Body = b
  emit('request-update', clone)
}
</script>

<template>
  <div class="flex h-full max-h-full min-w-0 max-w-full flex-col overflow-hidden py-2 relative">
    <div class="flex-0 flex items-stretch text-xs">
      <!-- Button for method (GET, POST etc.) -->
      <div class="flex-0 rounded-l-md border border-frost-1 bg-frost-1 py-2 hover:bg-frost-1/80">
        <span v-if="readonly" class="px-4">{{ selectedMethod }}</span>
        <Listbox v-else as="div" @update:model-value="updateMethod($event)" v-model="selectedMethod" class="relative m-0 p-0 text-left align-middle">
          <div class="relative">
            <ListboxButton
                class="focus:ring-none relative w-full cursor-pointer pl-4 pr-10 text-left align-middle shadow-sm focus:border-none focus:outline-none">
              <span class="truncate">{{ selectedMethod }}</span>
              <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                <ChevronDownIcon class="h-5 w-5" aria-hidden="true"/>
              </span>
            </ListboxButton>

            <transition
                leave-active-class="transition ease-in duration-100"
                leave-from-class="opacity-100"
                leave-to-class="opacity-0">
              <ListboxOptions
                  style="z-index: 11"
                  class="max-h-70 absolute z-10 mt-1 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none dark:bg-polar-night-4 sm:text-sm">
                <ListboxOption
                    as="template"
                    v-for="method in methods"
                    :key="method"
                    :value="method"
                    v-slot="{ active, selected }">
                  <li
                      :class="[
                      active ? 'bg-frost-1 text-white' : 'text-gray-900 dark:text-snow-storm-1',
                      'relative cursor-pointer select-none py-2 px-1',
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
        <input
            type="text"
            :readonly="readonly"
            :value="request.URL"
            @input="updateURL(($event.target as HTMLInputElement).value)"
            @change="updateURL(($event.target as HTMLInputElement).value)"
            class="my-0 h-full w-full border-none bg-transparent py-0 text-sm outline-none ring-0 focus:border-transparent focus:outline-none focus:ring-0"/>
      </div>

      <!-- default action button -->
      <button
          type="button"
          :class="[
          extraActions.length > 0 ? 'rounded-r-none' : 'rounded-r-md',
          'flex-0 bg-frost-4 px-4 font-semibold text-snow-storm shadow-sm',
          'my-0 py-0 text-center align-middle hover:bg-frost-4/80 focus:outline-none',
        ]"
          @click="emit('action', defaultAction)">
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

        <transition
            enter-active-class="transition ease-out duration-100"
            enter-from-class="transform opacity-0 scale-95"
            enter-to-class="transform opacity-100 scale-100"
            leave-active-class="transition ease-in duration-75"
            leave-from-class="transform opacity-100 scale-100"
            leave-to-class="transform opacity-0 scale-95">
          <MenuItems
              class="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-md bg-snow-storm-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none dark:bg-polar-night-4">
            <div class="py-1">
              <MenuItem v-slot="{ active }" v-for="action in extraActions" :key="action">
                <a
                    @click="emit('action', action)"
                    :class="[
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

      <div class="flex-0 -mr-1 flex pl-2 pt-1 text-right align-middle">
        <a
            class="cursor-pointer pt-1 text-gray-400 hover:text-polar-night-1 dark:hover:text-snow-storm-1"
            @click="toggleFullscreen(!fullscreen)">
          <ArrowsPointingOutIcon v-if="!fullscreen" class="h-5 w-5"/>
          <ArrowsPointingInIcon v-else class="h-5 w-5"/>
        </a>
        <a
            class="cursor-pointer text-gray-400 hover:text-polar-night-1 dark:hover:text-snow-storm-1"
            @click="emit('close')">
          <XMarkIcon class="h-7 w-7"/>
        </a>
      </div>
    </div>

    <!-- request -->
    <div :class="[request.Response ? 'h-1/2' : 'h-full', 'overflow-y-hidden']">
      <div class="flex flex-col h-full">

        <!-- request tabs (headers, body, etc. )-->
        <div class="flex-0 min-h-16 h-16 max-h-16 px-2">
          <div class="sm:hidden">
            <label for="tabs" class="sr-only">Select a tab</label>
            <select
                @change="selectRequestTab"
                id="tabs"
                name="tabs"
                class="block w-full rounded-md bg-polar-night-2 text-sm text-snow-storm-1 focus:border-frost focus:ring-frost">
              <option v-for="tab in requestTabs" :key="tab.id" :selected="tab.current" :value="tab.id">
                {{ tab.name }}
              </option>
            </select>
          </div>
          <div class="hidden sm:block">
            <div class="border-b dark:border-polar-night-4">
              <nav class="-mb-px flex space-x-8" aria-label="Tabs">
                <a
                    v-for="tab in requestTabs"
                    @click="switchRequestTab(tab.id)"
                    :key="tab.id"
                    :class="[
                tab.current
                  ? 'border-frost text-frost'
                  : 'border-transparent text-gray-400 hover:border-gray-500 hover:text-gray-200',
                'group cursor-pointer border-b-2 px-1 pt-5 pb-1 text-xs',
              ]"
                    :aria-current="tab.current ? 'page' : undefined">
                  <span>{{ tab.name }}</span>
                </a>
              </nav>
            </div>
          </div>
        </div>

        <div class="flex-1 min-w-0 w-full overflow-y-auto">
          <KeyValEditor
              v-if="requestTab === 'headers'"
              :data="request.Headers"
              :readonly="readonly"
              :key-suggestions="Object.keys(Headers)"
              @publish="updateHeaders($event)"/>
          <KeyValEditor
              v-else-if="requestTab === 'query'"
              :data="request.Query"
              :readonly="readonly"
              @publish="updateQuery($event)"/>
          <CodeEditor v-else-if="requestTab === 'body'" :code="request.Body" :readonly="readonly" @change="updateBody($event)"/>
        </div>
      </div>
    </div>

    <!-- response -->
    <div v-if="request.Response" class="overflow-y-hidden h-1/2 pt-4">
      <div class="flex flex-col h-full">

        <h3 class="text-left text-sm p-2 text-polar-night-4 bg-snow-storm-1 dark:text-snow-storm-1 dark:bg-polar-night-4">Response ({{ request.Response.StatusCode }})</h3>
        <div class="flex-0 min-h-16 h-16 max-h-16 px-2">
          <div class="sm:hidden">
            <label for="tabs" class="sr-only">Select a tab</label>
            <select
                @change="selectResponseTab"
                id="tabs"
                name="tabs"
                class="block w-full rounded-md bg-polar-night-2 text-sm text-snow-storm-1 focus:border-frost focus:ring-frost">
              <option v-for="tab in responseTabs" :key="tab.id" :selected="tab.current" :value="tab.id">
                {{ tab.name }}
              </option>
            </select>
          </div>
          <div class="hidden sm:block">
            <div class="border-b dark:border-polar-night-4">
              <nav class="-mb-px flex space-x-8" aria-label="Tabs">
                <a
                    v-for="tab in responseTabs"
                    @click="switchResponseTab(tab.id)"
                    :key="tab.id"
                    :class="[
                tab.current
                  ? 'border-frost text-frost'
                  : 'border-transparent text-gray-400 hover:border-gray-500 hover:text-gray-200',
                'group cursor-pointer border-b-2 px-1 pt-5 pb-1 text-xs',
              ]"
                    :aria-current="tab.current ? 'page' : undefined">
                  <span>{{ tab.name }}</span>
                </a>
              </nav>
            </div>
          </div>
        </div>

        <div class="flex-1 min-w-0 w-full overflow-y-auto">
          <KeyValEditor
              v-if="responseTab === 'headers'"
              :data="request.Response.Headers"
              :readonly="readonly"
              :key-suggestions="Object.keys(Headers)"
              @publish="updateResponseHeaders($event)"/>
          <CodeEditor v-else-if="responseTab === 'body'" :code="request.Response.Body" :readonly="readonly" @change="updateResponseBody($event)"/>
        </div>
      </div>
    </div>
  </div>


</template>
