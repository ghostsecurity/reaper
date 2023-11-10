<script lang="ts" setup>
import { PropType, ref, watch } from 'vue'
import { PlusIcon, FolderIcon, BeakerIcon } from '@heroicons/vue/20/solid'
import { saveAs } from 'file-saver'
import { useFileDialog } from '@vueuse/core/index'
import { WorkflowM } from '../../lib/api/workflow'
import { Workspace } from '../../lib/api/workspace'
import List from './WorkflowList.vue'
import InputBox from '../InputBox.vue'
import Editor from './WorkflowEditor.vue'
import Client from '../../lib/api/Client'
import { VarStorageM } from '../../lib/api/node'

const props = defineProps({
  ws: { type: Object as PropType<Workspace>, required: true },
  selectedWorkflowId: { type: String, required: false, default: '' },
  runningWorkflowId: { type: String, required: false, default: '' },
  statuses: { type: Object as PropType<Map<string, string>>, required: true },
  stdoutLines: { type: Array as PropType<string[]>, required: true },
  stderrLines: { type: Array as PropType<string[]>, required: true },
  activityLines: { type: Array as PropType<string[]>, required: true },
  client: { type: Object as PropType<Client>, required: true },
})

const safe = ref<Workspace>(JSON.parse(JSON.stringify(props.ws)))
watch(() => props.ws, ws => {
  if (ws) {
    safe.value = JSON.parse(JSON.stringify(props.ws)) as Workspace
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
const currentFlow = ref<WorkflowM | null>(safe.value.workflows.find(
  wf => wf.id === props.selectedWorkflowId,
) ?? null)

const emit = defineEmits(['select', 'save', 'run', 'stop', 'clean'])

function addWorkflow(name: string) {
  creating.value = false
  props.client.CreateWorkflow().then((w: WorkflowM | null) => {
    if (!w) {
      return
    }
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

function saveWorkspace(w: Workspace) {
  emit('save', w)
}

function saveWorkflow(w: WorkflowM) {
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

function importWorkflow() {
  const { files, open, reset, onChange } = useFileDialog()
  open({ multiple: false, directory: false, accept: '.atk' })
  onChange((files: FileList) => {
    if (files.length === 0) {
      return
    }
    const reader = new FileReader()
    reader.onload = function () {
      const w: WorkflowM = JSON.parse(reader.result as string)
      safe.value.workflows.push(w)
      saveWorkspace(safe.value)
      emit('select', w.id)
    }
    reader.readAsText(files[0])
  })
}

function exportWorkflow(id: string) {
  const wf = safe.value.workflows.find(mwf => mwf.id === id)
  if (!wf) {
    return
  }
  const data = new TextEncoder().encode(JSON.stringify(wf))
  saveAs(new Blob([data.buffer]), `${id}.atk`)
  emit('select', id)
}

</script>

<template>

  <div ref="root" class="flex h-full overflow-hidden">

    <InputBox v-if="creating" title="New Workflow" message="Enter the workflow name." @cancel="creating = false"
              @confirm="addWorkflow($event)"/>
    <div ref="leftPanel" class="box-border h-full w-48 shrink overflow-y-auto pr-2 text-right">
      <button type="button" @click="creating = true"
              class="mb-1 rounded-full bg-frost-4 p-1.5 text-white shadow-sm hover:bg-frost-1 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
        <PlusIcon class="h-5 w-5" aria-hidden="true"/>
      </button>
      <button type="button" @click="importWorkflow"
              class="mb-1 ml-1 rounded-full bg-frost-4 p-1.5 text-white shadow-sm hover:bg-frost-1 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
        <FolderIcon class="h-5 w-5" aria-hidden="true"/>
      </button>
      <List :selected="selectedWorkflowId" :flows="safe.workflows" @select="selectWorkflow($event)"
            @delete="deleteWorkflow" @rename="renameWorkflow"/>
    </div>

    <div ref="rightPanel"
         class="mx-2 box-border h-full w-[60%] grow px-2">
      <Editor v-if="currentFlow" :client="client" :flow="currentFlow" @save="saveWorkflow($event)"
              @run="emit('run', $event)"
              @stop="emit('stop', $event)"
              :running="runningWorkflowId===currentFlow.id"
              :statuses="statuses" :stdout-lines="stdoutLines" :stderr-lines="stderrLines"
              :activity-lines="activityLines"
              @clean="emit('clean', $event)"
              @export="exportWorkflow($event)"
      />
      <div v-else class="mt-16 flex flex-col items-center text-white/20">
        <BeakerIcon class="h-12 w-12"/>
        <h3 class="mt-2 text-sm font-bold">No Workflow Selected</h3>
        <p class="mt-1 text-sm">Select or create a workflow from the list.</p>
      </div>
    </div>

  </div>

</template>
