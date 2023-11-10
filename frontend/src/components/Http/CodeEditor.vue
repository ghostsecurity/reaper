<script lang="ts" setup>
import { watch, ref, onMounted, PropType } from 'vue'

import { DocumentDuplicateIcon, SparklesIcon, ArrowPathRoundedSquareIcon } from '@heroicons/vue/24/outline'
import Client from '../../lib/api/Client'

const props = defineProps({
  code: { type: String, required: true },
  readonly: { type: Boolean, required: true },
  http: { type: Boolean, default: false },
  mime: { type: String, default: 'text/plain' },
  client: { type: Object as PropType<Client>, required: true },
})

const buffer = ref(props.code)
const busy = ref(true)
const highlighted = ref('')
const sent = ref('')
const textarea = ref()
const pre = ref()

const emit = defineEmits(['change'])

watch(
  () => props.code,
  () => {
    buffer.value = props.code
    const element = textarea.value as HTMLTextAreaElement
    element.value = buffer.value
    updateCode()
  },
)

onMounted(() => {
  updateCode()
})

function setHighlighted(hl: string) {
  highlighted.value = hl
  busy.value = false
}

function syncScroll() {
  const tElement = textarea.value as HTMLTextAreaElement
  const pElement = pre.value as HTMLTextAreaElement
  pElement.scrollTop = tElement.scrollTop
  pElement.scrollLeft = tElement.scrollLeft
}

function updateCode() {
  busy.value = true

  emit('change', buffer.value)

  sent.value = buffer.value
  if (sent.value.length > 0 && sent.value[sent.value.length - 1] === '\n') {
    sent.value += ' '
  }

  if (props.http) {
    props.client.HighlightHTTP(sent.value).then((hl: string) => {
      setHighlighted(hl)
    })
  } else {
    props.client.HighlightBody(sent.value, props.mime).then((hl: string) => {
      setHighlighted(hl)
    })
  }
}

function onKeydown(e: KeyboardEvent) {
  switch (e.key) {
    case 'Tab':
      e.preventDefault()
      indent()
      break
    default:
      break
  }
}

function indent() {
  const element = textarea.value as HTMLTextAreaElement
  const start = element.selectionStart
  const end = element.selectionEnd
  const { value } = element
  const before = value.substring(0, start)
  const after = value.substring(end)
  const insert = '  '
  buffer.value = before + insert + after
  element.value = buffer.value
  element.selectionStart = start + insert.length
  element.selectionEnd = start + insert.length
  updateCode()
}

function copyToClipboard() {
  navigator.clipboard.writeText(buffer.value)
}

function formatCode() {
  props.client.FormatCode(buffer.value, props.mime).then((formatted: string) => {
    buffer.value = formatted
    updateCode()
  })
}

function encodeParams(rawParams: string) {
  const lines = rawParams.split('\n')
  const params = lines.map((line: string) => {
    const [key, value] = line.split(': ').map(str => str.trim())
    return `${encodeURIComponent(key)}=${encodeURIComponent(value)}`
  })
  buffer.value = params.join('&')
}

function decodeParams(encodedParams: string) {
  const params = new URLSearchParams(encodedParams)
  let result = ''
  params.forEach((value, key) => {
    result += `${key}: ${value}\n`
  })
  buffer.value = result.trim()
}

function transformCode() {
  const input = buffer.value.trim()

  // presence of & or = means it's an encoded string
  if (input.includes('&') || input.includes('=')) {
    decodeParams(input)
  } else {
    encodeParams(input)
  }
  updateCode()
}
</script>

<template>
  <div class="absolute flex h-7 w-full border border-polar-night-3">
    <button class="rounded px-1 text-snow-storm-1/70 hover:bg-polar-night-3 hover:text-snow-storm-1"
            @click="copyToClipboard">
      <DocumentDuplicateIcon class="h-6 w-6"
                             aria-hidden="true"/>
    </button>
    <button class="rounded px-1 text-snow-storm-1/70 hover:bg-polar-night-3 hover:text-snow-storm-1"
            @click="formatCode">
      <SparklesIcon class="h-6 w-6" aria-hidden="true"/>
    </button>
    <button class="rounded px-1 text-snow-storm-1/70 hover:bg-polar-night-3 hover:text-snow-storm-1"
            @click="transformCode">
      <ArrowPathRoundedSquareIcon class="h-6 w-6"
                                  aria-hidden="true"/>
    </button>
  </div>
  <div class="h-full w-full border border-polar-night-3 bg-neutral-800/50 p-8 px-1 pb-2">
    <div style="min-height: 100px;"
         :class="[
           'wrapper h-full w-full overflow-x-auto',
           busy ? 'wrapper plain text-left' : 'wrapper highlighted min-h-full text-left',
         ]">
      <pre ref="pre"
           class="h-full min-h-full w-full"
           aria-hidden="true"><code v-html="highlighted"></code></pre>
      <textarea class="h-full w-full"
                :readonly="readonly"
                spellcheck="false"
                ref="textarea"
                @input="updateCode"
                @scroll="syncScroll"
                @keydown="onKeydown"
                v-model="buffer"></textarea>
    </div>
  </div>
</template>

<style scoped>
.wrapper {
  position: relative;
}

.v-theme--dark .wrapper,
.v-theme--ghost .wrapper {
  border-right: 1px solid #444;
}

.v-theme--light .wrapper {
  border-right: 1px solid #ccc;
}

textarea,
pre {
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  white-space: pre;
  overflow-wrap: normal;
  overflow: auto;
  overflow-x: auto !important;
  padding: 0;
  border: none;
}

textarea,
code {
  font-size: 0.9em !important;
  font-family: monospace !important;
  line-height: 1.5em !important;
  tab-size: 2;
  word-spacing: 0;
  letter-spacing: 0;
  font-weight: 400;
  font-style: normal;
  font-variant: normal;
  text-rendering: optimizeLegibility;
  text-transform: none;
  text-align: left;
  text-indent: 0;
  text-shadow: none;
  text-decoration: none;
  text-decoration-style: solid;
  writing-mode: horizontal-tb;
  white-space: pre;
  font-feature-settings: normal;
  overflow-wrap: normal;
}

pre {
  z-index: 9;
  padding: 0 !important;
  margin: 0 !important;
}

.plain pre {
  display: none;
}

textarea {
  box-shadow: none;
  outline: none;
  white-space: pre;
  z-index: 10;
  resize: none;
  caret-color: white;
  background-color: transparent;
}

textarea:focus {
  outline: none !important;
}

.highlighted textarea {
  color: transparent;
}

.plain textarea {
  color: white;
}
</style>
