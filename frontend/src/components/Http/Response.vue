<script lang="ts" setup>
import { PropType } from "vue";
import { HttpResponse } from "../../lib/Http.js";
import Code from "./Code.vue";

const props = defineProps({
  response: { type: Object as PropType<HttpResponse>, required: true },
  readonly: { type: Boolean, required: true },
})

const emit = defineEmits(['change'])

</script>

<template>
  <div class="h-full min-h-full w-full max-w-full overflow-x-auto">
    <Code :code="response.Raw" @change="emit('change', $event)" :readonly="readonly" class="h-full" />
  </div>
</template>
