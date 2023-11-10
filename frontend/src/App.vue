<script lang="ts" setup>
import { onBeforeMount, reactive, ref } from 'vue'
import { FunnelIcon, FolderIcon, CogIcon, BriefcaseIcon } from '@heroicons/vue/24/outline'
import Client from './lib/api/Client'
import { Settings } from './lib/api/settings'
import setDarkMode from './lib/theme'
import { Criteria } from './lib/Criteria/Criteria'
import TreeStructure from './components/TreeStructure.vue'
import AppDashboard from './components/AppDashboard.vue'
import SettingsModal from './components/SettingsModal.vue'
import WorkspaceModal from './components/WorkspaceModal.vue'
import WorkspaceSelection from './components/WorkspaceSelection.vue'
import { HttpRequest } from './lib/api/packaging'
import { VersionInfo } from './lib/api/api'
import { Workspace, StructureNode, Group, Request } from './lib/api/workspace'
import { WorkflowM } from './lib/api/workflow'
import MessageDialog from './components/MessageDialog.vue'

const settings = reactive({} as Settings)
const currentWorkspace = reactive({} as Workspace)
const workspaces = ref([] as Workspace[])
const loadedSettings = ref(false)
const loadedVersion = ref(false)
const loadedWorkspaces = ref(false)
const hasWorkspace = ref(false)
const settingsVisible = ref(false)
const workspaceConfigVisible = ref(false)
const sidebar = ref('')
const nodes = ref([] as Array<StructureNode>)
const criteria = reactive(new Criteria(''))
const proxyStatus = ref(false)
const proxyAddress = ref('')
const proxyMessage = ref('Starting...')
const versionInfo = ref(<VersionInfo | null>null)
const savedRequestIds = ref([] as string[])
const alertMessage = ref('')

const client = new Client()

function resetSavedIDs() {
  const list = [] as string[]
  currentWorkspace.collection.groups.forEach((group: Group) => {
    group.requests.forEach((req: Request) => {
      list.push(req.inner.id)
    })
  })
  savedRequestIds.value = list
}

onBeforeMount(() => {
  client.Init().then(() => {
    client.OnEvent('NotifyUser', (msg: string) => {
      alertMessage.value = msg
    })
    client.OnEvent('Close', () => {
      document.location.reload()
    })
    client.GetSettings().then((stngs: Settings) => {
      Object.assign(settings, stngs)
      loadedSettings.value = true
      setDarkMode(settings.dark_mode)
      client.GetWorkspaces().then(spaces => {
        workspaces.value = spaces
        loadedWorkspaces.value = true
        client.GetVersionInfo().then((info: VersionInfo) => {
          versionInfo.value = info
          loadedVersion.value = true
          client.GetWorkspace().then((ws: Workspace | null) => {
            if (!ws) {
              return
            }
            prepareWorkspace(ws)
            nodes.value = ws.tree.root.children
            Object.assign(currentWorkspace, ws)
            hasWorkspace.value = true
          })
        })
      })
    })
    client.OnEvent('TreeUpdate', (n: Array<StructureNode>) => {
      nodes.value = n
    })
    client.OnEvent('ProxyStatusChange', (up: boolean, addr: string, msg: string) => {
      proxyStatus.value = up
      proxyAddress.value = addr
      proxyMessage.value = msg
    })
  })
})

function isLoaded(): boolean {
  return loadedSettings.value && loadedVersion.value && loadedWorkspaces.value
}

function closeSettings() {
  settingsVisible.value = false
}

function saveSettings(s: Settings) {
  client.SaveSettings(s)
  closeSettings()
}

function closeWorkspaceConfig() {
  workspaceConfigVisible.value = false
}

let saveTimeout = 0

function saveWorkspace(ws: Workspace) {
  Object.assign(currentWorkspace, ws)
  currentWorkspace.tree.root.children = nodes.value
  // buffer saves to 3 seconds after last activity
  clearTimeout(saveTimeout)
  saveTimeout = setTimeout(() => {
    client.SaveWorkspace(currentWorkspace)
  }, 1000) as unknown as number
  closeWorkspaceConfig()
}

