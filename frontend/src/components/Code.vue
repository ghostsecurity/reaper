

<script lang="ts">
import {defineComponent, PropType} from "vue";

import {EventsEmit, EventsOn} from "../../wailsjs/runtime";

export default /*#__PURE__*/ defineComponent({
  props: {
    code: {type: String, required: true},
    readonly: {type: Boolean, required: true},
    onchange: {type: Function as PropType<(raw: string) => void>, required: false},
  },
  data: function () {
    return {
      buffer: this.code,
      busy: true,
      highlighted: '',
      cancel: ()=> {},
      sent: '',
    }
  },
  watch: {
    code: function () {
      this.buffer = this.code;
      let textarea = (this.$refs['textarea'] as any) as HTMLTextAreaElement
      textarea.value = this.buffer
      this.updateCode()
    }
  },
  beforeMount: function () {

  },
  unmounted() {
    if (typeof this.cancel !== 'undefined') {
      this.cancel()
    }
  },
  mounted() {
    this.updateCode()
  },
  methods: {
    setHighlighted(highlighted: string){
      this.highlighted = highlighted
      this.busy = false;
      if (typeof this.cancel !== 'undefined') {
        this.cancel()
      }
    },
    syncScroll() {
      let textarea = (this.$refs['textarea'] as any) as HTMLTextAreaElement
      let pre = (this.$refs['pre'] as any) as HTMLElement
      pre.scrollTop = textarea.scrollTop;
      pre.scrollLeft = textarea.scrollLeft;
    },
    updateCode() {

      this.busy = true

      if(this.onchange !== undefined){
        this.onchange(this.buffer)
      }

      this.sent = this.buffer
      if (this.sent.length > 0 && this.sent[this.sent.length - 1] === '\n') {
        this.sent += ' '
      }

      this.cancel = EventsOn('OnHighlightResponse', (highlighted: string, original: string) => {
        if (original === this.sent) {
          this.setHighlighted(highlighted)
        }
      })

      EventsEmit("OnHighlightRequest", this.sent)
    }
  }
})
</script>

<template>
  <div v-bind:class="'overflow-x-auto ' + (busy?'h-full text-left wrapper plain':'h-full text-left wrapper highlighted min-h-full')">
    <pre ref="pre" class="h-full min-h-full"  aria-hidden="true"><code v-html="highlighted"></code></pre>
    <textarea :readonly="readonly" spellcheck="false" ref="textarea" @input="updateCode" @scroll="syncScroll" v-model="buffer"></textarea>
  </div>
</template>

<style scoped>
.wrapper {
  position: relative;
}
.v-theme--dark .wrapper, .v-theme--ghost .wrapper{
  border-right: 1px solid #444;
}
.v-theme--light .wrapper{
  border-right: 1px solid #ccc;
}
textarea, pre {
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  white-space: pre; /*nowrap;*/
  overflow-wrap: normal;
  overflow: auto;
  overflow-x: scroll !important;
  padding: 0;
  border: none;
}
textarea, pre, code{
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