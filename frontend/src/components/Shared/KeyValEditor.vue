<script lang="ts" setup>
import { PropType, computed, ref, watch } from 'vue'
import { TrashIcon } from '@heroicons/vue/20/solid'

const props = defineProps({
  data: {
    type: Object as PropType<{ [key: string]: Array<string> }>,
    required: true,
  },
  readonly: {
    type: Boolean,
    required: false,
    default: true,
  },
})

const copy = ref(copyData(props.data))

const emit = defineEmits(['publish'])
const rootKey = ref(0)

watch(
  () => props.data,
  () => {
    copy.value = copyData(props.data)
    rootKey.value += 1
  },
)

function copyData(data: { [key: string]: Array<string> }) {
  const c: { [key: string]: Array<string> } = {}
  Object.keys(data).forEach(key => {
    c[key] = [...data[key].sort()]
  })
  return c
}

interface Flat {
  key: string
  val: string
  id: number
  subId: number
}

const flat = computed(() => {
  const flattened: Array<Flat> = []
  let id = 0
  Object.keys(props.data)
    .sort()
    .forEach(key => {
      let subId = 0
      props.data[key].forEach(val => {
        flattened.push({ key, val, id, subId })
        subId += 1
      })
      id += 1
    })
  flattened.push({ key: '', val: '', id, subId: 0 })
  return flattened
})

function publish() {
  const publication = copyData(copy.value)
  delete publication['']
  emit('publish', publication)
}

function updateKey(oldKey: string, newKey: string, subId: number) {
  if (oldKey === newKey) {
    return
  }
  if (copy.value[newKey] === undefined) {
    copy.value[newKey] = []
  }
  if (copy.value[oldKey] !== undefined) {
    copy.value[newKey].push(copy.value[oldKey].slice(subId, 1)[0])
    copy.value[oldKey].splice(subId, 1)
    if (copy.value[oldKey].length === 0) {
      delete copy.value[oldKey]
    }
  } else {
    copy.value[newKey].push('')
  }
  publish()
}

function updateValue(key: string, value: string, subId: number) {
  if (copy.value[key] !== undefined) {
    copy.value[key][subId] = value
  } else {
    copy.value[key] = [value]
  }
  publish()
}

function deleteRow(key: string, subId: number) {
  copy.value[key].splice(subId, 1)
  publish()
}
</script>

<template>
  <div class="overflow-hidden">
    <div v-if="Object.keys(data).length === 0" class="w-full">no things</div>
    <table class="w-full" v-else>
      <tbody class="bg-white dark:bg-reaper-bg-dark">
        <tr v-for="row in flat" :key="row.id">
          <td class="text-left text-xs text-gray-500">
            <div class="max-w-s">
              <input
                type="text"
                :value="row.key"
                :readonly="readonly"
                @change="updateKey(row.key, ($event.target as HTMLInputElement).value, row.subId)"
                @keypress="updateKey(row.key, ($event.target as HTMLInputElement).value, row.subId)"
                autocomplete="off"
                autocapitalize="off"
                spellcheck="false"
                class="m-0 overflow-x-hidden border-none bg-transparent p-0 text-xs text-polar-night-1 ring-0 focus:border-none focus:text-xs focus:outline-none focus:ring-0 dark:text-snow-storm-1" />
            </div>
          </td>
          <td class="text-xs font-medium text-gray-900">
            <div class="max-w-s">
              <input
                type="text"
                :value="row.val"
                :readonly="readonly"
                @change="updateValue(row.key, ($event.target as HTMLInputElement).value, row.subId)"
                @keypress="updateValue(row.key, ($event.target as HTMLInputElement).value, row.subId)"
                autocomplete="off"
                autocapitalize="off"
                spellcheck="false"
                class="m-0 w-full max-w-sm truncate border-none bg-transparent p-0 text-xs text-polar-night-1 outline-none ring-0 hover:truncate focus:border-none focus:text-xs focus:outline-none focus:ring-0 dark:text-snow-storm-1" />
            </div>
          </td>
          <td class="">
            <a @click.stop.prevent="deleteRow(row.key, row.subId)" class="cursor-pointer text-gray-400">
              <TrashIcon class="h-4 w-4" />
            </a>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
