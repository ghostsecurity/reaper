<script lang="ts" setup>
import { computed, onBeforeMount, reactive, ref } from "vue";
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
  LoadWorkspace,
  DeleteWorkspace,
  SaveSettings, GenerateID,
  Confirm, Notify, Warn, Error
} from "../wailsjs/go/app/App";
import { FunnelIcon, FolderIcon, CogIcon, BriefcaseIcon } from '@heroicons/vue/24/outline'
import Structure from './components/Structure.vue'
import Dashboard from './components/Dashboard.vue'
import SettingsModal from './components/SettingsModal.vue'
import WorkspaceModal from './components/WorkspaceModal.vue'
import WorkspaceSelection from './components/WorkspaceSelection.vue'
import { HttpRequest } from "./lib/Http";

const settings = reactive(new Settings());
const currentWorkspace = reactive(new workspace.Workspace({}));
const workspaces = ref([] as workspace.Workspace[])
const loadedSettings = ref(false)
const loadedWorkspaces = ref(false)
const hasWorkspace = ref(false)
const settingsVisible = ref(false)
const workspaceConfigVisible = ref(false)
const sidebar = ref('')
const nodes = ref(Array<workspace.StructureNode>())
const criteria = reactive(new Criteria(""))
const proxyStatus = ref(false)
const proxyAddress = ref("")
const proxyMessage = ref("Starting...")

const savedRequestIds = computed(() => {
  let list = [] as string[];
  currentWorkspace.collection.groups.forEach((group) => {
    group.requests.forEach((req) => {
      list.push(req.inner.ID)
    })
  })
  return list
})

onBeforeMount(() => {
  GetSettings().then((stngs: Settings) => {
    Object.assign(settings, stngs);
    loadedSettings.value = true;
    setDarkMode(settings.DarkMode);
    GetWorkspaces().then((spaces) => {
      workspaces.value = spaces
      loadedWorkspaces.value = true
    })
  });
  EventsOn('TreeUpdate', (n: Array<workspace.StructureNode>) => {
    nodes.value = n
  })
  EventsOn("ProxyStatusChange", (up: boolean, addr: string, msg: string) => {
    proxyStatus.value = up
    proxyAddress.value = addr
    proxyMessage.value = msg
  })
})

function isLoaded(): boolean {
  return loadedSettings.value
}

function saveSettings(settings: Settings) {
  SaveSettings(settings)
  closeSettings()
}

function saveWorkspace(ws: workspace.Workspace) {
  Object.assign(currentWorkspace, ws)
  currentWorkspace.tree.root.children = nodes.value
  SaveWorkspace(ws)
  closeWorkspaceConfig()
}

function closeSettings() {
  settingsVisible.value = false
}

function showSettings() {
  settingsVisible.value = true
}

function showWorkspaceConfig() {
  workspaceConfigVisible.value = true
}

function closeWorkspaceConfig() {
  workspaceConfigVisible.value = false
}

function setSidebar(id: string) {
  sidebar.value = sidebar.value === id ? '' : id;
}

function onCriteriaChange(c: Criteria) {
  Object.assign(criteria, c)
}
function onStructureSelect(parts: Array<string>) {
  let query = "(host is " + parts[0] + " and path is /" + parts.slice(1).join("/") + ")"
  Object.assign(criteria, new Criteria(query))
}

function selectWorkspaceById(id: string) {
  LoadWorkspace(id).then((ws) => {
    selectWorkspace(ws)
  })
}

function selectWorkspace(ws: workspace.Workspace) {
  StopProxy().then(() => {
    SetWorkspace(ws).then(() => {
      StartProxy().then(() => {
        prepareWorkspace(ws)
        nodes.value = ws.tree.root.children
        Object.assign(currentWorkspace, ws)
        hasWorkspace.value = true
      })
    })
  })
  // TODO: handle errors here
}
function prepareWorkspace(ws: workspace.Workspace) {
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
}

function createWorkspace(ws: workspace.Workspace) {
  CreateWorkspace(ws).then((created: workspace.Workspace) => {
    selectWorkspace(created)
  })
}

