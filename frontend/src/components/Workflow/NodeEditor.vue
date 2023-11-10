<script lang="ts" setup>
import {useFileDialog} from '@vueuse/core'
import {computed, PropType, ref, watch} from 'vue'
import {XMarkIcon, FolderIcon} from '@heroicons/vue/20/solid'
import {NodeM} from '../../lib/api/workflow'
import {TransmissionM, Connector, VarStorageM} from '../../lib/api/node'
import {NodeType, ParentType, NodeTypeName, ChildType} from '../../lib/Workflows'
import IDE from '../Http/IDE.vue'
import KeyValEditor from '../KeyValEditor.vue'
import Client from "../../lib/api/Client";
import {HttpRequest, KeyValue} from "../../lib/api/packaging";

interface IMap<T> {
  [index: string]: T;
}

interface Choice {
  key: string
  options: IMap<string>
}

const props = defineProps({
  node: {type: Object as PropType<NodeM>, required: true},
  client: {type: Object as PropType<Client>, required: true},
})

const safe = ref<NodeM>(safeCopy(props.node))
watch(() => props.node, n => {
  if (n) {
    safe.value = safeCopy(n)
  }
})

const emit = defineEmits(['update', 'close'])

function safeCopy(n: NodeM): NodeM {
  const c = JSON.parse(JSON.stringify(n)) as NodeM
  if (!c.name) {
    c.name = NodeTypeName(c.type as NodeType)
  }
  return c
}

function publish() {
  emit('update', safe.value)
}

const staticInputs = computed(() => safe.value?.vars?.inputs?.filter((input: Connector) => {
  switch (input.type) {
    case ParentType.STRING:
      return true
    case ParentType.INT:
      return true
    case ParentType.LIST:
      return true
    case ParentType.CHOICE:
      return true
    case ParentType.REQUEST:
      return !input.linkable
    default:
      return !input.linkable
  }
}) || [])

function updateStringField(field: Connector, event: Event) {
  if (!safe.value?.vars?.static) {
    return
  }
  safe.value.vars.static[field.name].data = (event.target as HTMLInputElement).value
  publish()
}

function updateIntField(field: Connector, event: Event) {
  if (!safe.value?.vars?.static) {
    return
  }
  event.preventDefault()
  const el = (event.target as HTMLInputElement)
  const val = el.value
  const num = val.replace(/[^0-9]/g, '')
  safe.value.vars.static[field.name].data = parseInt(num, 10)
  el.value = num
  publish()
}

function isFieldChildType(field: Connector, type: ChildType) {
  const actual = safe.value.vars?.static[field.name]?.internal
  if (!actual) {
    return false
  }
  return actual === type
}

function updateListType(field: Connector, ev: Event) {
  if (!safe.value || !safe.value.vars?.static[field.name]) {
    return
  }
  const newType = parseInt((ev.target as HTMLSelectElement).value, 10)
  safe.value.vars.static[field.name] = {
    type: ParentType.LIST,
    internal: newType,
    data: createDefaultListData(newType),
  } as TransmissionM
  publish()
}

function createDefaultListData(t: ChildType) {
  switch (t) {
    case ChildType.NUMERIC_RANGE_LIST:
      return [0, 100]
    case ChildType.WORD_LIST:
      return ''
    default:
      return null
  }
}

function updateNumericRangeStart(field: Connector, ev: Event) {
  if (!safe.value || !safe.value.vars?.static[field.name]) {
    return
  }
  const val = parseInt((ev.target as HTMLInputElement).value, 10)
  safe.value.vars.static[field.name].data[0] = val
  publish()
}

function updateNumericRangeEnd(field: Connector, ev: Event) {
  if (!safe.value || !safe.value.vars?.static[field.name]) {
    return
  }
  const val = parseInt((ev.target as HTMLInputElement).value, 10)
  safe.value.vars.static[field.name].data[1] = val
  publish()
}

const requestActions = new Map<string, string>([])

