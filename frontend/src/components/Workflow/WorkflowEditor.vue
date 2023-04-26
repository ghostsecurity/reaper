<script lang="ts" setup>
import { PropType, ref, watch } from 'vue'
import { BeakerIcon } from '@heroicons/vue/20/solid'
import { workflow } from '../../../wailsjs/go/models'

const props = defineProps({
  flow: { type: Object as PropType<workflow.WorkflowM>, required: true },
})

const safe = ref<workflow.WorkflowM>(JSON.parse(JSON.stringify(props.flow)))
watch(() => props.flow, flow => {
  if (flow) {
    safe.value = JSON.parse(JSON.stringify(props.flow)) as workflow.WorkflowM
  }
})

const emit = defineEmits(['save'])

function saveWorkflow(f: workflow.WorkflowM) {
  emit('save', f)
}

function linkNodes(nodeA: string, connectorA: string, nodeB: string, connectorB: string) {
  if (!safe.value) {
    return
  }
  safe.value.links.push(new workflow.LinkM({
    from: new workflow.LinkDirectionM({
      node: nodeA,
      connector: connectorA,
    }),
    to: new workflow.LinkDirectionM({
      node: nodeB,
      connector: connectorB,
    }),
  }))
  saveWorkflow(safe.value)
}
</script>

<template>
  <div class="h-full">
    <div v-if="!safe" class="flex flex-col items-center mt-16">
      <BeakerIcon class="h-12 w-12"/>
      <h3 class="mt-2 text-sm font-bold">No Workflow Selected</h3>
      <p class="mt-1 text-sm">Select or create a workflow from the list.</p>
    </div>
    <div v-else class="h-full flex flex-col">
      {{ safe.name }}
      <div class="bg-polar-night-1a border border-polar-night-4 flex-auto my-2">
        diagram goes here
        <a @click="linkNodes('a', 'b', 'c', 'd')">link</a>
      </div>
    </div>
  </div>
</template>
