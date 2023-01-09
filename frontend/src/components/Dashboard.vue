

<script lang="ts">
import TrafficLog from "./TrafficLog.vue";
import {defineComponent, PropType} from "vue";
import {HttpRequest, HttpResponse} from "../lib/Http";
import HttpInspector from "./HttpInspector.vue";
import {EventsOn} from "../../wailsjs/runtime";
import RightPop from "./RightPop.vue";
import { BarsArrowDownIcon, CircleStackIcon, BeakerIcon, BriefcaseIcon} from "@heroicons/vue/20/solid";
import Search from "./Search.vue";
import {Criteria}  from "../lib/Criteria";
import WorkspaceMenu from "./WorkspaceMenu.vue";
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
  },
  emits: ['switchWorkspace'],
  data() {
    return {
      requests: Array<HttpRequest>(),
      req: null as HttpRequest | null,
      tabs: [
        { name: 'Log Stream', id: 'log', icon: BarsArrowDownIcon, current: true },
        { name: 'Saved Requests', id: 'saved', icon: CircleStackIcon, current: false },
        { name: 'Attack Workflows', id: 'workflows', icon: BeakerIcon, current: false },
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
      this.$emit("switchWorkspace")
    },
  },
  components: {WorkspaceMenu, Search, RightPop, HttpInspector, TrafficLog, BriefcaseIcon}
})
</script>

<template>
  <div class="min-h-16 h-16 max-h-16 flex p-2">
    <div class="flex-grow text-left">
      <Search :on-search="onSearch" :query="liveCriteria.Raw"/>
    </div>
    <div class="flex-shrink p-0 ml-2">
      <WorkspaceMenu :on-workspace-config="onWorkspaceConfig" :ws="ws" @switchWorkspace="switchWorkspace"/>
    </div>
  </div>
  <div class="min-h-16 h-16 max-h-16 px-2">
    <div class="sm:hidden">
      <label for="tabs" class="sr-only">Select a tab</label>
      <select @change="selectTab"  id="tabs" name="tabs" class="block w-full rounded-md text-gray-500 focus:border-frost focus:ring-frost">
        <option v-for="tab in tabs" :key="tab.id" :selected="tab.current" :value="tab.id">{{ tab.name }}</option>
      </select>
    </div>
    <div class="hidden sm:block">
      <div class="border-b dark:border-polar-night-4">
        <nav class="-mb-px flex space-x-8" aria-label="Tabs">
          <a v-for="tab in tabs" @click="switchTab(tab.id)" :key="tab.name" :class="[tab.current ? 'border-frost text-frost' : 'border-transparent text-gray-400 hover:text-gray-200 hover:border-gray-500', 'group inline-flex items-center py-4 px-1 border-b-2 font-medium text-sm']" :aria-current="tab.current ? 'page' : undefined">
            <component :is="tab.icon" :class="[tab.current ? 'text-frost' : 'text-gray-400 group-hover:text-gray-300', '-ml-0.5 mr-2 h-5 w-5']" aria-hidden="true" />
            <span>{{ tab.name }}</span>
          </a>
        </nav>
      </div>
    </div>
  </div>
  <div class="body overflow-y-auto px-2">
    <!-- TODO: show intercepted requests here in dialog -->
    <TrafficLog :key="liveCriteria.Raw" v-if="selectedTab() == 'log'" :proxy-address="proxyAddress" :requests="requests" :onSelect="examineRequest" :selected="req ? req.ID : -1" :criteria="liveCriteria"/>
    <RightPop :show="req !== null" :onRequestClose="clearRequest">
      <HttpInspector v-if="req" :request="req" />
    </RightPop>

    <div v-if="selectedTab() == 'saved'">
      TODO: Saved Requests
    </div>
    <div v-if="selectedTab() == 'workflows'">
      TODO: Attack Workflows
    </div>
  </div>
</template>

<style scoped>
.body{
  max-height: calc(100vh - 8rem);
}
a{
  cursor: pointer;
}
</style>