function updateRequestField(field: Connector, req: HttpRequest) {
  if (!safe.value || !safe.value.vars?.static[field.name]) {
    return
  }
  safe.value.vars.static[field.name].data = req
  publish()
}

function updateMapField(field: Connector, kvs: KeyValue[]) {
  if (!safe.value || !safe.value.vars?.static[field.name]) {
    return
  }
  const m = new Map<string, string>([])
  kvs.forEach(kv => {
    m.set(kv.key, kv.value)
  })
  safe.value.vars.static[field.name].data = Object.fromEntries(m)
  publish()
}

function keyValsFromMap(field: Connector): KeyValue[] {
  const data = safe.value.vars?.static[field.name]?.data
  if (data) {
    return Object.entries(data).map(([k, v]) => ({
      key: k,
      value: v,
    } as KeyValue))
  }
  return []
}

function keyValsFromChoice(field: Connector): KeyValue[] {
  const data = (safe.value.vars?.static[field.name]?.data as Choice).options
  if (data) {
    return Object.entries(data).map(([k, v]) => ({
      key: k,
      value: v,
    } as KeyValue)).sort()
  }
  return []
}

function updateBooleanField(field: Connector, ev: Event) {
  if (!safe.value || !safe.value.vars?.static[field.name]) {
    return
  }
  const val = (ev.target as HTMLInputElement).checked
  safe.value.vars.static[field.name].data = val
  publish()
}

function updateChoiceField(field: Connector, ev: Event) {
  if (!safe.value || !safe.value.vars?.static[field.name]) {
    return
  }
  const val = (ev.target as HTMLSelectElement).value
  safe.value.vars.static[field.name].data.key = val
  publish()
}

function getLabel(field: Connector) {
  const label = field.name.replace(/_/g, ' ')
  if (!field.description) {
    return label
  }
  return `${label} (${field.description})`
}

function updateWordList(field: Connector) {
  if (!safe.value || !safe.value.vars) {
    return
  }

  const {files, open, reset, onChange} = useFileDialog()
  open({multiple: false, directory: false})
  onChange((files: FileList) => {
    if (files.length === 0) {
      return
    }
    const reader = new FileReader();
    reader.onload = function () {
      if (!safe.value.vars) {
        safe.value.vars = {} as VarStorageM
      }
      safe.value.vars.static[field.name].data = reader.result
      publish()
    };
    reader.readAsText(files[0])
  })
}

</script>

