<script lang="ts" setup>
import { PropType } from "vue";
import { HttpRequest } from "../../lib/Http.js";
import Code from "./Code.vue";

const props = defineProps({
  request: { type: Object as PropType<HttpRequest>, required: true },
  readonly: { type: Boolean, required: true },
})

const emit = defineEmits(['change'])

</script>

<template>
  <div class="h-full min-h-full w-full max-w-full overflow-x-auto">
    <Code :code="request.Raw" @change="emit('change', $event)" :readonly="readonly" class="h-full" />
  </div>
</template>
