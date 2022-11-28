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
        return 'text-success'
      } else if (this.code >= 300 && this.code < 400) {
        return 'text-info'
      } else if (this.code >= 400 && this.code < 500) {
        return 'text-warning'
      } else if (this.code >= 500 && this.code < 600) {
        return 'text-error'
      }
      return ''
    }
  },
})
</script>

<template>
  <span v-if="code" v-bind:class="getClass()">{{ code }}</span>
  <span v-if="!code">-</span>
</template>