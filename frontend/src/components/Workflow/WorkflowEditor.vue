<script lang="ts" setup>
import { onActivated, onBeforeUpdate, onMounted, onUpdated, PropType, ref, watch } from 'vue'
import {
  PlayIcon,
  PlusIcon,
  Square2StackIcon,
  StopIcon,
  TrashIcon,
  CheckCircleIcon,
  PauseCircleIcon,
  ExclamationCircleIcon, BarsArrowDownIcon, BoltIcon, EyeSlashIcon,
  XCircleIcon,
  ArrowUpOnSquareIcon,
} from '@heroicons/vue/20/solid'
import { uuid } from 'vue-uuid'
import { workflow } from '../../../wailsjs/go/models'
import { CreateNode } from '../../../wailsjs/go/backend/App'
import NodeEditor from './NodeEditor.vue'
import { NodeType, NodeTypeName } from '../../lib/Workflows'
import Spinner from '../Shared/LoadingSpinner.vue'
import ScrollingOutput from '../Shared/ScrollingOutput.vue'

const props = defineProps({
  flow: { type: Object as PropType<workflow.WorkflowM>, required: true },
  running: { type: Boolean, required: false, default: false },
  statuses: { type: Object as PropType<Map<string, string>>, required: true },
  stdoutLines: { type: Array as PropType<string[]>, required: true },
  stderrLines: { type: Array as PropType<string[]>, required: true },
  activityLines: { type: Array as PropType<string[]>, required: true },
})

const availableNodeTypes = ref(<NodeType[]>[
  NodeType.REQUEST,
  NodeType.SENDER,
  NodeType.STATUS_FILTER,
  NodeType.FUZZER,
  NodeType.OUTPUT,
  NodeType.VARIABLES,
  NodeType.DELAY,
])

const linkColour = '#8FBCBB'
const linkStrokeWidth = '2'
const editingNode = ref(<workflow.NodeM | null>null)
const menuMode = ref('')
// used to prevent click + drag events overlapping
const mouseMoved = ref(false)

const safe = ref<workflow.WorkflowM>(JSON.parse(JSON.stringify(props.flow)))
let initial = true
watch(() => props.flow, flow => {
  if (flow) {
    safe.value = JSON.parse(JSON.stringify(props.flow)) as workflow.WorkflowM
    redraw()
  }
})

onMounted(redraw)
onActivated(redraw)

onUpdated(() => {
  if (initial) {
    redraw()
    initial = false
  }
})

const svg = ref(<HTMLElement | null>null)
const connector = ref(<HTMLElement | null>null)
const paths = ref(<string[]>[])

const emit = defineEmits(['save', 'run', 'stop', 'clean', 'export'])

function saveWorkflow(f: workflow.WorkflowM) {
  emit('save', f)
  redraw()
}

const tabs = ref([
  { name: 'Stdout', id: 'stdout', icon: BarsArrowDownIcon, current: true },
  { name: 'Stderr', id: 'stderr', icon: BarsArrowDownIcon, current: false },
  { name: 'Activity', id: 'activity', icon: BoltIcon, current: false },
])

function selectedTab(): string {
  return tabs.value.find(tab => tab.current)?.id || ''
}

function switchTab(id: string) {
  tabs.value = tabs.value.map(tab => {
    const updatedTab = tab
    updatedTab.current = updatedTab.id === id
    return updatedTab
  })
}

function selectTab(e: Event) {
  switchTab((e.target as HTMLSelectElement).value)
}

const curveOffset = ref(100)

function redraw() {
  if (!safe.value) {
    return
  }
  if (canvas.value && svg.value) {
    svg.value.style.height = `${canvas.value.clientHeight + canvas.value.scrollTop}px`
    svg.value.style.width = `${canvas.value.clientWidth + canvas.value.scrollLeft}px`
  }

  const newPaths: string[] = []
  let ok = false
  safe.value.links.forEach(link => {
    const fromNode = movers.value.get(link.from.node) as HTMLElement
    const toNode = movers.value.get(link.to.node) as HTMLElement
    if (!fromNode || !toNode) {
      return
    }

    const fromConn = fromNode.querySelector(`.connector.output[data-connector="${link.from.connector}"]`) as HTMLElement
    const toConn = toNode.querySelector(`.connector.input[data-connector="${link.to.connector}"]`) as HTMLElement
    if (!fromConn || !toConn) {
      return
    }

    ok = true

    const posA = {
      x: (fromConn.offsetLeft + fromConn.offsetWidth) + fromNode.offsetLeft,
      y: (fromConn.offsetTop + fromConn.offsetHeight / 2) + fromNode.offsetTop,
    }
    const posB = {
      x: toConn.offsetLeft + toNode.offsetLeft,
      y: (toConn.offsetTop + toConn.offsetHeight / 2) + toNode.offsetTop,
    }
    const path = `M${
      posA.x},${posA.y} `
        + `C${
          posA.x + curveOffset.value},${posA.y} ${
          posB.x - curveOffset.value},${posB.y} ${
          posB.x},${posB.y}`
    newPaths.push(path)
  })
  if (!ok) {
    initial = true
    return
  }
  paths.value = newPaths
}