<template>
  <div
      class="pointer-events-auto relative max-h-full overflow-y-auto rounded border border-polar-night-3 bg-polar-night-1 p-2 text-center">
    <button @click="emit('close')" class="absolute right-1 top-1">
      <XMarkIcon class="h-4 w-4"/>
    </button>
    {{ NodeTypeName(safe.type) }}
    <div class="relative mt-2 text-left">
      <!-- FORM BEGIN -->

      <!-- GLOBAL OPTIONS -->
      <div class="mt-2">
        <div class="sm:col-span-4">
          <label class="block text-sm font-medium capitalize leading-6 text-snow-storm-1">Name</label>
          <div class="mt-1">
            <div
                class="flex rounded-md bg-white/5 ring-1 ring-inset ring-white/10 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-500">
              <input type="text" autocomplete="off" autocapitalize="off" spellcheck="false"
                     v-model="safe.name"
                     @input="publish"
                     class="flex-1 border-0 bg-transparent px-2 py-1.5 text-snow-storm-1 focus:ring-0 sm:text-sm sm:leading-6"/>
            </div>
          </div>
        </div>
      </div>

      <div v-for="field in staticInputs" class="mt-2" :key="field.name">
        <div v-if="field.type === ParentType.STRING" class="sm:col-span-4">
          <label class="block text-sm font-medium capitalize leading-6 text-snow-storm-1">{{
              getLabel(field)
            }}</label>
          <div class="mt-1">
            <div
                class="flex rounded-md bg-white/5 ring-1 ring-inset ring-white/10 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-500">
              <input type="text" autocomplete="off" autocapitalize="off" spellcheck="false"
                     :value="safe.vars?.static[field.name].data"
                     @input="updateStringField(field, $event)"
                     class="flex-1 border-0 bg-transparent px-2 py-1.5 text-snow-storm-1 focus:ring-0 sm:text-sm sm:leading-6"/>
            </div>
          </div>
        </div>
        <div v-else-if="field.type === ParentType.INT" class="sm:col-span-4">
          <label class="block text-sm font-medium capitalize leading-6 text-snow-storm-1">{{
              getLabel(field)
            }}</label>
          <div class="mt-1">
            <div
                class="flex rounded-md bg-white/5 ring-1 ring-inset ring-white/10 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-500">
              <input type="number" autocomplete="off" autocapitalize="off"
                     spellcheck="false"
                     :value="safe.vars?.static[field.name].data"
                     @input="updateIntField(field, $event)"
                     class="flex-1 border-0 bg-transparent px-2 py-1.5 text-snow-storm-1 focus:ring-0 sm:text-sm sm:leading-6"/>
            </div>
          </div>
        </div>
        <div v-else-if="field.type === ParentType.LIST" class="sm:col-span-4">
          <label class="block text-sm font-medium capitalize leading-6 text-snow-storm-1">{{
              getLabel(field)
            }}</label>
          <div class="mt-1">
            <div
                class="flex rounded-md bg-white/5 ring-1 ring-inset ring-white/10 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-500">
              <select
                  @change="updateListType(field, $event)"
                  class="flex-1 border-0 bg-transparent px-2 py-1.5 text-snow-storm-1 focus:ring-0 sm:text-sm sm:leading-6">
                <option :selected="isFieldChildType(field, ChildType.NUMERIC_RANGE_LIST)"
                        :value="ChildType.NUMERIC_RANGE_LIST">
                  Numeric Range
                </option>
                <option :selected="isFieldChildType(field, ChildType.WORD_LIST)"
                        :value="ChildType.WORD_LIST">
                  Wordlist
                </option>
                <option :selected="isFieldChildType(field, ChildType.COMMA_SEP_LIST)"
                        :value="ChildType.COMMA_SEP_LIST">
                  Comma Separated Values
                </option>
              </select>
            </div>
          </div>
          <div class="mt-1">
            <div v-if="isFieldChildType(field, ChildType.NUMERIC_RANGE_LIST)">
              <div class="mt-1">
                <label class="block text-sm font-medium capitalize leading-6 text-snow-storm-1">Start</label>
                <div class="mt-1">
                  <div
                      class="flex rounded-md bg-white/5 ring-1 ring-inset ring-white/10 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-500">
                    <input type="number" autocomplete="off" autocapitalize="off"
                           spellcheck="false"
                           :value="safe.vars?.static[field.name].data[0]"
                           @input="updateNumericRangeStart(field, $event)"
                           class="flex-1 border-0 bg-transparent px-2 py-1.5 text-snow-storm-1 focus:ring-0 sm:text-sm sm:leading-6"/>
                  </div>
                </div>
              </div>
              <div class="mt-1">
                <label class="block text-sm font-medium capitalize leading-6 text-snow-storm-1">End (inclusive)</label>
                <div class="mt-1">
                  <div
                      class="flex rounded-md bg-white/5 ring-1 ring-inset ring-white/10 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-500">
                    <input type="number" autocomplete="off" autocapitalize="off"
                           spellcheck="false"
                           :value="safe.vars?.static[field.name].data[1]"
                           @input="updateNumericRangeEnd(field, $event)"
                           class="flex-1 border-0 bg-transparent px-2 py-1.5 text-snow-storm-1 focus:ring-0 sm:text-sm sm:leading-6"/>
                  </div>
                </div>
              </div>
            </div>
            <div v-else-if="isFieldChildType(field, ChildType.WORD_LIST)"
                 class="mt-2 text-snow-storm-1/80">
              <label class="block text-sm font-medium capitalize leading-6 text-snow-storm-1">Wordlist</label>
              <div class="mt-1 flex rounded border border-polar-night-4 bg-white/5 py-1 text-sm">
                <button @click="updateWordList(field)" class="mx-2 shrink">
                  <FolderIcon class="h-4 w-4 text-snow-storm-1"/>
                </button>
                <div class="grow cursor-pointer pt-1" @click="updateWordList(field)">
                  <p
                      v-if="safe.vars?.static[field.name].data">Wordlist Data</p>
                  <p class="italic" v-else>No file selected</p>
                </div>
              </div>

            </div>
            <div v-else-if="isFieldChildType(field, ChildType.COMMA_SEP_LIST)"
                 class="mt-2 text-snow-storm-1/80">
              <label class="block text-sm font-medium capitalize leading-6 text-snow-storm-1">Comma Separated
                Values</label>
              <div class="mt-1">
                <div
                    class="flex rounded-md bg-white/5 ring-1 ring-inset ring-white/10 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-500">
                  <input type="text" autocomplete="off" autocapitalize="off" spellcheck="false"
                         :value="safe.vars?.static[field.name].data"
                         @input="updateStringField(field, $event)"
                         class="flex-1 border-0 bg-transparent px-2 py-1.5 text-snow-storm-1 focus:ring-0 sm:text-sm sm:leading-6"/>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-else-if="field.type === ParentType.BOOLEAN" class="sm:col-span-4">
          <div class="relative flex items-start">
            <div class="flex h-6 items-center">
              <input :id="field.name" :name="field.name" type="checkbox"
                     :checked="safe.vars?.static[field.name].data"
                     @change="updateBooleanField(field, $event)"
                     class="ml-2 h-4 w-4 rounded bg-polar-night-4 text-frost-1 focus:text-frost-1"/>
            </div>
            <div class="ml-2 text-sm leading-6">
              <label :for="field.name" class="font-medium capitalize text-snow-storm-1">{{
                  getLabel(field)
                }}</label>
            </div>
          </div>
        </div>
        <div v-else-if="field.type === ParentType.CHOICE" class="sm:col-span-4">
          <div class="relative flex items-start py-1">
            <div class="flex-0 pr-2 pt-1 text-sm leading-6">
              <label :for="field.name" class="font-medium capitalize text-snow-storm-1">{{
                  getLabel(field)
                }}</label>
            </div>
            <select :id="field.name" :name="field.name" @change="updateChoiceField(field, $event)"
                    class="flex-1 border-0 bg-polar-night-2 px-2 py-1.5 text-snow-storm-1 focus:ring-0 sm:text-sm sm:leading-6">
              <option
                  v-bind:key="option.key"
                  v-for="option in keyValsFromChoice(field)"
                  :value="option.key"
                  :selected="(safe.vars?.static[field.name].data as Choice).key === option.key">{{ option.value }}
              </option>
            </select>
          </div>
        </div>
        <div v-else-if="field.type === ParentType.REQUEST" class="sm:col-span-4">
          <IDE :client="client" :request="safe.vars?.static[field.name].data" :actions="requestActions"
               :readonly="false" :show-buttons="false"
               @request-update="updateRequestField(field, $event)"/>
        </div>
        <div v-else-if="field.type === ParentType.MAP" class="sm:col-span-4">
          <label class="block text-sm font-medium capitalize leading-6 text-snow-storm-1">{{
              getLabel(field)
            }}</label>
          <KeyValEditor :data="keyValsFromMap(field)"
                        :readonly="false"
                        @publish="updateMapField(field, $event)"/>
        </div>
        <div v-else class="sm:col-span-4">
          <label class="block text-sm font-medium capitalize leading-6 text-snow-storm-1">{{ field.name }}</label>
          <div class="mt-1">
            <i>This value cannot be edited. This is a bug!</i>
          </div>
        </div>
      </div>
      <!-- FORM END -->
    </div>
  </div>
</template>
