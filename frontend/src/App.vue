<script lang="ts">
import {defineComponent} from "vue";
import {EventsOn} from '../wailsjs/runtime';
import Settings from "./lib/Settings";
import setDarkMode from "./lib/theme";
import {Criteria} from "./lib/Criteria";
import {workspace} from "../wailsjs/go/models";
import {
  CreateWorkspace,
  GetSettings,
  GetWorkspaces,
  StartProxy,
  StopProxy,
  SetWorkspace,
  SaveWorkspace,
  LoadWorkspace,
  DeleteWorkspace,
  SaveSettings, GenerateID,
  Confirm, Notify, Warn, Error
} from "../wailsjs/go/app/App";
import {HttpRequest} from "./lib/Http";

export default defineComponent({
  components: {},
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
  computed: {
    savedRequestIds(): string[] {
      let list = [] as string[];
      this.currentWorkspace.collection.groups.forEach((group) => {
        group.requests.forEach((req) => {
          list.push(req.inner.ID)
        })
      })
      return list
    },
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
      this.criteria = new Criteria(query);
    },
    selectWorkspace(ws: workspace.Workspace) {
      StopProxy().then(() => {
        SetWorkspace(ws).then(() => {
          StartProxy().then(() => {
            this.prepareWorkspace(ws)
            this.currentWorkspace = ws
            this.hasWorkspace = true
          })
        })
      })
      // TODO: handle errors here
    },
    prepareWorkspace(ws: workspace.Workspace) {
      // ensure our collection has at least one default group
      if (ws.collection.groups === null) {
        ws.collection.groups = []
      }
      if (ws.collection.groups.length === 0) {
        GenerateID().then((id) => {
          ws.collection.groups.push(new workspace.Group({
            id: id,
            name: "Default",
            requests: [],
          }))
        })
      }
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
    },
    setRequestGroup(request: workspace.Request, groupID: string, nextID: string) {
      const oldGroup = this.currentWorkspace.collection.groups.find((g) => {
        return g.requests.find((r: workspace.Request) => {
          return r.id === request.id
        }) as (workspace.Request | undefined) !== undefined
      })
      if (oldGroup !== undefined) {
        oldGroup.requests = oldGroup.requests.filter((item) => item.id !== request.id)
      }
      const group = this.currentWorkspace.collection.groups.find((g) => g.id === groupID)
      if (group === undefined) {
        return
      }
      let index = group.requests.findIndex((r) => r.id === nextID)
      if (index === -1) {
        index = 0
      } else {
        index++
      }
      group.requests.splice(index, 0, request)
      this.saveWorkspace(this.currentWorkspace)
    },
    createRequestGroup(name: string) {
      GenerateID().then((id) => {
        this.currentWorkspace.collection.groups.splice(0, 0, new workspace.Group({
          id: id,
          name: name,
          requests: [],
        }))
      })
    },
    saveRequest(request: HttpRequest, groupID: string) {
      let group = this.currentWorkspace.collection.groups.find((g) => g.id === groupID)
      if (group === undefined) {
        group = this.currentWorkspace.collection.groups[0]
      }
      GenerateID().then((id) => {
        let wrapped = new workspace.Request({id: id, name: ''})
        wrapped.inner = Object.assign({}, request)
        wrapped.inner.Response = null
        if (group) {
          group.requests.push(wrapped)
        }
        this.saveWorkspace(this.currentWorkspace)
      })
    },
    unsaveRequest(request: HttpRequest | workspace.Request) {
      let id = 'inner' in request ? request.inner.ID : (request as unknown as HttpRequest).ID
      const group = this.currentWorkspace.collection.groups.find((g) => {
        return g.requests.find((r: workspace.Request) => {
          return r.inner.ID === id
        }) as (workspace.Request | undefined) !== undefined
      })
      if (group !== undefined) {
        group.requests = group.requests.filter((item) => item.inner.ID !== id)
      }
      this.saveWorkspace(this.currentWorkspace)
    },
    reorderGroup(fromID: string, toID: string) {

      const group = this.currentWorkspace.collection.groups.find((g) => g.id === fromID)

      // remove from old position
      this.currentWorkspace.collection.groups = this.currentWorkspace.collection.groups.filter((g) => g.id !== fromID)

      // find new position
      let index = this.currentWorkspace.collection.groups.findIndex((g) => g.id === toID)

      if (index === -1) {
        index = 0
      }
      this.currentWorkspace.collection.groups.splice(index, 0, group as workspace.Group)
      this.saveWorkspace(this.currentWorkspace)
    },
    duplicateRequest(request: workspace.Request) {
      const group = this.currentWorkspace.collection.groups.find((g) => {
        return g.requests.find((r: workspace.Request) => {
          return r.id === request.id
        }) as (workspace.Request | undefined) !== undefined
      })
      if (group === undefined) {
        return
      }
      let dupName = request.name.endsWith(" (copy)") ? request.name : request.name + " (copy)"
      GenerateID().then((id) => {
        let wrapped = new workspace.Request({
          id: id,
          name: dupName,
        })
        wrapped.inner = Object.assign({}, request.inner)
        wrapped.inner.ID = id // unlink this from the original request
        group.requests.push(wrapped)
        this.saveWorkspace(this.currentWorkspace)
      })
    },
    deleteRequestGroup(groupId: string) {
      if (this.currentWorkspace.collection.groups.length < 2) {
        Warn("Deletion failed", "Cannot delete this group - there must be at least one group. Try renaming it instead.")
        return
      }
      let group = this.currentWorkspace.collection.groups.find((g) => g.id === groupId)
      if (group === undefined) {
        return
      }
      if (group.requests.length > 0) {
        Confirm("Confirm deletion", `The group '${group.name}' contains ${group.requests.length}. Are you sure you want to delete it?`).then((confirmed) => {
          if (confirmed) {
            this.currentWorkspace.collection.groups = this.currentWorkspace.collection.groups.filter((g) => g.id !== groupId)
            this.saveWorkspace(this.currentWorkspace)
          }
        })
      }
    },
    renameRequestGroup(groupId: string, name: string) {
      let group = this.currentWorkspace.collection.groups.find((g) => g.id === groupId)
      if (group === undefined) {
        return
      }
      group.name = name
    },
    renameRequest(requestId: string, name: string) {
      let request = this.currentWorkspace.collection.groups.flatMap((g) => g.requests).find((r) => r.id === requestId)
      if (request === undefined) {
        return
      }
      request.name = name
    },
  }
})
</script>

