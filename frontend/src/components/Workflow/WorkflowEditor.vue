<script lang="ts" setup>
import {onBeforeUpdate, onMounted, onUpdated, PropType, ref, watch} from 'vue'
import {BeakerIcon, TrashIcon, PlusIcon, ArrowPathIcon, PlayIcon} from '@heroicons/vue/20/solid'
import {workflow} from '../../../wailsjs/go/models'
import {CreateNode} from "../../../wailsjs/go/backend/App";
import NodeEditor from "./NodeEditor.vue";
import {NodeType, ParentType, NodeTypeName} from "../../lib/Workflows";
import ConfirmDialog from "../ConfirmDialog.vue";

const props = defineProps({
  flow: {type: Object as PropType<workflow.WorkflowM>, required: true},
})

const availableNodeTypes = ref(<NodeType[]>[
  NodeType.REQUEST,
  NodeType.SENDER,
  NodeType.STATUS_FILTER,
  NodeType.FUZZER,
  NodeType.OUTPUT,
])

const linkColour = "#8FBCBB"
const linkStrokeWidth = "2"
const editingNode = ref(<workflow.NodeM | null>null)
const menuMode = ref("")
const resetting = ref(false)

const safe = ref<workflow.WorkflowM>(JSON.parse(JSON.stringify(props.flow)))
let initial = true
watch(() => props.flow, flow => {
  if (flow) {
    safe.value = JSON.parse(JSON.stringify(props.flow)) as workflow.WorkflowM
    redraw()
  }
})

onMounted(() => {
  redraw()
})

onUpdated(() => {
  if (initial) {
    redraw()
    initial = false
  }
})

const svg = ref(<HTMLElement | null>null)
const connector = ref(<HTMLElement | null>null)
const paths = ref(<string[]>[])

const emit = defineEmits(['save'])

function saveWorkflow(f: workflow.WorkflowM) {
  emit('save', f)
  redraw()
}

function redraw() {
  if (!safe.value) {
    return
  }
  let newPaths: string[] = []
  let ok = true
  safe.value.links.forEach((link) => {

    const fromNode = movers.value.get(link.from.node) as HTMLElement
    const toNode = movers.value.get(link.to.node) as HTMLElement
    if (!fromNode || !toNode) {
      ok = false
      return
    }

    let fromConn = fromNode.querySelector(".connector.output[data-connector=\"" + link.from.connector + "\"]") as HTMLElement
    let toConn = toNode.querySelector(".connector.input[data-connector=\"" + link.to.connector + "\"]") as HTMLElement
    if (!fromConn || !toConn) {
      ok = false
      return
    }

    let posA = {
      x: (fromConn.offsetLeft + fromConn.offsetWidth) + fromNode.offsetLeft,
      y: (fromConn.offsetTop + fromConn.offsetHeight / 2) + fromNode.offsetTop
    };
    let posB = {
      x: toConn.offsetLeft + toNode.offsetLeft,
      y: (toConn.offsetTop + toConn.offsetHeight / 2) + toNode.offsetTop
    };
    let path =
        "M" +
        (posA.x) + "," + (posA.y) + " " +
        "C" +
        (posA.x + 100) + "," + (posA.y) + " " +
        (posB.x - 100) + "," + (posB.y) + " " +
        (posB.x) + "," + (posB.y);
    newPaths.push(path)
  })
  if (!ok) {
    initial = true
    return
  }
  paths.value = newPaths
}

function getRawPosition(id: string, i: number) {
  let def = {
    x: 50 + (((i % 4) * 300)),
    y: 50 + (((i - (i % 4)) / 4) * 100),
  }
  if (!safe.value || !safe.value.positioning) {
    return def
  }
  let pos = safe.value.positioning[id]
  if (!pos) {
    return def
  }
  return {
    x: pos.x,
    y: pos.y,
  }
}

function getPosition(id: string, i: number) {
  let pos = getRawPosition(id, i)
  return {
    left: pos.x + 'px',
    top: pos.y + 'px',
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
    x: x,
    y: y,
  })
  saveWorkflow(safe.value)
}

function addNode(t: number) {
  if (!safe.value) {
    return
  }
  CreateNode(t).then((n) => {
    safe.value.nodes.push(n)
    editingNode.value = n
    menuMode.value = ""
    saveWorkflow(safe.value)
  })
}

function reset() {
  safe.value.nodes = safe.value.nodes.filter((n) => {
    return n.readonly
  })
  safe.value.links = []
  saveWorkflow(safe.value)
  resetting.value = false
}