function switchWorkspace() {
  loadedWorkspaces.value = false
  hasWorkspace.value = false
  GetWorkspaces().then((spaces) => {
    workspaces.value = spaces
    loadedWorkspaces.value = true
  })
}

function editWorkspace(id: string) {
  LoadWorkspace(id).then((ws) => {
    Object.assign(currentWorkspace, ws)
    showWorkspaceConfig()
  })
}

function deleteWorkspace(id: string) {
  DeleteWorkspace(id).then(() => {
    GetWorkspaces().then((spaces) => {
      workspaces.value = spaces
      loadedWorkspaces.value = true
    })
  })
}

function setRequestGroup(request: workspace.Request, groupID: string, nextID: string) {
  const oldGroup = currentWorkspace.collection.groups.find((g) => {
    return g.requests.find((r: workspace.Request) => {
      return r.id === request.id
    }) as (workspace.Request | undefined) !== undefined
  })
  if (oldGroup !== undefined) {
    oldGroup.requests = oldGroup.requests.filter((item) => item.id !== request.id)
  }
  const group = currentWorkspace.collection.groups.find((g) => g.id === groupID)
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
  saveWorkspace(currentWorkspace)
}

function createRequestGroup(name: string) {
  GenerateID().then((id) => {
    currentWorkspace.collection.groups.splice(0, 0, new workspace.Group({
      id: id,
      name: name,
      requests: [],
    }))
  })
}

function saveRequest(request: HttpRequest, groupID: string) {
  let group = currentWorkspace.collection.groups.find((g) => g.id === groupID)
  if (group === undefined) {
    group = currentWorkspace.collection.groups[0]
  }
  GenerateID().then((id) => {
    let wrapped = new workspace.Request({ id: id, name: '' })
    wrapped.inner = Object.assign({}, request)
    wrapped.inner.Response = null
    if (group) {
      group.requests.push(wrapped)
    }
    saveWorkspace(currentWorkspace)
  })
}

function unsaveRequest(request: HttpRequest | workspace.Request) {
  let id = 'inner' in request ? request.inner.ID : (request as unknown as HttpRequest).ID
  const group = currentWorkspace.collection.groups.find((g) => {
    return g.requests.find((r: workspace.Request) => {
      return r.inner.ID === id
    }) as (workspace.Request | undefined) !== undefined
  })
  if (group !== undefined) {
    group.requests = group.requests.filter((item) => item.inner.ID !== id)
  }
  saveWorkspace(currentWorkspace)
}

function reorderGroup(fromID: string, toID: string) {

  const group = currentWorkspace.collection.groups.find((g) => g.id === fromID)

  // remove from old position
  currentWorkspace.collection.groups = currentWorkspace.collection.groups.filter((g) => g.id !== fromID)

  // find new position
  let index = currentWorkspace.collection.groups.findIndex((g) => g.id === toID)

  if (index === -1) {
    index = 0
  }
  currentWorkspace.collection.groups.splice(index, 0, group as workspace.Group)
  saveWorkspace(currentWorkspace)
}

function duplicateRequest(request: workspace.Request) {
  const group = currentWorkspace.collection.groups.find((g) => {
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
    saveWorkspace(currentWorkspace)
  })
}

function deleteRequestGroup(groupId: string) {
  if (currentWorkspace.collection.groups.length < 2) {
    Warn("Deletion failed", "Cannot delete this group - there must be at least one group. Try renaming it instead.")
    return
  }
  let group = currentWorkspace.collection.groups.find((g) => g.id === groupId)
  if (group === undefined) {
    return
  }
  if (group.requests.length > 0) {
    Confirm("Confirm deletion", `The group '${group.name}' contains ${group.requests.length}. Are you sure you want to delete it?`).then((confirmed) => {
      if (confirmed) {
        currentWorkspace.collection.groups = currentWorkspace.collection.groups.filter((g) => g.id !== groupId)
        saveWorkspace(currentWorkspace)
      }
    })
  }
}

