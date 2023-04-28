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
const handle = ref()
const resizing = ref(false)

const creating = ref(false)
const currentFlow = ref<workflow.WorkflowM | null>(safe.value.workflows.find(
    wf => wf.id === props.selectedWorkflowId) ?? null
)

const emit = defineEmits(['select', 'save', 'run', 'stop'])

onMounted(() => {
  root.value.addEventListener('mousemove', (e: MouseEvent) => {
    if (!resizing.value) {
      return
    }
    if (e.buttons === 0) {
      resizing.value = false
      return
    }

    // Get offset
    const containerOffsetLeft = root.value.offsetLeft

    // Get x-coordinate of pointer relative to container
    const pointerRelativeXpos = e.clientX - containerOffsetLeft

    // Arbitrary minimum width set on box A, otherwise its inner content will collapse to width of 0
    const boxAminWidth = 475

    rightPanel.value.style.width = `${Math.min(
        Math.max(400, root.value.offsetWidth - (pointerRelativeXpos + 10)), // 8px padding + 2px border
        root.value.offsetWidth - boxAminWidth,
    )}px`
    rightPanel.value.style.flexGrow = 0
    rightPanel.value.style.flexShrink = 0
  })
  root.value.addEventListener('mouseup', () => {
    resizing.value = false
  })
})

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
</script>

<template>

  <div ref="root" class="flex h-full overflow-x-hidden">

    <InputBox v-if="creating" title="New Workflow" message="Enter the workflow name." @cancel="creating = false"
              @confirm="addWorkflow($event)"/>
    <div ref="leftPanel" class="box-border flex-auto overflow-y-auto h-full text-right pr-2">
      <button type="button" @click="creating = true"
              class="mb-1 rounded-full bg-frost-4 p-1.5 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
        <PlusIcon class="h-5 w-5" aria-hidden="true"/>
      </button>
      <List :selected="selectedWorkflowId" :flows="safe.workflows" @select="selectWorkflow($event)"
            @delete="deleteWorkflow"/>
    </div>

    <div v-if="currentFlow" @mousedown.prevent="resizing = true" ref="handle"
         class="w-0.5 flex-none cursor-ew-resize bg-gray-500 dark:bg-polar-night-4"></div>

    <div v-if="currentFlow" ref="rightPanel"
         class="mx-2 box-border h-full flex-auto overflow-hidden px-2 w-[60%]">
      <Editor :flow="currentFlow" @save="saveWorkflow($event)" @run="emit('run', $event)"
              @stop="emit('stop', $event)"
              :running="runningWorkflowId===currentFlow.id"/>
    </div>

  </div>

</template>