const canvas = ref(<HTMLDivElement | null>null)
let dragId = ""
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
    x: x,
    y: y,
  }
}

function dragStart(id: string, ev: MouseEvent) {
  if (!canvas.value) {
    return
  }
  dragId = id
  let el = ev.target as HTMLDivElement
  let offset = getOffsetFrom(el, "mover")
  offsetX = ev.offsetX + offset.x
  offsetY = ev.offsetY + offset.y
}

function dragEnd(ev: MouseEvent) {
  if (!canvas.value) {
    return
  }
  dragId = ""
  if (currentLink) {
    endLinking(ev)
  }
}

function drag(ev: MouseEvent) {
  if (dragId === "" || !canvas.value) {
    if (currentLink) {
      moveLink(ev)
    }
    return
  }
  ev.preventDefault()
  let x = ev.clientX - canvas.value.offsetLeft
  let y = ev.clientY - canvas.value.offsetTop
  let el = movers.value.get(dragId) as HTMLElement
  if (!el) {
    return
  }
  el.style.left = (x - offsetX) + 'px'
  el.style.top = (y - offsetY) + 'px'
  setPosition(dragId, x - offsetX, y - offsetY)
}

//const movers = ref([] as any)
const movers = ref(new Map<string, any>())

// Make sure to reset the refs before each update.
onBeforeUpdate(() => {
  movers.value = new Map<string, any>()
});

let currentLink: workflow.LinkM | null = null
let currentLinkNode = ""
let currentLinkConnector = ""

let linkSearchMode = '';

function startLinkFromInput(nodeId: string, connector: string) {
  if (!safe.value) {
    return
  }
  linkSearchMode = 'output'
  currentLinkNode = nodeId
  currentLinkConnector = connector
  const existing = safe.value.links.find((l) => {
    return l.to.node == nodeId && l.to.connector == connector
  })
  unlinkAnyFromInput(nodeId, connector)
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
        node: "",
        connector: "",
      }),
    })
  } else {
    currentLink = new workflow.LinkM({
      from: new workflow.LinkDirectionM({
        node: "",
        connector: "",
      }),
      to: new workflow.LinkDirectionM({
        node: nodeId,
        connector: connector,
      }),
    })
  }
}

