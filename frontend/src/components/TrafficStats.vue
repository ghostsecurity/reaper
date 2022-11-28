<script lang="ts">
import {defineComponent, PropType} from 'vue'
import {HttpRequest, HttpResponse} from '../lib/Http.js';

export default /*#__PURE__*/ defineComponent({
  props: {
    requests: {type: Array<HttpRequest>},
  },
})
</script>

<template>
  <div class="d-flex">
    <div class="v-col-6">
     lol
    </div>
    <div class="v-col-6">
      some stats
    </div>
  </div>
</template>

<style scoped>
</style>

