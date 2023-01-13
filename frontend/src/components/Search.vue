<script lang="ts" setup>
import { Criteria } from "../lib/Criteria";
import { MagnifyingGlassIcon } from "@heroicons/vue/20/solid";

const props = defineProps({
  query: {
    type: String,
    required: true,
  }
})

const emit = defineEmits(['search'])

function search(evt: Event) {
  emit('search', new Criteria((evt.target as HTMLInputElement).value))
}
</script>

<template>
  <div class="relative w-full text-gray-400 focus-within:text-gray-200 inline-block">
    <span class="absolute p-1">
      <button type="submit" class="p-1 focus:outline-none focus:shadow-outline">
        <MagnifyingGlassIcon class="w-6 h-6" />
      </button>
    </span>
    <input :value="props.query" @keyup="search" @change="search" @search="search" type="search" autocomplete="off"
      autocapitalize="off" spellcheck="false"
      class="w-full py-2 text-sm text-snow-storm-1 bg-polar-night-2 focus:bg-polar-night-4 focus:text-snow-storm-3 focus:outline-none rounded-md pl-10"
      placeholder="Search...">
  </div>
</template>
