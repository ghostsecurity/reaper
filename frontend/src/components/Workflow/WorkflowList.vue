<script lang="ts" setup>
import { PropType } from 'vue'
import { BeakerIcon } from '@heroicons/vue/20/solid'
import { workflow } from '../../../wailsjs/go/models'

defineProps({
  flows: { type: Array as PropType<workflow.WorkflowM[]>, required: true },
  selected: { type: String, required: true },
})

const emit = defineEmits(['select'])
</script>

<template>
  <div ref="root" class="flex overflow-x-hidden">
    <div v-if="!flows || flows.length === 0" class="pt-8 pl-8 text-center text-frost-3 w-full">
      <div class="flex flex-col items-center">
        <BeakerIcon class="h-12 w-12"/>
        <h3 class="mt-2 text-sm font-bold">No Workflows</h3>
        <p class="mt-1 text-sm">Create a workflow using the '+' button above.</p>
      </div>
    </div>
    <ul v-else class="w-full block flex-auto text-left">
      <li v-for="flow in flows" :key="flow.id"
          class="w-full block">
        <a @click="emit('select', flow.id)"
           :class="[flow.id === selected ? 'border-frost-2 bg-polar-night-1 border-b': 'border-polar-night-3 hover:border-polar-night-4', 'border-t w-full block my-1 py-1 cursor-pointer']">
          {{ flow.name }}
          <p class="text-polar-night-4">Something</p>
        </a>
      </li>
    </ul>
  </div>
</template>
