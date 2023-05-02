<script lang="ts" setup>
import {onMounted, PropType, ref, watch} from 'vue'
import {PlusIcon} from '@heroicons/vue/20/solid'
import {workflow, workspace} from '../../../wailsjs/go/models'
import List from './WorkflowList.vue'
import InputBox from '../InputBox.vue'
import Editor from './WorkflowEditor.vue'
import {CreateWorkflow} from "../../../wailsjs/go/backend/App";

const props = defineProps({
  ws: {type: Object as PropType<workspace.Workspace>, required: true},
  selectedWorkflowId: {type: String, required: false, default: ''},
  runningWorkflowId: {type: String, required: false, default: ''},
  statuses: {type: Object as PropType<Map<string, string>>, required: true},
  stdoutLines: {type: Array as PropType<string[]>, required: true},
  stderrLines: {type: Array as PropType<string[]>, required: true},
  activityLines: {type: Array as PropType<string[]>, required: true},
})

const safe = ref<workspace.Workspace>(JSON.parse(JSON.stringify(props.ws)))
watch(() => props.ws, ws => {
  if (ws) {
    safe.value = JSON.parse(JSON.stringify(props.ws)) as workspace.Workspace
  }
})
watch(() => props.selectedWorkflowId, id => {
  if (id) {
    const index = safe.value.workflows.findIndex(wf => wf.id === id)
    if (index === -1) {
      currentFlow.value = null
      return
    }
    currentFlow.value = safe.value.workflows[index]
  } else {
    currentFlow.value = null
  }
})

const root = ref()
const leftPanel = ref()
const rightPanel = ref()

const creating = ref(false)
const currentFlow = ref<workflow.WorkflowM | null>(safe.value.workflows.find(
    wf => wf.id === props.selectedWorkflowId) ?? null
)

const emit = defineEmits(['select', 'save', 'run', 'stop'])

function addWorkflow(name: string) {
  creating.value = false
  CreateWorkflow().then((w) => {
    w.name = name
    safe.value.workflows.push(w)
    saveWorkspace(safe.value)
    emit('select', w.id)
  })
}

function deleteWorkflow(id: string) {
  const index = safe.value.workflows.findIndex(wf => wf.id === id)
  if (index === -1) {
    return
  }
  if (props.selectedWorkflowId === id) {
    selectWorkflow('')
  }
  safe.value.workflows.splice(index, 1)
  saveWorkspace(safe.value)
}

function saveWorkspace(w: workspace.Workspace) {
  emit('save', w)
}

function saveWorkflow(w: workflow.WorkflowM) {
  const index = safe.value.workflows.findIndex(wf => wf.id === w.id)
  if (index === -1) {
    return
  }
  safe.value.workflows[index] = w
  saveWorkspace(safe.value)
}

function selectWorkflow(id: string) {
  emit('select', id)
}

function renameWorkflow(id: string, name: string) {
  const index = safe.value.workflows.findIndex(wf => wf.id === id)
  if (index === -1) {
    return
  }
  safe.value.workflows[index].name = name
  saveWorkspace(safe.value)
  emit('select', id)
}


</script>

<template>

  <div ref="root" class="flex h-full overflow-x-hidden">

    <InputBox v-if="creating" title="New Workflow" message="Enter the workflow name." @cancel="creating = false"
              @confirm="addWorkflow($event)"/>
    <div ref="leftPanel" class="box-border flex-shrink overflow-y-auto w-64 h-full text-right pr-2">
      <button type="button" @click="creating = true"
              class="mb-1 rounded-full bg-frost-4 p-1.5 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
        <PlusIcon class="h-5 w-5" aria-hidden="true"/>
      </button>
      <List :selected="selectedWorkflowId" :flows="safe.workflows" @select="selectWorkflow($event)"
            @delete="deleteWorkflow" @rename="renameWorkflow"/>
    </div>

    <div v-if="currentFlow" ref="rightPanel"
         class="mx-2 box-border h-full flex-grow overflow-hidden px-2 w-[60%]">
      <Editor :flow="currentFlow" @save="saveWorkflow($event)" @run="emit('run', $event)"
              @stop="emit('stop', $event)"
              :running="runningWorkflowId===currentFlow.id"
              :statuses="statuses" :stdout-lines="stdoutLines" :stderr-lines="stderrLines"
              :activity-lines="activityLines"
      />
    </div>

  </div>

</template>
