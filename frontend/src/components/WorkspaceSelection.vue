<script lang="ts" setup>
import { PropType, reactive, ref } from "vue";
import { workspace } from "../../wailsjs/go/models";
import { PlusCircleIcon, BriefcaseIcon, NoSymbolIcon, PlusIcon } from "@heroicons/vue/20/solid";
import ScopeEditor from "./ScopeEditor.vue";
import ConfirmDialog from "./ConfirmDialog.vue";
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'

const props = defineProps({
  workspaces: {
    type: Array as PropType<workspace.Workspace[]>,
    required: true,
  },
})

const showDelete = ref(false)
const toDelete = ref("")
const creating = ref(false)
const ws = reactive(new workspace.Workspace({
  name: "",
  scope: new workspace.Scope({
    include: [],
    exclude: [],
  }),
}))

const emit = defineEmits(['select', 'create', 'edit', 'delete'])

function selectWorkspace(workspaceId: string) {
  emit('select', workspaceId)
}

function createWorkspace() {
  creating.value = true;
}

function selectNewWorkspace() {
  if (ws.name === "") {
    ws.name = "Untitled Workspace";
  }
  creating.value = false;
  emit('create', ws)
}

function setScope(scope: workspace.Scope) {
  ws.scope = scope
}

function editWorkspace(id: string) {
  emit('edit', id)
}

function deleteWorkspace(id: string) {
  toDelete.value = id;
  showDelete.value = true;
}

function confirmDelete() {
  showDelete.value = false;
  if (toDelete.value !== "") {
    emit('delete', toDelete.value)
  }
}

function cancelDelete() {
  showDelete.value = false;
}
</script>

