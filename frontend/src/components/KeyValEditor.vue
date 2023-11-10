<script lang="ts" setup>
import {PropType, ref, watch} from 'vue'
import {TrashIcon} from '@heroicons/vue/20/solid'
import {KeyValue} from '../lib/api/packaging'
import AutocompleteInput from './Shared/AutocompleteInput.vue'

const props = defineProps({
  data: {
    type: Array as PropType<KeyValue[]>,
    required: true,
  },
  readonly: {
    type: Boolean,
    required: false,
    default: true,
  },
  emptyMessage: {
    type: String,
    required: false,
    default: 'No values found.',
  },
  keySuggestions: {
    type: Array as PropType<string[]>,
    required: false,
    default: () => [],
  },
})

watch(
    () => props.data,
    newVal => {
      copy.value = newVal.map(kv => ({
        key: kv.key,
        value: kv.value,
      }))
    },
)

const emit = defineEmits(['publish'])

const copy = ref(
    props.data.map(kv => ({
      key: kv.key,
      value: kv.value,
    })),
)

function publish() {
  emit('publish', copy.value)
}

function updateKey(index: number, key: string) {
  if (props.readonly) {
    return
  }
  if (index === copy.value.length) {
    copy.value.push({key: '', value: ''})
  }
  copy.value[index].key = key
  publish()
}

function updateValue(index: number, value: string) {
  if (props.readonly) {
    return
  }
  if (index === copy.value.length) {
    copy.value.push({key: '', value: ''})
  }
  copy.value[index].value = value
  publish()
}

function deleteRow(index: number) {
  if (props.readonly) {
    return
  }
  copy.value.splice(index, 1)
  publish()
}

function localParams(c: KeyValue[]): KeyValue[] {
  if (props.readonly) {
    return c
  }
  const extra: KeyValue[] = c.map(kv => ({
    key: kv.key,
    value: kv.value,
  }))
  extra.push({key: '', value: ''})
  return extra
}
</script>

<template>
  <div>
    <div v-if="readonly && Object.keys(data).length === 0" class="relative w-full">{{ emptyMessage }}</div>
    <table class="w-full" v-else>
      <tbody class="bg-white dark:bg-reaper-bg-dark">
      <tr v-for="(row, index) in localParams(copy)" :key="index">
        <td class="min-w-200 w-2/5 border border-snow-storm-3 px-3 py-2 text-left text-xs dark:border-polar-night-4">
          <AutocompleteInput
              @change="updateKey(index, $event)"
              :value="row.key"
              :readonly="readonly"
              :suggestions="keySuggestions"
              :left="true"/>
        </td>
        <td class="border border-snow-storm-3 px-3 py-2 text-left text-xs dark:border-polar-night-4">
          <AutocompleteInput
              @change="updateValue(index, $event)"
              :value="row.value"
              :readonly="readonly"
              :suggestions="[]"/>
        </td>
        <td
            v-if="!readonly && index !== copy.length"
            class="w-[1%] border border-snow-storm-3 px-2 text-center text-xs dark:border-polar-night-4">
          <a v-if="index != copy.length" @click.stop.prevent="deleteRow(index)" class="cursor-pointer text-gray-400">
            <TrashIcon class="m-auto h-4 w-4"/>
          </a>
        </td>
      </tr>
      </tbody>
    </table>
  </div>
</template>
