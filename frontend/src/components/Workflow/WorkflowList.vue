<script lang="ts" setup>
import { ref, PropType } from 'vue'
import { BeakerIcon, TrashIcon, PencilSquareIcon } from '@heroicons/vue/20/solid'
import { workflow } from '../../../wailsjs/go/models'
import ConfirmDialog from '../ConfirmDialog.vue'
import InputBox from '../InputBox.vue'

const props = defineProps({
  flows: { type: Array as PropType<workflow.WorkflowM[]>, required: true },
  selected: { type: String, required: true },
})

const deleting = ref('')
const renaming = ref('')

const emit = defineEmits(['select', 'delete', 'rename'])

function renameWorkflow(name: string) {
  const index = props.flows.findIndex(wf => wf.id === renaming.value)
  renaming.value = ''
  if (index === -1) {
    return
  }
  emit('rename', props.flows[index].id, name)
}
</script>

<template>
  <InputBox v-if="renaming" @cancel="renaming = ''" @confirm="renameWorkflow"
            title="Rename workflow" message="Enter a new name for the workflow"
            :initial="flows.find(wf => wf.id === renaming)?.name ?? ''"
  />
  <ConfirmDialog title="Delete Workflow" cancel="Cancel" confirm="Delete"
                 message="Are you sure you want to delete this workflow?" :show="!!deleting"
                 @confirm="emit('delete', deleting);deleting=''"
                 @cancel="deleting = ''"/>
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
          :class="[flow.id === selected ? 'border-frost-2 bg-polar-night-1 border-b': 'border-polar-night-3 hover:border-polar-night-4', 'w-full block flex border-t']">
        <a @click="emit('select', flow.id)"
           class="block my-1 pl-2 cursor-pointer flex-grow">
          {{ flow.name }}
          <p class="text-polar-night-4">Something</p>
        </a>
        <button class="flex-shrink pr-2" @click="renaming = flow.id">
          <PencilSquareIcon class="h-5 w-5 text-polar-night-4" aria-hidden="true"/>
        </button>
        <button class="flex-shrink pr-2" @click="deleting = flow.id">
          <TrashIcon class="h-5 w-5 text-polar-night-4" aria-hidden="true"/>
        </button>
      </li>
    </ul>
  </div>
</template>