function showSettings() {
  settingsVisible.value = true
}

function showWorkspaceConfig() {
  workspaceConfigVisible.value = true
}

function setSidebar(id: string) {
  sidebar.value = sidebar.value === id ? '' : id
}

function onCriteriaChange(c: Criteria) {
  Object.assign(criteria, c)
}

function setQuery(query: string) {
  Object.assign(criteria, new Criteria(query))
}

function onStructureSelect(parts: Array<string>) {
  const query = `(host is ${parts[0]} and path is /${parts.slice(1).join('/')})`
  setQuery(query)
}

function prepareWorkspace(ws: Workspace) {
  // ensure our collection has at least one default group
  if (ws.collection.groups === null) {
    ws.collection.groups = [] /* eslint-disable-line */
  }
  if (ws.workflows === null) {
    ws.workflows = [] /* eslint-disable-line */
  }
  if (ws.collection.groups.length === 0) {
    client.GenerateID().then(id => {
      ws.collection.groups.push(
          {
            id,
            name: 'Default',
            requests: [],
          } as Group,
      )
      resetSavedIDs()
    })
  }
  return ws
}

function selectWorkspace(ws: Workspace) {
  client.StopProxy().then(() => {
    client.SetWorkspace(ws).then(() => {
      client.StartProxy().then(() => {
        prepareWorkspace(ws)
        nodes.value = ws.tree.root.children
        Object.assign(currentWorkspace, ws)
        hasWorkspace.value = true
      })
    })
  })
  // TODO: handle errors here
}

function selectWorkspaceById(id: string) {
  client.LoadWorkspace(id).then((ws: Workspace | null) => {
    if (ws === null) {
      return
    }
    selectWorkspace(ws)
  })
}

function createWorkspace(ws: Workspace) {
  client.CreateWorkspace(ws).then((created: Workspace | null) => {
    if (created === null) {
      return
    }
    selectWorkspace(created)
  })
}

function switchWorkspace() {
  loadedWorkspaces.value = false
  hasWorkspace.value = false
  client.GetWorkspaces().then(spaces => {
    workspaces.value = spaces
    loadedWorkspaces.value = true
  })
}

function editWorkspace(id: string) {
  client.LoadWorkspace(id).then(ws => {
    Object.assign(currentWorkspace, ws)
    showWorkspaceConfig()
  })
}

function deleteWorkspace(id: string) {
  client.DeleteWorkspace(id).then(() => {
    client.GetWorkspaces().then(spaces => {
      workspaces.value = spaces
      loadedWorkspaces.value = true
    })
  })
}

function setRequestGroup(request: Request, groupID: string, nextID: string) {
  const oldGroup = currentWorkspace.collection.groups.find(
    // TODO: maybe this can be cleaned up?
    // eslint-disable-next-line
      (g: Group) => (g.requests.find((r: Request) => r.id === request.id) as Request | undefined) !== undefined,
  )
  if (oldGroup !== undefined) {
    oldGroup.requests = oldGroup.requests.filter((item: Request) => item.id !== request.id)
  }
  const group = currentWorkspace.collection.groups.find((g: Group) => g.id === groupID)
  if (group === undefined) {
    return
  }
  let index = group.requests.findIndex((r: Request) => r.id === nextID)
  if (index === -1) {
    index = 0
  } else {
    index += 1
  }
  group.requests.splice(index, 0, request)
  saveWorkspace(currentWorkspace)
}

function createRequestGroup(name: string) {
  client.GenerateID().then(id => {
    currentWorkspace.collection.groups.splice(
      0,
      0,
        {
          id,
          name,
          requests: [],
        } as Group,
    )
  })
}