<script lang="ts" setup>
import {FunnelIcon, FolderIcon, CogIcon, BriefcaseIcon} from '@heroicons/vue/24/outline'
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
                        :onWorkspaceCreate="createWorkspace" :onWorkspaceConfig="editWorkspace"
                        :onWorkspaceDelete="deleteWorkspace"/>
    <WorkspaceModal :show="isLoaded() && workspaceConfigVisible" :onRequestClose="closeWorkspaceConfig"
                    :onSave="saveWorkspace" :workspace="currentWorkspace"/>
  </div>
  <div v-else class="h-full">
    <SettingsModal :show="isLoaded() && settingsVisible" :onRequestClose="closeSettings" :onSave="saveSettings"
                   :settings="settings"/>
    <WorkspaceModal :show="isLoaded() && workspaceConfigVisible" :onRequestClose="closeWorkspaceConfig"
                    :onSave="saveWorkspace" :workspace="currentWorkspace"/>
    <div class="fixed pt-1 w-10 bg-polar-night-1a h-full">
      <button
          :class="'text-snow-storm-1 hover:bg-polar-night-3 p-1 rounded ' + (sidebar === 'structure' ? 'bg-polar-night-4' : '')"
          @click="setSidebar('structure')">
        <FolderIcon class="h-6 w-6" aria-hidden="true" title="Structure"/>
      </button>
      <button
          :class="'text-snow-storm-1 hover:bg-polar-night-3 p-1 rounded ' + (sidebar === 'scope' ? 'bg-polar-night-4' : '')"
          @click="setSidebar('scope')">
        <FunnelIcon class="h-6 w-6" aria-hidden="true" title="Scope"/>
      </button>
      <div class="absolute bottom-0 left-1">
        <button class="text-snow-storm-1 hover:bg-polar-night-3 p-1 rounded" title="Workspace"
                @click="showWorkspaceConfig">
          <BriefcaseIcon class="h-6 w-6" aria-hidden="true" title="Workspace"/>
        </button>
        <button class="text-snow-storm-1 hover:bg-polar-night-3 p-1 rounded" title="Settings" @click="showSettings">
          <CogIcon class="h-6 w-6" aria-hidden="true"/>
        </button>
      </div>
    </div>

    <div class="flex ml-10 h-full">
      <div
          :class="['sidebar', 'overflow-auto', 'pr-12', 'border-l-2', 'border-polar-night-1', 'relative', 'py-1', 'h-screen', 'bg-polar-night-1a', 'flex-none', 'w-fit', 'min-w-[10%]', 'max-w-[25%]', (sidebar !== '' ? '' : 'hidden')]">
        <Structure v-if="sidebar === 'structure'" :expanded="true" :nodes="nodes" :on-select="onStructureSelect"/>
        <p v-else>
          not implemented yet
        </p>
      </div>
      <div class="flex-1 w-3/4 h-full">
        <Dashboard @save-request="saveRequest" @unsave-request="unsaveRequest" @request-group-change="setRequestGroup"
                   @request-group-create="createRequestGroup" @switch-workspace="switchWorkspace" :criteria="criteria"
                   :onCriteriaChange="onCriteriaChange" :proxy-address="'127.0.0.1:' + settings.ProxyPort"
                   :onWorkspaceConfig="showWorkspaceConfig" :ws="currentWorkspace"
                   @group-order-change="reorderGroup"
                   @duplicate-request="duplicateRequest"
                   @request-group-delete="deleteRequestGroup"
                   @request-group-rename="renameRequestGroup"
                   @request-rename="renameRequest"
                   :saved-request-ids="savedRequestIds"/>
      </div>
    </div>
  </div>

</template>

<style scoped>
.sidebar {
  resize: horizontal;
}
</style>