function getRawPosition(id: string, i: number) {
  const def = {
    x: 50 + (((i % 4) * 300)),
    y: 50 + (((i - (i % 4)) / 4) * 100),
  }
  if (!safe.value || !safe.value.positioning) {
    return def
  }
  const pos = safe.value.positioning[id]
  if (!pos) {
    return def
  }
  return {
    x: pos.x,
    y: pos.y,
  }
}

function getPosition(id: string, i: number) {
  const pos = getRawPosition(id, i)
  return {
    left: `${pos.x}px`,
    top: `${pos.y}px`,
  }
}

function setPosition(id: string, x: number, y: number) {
  if (!safe.value) {
    return
  }
  if (!safe.value.positioning) {
    safe.value.positioning = {}
  }
  safe.value.positioning[id] = new workflow.Position({
    x,
    y,
  })
  saveWorkflow(safe.value)
}

function addNode(t: number) {
  if (!safe.value) {
    return
  }
  CreateNode(t).then(n => {
    safe.value.nodes.push(n)
    editingNode.value = n
    menuMode.value = ''
    saveWorkflow(safe.value)
  })
}

const canvas = ref(<HTMLDivElement | null>null)
let dragId = ''
let offsetX = 0
let offsetY = 0

function getOffsetFrom(el: HTMLElement, className: string) {
  let x = 0
  let y = 0
  while (el && (!el.className || !el.className.includes || !el.className.includes(className))) {
    if (typeof el.offsetLeft === 'number') {
      x += el.offsetLeft
      y += el.offsetTop
    }
    if (!el.offsetParent) {
      break
    }
    el = el.offsetParent as HTMLElement
  }
  return {
    x,
    y,
  }
}

function dragStart(id: string, ev: MouseEvent) {
  mouseMoved.value = false
  if (!canvas.value) {
    return
  }
  dragId = id
  const el = ev.target as HTMLDivElement
  const offset = getOffsetFrom(el, 'mover')
  offsetX = ev.offsetX + offset.x
  offsetY = ev.offsetY + offset.y
}

function dragEnd() {
  if (!canvas.value) {
    return
  }
  dragId = ''
  if (currentLink) {
    endLinking()
  }
}

function drag(ev: MouseEvent) {
  mouseMoved.value = true
  if (dragId === '' || !canvas.value) {
    if (currentLink) {
      moveLink(ev)
    }
    return
  }
  ev.preventDefault()
  const off = getOffsetFrom(ev.target as HTMLElement, 'canvas')
  const x = ev.offsetX + off.x - canvas.value.offsetLeft
  const y = ev.offsetY + off.y - canvas.value.offsetTop
  const el = movers.value.get(dragId) as HTMLElement
  if (!el) {
    return
  }
  el.style.left = `${x - offsetX}px`
  el.style.top = `${y - offsetY}px`
  setPosition(dragId, x - offsetX, y - offsetY)
}

// const movers = ref([] as any)
const movers = ref(new Map<string, Element | HTMLElement | null>())

// Make sure to reset the refs before each update.
onBeforeUpdate(() => {
  movers.value = new Map<string, Element | HTMLElement | null>()
})

let currentLink: workflow.LinkM | null = null
let currentLinkNode = ''
let currentLinkConnector = ''

let linkSearchMode = ''

function startLinkFromInput(nodeId: string, conn: string) {
  if (!safe.value) {
    return
  }
  linkSearchMode = 'output'
  currentLinkNode = nodeId
  currentLinkConnector = conn
  const existing = safe.value.links.find(l => l.to.node === nodeId && l.to.connector === conn)
  unlinkAnyFromInput(nodeId, conn)
  if (existing) {
    linkSearchMode = 'input'
    currentLinkNode = existing.from.node
    currentLinkConnector = existing.from.connector
    currentLink = new workflow.LinkM({
      from: new workflow.LinkDirectionM({
        node: existing.from.node,
        connector: existing.from.connector,
      }),
      to: new workflow.LinkDirectionM({
        node: '',
        connector: '',
      }),
    })
  } else {
    currentLink = new workflow.LinkM({
      from: new workflow.LinkDirectionM({
        node: '',
        connector: '',
      }),
      to: new workflow.LinkDirectionM({
        node: nodeId,
        connector: conn,
      }),
    })
  }
}

