<script lang="ts">
import {defineComponent} from "vue";

export default /*#__PURE__*/ defineComponent({
  props: {
    code: {type: Number, required: false},
  },
  methods: {
    getClass() {
      if (this.code === undefined) {
        return ''
      }
      if (this.code >= 200 && this.code < 300) {
        return 'bg-aurora-4'
      } else if (this.code >= 300 && this.code < 400) {
        return 'bg-aurora-3'
      } else if (this.code >= 400 && this.code < 500) {
        return 'bg-aurora-2'
      } else if (this.code >= 500 && this.code < 600) {
        return 'bg-aurora-1'
      }
      return ''
    }
  },
})
</script>

<template>
  <span v-if="code" v-bind:class="'px-2 inline-flex rounded-full '+getClass()">{{ code }}</span>
  <span v-else></span>
</template>