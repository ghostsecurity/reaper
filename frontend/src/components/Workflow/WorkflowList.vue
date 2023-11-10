<script lang="ts" setup>
import { ref, PropType } from 'vue'
import { TrashIcon, PencilSquareIcon } from '@heroicons/vue/20/solid'
import { WorkflowM } from '../../lib/api/workflow'
import ConfirmDialog from '../ConfirmDialog.vue'
import InputBox from '../InputBox.vue'

const props = defineProps({
  flows: { type: Array as PropType<WorkflowM[]>, required: true },
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
    <div v-if="!flows || flows.length === 0" class="w-full pl-8 pt-8 text-center text-frost-3">
      <div class="flex flex-col items-center">
        <h3 class="mt-2 text-sm font-bold">No Workflows</h3>
        <p class="mt-1 text-sm">Create a workflow using the '+' button above.</p>
      </div>
    </div>
    <ul v-else class="block w-full flex-auto text-left">
      <li v-for="flow in flows" :key="flow.id" class="group"
          :class="[flow.id === selected ? 'border-b border-frost-2 bg-polar-night-1': 'border-polar-night-3 hover:border-polar-night-4', 'flex w-full border-t']">
        <a @click="emit('select', flow.id)"
           class="my-1 block grow cursor-pointer truncate pl-2">
          {{ flow.name }}
          <div class="flex items-center justify-between">
            <div class="py-1 text-xs text-polar-night-4">{{ flow.id.substring(0, 8) }}</div>
            <div class="hidden group-hover:flex">
              <button class="shrink pr-2" @click="renaming = flow.id">
                <PencilSquareIcon class="h-5 w-5 text-polar-night-4 hover:text-frost-1" aria-hidden="true"/>
              </button>
              <button class="shrink pr-2" @click="deleting = flow.id">
                <TrashIcon class="h-4 w-4 text-polar-night-4 hover:text-aurora-1" aria-hidden="true"/>
              </button>
            </div>
          </div>
        </a>
      </li>
    </ul>
  </div>
</template>