function startLinkFromOutput(nodeId: string, conn: string) {
  if (!safe.value) {
    return
  }
  linkSearchMode = 'input'
  currentLinkNode = nodeId
  currentLinkConnector = conn
  const existing = safe.value.links.find(l => l.from.node === nodeId && l.from.connector === conn)
  unlinkAnyFromOutput(nodeId, conn)
  if (existing) {
    linkSearchMode = 'output'
    currentLinkNode = existing.to.node
    currentLinkConnector = existing.to.connector
    currentLink = new workflow.LinkM({
      from: new workflow.LinkDirectionM({
        node: '',
        connector: '',
      }),
      to: new workflow.LinkDirectionM({
        node: existing.to.node,
        connector: existing.to.connector,
      }),
    })
  } else {
    currentLink = new workflow.LinkM({
      from: new workflow.LinkDirectionM({
        node: nodeId,
        connector: conn,
      }),
      to: new workflow.LinkDirectionM({
        node: '',
        connector: '',
      }),
    })
  }
}

// x and y are offsets into the canvas
function findConnectorAt(x: number, y: number, className: string): HTMLElement | null {
  if (!canvas.value) {
    return null
  }
  const magnetism = 25
  const connectors = canvas.value.querySelectorAll(`.connector.${className}`)
  for (let i = 0; i < connectors.length; i += 1) {
    const el = connectors[i] as HTMLElement
    if (!el.parentElement || !canvas.value) {
      continue
    }
    const offset = getOffsetFrom(el, 'canvas')
    const rect = el.parentElement.getBoundingClientRect()
    const area = {
      x1: offset.x - magnetism,
      y1: offset.y,
      x2: offset.x + rect.width + magnetism,
      y2: offset.y + rect.height,
    }
    if (x >= area.x1 && x <= area.x2 && y >= area.y1 && y <= area.y2) {
      return el
    }
  }
  return null
}

function moveLink(ev: MouseEvent) {
  if (!currentLink || !connector.value || !canvas.value) {
    return
  }

  let selector = '.connector.'
  if (linkSearchMode === 'input') {
    selector += 'output'
  } else {
    selector += 'input'
  }

  const nodeA = movers.value.get(currentLinkNode) as HTMLElement
  const connectorA = nodeA.querySelector(`${selector}[data-connector="${currentLinkConnector}"]`) as HTMLElement

  let posA = {
    x: (connectorA.offsetLeft + connectorA.offsetWidth) + nodeA.offsetLeft - canvas.value.offsetLeft,
    y: (connectorA.offsetTop + connectorA.offsetHeight / 2) + nodeA.offsetTop - canvas.value.offsetTop,
  }

  const off = getOffsetFrom(ev.target as HTMLElement, 'canvas')
  let posB = {
    x: ev.offsetX + off.x,
    y: ev.offsetY + off.y,
  }

  const connectorB = findConnectorAt(posB.x, posB.y, linkSearchMode)
  if (connectorB) {
    let nodeB = connectorB.parentElement as HTMLElement
    while (nodeB && !nodeB.className.includes('mover')) {
      nodeB = nodeB.parentElement as HTMLElement
    }
    posB = {
      x: connectorB.offsetLeft + nodeB.offsetLeft,
      y: connectorB.offsetTop + nodeB.offsetTop + (connectorB.offsetHeight / 2),
    }
    const node = connectorB.getAttribute('data-node')
    const conn = connectorB.getAttribute('data-connector')
    if (node && conn) {
      if (linkSearchMode === 'input') {
        currentLink.to.node = node
        currentLink.to.connector = conn
      } else {
        currentLink.from.node = node
        currentLink.from.connector = conn
      }
      if (canLink(currentLink)) {
        connector.value.setAttribute('stroke', 'green')
      } else {
        connector.value.setAttribute('stroke', 'red')
      }
    }
  } else if (linkSearchMode === 'input') {
    currentLink.to.node = ''
    currentLink.to.connector = ''
  } else {
    currentLink.from.node = ''
    currentLink.from.connector = ''
  }

  if (linkSearchMode !== 'input') {
    const buf = posA
    posA = posB
    posB = buf
  }

  const dStr = `M${
    posA.x},${posA.y} `
      + `C${
        posA.x + curveOffset.value},${posA.y} ${
        posB.x - curveOffset.value},${posB.y} ${
        posB.x},${posB.y}`
  connector.value.setAttribute('d', dStr)
  if (!connectorB) {
    connector.value.setAttribute('stroke', 'blue')
  }
}

