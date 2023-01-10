



<script lang="ts">
import { defineComponent } from "vue";
import { EventsOn } from '../wailsjs/runtime';
import Settings from "./lib/Settings";
import setDarkMode from "./lib/theme";
import { Criteria } from "./lib/Criteria";
import { workspace } from "../wailsjs/go/models";
import {
  CreateWorkspace,
  GetSettings,
  GetWorkspaces,
  StartProxy,
  StopProxy,
  SetWorkspace,
  SaveWorkspace,
  SaveSettings,
  LoadWorkspace,
  DeleteWorkspace,
} from "../wailsjs/go/app/App";


export default defineComponent({
  components: {
  },
  data: () => ({
    settings: new Settings(),
    currentWorkspace: new workspace.Workspace({}),
    workspaces: [] as workspace.Workspace[],
    loadedSettings: false,
    loadedWorkspaces: false,
    hasWorkspace: false,
    settingsVisible: false,
    workspaceConfigVisible: false,
    sidebar: '',
    nodes: Array<workspace.StructureNode>(),
    criteria: new Criteria(""),
    proxyStatus: false,
    proxyAddress: "",
    proxyMessage: "Starting...",
  }),
  beforeMount() {
    GetSettings().then((settings: Settings) => {
      this.settings = settings;
      this.loadedSettings = true;
      setDarkMode(this.settings.DarkMode);
      GetWorkspaces().then((workspaces) => {
        this.workspaces = workspaces
        this.loadedWorkspaces = true
      })
    });
    EventsOn('TreeUpdate', (nodes: Array<workspace.StructureNode>) => {
      this.nodes = nodes
    })
    EventsOn("ProxyStatusChange", (up: boolean, addr: string, msg: string) => {
      this.proxyStatus = up
      this.proxyAddress = addr
      this.proxyMessage = msg
    })
  },
  methods: {
    isLoaded(): boolean {
      return this.loadedSettings
    },
    saveSettings(settings: Settings) {
      SaveSettings(settings)
    },
    saveWorkspace(ws: workspace.Workspace) {
      SaveWorkspace(ws)
    },
    closeSettings() {
      this.settingsVisible = false
    },
    showSettings() {
      this.settingsVisible = true
    },
    showWorkspaceConfig() {
      this.workspaceConfigVisible = true
    },
    closeWorkspaceConfig() {
      this.workspaceConfigVisible = false
    },
    setSidebar(id: string) {
      this.sidebar = this.sidebar === id ? '' : id;
    },
    onCriteriaChange(criteria: Criteria) {
      this.criteria = criteria
    },
    onStructureSelect(parts: Array<string>) {
      let query = "(host is " + parts[0] + " and path is /" + parts.slice(1).join("/") + ")"
      let crit = new Criteria(query);
      this.criteria = crit;
    },
    selectWorkspace(ws: workspace.Workspace) {
      StopProxy().then(() => {
        SetWorkspace(ws).then(() => {
          StartProxy().then(() => {
            this.currentWorkspace = ws
            this.hasWorkspace = true
          })
        })
      })
      // TODO: handle errors here
    },
    createWorkspace(ws: workspace.Workspace) {
      CreateWorkspace(ws).then((created: workspace.Workspace) => {
        this.selectWorkspace(created)
      })
    },
    switchWorkspace() {
      this.loadedWorkspaces = false
      this.hasWorkspace = false
      GetWorkspaces().then((workspaces) => {
        this.workspaces = workspaces
        this.loadedWorkspaces = true
      })
    },
    editWorkspace(id: string) {
      LoadWorkspace(id).then((ws) => {
        this.currentWorkspace = ws
        this.showWorkspaceConfig()
      })
    },
    deleteWorkspace(id: string) {
      DeleteWorkspace(id).then(() => {
        GetWorkspaces().then((workspaces) => {
          this.workspaces = workspaces
          this.loadedWorkspaces = true
        })
      })
    }
  }

})
</script>

<script lang="ts" setup>
import { FunnelIcon, FolderIcon, CogIcon, BriefcaseIcon } from '@heroicons/vue/24/outline'
import Structure from "./components/Structure.vue";
import Dashboard from "./components/Dashboard.vue";
import SettingsModal from './components/SettingsModal.vue'
import WorkspaceModal from './components/WorkspaceModal.vue'
import WorkspaceSelection from "./components/WorkspaceSelection.vue";
</script>

<template>
  <div v-if="!isLoaded()">
    Loading...
  </div>
  <div v-else-if="!hasWorkspace">
    <WorkspaceSelection :workspaces="workspaces" :onWorkspaceSelect="selectWorkspace"
      :onWorkspaceCreate="createWorkspace" :onWorkspaceConfig="editWorkspace" :onWorkspaceDelete="deleteWorkspace" />
    <WorkspaceModal :show="isLoaded() && workspaceConfigVisible" :onRequestClose="closeWorkspaceConfig"
      :onSave="saveWorkspace" :workspace="currentWorkspace" />
  </div>
  <div v-else>
    <SettingsModal :show="isLoaded() && settingsVisible" :onRequestClose="closeSettings" :onSave="saveSettings"
      :settings="settings" />
    <WorkspaceModal :show="isLoaded() && workspaceConfigVisible" :onRequestClose="closeWorkspaceConfig"
      :onSave="saveWorkspace" :workspace="currentWorkspace" />
    <div class="fixed pt-1 w-10 bg-polar-night-1a h-full">
      <button
        :class="'text-snow-storm-1 hover:bg-polar-night-3 p-1 rounded ' + (sidebar === 'structure' ? 'bg-polar-night-4' : '')"
        @click="setSidebar('structure')">
        <FolderIcon class="h-6 w-6" aria-hidden="true" title="Structure" />
      </button>
      <button
        :class="'text-snow-storm-1 hover:bg-polar-night-3 p-1 rounded ' + (sidebar === 'scope' ? 'bg-polar-night-4' : '')"
        @click="setSidebar('scope')">
        <FunnelIcon class="h-6 w-6" aria-hidden="true" title="Scope" />
      </button>
      <div class="absolute bottom-0 left-1">
        <button class="text-snow-storm-1 hover:bg-polar-night-3 p-1 rounded" title="Workspace"
          @click="showWorkspaceConfig">
          <BriefcaseIcon class="h-6 w-6" aria-hidden="true" title="Workspace" />
        </button>
        <button class="text-snow-storm-1 hover:bg-polar-night-3 p-1 rounded" title="Settings" @click="showSettings">
          <CogIcon class="h-6 w-6" aria-hidden="true" />
        </button>
      </div>
    </div>

    <div class="flex ml-10 h-full">
      <div
        :class="['sidebar', 'overflow-auto', 'pr-12', 'border-l-2', 'border-polar-night-1', 'relative', 'py-1', 'h-screen', 'bg-polar-night-1a', 'flex-none', 'w-fit', 'min-w-[10%]', 'max-w-[25%]', (sidebar !== '' ? '' : 'hidden')]">
        <Structure v-if="sidebar === 'structure'" :expanded="true" :nodes="nodes" :on-select="onStructureSelect" />
        <p v-else>
          not implemented yet
        </p>
      </div>
      <div class="flex-1 w-3/4">
        <Dashboard @switchWorkspace="switchWorkspace" :criteria="criteria" :onCriteriaChange="onCriteriaChange"
          :proxy-address="'127.0.0.1:' + settings.ProxyPort" :onWorkspaceConfig="showWorkspaceConfig"
          :ws="currentWorkspace" />
      </div>
    </div>
  </div>

</template>

<style scoped>
.sidebar {
  resize: horizontal;
}
</style>