function saveRequest(request: HttpRequest, groupID: string) {
  let group = currentWorkspace.collection.groups.find((g: Group) => g.id === groupID)
  if (!group) {
    // TODO: lint fix?
    ;[group] = currentWorkspace.collection.groups // eslint-disable-line
  }
  client.GenerateID().then(id => {
    const wrapped = { id, name: '' } as Request
    wrapped.inner = JSON.parse(JSON.stringify(request))
    wrapped.inner.response = null
    if (group) {
      if (!group.requests) {
        group.requests = [] /* eslint-disable-line */
      }
      group.requests.push(wrapped)
    }
    saveWorkspace(currentWorkspace)
    resetSavedIDs()
  })
}

function unsaveRequest(request: HttpRequest | Request) {
  const id = 'inner' in request ? request.inner.id : (request as unknown as HttpRequest).id
  const group = currentWorkspace.collection.groups.find(
    (g: Group) => (g.requests.find((r: Request) => r.inner.id === id) as Request | undefined) !== undefined,
  )
  if (group) {
    group.requests = group.requests.filter((item: Request) => item.inner.id !== id)
  }
  saveWorkspace(currentWorkspace)
  resetSavedIDs()
}

function updateRequest(request: HttpRequest) {
  for (let i = 0; i < currentWorkspace.collection.groups.length; i += 1) {
    const group = currentWorkspace.collection.groups[i]
    for (let j = 0; j < group.requests.length; j += 1) {
      if (group.requests[j].inner.id === request.id) {
        group.requests[j].inner = request
        currentWorkspace.collection.groups.splice(i, 1, group)
        saveWorkspace(currentWorkspace)
        return
      }
    }
  }
}

function reorderGroup(fromID: string, toID: string) {
  const group = currentWorkspace.collection.groups.find((g: Group) => g.id === fromID)

  // remove from old position
  currentWorkspace.collection.groups = currentWorkspace.collection.groups.filter((g: Group) => g.id !== fromID)

  // find new position
  let index = currentWorkspace.collection.groups.findIndex((g: Group) => g.id === toID)

  if (index === -1) {
    index = 0
  }
  currentWorkspace.collection.groups.splice(index, 0, group as Group)
  saveWorkspace(currentWorkspace)
}

function duplicateRequest(request: Request) {
  const group = currentWorkspace.collection.groups.find(
    (g: Group) =>
    // TODO: maybe this can be cleaned up?
    // eslint-disable-next-line
          (g.requests.find((r: Request) => r.id === request.id) as Request | undefined) !== undefined,
  )
  if (group === undefined) {
    return
  }
  const dupName = request.name.endsWith(' (copy)') ? request.name : `${request.name} (copy)`
  client.GenerateID().then(id => {
    const wrapped = {
      id,
      name: dupName,
    } as Request
    wrapped.inner = { ...request.inner }
    wrapped.inner.id = id // unlink this from the original request
    group.requests.push(wrapped)
    saveWorkspace(currentWorkspace)
  })
}

function deleteRequestGroup(groupId: string) {
  if (currentWorkspace.collection.groups.length < 2) {
    alertMessage.value = 'Cannot delete this group - there must be at least one group. Try renaming it instead.'
    return
  }
  const group = currentWorkspace.collection.groups.find((g: Group) => g.id === groupId)
  if (group === undefined) {
    return
  }
  if (group.requests.length > 0) {
    if (window.confirm(
      `The group '${group.name}' contains ${group.requests.length}. Are you sure you want to delete it?`,
    )) {
      currentWorkspace.collection.groups = currentWorkspace.collection.groups.filter((g: Group) => g.id !== groupId)
      saveWorkspace(currentWorkspace)
    }
  }
}

function renameRequestGroup(groupId: string, name: string) {
  const group = currentWorkspace.collection.groups.find((g: Group) => g.id === groupId)
  if (!group) {
    return
  }
  group.name = name
}

function renameRequest(requestId: string, name: string) {
  const request = currentWorkspace.collection.groups.flatMap((g: Group) => g.requests)
    .find((r: Request) => r.id === requestId)
  if (!request) {
    return
  }
  request.name = name
}

const workflowId = ref('')

