<script lang="ts" setup>
import { onBeforeMount, onMounted, PropType, reactive, ref, watch } from 'vue'
import { BarsArrowDownIcon, BeakerIcon, StarIcon, HandRaisedIcon } from '@heroicons/vue/20/solid'
import { EventsEmit, EventsOn } from '../../wailsjs/runtime'
import { HttpRequest, HttpResponse } from '../lib/Http'
import { Criteria } from '../lib/Criteria/Criteria'
import { workspace } from '../../wailsjs/go/models'
import RequestList from './Http/RequestList.vue'
import GroupedRequestList from './Http/GroupedRequestList.vue'
import WorkspaceMenu from './WorkspaceMenu.vue'
import Search from './SearchInput.vue'
import IDE from './Http/IDE.vue'
import RequestInterceptor from './RequestInterceptor.vue'
import GUI from './Workflow/WorkflowGUI.vue'

const props = defineProps({
  ws: { type: Object as PropType<workspace.Workspace>, required: true },
  criteria: { type: Object as PropType<Criteria>, required: true },
  proxyAddress: { type: String, required: true },
  savedRequestIds: { type: Array as PropType<string[]>, required: false, default: () => [] },
})

const emit = defineEmits([
  'switch-workspace',
  'save-request',
  'unsave-request',
  'request-group-change',
  'request-group-create',
  'group-order-change',
  'group-expand',
  'duplicate-request',
  'request-group-delete',
  'request-group-rename',
  'request-rename',
  'criteria-change',
  'workspace-edit',
  'workspace-save',
  'copy-request-as-curl',
  'send-request',
  'update-request',
])

const requests = ref<HttpRequest[]>([])
const req = ref<HttpRequest | null>(null)
const tabs = ref([
  { name: 'Log Stream', id: 'log', icon: BarsArrowDownIcon, current: true, count: 0 },
  { name: 'Intercepted Requests', id: 'intercepted', icon: HandRaisedIcon, current: false, count: 0 },
  { name: 'Saved Requests', id: 'saved', icon: StarIcon, current: false, count: 0 },
  { name: 'Attack Workflows', id: 'workflows', icon: BeakerIcon, current: false, count: 0 },
])
const liveCriteria = reactive(props.criteria)
const reqReadOnly = ref(true)
const fullscreenIDE = ref(false)

const root = ref()
const leftPanel = ref()
const rightPanel = ref()
const handle = ref()
const resizing = ref(false)

let sendingReq: HttpRequest | null = null

const readOnlyActions = new Map<string, string>([
  ['send', 'Resend'],
  ['copy-curl', 'Copy as curl'],
])

const writeActions = new Map<string, string>([
  ['send', 'Send'],
  ['copy-curl', 'Copy as curl'],
])

const interceptionCount = ref(0)
const interceptedRequest = ref<HttpRequest | null>(null)
const sentInterceptedRequest = ref<HttpRequest | null>(null)

function ideAction(action: string) {
  switch (action) {
    case 'send':
      sendingReq = req.value
      emit('send-request', req.value)
      break
    case 'copy-curl':
      break
    default:
      throw new Error(`unknown action ${action}`)
  }
}

function toggleFullscreenIDE() {
  fullscreenIDE.value = !fullscreenIDE.value
  setDefaultSizing()
}

watch(
  () => props.criteria,
  criteria => {
    Object.assign(liveCriteria, criteria)
  },
)

function compareRequests(a: HttpRequest, b: HttpRequest) {
  if (a.Method !== b.Method) {
    return false
  }
  if (a.URL !== b.URL) {
    return false
  }
  if (a.Body !== b.Body) {
    return false
  }
  return true
}

