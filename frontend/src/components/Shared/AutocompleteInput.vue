<script lang="ts" setup>
import { PropType, ref, watch, computed } from 'vue'

const props = defineProps({
  value: { type: String, required: false, default: '' },
  suggestions: { type: Array as PropType<string[]>, required: false, default: () => [] },
  readonly: { type: Boolean, required: false, default: false },
  left: { type: Boolean, required: false, default: false },
})

const liveValue = ref(props.value)
const showAutocomplete = ref(false)
const input = ref()
const index = ref(0)

watch(
  () => props.value,
  () => {
    liveValue.value = props.value
  },
)

const relevantSuggestions = computed(() => props.suggestions.filter((sugg: string) => {
  if (liveValue.value === '') {
    return true
  }
  return sugg.toLowerCase().startsWith(liveValue.value.toLowerCase()) && sugg !== liveValue.value
}))

const emit = defineEmits(['change'])

function onChange(event: Event) {
  changeValue((event.target as HTMLInputElement).value)
  showAutocomplete.value = true
}
function changeValue(val: string) {
  liveValue.value = val
  emit('change', liveValue.value)
  input.value.focus()
  index.value = 0
}
function hasSuggestions(): boolean {
  return showAutocomplete.value && relevantSuggestions.value.length > 0
}
function onKeydown(e: KeyboardEvent) {
  if (!hasSuggestions()) {
    return
  }
  switch (e.key) {
  case 'Tab':
    e.preventDefault()
    if (index.value < relevantSuggestions.value.length - 1) {
      index.value += 1
    }
    break
  case 'Escape':
    e.preventDefault()
    showAutocomplete.value = false
    break
  case 'ArrowUp':
    e.preventDefault()
    if (index.value > 0) {
      index.value -= 1
    }
    break
  case 'ArrowDown':
    e.preventDefault()
    if (index.value < relevantSuggestions.value.length - 1) {
      index.value += 1
    }
    break
    case 'Enter':
    e.preventDefault()
    changeValue(relevantSuggestions.value[index.value])
    break
  }
}

function highlightedPortion(sugg: string): string {
  return sugg.substring(0, liveValue.value.length)
}
function remainingPortion(sugg: string): string {
  return sugg.substring(liveValue.value.length)
}
</script>

<template>
  <div class="relative">
    <input
      ref="input"
      type="text"
      :value="liveValue"
      :readonly="readonly"
      @focus="showAutocomplete = !readonly"
      @blur="showAutocomplete = false"
      @input="onChange"
      @keydown="onKeydown"
      autocomplete="off"
      autocapitalize="off"
      spellcheck="false"
      class="m-0 w-full truncate border-none bg-transparent p-0 text-xs text-polar-night-1 outline-none ring-0 hover:truncate focus:border-none focus:text-xs focus:outline-none focus:ring-0 dark:text-snow-storm-1" />
    <div
      :class="[
        hasSuggestions() ? '' : 'hidden',
        'absolute z-10 mt-2 w-56 origin-top-left rounded-md bg-white dark:bg-polar-night-4 shadow-lg ring-1',
        'ring-black ring-opacity-5 focus:outline-none overflow-y-auto',
      ]"
      :style="{
        left: left ? '0' : 'auto',
        right: left ? 'auto' : '0',
        'max-height': '240px',
      }"
      >
      <div class="py-1" :key="liveValue">
        <div v-for="(suggestion, i) in relevantSuggestions" :key="suggestion">
          <a
            href="#"
            @focus="input.focus()"
            @mousedown.prevent.stop="changeValue(suggestion)"
            @click.prevent.stop
            @mouseup.prevent.stop
            @mouseover="index = i"
            :class="[
              i === index ? 'bg-aurora-5' : '',
              'flex cursor-pointer px-4 py-2 text-sm text-gray-700 dark:text-snow-storm-1 hover:bg-aurora-5',
            ]">
            <span class="bg-aurora-4/25">{{ highlightedPortion(suggestion) }}</span>
            {{ remainingPortion(suggestion) }}
          </a>
        </div>
      </div>
    </div>
  </div>
</template>