<template>
  <div v-if="creating" class="max-w-2xl m-auto pt-4 text-left overflow-y-auto max-h-screen">
    <div class="mt-8">
      <label for="name" class="block text-sm font-medium text-snow-storm">Workspace Name</label>
      <div class="relative mt-1 rounded-md shadow-sm">
        <input autofocus autocomplete="off" autocapitalize="off" spellcheck="false" v-model="ws.name" type="text"
          name="name" id="name" class="block w-full rounded-md bg-polar-night-4 pr-10 focus:outline-none sm:text-sm"
          aria-invalid="true" aria-describedby="name-error" />
      </div>
    </div>
    <div class="mt-8 italic text-gray-400 text-sm">
      Use the rules below to define what is in scope for your workspace. If you continue without defining any rules, all
      requests will be treated as being in scope.
    </div>
    <div class="mt-4">
      <ScopeEditor :scope="ws.scope" @save="setScope" :allowSimpleView="true" />
    </div>

    <div class="divide-y divide-gray-200 pt-6 text-right">
      <div class="pb-4">
        <button @click="selectNewWorkspace"
          class="inline-flex items-center rounded border border-transparent bg-aurora-4 px-2.5 py-1.5 text-xs font-medium text-snow-storm-3 shadow-sm hover:bg-aurora-5 focus:outline-none">
          Create
        </button>
        <button @click="creating = false"
          class="ml-2 inline-flex items-center rounded border border-transparent bg-aurora-1 px-2.5 py-1.5 text-xs font-medium text-snow-storm-3 shadow-sm hover:bg-aurora-5 focus:outline-none">
          Cancel
        </button>
      </div>
    </div>
  </div>
  <div v-else class="h-screen">
    <div class="flex h-full">
      <div class="flex-1 flex h-full">
        <div class="text-center pl-8 m-auto">
          <div class="flex flex-col items-center">
            <div class="flex-shrink-0">
              <PlusCircleIcon class="h-12 w-12 text-frost-3" />
            </div>
            <p class="mt-1 text-base font-medium text-frost-3">
              Create new workspace
            </p>
            <p class="mt-1 text-sm text-gray-500 italic">A whole new world, a new fantastic point of view.</p>
          </div>
          <div class="mt-6">
            <button @click="createWorkspace" type="button"
              class="inline-flex items-center rounded-md border border-transparent bg-frost-3 px-4 py-2 text-sm font-medium text-snow-storm-1 shadow-sm hover:bg-aurora-5 focus:outline-none focus:ring-2 focus:ring-aurora-5 focus:ring-offset-2">
              <PlusIcon class="-ml-1 mr-2 h-5 w-5" aria-hidden="true" />
              New Workspace
            </button>
          </div>
        </div>
      </div>
      <div class="flex-1 border-l border-polar-night-4 text-left px-2 max-h-full overflow-y-auto h-full">
        <div v-if="props.workspaces.length === 0" class="flex-1 flex h-full">
          <div class="text-center pl-8 m-auto">
            <div class="flex flex-col items-center">
              <div class="flex-shrink-0">
                <NoSymbolIcon class="h-12 w-12 text-gray-500" />
              </div>
              <p class="mt-1 text-base font-medium text-gray-500">
                No existing workspaces found
              </p>
              <p class="mt-1 text-sm text-gray-500 italic">
                Workspaces you create will show up here.
              </p>
            </div>
          </div>
        </div>
        <ul v-else>
          <li v-for="ws in props.workspaces" :key="ws.id">
            <div class="flex items-center px-2 py-4 hover:bg-frost-4 group">
              <div class="flex-shrink-0">
                <BriefcaseIcon class="h-10 w-10 text-frost-3 group-hover:text-snow-storm-3" />
              </div>
              <a @click="selectWorkspace(ws.id)"
                class="flex-grow items-center px-2 py-4 hover:bg-frost-4 group cursor-pointer">
                <div class="ml-4">
                  <p class="text-base font-medium text-snow-storm-3">
                    {{ ws.name }}
                  </p>
                  <p class="text-sm text-polar-night-4 italic">{{ ws.id }}</p>
                </div>
              </a>
              <Menu as="div" class="relative inline-block text-left">
                <div>
                  <MenuButton
                    class="inline-flex w-full justify-center rounded-md dark:text-snow-storm-1 px-4 py-2 text-sm font-medium text-gray-700 shadow-sm">
                    <svg class="w-4 h-4" aria-hidden="true" fill="currentColor" viewBox="0 0 20 20"
                      xmlns="http://www.w3.org/2000/svg">
                      <path
                        d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z">
                      </path>
                    </svg>
                  </MenuButton>
                </div>

                <transition enter-active-class="transition ease-out duration-100"
                  enter-from-class="transform opacity-0 scale-95" enter-to-class="transform opacity-100 scale-100"
                  leave-active-class="transition ease-in duration-75" leave-from-class="transform opacity-100 scale-100"
                  leave-to-class="transform opacity-0 scale-95">
                  <MenuItems
                    class="absolute right-0 z-10 mt-2 w-35 origin-top-right rounded-md dark:text-snow-storm-1 bg-white dark:bg-gray-700 shadow-lg">
                    <div class="py-1">
                      <MenuItem v-slot="{ active }">
                      <a @click="editWorkspace(ws.id)"
                        :class="[active ? 'bg-gray-100 text-gray-900 dark:bg-gray-600 dark:text-snow-storm-1' : 'text-gray-700 dark:text-snow-storm-1', 'block px-4 py-2 text-sm cursor-pointer']">Edit</a>
                      </MenuItem>
                      <MenuItem v-slot="{ active }">
                      <a @click="deleteWorkspace(ws.id)"
                        :class="[active ? 'bg-gray-100 text-gray-900 dark:bg-gray-600 dark:text-snow-storm-1' : 'text-gray-700 dark:text-snow-storm-1', 'block px-4 py-2 text-sm cursor-pointer']">Delete</a>
                      </MenuItem>
                    </div>
                  </MenuItems>
                </transition>
              </Menu>


            </div>
          </li>
        </ul>
      </div>
    </div>
  </div>

  <ConfirmDialog :show-delete="showDelete" @confirm="confirmDelete" @cancel="cancelDelete" :confirm="`Delete`"
    :cancel="`Cancel`" :title="`Delete Workspace`" :message="`Are you sure you want to delete this workspace?`" />

</template>