function createWorkflowFromRequest(request: HttpRequest) {
  client.CreateWorkflowFromRequest(request).then((w: WorkflowM | null) => {
    if (w === null) {
      return
    }
    currentWorkspace.workflows.push(w)
    saveWorkspace(currentWorkspace)
    workflowId.value = w.id
  })
}

function sendRequest(request: HttpRequest) {
  client.SendRequest(request)
}
</script>

<template>
  <div class="absolute " v-if="alertMessage">
    <div class="rounded bg-polar-night-1a p-2 text-snow-storm-1">
      <MessageDialog :message="alertMessage" @close="alertMessage = ''"/>
    </div>
  </div>
  <div v-if="!isLoaded()">Loading...</div>
  <div v-else-if="!hasWorkspace">
    <WorkspaceSelection :workspaces="workspaces" @select="selectWorkspaceById" @create="createWorkspace"
                        @edit="editWorkspace" @delete="deleteWorkspace"/>
    <WorkspaceModal :show="isLoaded() && workspaceConfigVisible" @close="closeWorkspaceConfig" @save="saveWorkspace"
                    :ws="currentWorkspace"/>
  </div>
  <div v-else class="h-full">
    <SettingsModal :client="client" :show="isLoaded() && settingsVisible" @close="closeSettings" @save="saveSettings"
                   :settings="settings"
                   :version="versionInfo"/>
    <WorkspaceModal :show="isLoaded() && workspaceConfigVisible" @close="closeWorkspaceConfig" @save="saveWorkspace"
                    :ws="currentWorkspace"/>
    <div class="fixed h-full w-10 bg-polar-night-1a pt-1">
      <button :class="
        'rounded p-1 text-snow-storm-1 hover:bg-polar-night-3 ' + (sidebar === 'structure' ? 'bg-polar-night-4' : '')
      " @click="setSidebar('structure')">
        <FolderIcon class="h-6 w-6" aria-hidden="true" title="Structure"/>
      </button>
      <button :class="
        'rounded p-1 text-snow-storm-1 hover:bg-polar-night-3 ' + (sidebar === 'scope' ? 'bg-polar-night-4' : '')
      " @click="setSidebar('scope')">
        <FunnelIcon class="h-6 w-6" aria-hidden="true" title="Scope"/>
      </button>
      <div class="absolute bottom-0 left-1">
        <button class="rounded p-1 text-snow-storm-1 hover:bg-polar-night-3" title="Workspace"
                @click="showWorkspaceConfig">
          <BriefcaseIcon class="h-6 w-6" aria-hidden="true" title="Workspace"/>
        </button>
        <button class="rounded p-1 text-snow-storm-1 hover:bg-polar-night-3" title="Settings" @click="showSettings">
          <CogIcon class="h-6 w-6" aria-hidden="true"/>
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
        <TreeStructure v-if="sidebar === 'structure'" :expanded="true" :nodes="nodes" @select="onStructureSelect"/>
        <p v-else>not implemented yet</p>
      </div>
      <div class="h-full w-3/4 flex-1">
        <AppDashboard :client="client" :criteria="criteria" :proxy-address="'127.0.0.1:' + settings.proxy_port"
                      :ws="currentWorkspace"
                      :saved-request-ids="savedRequestIds" @save-request="saveRequest" @unsave-request="unsaveRequest"
                      @request-group-change="setRequestGroup" @request-group-create="createRequestGroup"
                      @switch-workspace="switchWorkspace" @criteria-change="onCriteriaChange"
                      @workspace-edit="showWorkspaceConfig"
                      @workspace-save="saveWorkspace" @group-order-change="reorderGroup"
                      @duplicate-request="duplicateRequest"
                      @request-group-delete="deleteRequestGroup" @request-group-rename="renameRequestGroup"
                      @request-rename="renameRequest" @send-request="sendRequest" @update-request="updateRequest"
                      @create-workflow-from-request="createWorkflowFromRequest" :current-workflow-id="workflowId"/>
      </div>
    </div>
  </div>
</template>
