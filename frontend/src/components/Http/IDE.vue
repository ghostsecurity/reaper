<script lang="ts" setup>
import { XMarkIcon, ChevronDownIcon } from '@heroicons/vue/20/solid'
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

const props = defineProps({
  request: { type: Object as PropType<HttpRequest>, required: true },
  actions: {
    type: Object as PropType<Map<string, string>>,
    required: false,
    default: new Map<string, string>([['send', 'Send']]),
  },
  readonly: { type: Boolean, default: true },
})

const emit = defineEmits(['action', 'close'])
const defaultAction = ref(props.actions.keys().next().value)
const extraActions = ref([...props.actions.keys()].filter(key => key !== defaultAction.value))

watch(
  () => props.actions,
  () => {
    defaultAction.value = props.actions.keys().next().value
    extraActions.value = [...props.actions.keys()].filter(key => key !== defaultAction.value)
  },
)

const methods = [
  'GET',
  'POST',
  'PUT',
  'PATCH',
  'DELETE',
  'HEAD',
  'OPTIONS',
  'TRACE',
]

const requestTabs = ref([
  {
    name: 'Headers',
    current: true,
  },
  {
    name: 'Query',
  },
  {
    name: 'Body',
  },
  {
    name: 'Auth',
  },
])

const selectedMethod = ref(methods[0])

const requestTab = computed(() => requestTabs.value.find(tab => tab.current)?.name)

function switchTab(name: string) {
  requestTabs.value = requestTabs.value.map(tab => {
    const updatedTab = tab
    updatedTab.current = updatedTab.name === name
    return updatedTab
  })
}

function selectTab(e: Event) {
  switchTab((e.target as HTMLSelectElement).value)
}

</script>

<template>
  <div class="max-h-full h-full flex flex-col pb-2">

    <div class="flex-0 flex border-b border-snow-storm-1 dark:border-polar-night-4 py-1.5 mb-2">
      <div class="flex-1"><!-- title? --></div>
      <div class="flex-0">
        <a class="cursor-pointer" @click="emit('close')">
          <XMarkIcon class="h-6 w-6" />
        </a>
      </div>

    </div>

    <div class="flex-0 flex">

      <!-- Button for method (GET, POST etc.) -->
      <Listbox as="div" v-model="selectedMethod"
        class="relative inline-block text-left flex-0 border border-frost-1 bg-frost-1 rounded-l-md">
        <div class="relative">
          <ListboxButton
            class="relative w-full min-w-[7rem] cursor-pointer h-full pt-2.5 pb-2 pl-4 pr-10 text-left shadow-sm focus:border-none focus:outline-none focus:ring-none text-sm">
            <span class="block truncate">{{ selectedMethod }}</span>
            <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
              <ChevronDownIcon class="h-5 w-5" aria-hidden="true" />
            </span>
          </ListboxButton>

          <transition leave-active-class="transition ease-in duration-100" leave-from-class="opacity-100"
            leave-to-class="opacity-0">
            <ListboxOptions style="z-index: 11"
              class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white dark:bg-polar-night-4 py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
              <ListboxOption as="template" v-for="method in methods" :key="method" :value="method"
                v-slot="{ active, selected }">
                <li :class="[
                  active ?
                    'text-white bg-frost-1' :
                    'text-gray-900 dark:text-snow-storm-1',
                  'relative cursor-pointer select-none py-2 px-4'
                ]">
                  <span :class="[
                    selected ?
                      'font-semibold' :
                      'font-normal', 'block truncate text-sm'
                  ]">{{ method }}</span>
                </li>
              </ListboxOption>
            </ListboxOptions>
          </transition>
        </div>
      </Listbox>

      <!-- Input for full URL -->
      <div class="flex-1 inline-block w-full bg-snow-storm-1 dark:bg-polar-night-4">
        <input type="text" :readonly="readonly" :value="request.URL"
          class="w-full h-full text-sm bg-transparent border-none ring-0 outline-none focus:outline-none focus:border-transparent focus:ring-0" />
      </div>

      <button type="button" :class="[
        extraActions.length > 0 ?
          'rounded-r-none' :
          'rounded-r-md',
        'flex-0 items-center bg-frost-4 px-4 py-1.5 text-sm font-semibold text-snow-storm shadow-sm',
        'hover:bg-frost-4/80 focus:outline-none'
      ]" @click="emit('action', defaultAction)">
        {{ actions.get(defaultAction) }}
      </button>

      <!-- Button/dropdown for Send/Export etc. -->
      <Menu v-if="extraActions.length > 0" as="div" class="relative inline-block text-left flex-0 p-0 m-0">
        <div class="border border-frost-4 bg-frost-4 hover:dark:bg-frost-4/80 rounded-r-md overflow-hidden">
          <MenuButton class="text-sm w-full justify-center px-0 py-1.5 font-medium text-gray-700 focus:outline-none">
            <ChevronDownIcon class="h-7 w-7" aria-hidden="true" />
          </MenuButton>
        </div>

        <transition enter-active-class="transition ease-out duration-100"
          enter-from-class="transform opacity-0 scale-95" enter-to-class="transform opacity-100 scale-100"
          leave-active-class="transition ease-in duration-75" leave-from-class="transform opacity-100 scale-100"
          leave-to-class="transform opacity-0 scale-95">
          <MenuItems
            class="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-md bg-snow-storm-1 dark:bg-polar-night-4 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
            <div class="py-1">
              <MenuItem v-slot="{ active }" v-for="action in extraActions" :key="action">
              <a @click="emit('action', action)" :class="[
                active ?
                  'bg-gray-100 text-gray-900' :
                  'text-gray-700',
                'block px-4 py-2 text-sm cursor-pointer',
              ]">
                {{ actions.get(action) }}
              </a>
              </MenuItem>
            </div>
          </MenuItems>
        </transition>
      </Menu>
    </div>

    <!-- request tabs (headers, body, etc. )-->
    <div class="flex-0 min-h-16 h-16 max-h-16 px-2">
      <div class="sm:hidden">
        <label for="tabs" class="sr-only">Select a tab</label>
        <select @change="selectTab" id="tabs" name="tabs"
          class="block w-full text-sm rounded-md bg-polar-night-2 text-snow-storm-1 focus:border-frost focus:ring-frost">
          <option v-for="tab in requestTabs" :key="tab.name" :selected="tab.current" :value="tab.name">{{ tab.name }}
          </option>
        </select>
      </div>
      <div class="hidden sm:block">
        <div class="border-b dark:border-polar-night-4">
          <nav class="-mb-px flex space-x-8" aria-label="Tabs">
            <a v-for="tab in requestTabs" @click="switchTab(tab.name)" :key="tab.name" :class="[
              tab.current
                ? 'border-frost text-frost'
                : 'border-transparent text-gray-400 hover:border-gray-500 hover:text-gray-200',
              'group inline-flex cursor-pointer items-center border-b-2 pt-5 pb-1 px-1 text-xs font-medium',
            ]" :aria-current="tab.current ? 'page' : undefined">
              <span>{{ tab.name }}</span>
            </a>
          </nav>
        </div>
      </div>
    </div>

    <div class="flex-auto">
      <div v-if="requestTab === 'Headers'">
        headers here
      </div>
      <div v-else-if="requestTab === 'Body'" class="h-full">
        <CodeEditor :code="request.Body" :readonly="readonly" />
      </div>
    </div>

  </div>
</template>
