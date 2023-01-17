<script lang="ts" setup>
import { XMarkIcon, ChevronDownIcon, ArrowsPointingOutIcon, ArrowsPointingInIcon } from '@heroicons/vue/20/solid'
import { PropType, ref, computed, watch } from 'vue'
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
import { HttpRequest } from '../../lib/Http'
import CodeEditor from './CodeEditor.vue'
import KeyValEditor from '../Shared/KeyValEditor.vue'

const props = defineProps({
  request: { type: Object as PropType<HttpRequest>, required: true },
  actions: {
    type: Object as PropType<Map<string, string>>,
    required: false,
    default: new Map<string, string>([['send', 'Send']]),
  },
  readonly: { type: Boolean, default: true },
  fullscreen: { type: Boolean, default: false },
})

const emit = defineEmits(['action', 'close', 'fullscreen', 'headers-update'])
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

const selectedMethod = ref(methods[0])

const requestTab = computed(() => requestTabs.value.find(tab => tab.current)?.id)

function switchTab(id: string) {
  requestTabs.value = requestTabs.value.map(tab => {
    const updatedTab = tab
    updatedTab.current = updatedTab.id === id
    return updatedTab
  })
}

function selectTab(e: Event) {
  switchTab((e.target as HTMLSelectElement).value)
}

function toggleFullscreen(on: boolean) {
  emit('fullscreen', on)
}
</script>

<template>
  <div class="flex h-full max-h-full min-w-0 max-w-full flex-col overflow-hidden py-2">
    <div class="flex-0 flex items-stretch text-xs">
      <!-- Button for method (GET, POST etc.) -->
      <div class="flex-0 rounded-l-md border border-frost-1 bg-frost-1 py-2 hover:bg-frost-1/80">
        <span v-if="readonly" class="px-4">{{ selectedMethod }}</span>
        <Listbox v-else as="div" v-model="selectedMethod" class="relative m-0 p-0 text-left align-middle">
          <div class="relative">
            <ListboxButton
              class="focus:ring-none relative w-full cursor-pointer pl-4 pr-10 text-left align-middle shadow-sm focus:border-none focus:outline-none">
              <span class="truncate">{{ selectedMethod }}</span>
              <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                <ChevronDownIcon class="h-5 w-5" aria-hidden="true" />
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
          class="my-0 h-full w-full border-none bg-transparent py-0 text-sm outline-none ring-0 focus:border-transparent focus:outline-none focus:ring-0" />
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
            <ChevronDownIcon class="h-5 w-5" aria-hidden="true" />
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
          <ArrowsPointingOutIcon v-if="!fullscreen" class="h-5 w-5" />
          <ArrowsPointingInIcon v-else class="h-5 w-5" />
        </a>
        <a
          class="cursor-pointer text-gray-400 hover:text-polar-night-1 dark:hover:text-snow-storm-1"
          @click="emit('close')">
          <XMarkIcon class="h-7 w-7" />
        </a>
      </div>
    </div>

    <!-- request tabs (headers, body, etc. )-->
    <div class="flex-0 min-h-16 h-16 max-h-16 px-2">
      <div class="sm:hidden">
        <label for="tabs" class="sr-only">Select a tab</label>
        <select
          @change="selectTab"
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
              @click="switchTab(tab.id)"
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

    <div class="min-w-0">
      <KeyValEditor
        v-if="requestTab === 'headers'"
        :data="request.Headers"
        :readonly="readonly"
        @publish="emit('headers-update', $event)" />
      <KeyValEditor v-else-if="requestTab === 'query'" :data="request.Query" :readonly="readonly" />
      <CodeEditor v-else-if="requestTab === 'body'" :code="request.Body" :readonly="readonly" />
    </div>
  </div>
</template>
