<script lang="ts" setup>
import { onBeforeMount, PropType, reactive, ref, watch } from 'vue'
import { BarsArrowDownIcon, BeakerIcon, StarIcon } from '@heroicons/vue/20/solid'
import { EventsOn } from '../../wailsjs/runtime'
import { HttpRequest, HttpResponse } from '../lib/Http'
import { Criteria } from '../lib/Criteria/Criteria'
import { workspace } from '../../wailsjs/go/models'
import RequestList from './Http/RequestList.vue'
import GroupedRequestList from './Http/GroupedRequestList.vue'
import WorkspaceMenu from './WorkspaceMenu.vue'
import Search from './SearchInput.vue'
import IDE from './Http/IDE.vue'

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
])

const requests = ref<HttpRequest[]>([])
const req = ref<HttpRequest | null>(null)
const tabs = ref([
  { name: 'Log Stream', id: 'log', icon: BarsArrowDownIcon, current: true },
  { name: 'Saved Requests', id: 'saved', icon: StarIcon, current: false },
  { name: 'Attack Workflows', id: 'workflows', icon: BeakerIcon, current: false },
])
const liveCriteria = reactive(props.criteria)
const reqReadOnly = ref(true)

watch(() => props.criteria, (criteria) => {
  Object.assign(liveCriteria, criteria)
})

onBeforeMount(() => {
  EventsOn('HttpRequest', (data: HttpRequest) => {
    requests.value.push(data)
  })
  EventsOn('HttpResponse', (response: HttpResponse) => {
    for (let i = 0; i < requests.value.length; i += 1) {
      if (requests.value[i].ID === response.ID) {
        requests.value[i].Response = response
        break
      }
    }
  })
})

function examineRequest(request: HttpRequest, readonly: boolean) {
  req.value = request
  reqReadOnly.value = readonly
}

function clearRequest() {
  req.value = null
}

function switchTab(id: string) {
  tabs.value = tabs.value.map((tab) => {
    const updatedTab = tab
    updatedTab.current = updatedTab.id === id
    return updatedTab
  })
}

function selectTab(e: Event) {
  switchTab((e.target as HTMLSelectElement).value)
}

function selectedTab(): string {
  return tabs.value.find((tab) => tab.current)?.id || ''
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
</script>

<template>
  <div class="min-h-16 h-16 max-h-16 flex p-2">
    <div class="flex-grow text-left">
      <Search @search="onSearch" :query="liveCriteria.Raw" />
    </div>
    <div class="flex-shrink p-0 ml-2">
      <WorkspaceMenu @edit="emit('workspace-edit')" :ws="ws" @switch="switchWorkspace" />
    </div>
  </div>
  <div class="min-h-16 h-16 max-h-16 px-2">
    <div class="sm:hidden">
      <label for="tabs" class="sr-only">Select a tab</label>
      <select @change="selectTab" id="tabs" name="tabs"
        class="bg-polar-night-2 text-snow-storm-1 block w-full rounded-md focus:border-frost focus:ring-frost">
        <option v-for="tab in tabs" :key="tab.id" :selected="tab.current" :value="tab.id">{{ tab.name }}</option>
      </select>
    </div>
    <div class="hidden sm:block">
      <div class="border-b dark:border-polar-night-4">
        <nav class="-mb-px flex space-x-8" aria-label="Tabs">
          <a v-for="tab in tabs" @click="switchTab(tab.id)" :key="tab.name" :class="[
          tab.current ?
            'border-frost text-frost' :
            'border-transparent text-gray-400 hover:text-gray-200 hover:border-gray-500',
          'group inline-flex items-center py-4 px-1 border-b-2 font-medium text-sm cursor-pointer']"
            :aria-current="tab.current ? 'page' : undefined">
            <component :is="tab.icon" :class="[
              tab.current ?
                'text-frost' :
                'text-gray-400 group-hover:text-gray-300',
              '-ml-0.5 mr-2 h-5 w-5'
            ]" aria-hidden="true" />
            <span>{{ tab.name }}</span>
          </a>
        </nav>
      </div>
    </div>
  </div>
  <div class="px-2 flex h-full">
    <div class="flex-1">
      <RequestList @save-request="saveRequest" @unsave-request="unsaveRequest" :saved-request-ids="savedRequestIds"
        :key="liveCriteria.Raw" v-if="selectedTab() === 'log'"
        :empty-message="'Reaper is ready to receive requests at ' + proxyAddress" :requests="requests"
        @select="examineRequest($event, true)" :selected="req ? req.ID : ''" :criteria="liveCriteria" />
      <GroupedRequestList :key="liveCriteria.Raw" v-if="selectedTab() === 'saved'"
        :groups="ws.collection.groups ? ws.collection.groups : []" @select="examineRequest($event, false)"
        :selected="req ? req.ID : ''" :criteria="liveCriteria" :empty-title="'No saved requests'"
        :empty-message="'Save some requests from the log stream to access them here'" :empty-icon="StarIcon"
        @request-group-change="setRequestGroup" @request-group-create="createRequestGroup"
        @group-order-change="reorderGroup" @unsave-request="unsaveRequest" @duplicate-request="duplicateRequest"
        @request-group-delete="deleteGroup" @request-rename="renameRequest" @request-group-rename="renameGroup" />
      <div v-if="selectedTab() === 'workflows'">
        TODO: Attack Workflows
      </div>
    </div>

    <!-- TODO: intercept stuff here -->
    <div v-if="false" class="flex-0 w-[50%] min-w- pl-2 ml-2 border-l border-polar-night-2 h-full">
      intercept ui here
    </div>

    <div v-else-if="req" class="flex-0 w-[50%] min-w- pl-2 ml-2 border-l border-polar-night-2 h-full">
      <IDE :request="req" :readonly="reqReadOnly" @close="clearRequest" />
    </div>
  </div>
</template>

<style scoped>
.body {
  max-height: calc(100vh - 8rem);
}
</style>