function canLink(link: workflow.LinkM): boolean {
  if (!safe.value) {
    return false
  }
  if (link.from.node === link.to.node) {
    return false
  }
  const fromNode = safe.value.nodes.find(n => n.id === link.from.node)
  const toNode = safe.value.nodes.find(n => n.id === link.to.node)
  if (!fromNode || !toNode || !fromNode.vars || !toNode.vars) {
    return false
  }
  const output = fromNode.vars.outputs.find(o => o.name === link.from.connector)
  const input = toNode.vars.inputs.find(i => i.name === link.to.connector)
  if (!output || !input) {
    return false
  }
  if ((input.type & output.type) !== input.type) {
    return false
  }
  return true
}

function endLinking() {
  if (!currentLink || !connector.value || !canvas.value || !safe.value) {
    return
  }

  if (canLink(currentLink)) {
    safe.value.links.push(currentLink)
    saveWorkflow(safe.value)
  }

  connector.value.setAttribute('d', '')
  currentLink = null
  currentLinkConnector = ''
  currentLinkNode = ''
}

function unlinkAnyFromOutput(node: string, conn: string) {
  if (!safe.value) {
    return
  }
  const index = safe.value.links.findIndex(l => l.from.node === node && (l.from.connector === conn || conn === ''))
  if (index >= 0) {
    safe.value.links.splice(index, 1)
    saveWorkflow(safe.value)
  }
}

function unlinkAnyFromInput(node: string, conn: string) {
  if (!safe.value) {
    return
  }
  const index = safe.value.links.findIndex(l => l.to.node === node && (l.to.connector === conn || conn === ''))
  if (index >= 0) {
    safe.value.links.splice(index, 1)
    saveWorkflow(safe.value)
  }
}

function editNode(id: string) {
  editingNode.value = safe.value.nodes.find(n => n.id === id) as workflow.NodeM
  menuMode.value = ''
}

function duplicateNode(id: string) {
  if (!safe.value) {
    return
  }
  const node = safe.value.nodes.find(n => n.id === id)
  if (!node) {
    return
  }
  const newNode = JSON.parse(JSON.stringify(node)) as workflow.NodeM
  newNode.id = uuid.v4()
  safe.value.nodes.push(newNode)
  editingNode.value = newNode
  saveWorkflow(safe.value)
}

function deleteNode(id: string) {
  if (!safe.value) {
    return
  }
  if (editingNode.value && editingNode.value.id === id) {
    editingNode.value = null
  }
  const index = safe.value.nodes.findIndex(n => n.id === id)
  if (index < 0) {
    return
  }
  safe.value.nodes.splice(index, 1)
  unlinkAnyFromInput(id, '')
  unlinkAnyFromOutput(id, '')
  saveWorkflow(safe.value)
}

function updateNode(n: workflow.NodeM) {
  if (!safe.value) {
    return
  }
  const index = safe.value.nodes.findIndex(node => node.id === n.id)
  if (index < 0) {
    return
  }
  safe.value.nodes[index] = n
  saveWorkflow(safe.value)
}

function setMenu(n: string) {
  editingNode.value = null
  if (menuMode.value === n) {
    menuMode.value = ''
    return
  }
  menuMode.value = n
}

function getStatusClass(id: string): string {
  switch (props.statuses.get(id)) {
    case 'pending':
      return 'border-grey-500'
    case 'running':
      return 'border-aurora-3'
    case 'error':
      return 'border-aurora-1'
    case 'aborted':
      return 'border-aurora-5'
    case 'success':
      return 'border-aurora-4'
    default:
      return 'bg-polar-night-4'
  }
}

function trackMover(id: string, el: any) {
  if (!el || !(el instanceof HTMLElement)) {
    return
  }
  movers.value.set(id, el)
}

</script>

