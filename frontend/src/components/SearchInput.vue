<script lang="ts" setup>
import { MagnifyingGlassIcon } from '@heroicons/vue/20/solid'
import { Criteria } from '../lib/Criteria/Criteria'

const props = defineProps({
  query: {
    type: String,
    required: true,
  },
})

const emit = defineEmits(['search'])

function search(evt: Event) {
  emit('search', new Criteria((evt.target as HTMLInputElement).value))
}
</script>

<template>
  <div class="relative inline-block w-full text-gray-400">
    <span class="absolute p-1">
      <button type="submit" class="focus:shadow-outline p-1 focus:outline-none">
        <MagnifyingGlassIcon class="h-6 w-6" />
      </button>
    </span>
    <input
      :value="props.query"
      @keyup="search"
      @change="search"
      @search="search"
      type="search"
      autocomplete="off"
      autocapitalize="off"
      spellcheck="false"
      class="block w-full rounded-md border-2 bg-snow-storm-1 pl-10 text-polar-night-1 focus:border-frost-4 focus:outline-none focus:ring-frost-4 dark:bg-polar-night-4 dark:text-snow-storm-3 sm:text-sm"
      placeholder="Search..." />
  </div>
</template>
