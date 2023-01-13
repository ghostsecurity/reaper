<script lang="ts" setup>
import { watch, ref, onMounted } from 'vue'

import { HighlightCode } from '../../../wailsjs/go/app/App'

const props = defineProps({
  code: { type: String, required: true },
  readonly: { type: Boolean, required: true },
})

const buffer = ref(props.code)
const busy = ref(true)
const highlighted = ref('')
const sent = ref('')
const textarea = ref()
const pre = ref()

const emit = defineEmits(['change'])

watch(() => props.code, () => {
  buffer.value = props.code
  const element = (textarea.value as HTMLTextAreaElement)
  element.value = buffer.value
  updateCode()
})

onMounted(() => {
  updateCode()
})

function setHighlighted(hl: string) {
  highlighted.value = hl
  busy.value = false
}

function syncScroll() {
  const tElement = (textarea.value as HTMLTextAreaElement)
  const pElement = (pre.value as HTMLTextAreaElement)
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

  HighlightCode(sent.value).then((hl: string) => {
    setHighlighted(hl)
  })
}
</script>

<template>
  <div :class="[
    'overflow-x-auto',
    busy ?
      'h-full text-left wrapper plain' :
      'h-full text-left wrapper highlighted min-h-full',
  ]">
    <pre ref="pre" class="h-full min-h-full" aria-hidden="true"><code v-html="highlighted"></code></pre>
    <textarea :readonly="readonly" spellcheck="false" ref="textarea" @input="updateCode" @scroll="syncScroll"
      v-model="buffer"></textarea>
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
  /*nowrap;*/
  overflow-wrap: normal;
  overflow: auto;
  overflow-x: scroll !important;
  padding: 0;
  border: none;
}

textarea,
pre,
code {
  font-size: 12pt !important;
  font-family: monospace !important;
  line-height: 20pt !important;
  tab-size: 2;
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
