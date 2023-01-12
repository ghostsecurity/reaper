<script lang="ts">
import {defineComponent, PropType} from "vue";
import {HttpRequest, HttpResponse} from "../lib/Http";
import {EventsOn} from "../../wailsjs/runtime";
import {BarsArrowDownIcon, BeakerIcon, StarIcon} from "@heroicons/vue/20/solid";
import {Criteria} from "../lib/Criteria";
import {workspace} from "../../wailsjs/go/models";

export default /*#__PURE__*/ defineComponent({
  name: "Dashboard",
  props: {
    ws: {
      type: Object as PropType<workspace.Workspace>,
      required: true
    },
    criteria: {
      type: Object as PropType<Criteria>,
      required: true,
    },
    onCriteriaChange: {
      type: Function as PropType<(criteria: Criteria) => void>,
      required: false,
    },
    proxyAddress: {
      type: String,
      required: true,
    },
    onWorkspaceConfig: {
      type: Function as PropType<() => void>,
      required: true,
    },
    savedRequestIds: {
      type: Array as PropType<string[]>,
      required: false,
      default: [],
    }
  },
  emits: [
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
  ],
  data() {
    return {
      requests: Array<HttpRequest>(),
      req: null as HttpRequest | null,
      tabs: [
        {name: 'Log Stream', id: 'log', icon: BarsArrowDownIcon, current: true},
        {name: 'Saved Requests', id: 'saved', icon: StarIcon, current: false},
        {name: 'Attack Workflows', id: 'workflows', icon: BeakerIcon, current: false},
      ],
      liveCriteria: this.criteria,
    }
  },
  watch: {
    criteria: {
      handler: function () {
        this.liveCriteria = this.criteria
      },
      immediate: true,
    },
  },
  beforeMount() {
    EventsOn("HttpRequest", (data) => {
      this.requests.push(data as HttpRequest);
    });
    EventsOn("HttpResponse", (response: HttpResponse) => {
      for (let i = 0; i < this.requests.length; i++) {
        if (this.requests[i].ID === response.ID) {
          this.requests[i].Response = response;
          break;
        }
      }
    });
  },
  methods: {
    StarIcon,
    examineRequest(request: HttpRequest) {
      this.req = request
    },
    clearRequest() {
      this.req = null
    },
    selectTab(e: Event) {
      this.switchTab((e.target as HTMLSelectElement).value)
    },
    switchTab(id: any) {
      this.tabs.forEach(tab => {
        tab.current = tab.id === id
      })
    },
    selectedTab(): string {
      return this.tabs.find(tab => tab.current)?.id || ''
    },
    onSearch(crit: Criteria) {
      this.liveCriteria = crit;
      if (this.onCriteriaChange) {
        this.onCriteriaChange(crit)
      }
    },
    switchWorkspace() {
      this.$emit("switch-workspace")
    },
    unsaveRequest(req: HttpRequest | workspace.Request) {
      this.$emit("unsave-request", req)
    },
    duplicateRequest(req: workspace.Request) {
      this.$emit("duplicate-request", req)
    },
    saveRequest(req: HttpRequest, groupID: string) {
      this.$emit("save-request", req, groupID)
    },
    clearLog() {
      // TODO: call this from clear button
      this.clearRequest()
      this.requests = []
    },
    setRequestGroup(request: workspace.Request, groupID: string, nextID: string) {
      this.$emit('request-group-change', request, groupID, nextID)
    },
    createRequestGroup(name: string) {
      this.$emit('request-group-create', name)
    },
    reorderGroup(fromID: string, toID: string) {
      this.$emit('group-order-change', fromID, toID)
    },
    deleteGroup(groupId: string) {
      this.$emit("request-group-delete", groupId)
    },
    renameGroup(groupId: string, name: string) {
      this.$emit("request-group-rename", groupId, name)
    },
    renameRequest(requestId: string, name: string) {
      this.$emit("request-rename", requestId, name)
    },
  },
})
</script>

<script lang="ts" setup>
import RequestList from "./RequestList.vue";
import GroupedRequestList from "./GroupedRequestList.vue";
import WorkspaceMenu from "./WorkspaceMenu.vue";
import Search from "./Search.vue";
import RightPop from "./RightPop.vue";
import HttpInspector from "./HttpInspector.vue";
</script>

<template>
  <div class="min-h-16 h-16 max-h-16 flex p-2">
    <div class="flex-grow text-left">
      <Search :on-search="onSearch" :query="liveCriteria.Raw"/>
    </div>
    <div class="flex-shrink p-0 ml-2">
      <WorkspaceMenu :on-workspace-config="onWorkspaceConfig" :ws="ws" @switch-workspace="switchWorkspace"/>
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
          <a v-for="tab in tabs" @click="switchTab(tab.id)" :key="tab.name"
             :class="[tab.current ? 'border-frost text-frost' : 'border-transparent text-gray-400 hover:text-gray-200 hover:border-gray-500', 'group inline-flex items-center py-4 px-1 border-b-2 font-medium text-sm cursor-pointer']"
             :aria-current="tab.current ? 'page' : undefined">
            <component :is="tab.icon"
                       :class="[tab.current ? 'text-frost' : 'text-gray-400 group-hover:text-gray-300', '-ml-0.5 mr-2 h-5 w-5']"
                       aria-hidden="true"/>
            <span>{{ tab.name }}</span>
          </a>
        </nav>
      </div>
    </div>
  </div>
  <div class="px-2">
    <!-- TODO: show intercepted requests here in dialog -->
    <RequestList @save-request="saveRequest" @unsave-request="unsaveRequest" :saved-request-ids="savedRequestIds"
                 :key="liveCriteria.Raw" v-if="selectedTab() === 'log'"
                 :empty-message="'Reaper is ready to receive requests at ' + proxyAddress" :requests="requests"
                 :onSelect="examineRequest" :selected="req ? req.ID : ''" :criteria="liveCriteria"/>
    <RightPop :show="req !== null" :onRequestClose="clearRequest">
      <HttpInspector v-if="req" :request="req"/>
    </RightPop>

    <GroupedRequestList :key="liveCriteria.Raw" v-if="selectedTab() === 'saved'"
                        :groups="ws.collection.groups ? ws.collection.groups : []" :onSelect="examineRequest"
                        :selected="req ? req.ID : ''" :criteria="liveCriteria" :empty-title="'No saved requests'"
                        :empty-message="'Save some requests from the log stream to access them here'"
                        :empty-icon="StarIcon"
                        @request-group-change="setRequestGroup"
                        @request-group-create="createRequestGroup"
                        @group-order-change="reorderGroup"
                        @unsave-request="unsaveRequest"
                        @duplicate-request="duplicateRequest"
                        @request-group-delete="deleteGroup"
                        @request-rename="renameRequest"
                        @request-group-rename="renameGroup"
    />

    <div v-if="selectedTab() === 'workflows'">
      TODO: Attack Workflows
    </div>
  </div>
</template>

<style scoped>
.body {
  max-height: calc(100vh - 8rem);
}
</style>