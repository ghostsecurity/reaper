<script lang="ts" setup>
import { PropType } from "vue";
import { HttpResponse } from "../../lib/Http.js";
import Code from "./Code.vue";

const props = defineProps({
  response: { type: Object as PropType<HttpResponse>, required: true },
  readonly: { type: Boolean, required: true },
  onchange: { type: Function as PropType<(raw: string) => void>, required: false },
})
</script>

<template>
  <div class="h-full min-h-full w-full max-w-full overflow-x-auto">
    <Code :code="response.Raw" :onchange="onchange" :readonly="readonly" class="h-full" />
  </div>
</template>

<style scoped>
tbody {
  max-width: 100% !important;
}

th {
  font-weight: normal !important;
  max-width: 15vw !important;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

td {
  max-width: 15vw !important;
  text-align: right;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;

}

td input {
  width: 100%;
}
</style>