onBeforeMount(() => {
  EventsOn('HttpRequest', (data: HttpRequest) => {
    requests.value.push(data)
    // TODO: better way to identify the request we're sending
    if (sendingReq !== null && compareRequests(data, sendingReq)) {
      sendingReq = data
    }
  })
  EventsOn('HttpResponse', (response: HttpResponse) => {
    if (sentInterceptedRequest.value !== null && sentInterceptedRequest.value.ID === response.ID) {
      sentInterceptedRequest.value.Response = response
    }
    for (let i = 0; i < requests.value.length; i += 1) {
      if (requests.value[i].ID === response.ID) {
        requests.value[i].Response = response
        const r = requests.value[i]
        if (sendingReq !== null && sendingReq.ID === r.ID) {
          sendingReq = null
          if (req.value !== null) {
            req.value.Response = response
          }
        }
        break
      }
    }
  })
  EventsOn('InterceptedRequestQueueChange', (count: number) => {
    for (let i = 0; i < tabs.value.length; i += 1) {
      if (tabs.value[i].id === 'intercepted') {
        tabs.value[i].count = count
        interceptionCount.value = count
        break
      }
    }
  })
  EventsOn('InterceptedRequest', (request: HttpRequest) => {
    interceptedRequest.value = request
  })
})

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

    // Resize box A
    // * 8px is the left/right spacing between .handler and its inner pseudo-element
    // * Set flex-grow to 0 to prevent it from growing
    // leftPanel.value.style.width = `${Math.max(boxAminWidth, pointerRelativeXpos - 8)}px`
    // leftPanel.value.style.flexGrow = 0
    // leftPanel.value.style.flexShrink = 0
  })
  root.value.addEventListener('mouseup', () => {
    resizing.value = false
  })
})

function examineRequest(request: HttpRequest, readonly: boolean) {
  reqReadOnly.value = readonly
  req.value = request
}

function clearRequest() {
  fullscreenIDE.value = false
  req.value = null
  setDefaultSizing()
}