<template>
  <div class="w-full h-full box-border relative">
    <div class="w-full h-full flex flex-col overflow-hidden">
      <div class="canvas border border-polar-night-4 stripy grow-[3] relative w-full h-full overflow-auto"
           @mousemove="drag" @mouseup="dragEnd"
           @scroll="redraw"
           @resize="redraw"
           ref="canvas"
      >
        <svg ref="svg" xmlns="http://www.w3.org/2000/svg" class="absolute w-full h-full"
             @click="setMenu('')"
        >
          <path ref="connector" fill="none" stroke="" :stroke-width="linkStrokeWidth"/>
          <path v-for="path in paths" :d="path" fill="none" :stroke="linkColour" :stroke-width="linkStrokeWidth"
                :key="path"/>
        </svg>
        <div v-for="(node, i) in safe.nodes" :key="node.id"
             class="mover absolute"
             :ref="(el) => trackMover(node.id, el)"
             :style="getPosition(node.id, i)"
             draggable="true" @mousedown="dragStart(node.id, $event)"
        >
          <div class="flex items-center">
            <div class="flex-shrink">
              <div @mousedown.prevent.stop="startLinkFromInput(node.id, input.name)"
                   v-for="input in node.vars?.inputs?.filter((inp) => inp.linkable)"
                   :key="input.name"
                   class="flex items-center pr-2 py-0 my-0 leading-4 group">
                <div class="flex-grow mr-2 opacity-60 text-right">
                  {{ input.name }}
                </div>
                <div :data-node="node.id" :data-connector="input.name"
                     class="connector input flex-shrink w-2 h-2 border border-frost rounded-full group-hover:border-4">
                </div>
              </div>
            </div>
            <div class="flex-grow group">
              <div
                  @click="!mouseMoved && editNode(node.id)"
                  :class="[editingNode && editingNode.id == node.id ? 'border-snow-storm-1 bg-polar-night-2' : (node.type == NodeType.START ? 'border-aurora-4 bg-aurora-4/25' : 'bg-polar-night-2 border-frost-1/50'),  'px-2 py-4 border rounded cursor-move relative',  getStatusClass(node.id)]"
                  style="min-width:90px">
                {{ node.name }}
                <div v-if="node.type != NodeType.START && NodeTypeName(node.type) !== node.name"
                     class="absolute ml-1 left-0 top-0 py-0.5 p-2 text-xs bg-polar-night-2 rounded-md text-gray-400 italic"
                     style="margin-top: -0.5rem;">
                  {{ NodeTypeName(node.type) }}
                </div>
                <div>
                  <div class="absolute left-1 bottom-1">
                    <Spinner v-if="running && props.statuses.get(node.id) === 'running'"/>
                    <CheckCircleIcon v-else-if="props.statuses.get(node.id) === 'success'"
                                     class="w-4 h-4 mr-2 text-aurora-4"/>
                    <XCircleIcon v-else-if="props.statuses.get(node.id) === 'error'"
                                 class="w-4 h-4 mr-2 text-aurora-1"/>
                    <PauseCircleIcon v-else-if="props.statuses.get(node.id) === 'pending'"
                                     class="w-4 h-4 mr-2 text-aurora-3"/>
                    <ExclamationCircleIcon v-else-if="props.statuses.get(node.id) === 'aborted'"
                                           class="w-4 h-4 mr-2 text-gray-400"/>
                  </div>
                </div>
                <div v-if="!node.readonly && !dragId" @mousedown.prevent.stop
                     class="absolute top-0 right-0 invisible hover:visible group-hover:visible text-snow-storm-1">
                  <button @click.prevent.stop="duplicateNode(node.id)" class="group/btn px-0.5">
                    <Square2StackIcon class="h-4 w-4 group-hover/btn:text-frost-1"/>
                  </button>
                  <button @click.prevent.stop="deleteNode(node.id)" class="group/btn px-0.5">
                    <TrashIcon class="h-4 w-4 group-hover/btn:text-aurora-2"/>
                  </button>
                </div>
              </div>
            </div>
            <div class="flex-shrink">
              <div @mousedown.prevent.stop="startLinkFromOutput(node.id, output.name)"
                   v-for="output in node.vars?.outputs"
                   :key="output.name"
                   class="flex items-center pl-2 py-0 my-0 leading-4 group">
                <div :data-node="node.id" :data-connector="output.name"
                     class="connector output flex-shrink w-2 h-2 border border-frost rounded-full group-hover:border-4">
                </div>
                <div class="flex-grow ml-2 opacity-60 text-right">
                  {{ output.name }}
                </div>

              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="absolute right-5 top-5 text-right pointer-events-none">
        <div class="pointer-events-auto">
          <button type="button" @click="setMenu('add')"
                  :class="[menuMode==='add'?'bg-frost-1':'bg-frost-4', 'mb-1 rounded-full  p-1.5 text-white shadow-sm hover:bg-frost-1 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 mx-0.5']">
            <PlusIcon class="h-5 w-5" aria-hidden="true"/>
          </button>
        </div>
        <div v-if="editingNode" class="mt-1 h-full overflow-y-hidden pointer-events-none">
          <NodeEditor :node="editingNode" @update="updateNode" @close="editingNode = null; menuMode = ''"/>
        </div>
        <div v-else-if="menuMode==='add'"
             class="mt-1 border rounded border-snow-storm-1 relative bg-polar-night-2 pointer-events-auto">
          <button v-for="t in availableNodeTypes" :key="t" @click="addNode(t);setMenu('')"
                  class="w-full border border-polar-night-4 py-1 bg-polar-night-2 hover:bg-polar-night-4">
            {{ NodeTypeName(t) }}
          </button>
        </div>
      </div>
      <div
          class="flex flex-col box-border w-full box-border relative border-x border-b border-polar-night-3 text-left h-[25%]"
          style="flex: 0 0 auto">
        <div class="flex-shrink p-1">
          <button :disabled="running"
                  :class="['bg-polar-night-4 rounded-md p-2 mr-0.5', !running ? 'text-snow-storm-1 hover:text-frost-1' : 'text-snow-storm-1/20']"
                  @click="emit('run', safe.id)">
            <PlayIcon class="h-5 w-5" aria-hidden="true"/>
          </button>
          <button :disabled="!running"
                  :class="['bg-polar-night-4 rounded-md p-2 mx-0.5', running ? 'text-snow-storm-1 hover:text-frost-1' : 'text-snow-storm-1/20']"
                  @click="emit('stop', safe.id)">
            <StopIcon class="h-5 w-5" aria-hidden="true"/>
          </button>
          <button :disabled="running"
                  :class="['bg-polar-night-4 rounded-md p-2 mx-0.5', !running ? 'text-snow-storm-1 hover:text-frost-1' : 'text-snow-storm-1/20']"
                  @click="emit('clean', safe.id)">
            <EyeSlashIcon class="h-5 w-5" aria-hidden="true"/>
          </button>
          <button class="bg-polar-night-4 rounded-md p-2 mx-0.5 text-snow-storm-1 hover:text-frost-1"
                  @click="emit('export', safe.id)">
            <ArrowUpOnSquareIcon class="h-5 w-5" aria-hidden="true"/>
          </button>
        </div>

        <div class="flex-shrink px-1">
          <div class="sm:hidden">
            <label for="tabs" class="sr-only">Select a tab</label>
            <!-- Use an "onChange" listener to redirect the user to the selected tab URL. -->
            <select id="tabs" name="tabs" @change="selectTab($event)"
                    class="block w-full rounded-md border-gray-300 py-2 pl-3 pr-10 text-base focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm">
              <option v-for="tab in tabs" :key="tab.name" :selected="tab.current">{{ tab.name }}</option>
            </select>
          </div>
          <div class="hidden sm:block">
            <div>
              <nav class="-mb-px flex space-x-8 pt-2" aria-label="Tabs">
                <a v-for="tab in tabs" :key="tab.name" @click="switchTab(tab.id)"
                   :class="[tab.current
                    ? 'border-frost text-frost'
                    : 'border-transparent text-gray-400 hover:border-gray-500 hover:text-gray-200',
                  'group inline-flex cursor-pointer items-center border-b-2 py-2 px-1 text-sm font-medium']"
                   :aria-current="tab.current ? 'page' : undefined">
                  <component :is="tab.icon" :class="[
                    tab.current ? 'text-frost' : 'text-gray-400 group-hover:text-gray-300',
                    '-ml-0.5 mr-2 h-5 w-5',
                  ]" aria-hidden="true"/>
                  <span>{{ tab.name }}</span></a>
              </nav>
            </div>
          </div>
        </div>

        <div class="flex-grow overflow-hidden box-border border-t border-polar-night-3">
          <ScrollingOutput v-if="selectedTab() === 'stdout'" :lines="stdoutLines"/>
          <ScrollingOutput v-else-if="selectedTab() === 'stderr'" :lines="stderrLines"/>
          <ScrollingOutput v-else-if="selectedTab() === 'activity'" :lines="activityLines"/>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
.stripy {
  background: repeating-linear-gradient(
      45deg,
      #242933,
      #242933 10px,
      #212630 10px,
      #212630 20px
  );
}
</style>