function renameRequestGroup(groupId: string, name: string) {
  let group = currentWorkspace.collection.groups.find((g) => g.id === groupId)
  if (group === undefined) {
    return
  }
  group.name = name
}

function renameRequest(requestId: string, name: string) {
  let request = currentWorkspace.collection.groups.flatMap((g) => g.requests).find((r) => r.id === requestId)
  if (request === undefined) {
    return
  }
  request.name = name
}
</script>

<template>
  <div v-if="!isLoaded()">Loading...</div>
  <div v-else-if="!hasWorkspace">
    <WorkspaceSelection :workspaces="workspaces" @select="selectWorkspaceById" @create="createWorkspace"
      @edit="editWorkspace" @delete="deleteWorkspace" />
    <WorkspaceModal :show="isLoaded() && workspaceConfigVisible" @close="closeWorkspaceConfig" @save="saveWorkspace"
      :workspace="currentWorkspace" />
  </div>
  <div v-else class="h-full">
    <SettingsModal :show="isLoaded() && settingsVisible" @close="closeSettings" @save="saveSettings"
      :settings="settings" />
    <WorkspaceModal :show="isLoaded() && workspaceConfigVisible" @close="closeWorkspaceConfig" @save="saveWorkspace"
      :workspace="currentWorkspace" />
    <div class="fixed h-full w-10 bg-polar-night-1a pt-1">
      <button :class="
        'rounded p-1 text-snow-storm-1 hover:bg-polar-night-3 ' + (sidebar === 'structure' ? 'bg-polar-night-4' : '')
      " @click="setSidebar('structure')">
        <FolderIcon class="h-6 w-6" aria-hidden="true" title="Structure" />
      </button>
      <button :class="
        'rounded p-1 text-snow-storm-1 hover:bg-polar-night-3 ' + (sidebar === 'scope' ? 'bg-polar-night-4' : '')
      " @click="setSidebar('scope')">
        <FunnelIcon class="h-6 w-6" aria-hidden="true" title="Scope" />
      </button>
      <div class="absolute bottom-0 left-1">
        <button class="rounded p-1 text-snow-storm-1 hover:bg-polar-night-3" title="Workspace"
          @click="showWorkspaceConfig">
          <BriefcaseIcon class="h-6 w-6" aria-hidden="true" title="Workspace" />
        </button>
        <button class="rounded p-1 text-snow-storm-1 hover:bg-polar-night-3" title="Settings" @click="showSettings">
          <CogIcon class="h-6 w-6" aria-hidden="true" />
        </button>
      </div>
    </div>

    <div class="ml-10 flex h-full">
      <div :class="[
        'sidebar',
        'resize-x',
        'overflow-auto',
        'pr-12',
        'border-l-2',
        'border-polar-night-1',
        'relative',
        'py-1',
        'h-screen',
        'bg-polar-night-1a',
        'flex-none',
        'w-fit',
        'min-w-[10%]',
        'max-w-[25%]',
        sidebar !== '' ? '' : 'hidden',
      ]">
        <Structure v-if="sidebar === 'structure'" :expanded="true" :nodes="nodes" @select="onStructureSelect" />
        <p v-else>not implemented yet</p>
      </div>
      <div class="h-full w-3/4 flex-1">
        <Dashboard @save-request="saveRequest" @unsave-request="unsaveRequest" @request-group-change="setRequestGroup"
          @request-group-create="createRequestGroup" @switch-workspace="switchWorkspace" :criteria="criteria"
          @criteria-change="onCriteriaChange" :proxy-address="'127.0.0.1:' + settings.ProxyPort"
          @workspace-edit="showWorkspaceConfig" :ws="currentWorkspace" @group-order-change="reorderGroup"
          @duplicate-request="duplicateRequest" @request-group-delete="deleteRequestGroup"
          @request-group-rename="renameRequestGroup" @request-rename="renameRequest"
          :saved-request-ids="savedRequestIds" />
      </div>
    </div>
  </div>
</template>