function setDefaultSizing() {
  if (leftPanel.value !== null) {
    leftPanel.value.style.flexGrow = 1
    leftPanel.value.style.flexShrink = 1
  }
  if (rightPanel.value !== null) {
    rightPanel.value.style.flexGrow = 1
    rightPanel.value.style.flexShrink = 1
  }
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

function selectedTab(): string {
  return tabs.value.find(tab => tab.current)?.id || ''
}

function onSearch(crit: Criteria) {
  Object.assign(liveCriteria, crit)
  emit('criteria-change', crit)
}

function switchWorkspace() {
  emit('switch-workspace')
}

function unsaveRequest(r: HttpRequest | workspace.Request) {
  emit('unsave-request', r)
}

function duplicateRequest(r: workspace.Request) {
  emit('duplicate-request', r)
}

function saveRequest(r: HttpRequest, groupID: string) {
  emit('save-request', r, groupID)
}

function setRequestGroup(request: workspace.Request, groupID: string, nextID: string) {
  emit('request-group-change', request, groupID, nextID)
}

function createRequestGroup(name: string) {
  emit('request-group-create', name)
}

function reorderGroup(fromID: string, toID: string) {
  emit('group-order-change', fromID, toID)
}

function deleteGroup(groupId: string) {
  emit('request-group-delete', groupId)
}

function renameGroup(groupId: string, name: string) {
  emit('request-group-rename', groupId, name)
}

function renameRequest(requestId: string, name: string) {
  emit('request-rename', requestId, name)
}

function updateRequest(e: HttpRequest) {
  if (req.value !== null && req.value.ID === e.ID) {
    req.value = e
  }
  emit('update-request', e)
}

function dropInterceptedRequest(request: HttpRequest) {
  if (interceptedRequest.value === null) {
    return
  }
  interceptedRequest.value = null
  EventsEmit('InterceptedRequestDrop', request)
}

function sendInterceptedRequest(request: HttpRequest) {
  if (interceptedRequest.value === null) {
    return
  }
  sentInterceptedRequest.value = request
  interceptedRequest.value = null
  EventsEmit('InterceptedRequestChange', request)
}

function closeInterceptedRequest() {
  sentInterceptedRequest.value = null
}
</script>

<template>
  <div ref="root" class="flex h-full overflow-x-hidden">
    <!-- main content with search, tabs etc. -->
    <div v-if="!fullscreenIDE" ref="leftPanel" class="box-border flex-auto">
      <div class="flex h-full flex-col px-2">
        <div class="min-h-16 flex h-16 max-h-16 flex-shrink py-2">
          <div class="flex-grow text-left">
            <Search @search="onSearch" :query="liveCriteria.Raw"/>
          </div>
          <div class="ml-2 flex-shrink p-0">
            <WorkspaceMenu @edit="emit('workspace-edit')" :ws="ws" @switch="switchWorkspace"/>
          </div>
        </div>
        <div class="min-h-16 h-16 max-h-16 flex-shrink">
          <div class="sm:hidden">
            <label for="tabs" class="sr-only">Select a tab</label>
            <select @change="selectTab" id="tabs" name="tabs"
                    class="block w-full rounded-md bg-polar-night-2 text-snow-storm-1 focus:border-frost focus:ring-frost">
              <option v-for="tab in tabs" :key="tab.id" :selected="tab.current" :value="tab.id">{{ tab.name }}</option>
            </select>
          </div>
          <div class="hidden sm:block">
            <div class="border-b dark:border-polar-night-4">
              <nav class="-mb-px flex space-x-8" aria-label="Tabs">
                <a v-for="tab in tabs" @click="switchTab(tab.id)" :key="tab.name" :class="[
                  tab.current
                    ? 'border-frost text-frost'
                    : 'border-transparent text-gray-400 hover:border-gray-500 hover:text-gray-200',
                  'group inline-flex cursor-pointer items-center border-b-2 py-4 px-1 text-sm font-medium',
                ]" :aria-current="tab.current ? 'page' : undefined">
                  <component :is="tab.icon" :class="[
                    tab.current ? 'text-frost' : 'text-gray-400 group-hover:text-gray-300',
                    '-ml-0.5 mr-2 h-5 w-5',
                  ]" aria-hidden="true"/>
                  <span>{{ tab.name }}</span>
                  <span v-if="tab.count > 0"
                        class="bg-indigo-100 text-indigo-600 ml-3 hidden rounded-full py-0.5 px-2.5 text-xs font-medium md:inline-block">{{
                      tab.count
                    }}</span>
                </a>
              </nav>
            </div>
          </div>
        </div>
        <div class="my-2 flex-auto overflow-y-hidden">
          <RequestList @save-request="saveRequest" @unsave-request="unsaveRequest" :saved-request-ids="savedRequestIds"
                       :key="liveCriteria.Raw" v-if="selectedTab() === 'log'"
                       :empty-message="'Reaper is ready to receive requests at ' + proxyAddress" :requests="requests"
                       @select="examineRequest($event, true)" :selected="req ? req.ID : ''" :criteria="liveCriteria"/>
          <GroupedRequestList :key="liveCriteria.Raw" v-if="selectedTab() === 'saved'"
                              :groups="ws.collection.groups ? ws.collection.groups : []"
                              @select="examineRequest($event, false)"
                              :selected="req ? req.ID : ''" :criteria="liveCriteria" :empty-title="'No saved requests'"
                              :empty-message="'Save some requests from the log stream to access them here'"
                              :empty-icon="StarIcon"
                              @request-group-change="setRequestGroup" @request-group-create="createRequestGroup"
                              @group-order-change="reorderGroup" @unsave-request="unsaveRequest"
                              @duplicate-request="duplicateRequest"
                              @request-group-delete="deleteGroup" @request-rename="renameRequest"
                              @request-group-rename="renameGroup"/>
          <GUI v-if="selectedTab() === 'workflows'" :ws="ws" @save="emit('workspace-save', $event)"/>
          <RequestInterceptor v-if="selectedTab() === 'intercepted'" :request="interceptedRequest"
                              :previous="sentInterceptedRequest"
                              :count="interceptionCount"
                              @drop="dropInterceptedRequest" @send="sendInterceptedRequest"
                              @close="closeInterceptedRequest"/>
        </div>
      </div>
    </div>

    <!-- resize handle -->
    <div v-if="req && !fullscreenIDE" @mousedown.prevent="resizing = true" ref="handle"
         class="w-0.5 flex-none cursor-ew-resize bg-gray-500 dark:bg-polar-night-4"></div>

    <!-- request viewer/editor -->
    <div v-if="req" ref="rightPanel" class="mx-2 box-border h-full flex-auto overflow-hidden px-2">
      <IDE :request="req" :readonly="reqReadOnly" @action="ideAction" @close="clearRequest" :fullscreen="fullscreenIDE"
           @fullscreen="toggleFullscreenIDE" @request-update="updateRequest($event)"
           :actions="reqReadOnly ? readOnlyActions : writeActions"/>
    </div>
  </div>
</template>