function startLinkFromOutput(nodeId: string, connector: string) {
  if (!safe.value) {
    return
  }
  linkSearchMode = 'input'
  currentLinkNode = nodeId
  currentLinkConnector = connector
  const existing = safe.value.links.find((l) => {
    return l.from.node == nodeId && l.from.connector == connector
  })
  unlinkAnyFromOutput(nodeId, connector)
  if (existing) {
    linkSearchMode = 'output'
    currentLinkNode = existing.to.node
    currentLinkConnector = existing.to.connector
    currentLink = new workflow.LinkM({
      from: new workflow.LinkDirectionM({
        node: "",
        connector: "",
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
        connector: connector,
      }),
      to: new workflow.LinkDirectionM({
        node: "",
        connector: "",
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
  let connectors = canvas.value.querySelectorAll(".connector." + className)
  for (let i = 0; i < connectors.length; i++) {
    let el = connectors[i] as HTMLElement
    if (!el.parentElement || !canvas.value) {
      continue
    }
    let rect = el.parentElement.getBoundingClientRect()
    let area = {
      x1: rect.left - canvas.value.offsetLeft - magnetism,
      y1: rect.top - canvas.value.offsetTop,
      x2: rect.right - canvas.value.offsetLeft + magnetism,
      y2: rect.bottom - canvas.value.offsetTop,
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

  let selector = ".connector."
  if (linkSearchMode === 'input') {
    selector += 'output'
  } else {
    selector += 'input'
  }

  let nodeA = movers.value.get(currentLinkNode) as HTMLElement
  let connectorA = nodeA.querySelector(selector + "[data-connector=\"" + currentLinkConnector + "\"]") as HTMLElement

  let posA = {
    x: (connectorA.offsetLeft + connectorA.offsetWidth) + nodeA.offsetLeft,
    y: (connectorA.offsetTop + connectorA.offsetHeight / 2) + nodeA.offsetTop
  };

  let posB = {
    x: ev.clientX - canvas.value.offsetLeft,
    y: ev.clientY - canvas.value.offsetTop
  }

  let connectorB = findConnectorAt(posB.x, posB.y, linkSearchMode)
  if (connectorB) {
    let nodeB = connectorB.parentElement as HTMLElement
    while (nodeB && !nodeB.className.includes("mover")) {
      nodeB = nodeB.parentElement as HTMLElement
    }
    posB = {
      x: connectorB.offsetLeft + nodeB.offsetLeft,
      y: connectorB.offsetTop + nodeB.offsetTop + (connectorB.offsetHeight / 2)
    }
    let node = connectorB.getAttribute("data-node")
    let conn = connectorB.getAttribute("data-connector")
    if (node && conn) {
      if (linkSearchMode === 'input') {
        currentLink.to.node = node
        currentLink.to.connector = conn
      } else {
        currentLink.from.node = node
        currentLink.from.connector = conn
      }
      if (canLink(currentLink)) {
        connector.value.setAttribute("stroke", "green")
      } else {
        connector.value.setAttribute("stroke", "red")
      }
    }
  } else {
    if (linkSearchMode === 'input') {
      currentLink.to.node = ""
      currentLink.to.connector = ""
    } else {
      currentLink.from.node = ""
      currentLink.from.connector = ""
    }
  }

  if (posA.x > posB.x) {
    let buf = posA
    posA = posB
    posB = buf
  }

  let dStr =
      "M" +
      (posA.x) + "," + (posA.y) + " " +
      "C" +
      (posA.x + 100) + "," + (posA.y) + " " +
      (posB.x - 100) + "," + (posB.y) + " " +
      (posB.x) + "," + (posB.y);
  connector.value.setAttribute("d", dStr);
  if (!connectorB) {
    connector.value.setAttribute("stroke", "blue")
  }
}


function canLink(link: workflow.LinkM): boolean {
  if (!safe.value) {
    return false
  }
  if (link.from.node == link.to.node) {
    return false
  }
  let fromNode = safe.value.nodes.find((n) => {
    return n.id == link.from.node
  })
  let toNode = safe.value.nodes.find((n) => {
    return n.id == link.to.node
  })
  if (!fromNode || !toNode || !fromNode.vars || !toNode.vars) {
    return false
  }
  let output = fromNode.vars.outputs.find((o) => {
    return o.name == link.from.connector
  })
  let input = toNode.vars.inputs.find((i) => {
    return i.name == link.to.connector
  })
  if (!output || !input) {
    return false
  }
  if ((input.type & output.type) !== input.type) {
    return false
  }
  return true
}

function endLinking(ev: MouseEvent) {
  if (!currentLink || !connector.value || !canvas.value || !safe.value) {
    return
  }

  if (canLink(currentLink)) {
    safe.value.links.push(currentLink)
    saveWorkflow(safe.value)
  }

  connector.value.setAttribute("d", "");
  currentLink = null
  currentLinkConnector = ""
  currentLinkNode = ""
}

function unlinkAnyFromOutput(node: string, connector: string) {
  if (!safe.value) {
    return
  }
  let index = safe.value.links.findIndex((l) => {
    return l.from.node === node && (l.from.connector === connector || connector === "")
  })
  if (index >= 0) {
    safe.value.links.splice(index, 1)
    saveWorkflow(safe.value)

  }
}

function unlinkAnyFromInput(node: string, connector: string) {
  if (!safe.value) {
    return
  }
  let index = safe.value.links.findIndex((l) => {
    return l.to.node === node && (l.to.connector === connector || connector === "")
  })
  if (index >= 0) {
    safe.value.links.splice(index, 1)
    saveWorkflow(safe.value)
  }
}

function editNode(id: string) {
  editingNode.value = safe.value.nodes.find((n) => {
    return n.id == id
  }) as workflow.NodeM
  menuMode.value = ""
}

function deleteNode(id: string) {
  if (!safe.value) {
    return
  }
  const index = safe.value.nodes.findIndex((n) => {
    return n.id == id
  })
  if (index < 0) {
    return
  }
  safe.value.nodes.splice(index, 1)
  unlinkAnyFromInput(id, "")
  unlinkAnyFromOutput(id, "")
  saveWorkflow(safe.value)
}

function updateNode(n: workflow.NodeM) {
  if (!safe.value) {
    return
  }
  const index = safe.value.nodes.findIndex((node) => {
    return node.id == n.id
  })
  if (index < 0) {
    return
  }
  safe.value.nodes[index] = n
  saveWorkflow(safe.value)
}

function setMenu(n: string) {
  editingNode.value = null
  if (menuMode.value === n) {
    menuMode.value = ""
    return
  }
  menuMode.value = n
}

function requestReset() {
  resetting.value = true
}

</script>

<template>
  <div class="h-full">
    <ConfirmDialog :show="resetting" title="Reset Workflow" cancel="Cancel" confirm="Reset"
                   message="Are you sure you want to reset the workflow? This will remove all nodes except the 'Start'."
                   @confirm="reset" @cancel="resetting = false"/>
    <div v-if="!safe" class="flex flex-col items-center mt-16">
      <BeakerIcon class="h-12 w-12"/>
      <h3 class="mt-2 text-sm font-bold">No Workflow Selected</h3>
      <p class="mt-1 text-sm">Select or create a workflow from the list.</p>
    </div>
    <div v-else class="h-full flex flex-col">
      <div class="border border-polar-night-4 flex-auto my-2 relative stripy w-full h-full overflow-auto"
           @mousemove="drag" @mouseup="dragEnd"
           ref="canvas"
      >
        <svg ref="svg" xmlns="http://www.w3.org/2000/svg" class="absolute w-full h-full">
          <path ref="connector" fill="none" stroke="" :stroke-width="linkStrokeWidth"/>
          <path v-for="path in paths" :d="path" fill="none" :stroke="linkColour" :stroke-width="linkStrokeWidth"
                :key="path"/>
        </svg>
        <div v-for="(node, i) in safe.nodes" :key="node.id"
             class="mover absolute"
             :ref="el => { movers.set(node.id, el) }"
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
                  @click="editNode(node.id)"
                  :class="[editingNode && editingNode.id == node.id ? 'border-snow-storm-1' : 'border-polar-night-4',  'px-2 py-4 border rounded cursor-move bg-polar-night-2 relative']"
                  style="min-width:90px">
                {{ node.name }}
                <div v-if="!node.readonly" @mousedown.prevent.stop
                     class="absolute top-0 right-0 invisible hover:visible group-hover:visible text-snow-storm-1">
                  <button @click="deleteNode(node.id)" class="group/btn px-0.5">
                    <TrashIcon class="h-4 w-4 group-hover/btn:text-aurora-2"/>
                  </button>
                </div>
              </div>
              <div
                  class="absolute ml-1 top-0 py-0.5 p-2 text-xs bg-polar-night-2 rounded-md text-gray-400 italic"
                  style="margin-top: -0.5rem;">
                {{ NodeTypeName(node.type) }}
              </div>


            </div>
            <div class="flex-shrink">
              <div @mousedown.prevent.stop="startLinkFromOutput(node.id, output.name)"
                   v-for="(output, j) in node.vars?.outputs"
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
        <div class="absolute right-5 top-5 text-right w-96">
          <button type="button" @click="setMenu('add')"
                  :class="[menuMode==='add'?'bg-frost-1':'bg-frost-4', 'mb-1 rounded-full  p-1.5 text-white shadow-sm hover:bg-frost-1 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 mx-0.5']">
            <PlusIcon class="h-5 w-5" aria-hidden="true"/>
          </button>
          <button type="button" @click="setMenu('run')"
                  :class="[menuMode==='run'?'bg-frost-1':'bg-frost-4', 'mb-1 rounded-full  p-1.5 text-white shadow-sm hover:bg-frost-1 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 mx-0.5']">
            <PlayIcon class="h-5 w-5" aria-hidden="true"/>
          </button>
          <button type="button" @click="setMenu('');requestReset()"
                  class="bg-frost-4 mb-1 rounded-full  p-1.5 text-white shadow-sm hover:bg-frost-1 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 mx-0.5">
            <ArrowPathIcon class="h-5 w-5" aria-hidden="true"/>
          </button>
          <div v-if="editingNode" class="border rounded border-snow-storm-1 mt-2">
            <NodeEditor :node="editingNode" @update="updateNode" @close="editingNode = null; menuMode = ''"/>
          </div>
          <div v-else-if="menuMode==='add'" class="border rounded border-snow-storm-1 relative bg-polar-night-2">
            <button v-for="t in availableNodeTypes" @click="addNode(t);setMenu('')"
                    class="w-full border border-polar-night-4 py-1 bg-polar-night-2 hover:bg-polar-night-4">
              {{ NodeTypeName(t) }}
            </button>
            <button @click="reset" class="w-full border-b border-frost-4 py-1 bg-polar-night-2">Reset</button>
          </div>
        </div>
      </div>
      <div v-if="menuMode==='run'">
        lol
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