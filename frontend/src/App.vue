



<script lang="ts">
import {defineComponent} from "vue";
import { EventsOn, EventsEmit } from '../wailsjs/runtime';
import Settings from "./lib/Settings";
import setDarkMode from "./lib/theme";
import StructureNode from "./lib/StructureNode";
import {Criteria}  from "./lib/Criteria";
import {Workspace} from "./lib/Workspace";


export default defineComponent({
  components: {
  },
  data: () => ({
    settings: new Settings,
    workspace: {} as Workspace,
    loaded: false,
    settingsVisible: false,
    workspaceConfigVisible: false,
    sidebar: '',
    nodes: Array<StructureNode>(),
    criteria: new Criteria(""),
  }),
  beforeMount() {
    EventsOn('OnSettingsLoad', (settings: Settings) => {
      this.settings = settings
      this.loaded = true
      setDarkMode(settings.DarkMode)
    })
    EventsOn('OnTreeUpdate', (nodes: Array<StructureNode>) => {
      this.nodes = nodes
    })
    EventsEmit("OnAppReady");
  },
  methods: {
    saveSettings(settings: Settings) {
      EventsEmit("OnSettingsSave", settings)
    },
    saveWorkspace(workspace: Workspace) {
      EventsEmit("OnWorkspaceSave", workspace)
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
    onStructureSelect(parts: Array<string>){
      let query = "(host is "  + parts[0] + " and path is /" + parts.slice(1).join("/") + ")"
      let crit = new Criteria(query);
      this.criteria = crit;
    }
  }

})
</script>

<script lang="ts" setup>
import { FunnelIcon, FolderIcon, CogIcon,BriefcaseIcon } from '@heroicons/vue/24/outline'
import Structure from "./components/Structure.vue";
import Dashboard from "./components/Dashboard.vue";
import SettingsModal from './components/SettingsModal.vue'
import WorkspaceModal from './components/WorkspaceModal.vue'
import workspace from "./components/Workspace.vue";
</script>

<template>
  <SettingsModal :show="loaded && settingsVisible" :onRequestClose="closeSettings" :onSave="saveSettings" :settings="settings" />
  <WorkspaceModal :show="loaded && workspaceConfigVisible" :onRequestClose="closeWorkspaceConfig" :onSave="saveWorkspace" :workspace="workspace" />
  <div class="fixed pt-1 w-10 bg-polar-night-1a h-full">
    <button :class="'text-snow-storm-1 hover:bg-polar-night-3 p-1 rounded '+(sidebar === 'structure' ? 'bg-polar-night-4' : '')" @click="setSidebar('structure')">
      <FolderIcon class="h-6 w-6" aria-hidden="true" title="Structure" />
    </button>
    <button :class="'text-snow-storm-1 hover:bg-polar-night-3 p-1 rounded '+(sidebar === 'scope' ? 'bg-polar-night-4' : '')" @click="setSidebar('scope')">
      <FunnelIcon class="h-6 w-6" aria-hidden="true" title="Scope" />
    </button>
  <div class="absolute bottom-0 left-1">
    <button class="text-snow-storm-1 hover:bg-polar-night-3 p-1 rounded" title="Workspace" @click="showWorkspaceConfig">
      <BriefcaseIcon class="h-6 w-6" aria-hidden="true" title="Workspace" />
    </button>
    <button class="text-snow-storm-1 hover:bg-polar-night-3 p-1 rounded" title="Settings" @click="showSettings">
      <CogIcon class="h-6 w-6" aria-hidden="true" />
    </button>
  </div>
  </div>

  <div class="flex ml-10 h-full">
  <div :class="['sidebar', 'overflow-auto', 'pr-12', 'border-l-2', 'border-polar-night-1', 'relative', 'py-1', 'h-full', 'bg-polar-night-1a', 'flex-none', 'w-fit', 'min-w-[10%]', 'max-w-[25%]', (sidebar !== '' ? '' : 'hidden')]">
    <Structure v-if="sidebar === 'structure'" :expanded="true" :nodes="nodes" :on-select="onStructureSelect"/>
    <p v-else>
      not implemented yet
    </p>
  </div>
  <div class="flex-1 w-3/4">
   <Dashboard :criteria="criteria" :onCriteriaChange="onCriteriaChange" :proxy-address="'127.0.0.1:' + settings.ProxyPort" :onWorkspaceConfig="showWorkspaceConfig"/>
  </div>
  </div>

</template>

<style scoped>
.sidebar {
  resize: horizontal;
}
</style>
