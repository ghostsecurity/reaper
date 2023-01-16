<script lang="ts" setup>
import { watch, ref, onMounted } from 'vue'

import { HighlightHTTP, HighlightBody } from '../../../wailsjs/go/app/App'

const props = defineProps({
  code: { type: String, required: true },
  readonly: { type: Boolean, required: true },
  http: { type: Boolean, default: false },
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
    HighlightHTTP(sent.value).then((hl: string) => {
      setHighlighted(hl)
    })
  } else {
    // TODO: content type
    HighlightBody(sent.value, 'application/json').then((hl: string) => {
      setHighlighted(hl)
    })
  }
}
</script>

<template>
  <div :class="[
    'overflow-x-auto',
    busy ? 'wrapper plain h-full text-left' : 'wrapper highlighted h-full min-h-full text-left',
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
  overflow-wrap: normal;
  overflow: auto;
  overflow-x: auto !important;
  padding: 0;
  border: none;

}

textarea,
code {
  font-size: 1.05em !important;
  font-family: monospace !important;
  line-height: 1.2em !important;
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
  text-decoration-line: none;
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
