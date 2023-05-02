<script lang="ts" setup>
import {HandRaisedIcon, PencilSquareIcon, ViewfinderCircleIcon} from '@heroicons/vue/20/solid'
import {reactive, ref, PropType} from 'vue'
import ScopeEditor from './ScopeEditor.vue'
import {workspace} from '../../wailsjs/go/models'

const props = defineProps({
  ws: {
    type: Object as PropType<workspace.Workspace>,
    required: true,
  },
})

const openTab = ref('overview')
const tabs = [
  {name: 'Overview', icon: PencilSquareIcon, id: 'overview'},
  {name: 'Scope', icon: ViewfinderCircleIcon, id: 'scope'},
  {name: 'Interception', icon: HandRaisedIcon, id: 'interception'},
]
const modifiedWorkspace = reactive(props.ws)

const emit = defineEmits(['save', 'cancel'])

function saveWorkspace() {
  emit('save', modifiedWorkspace)
}

function setWorkspaceName(event: Event) {
  let name = (event.target as HTMLInputElement).value
  if (name === '') {
    name = 'Untitled Workspace'
  }
  modifiedWorkspace.name = name
}

function setScope(scope: workspace.Scope) {
  modifiedWorkspace.scope = scope
}

function setInterceptionScope(scope: workspace.Scope) {
  modifiedWorkspace.interception_scope = scope
}

function cancel() {
  emit('cancel')
}

function toggleTab(tabId: string) {
  openTab.value = tabId
}
</script>

<template>
  <div>
    <main class="relative text-left">
      <div class="mx-auto max-w-screen-[100%] px-4 pb-6 sm:px-6 lg:px-8 lg:pb-16">
        <div
            class="overflow-hidden rounded-lg bg-snow-storm dark:bg-polar-night shadow text-polar-night dark:text-snow-storm">
          <div class="lg:grid lg:grid-cols-12 lg:divide-y-0 lg:divide-x divide-snow-storm-3 dark:divide-polar-night-3">
            <aside class="py-6 lg:col-span-3">
              <nav class="space-y-1">
                <a
                    @click="toggleTab(tab.id)"
                    v-for="tab in tabs"
                    :key="tab.name"
                    :class="[
                    tab.id === openTab
                      ? 'bg-polar-night-4 border-frost-3'
                      : 'border-transparent hover:bg-polar-night-3',
                    'group border-l-4 px-3 py-2 flex items-center text-sm font-medium cursor-pointer',
                  ]"
                    :aria-current="tab.id === openTab ? 'page' : undefined">
                  <component :is="tab.icon" :class="['flex-shrink-0 -ml-1 mr-3 h-6 w-6']" aria-hidden="true"/>
                  <span class="truncate">{{ tab.name }}</span>
                </a>
              </nav>
            </aside>

            <div class="lg:col-span-9">
              <!-- Workspace overview -->
              <div :class="{ hidden: 'overview' !== openTab }">
                <div class="py-6 px-4 sm:p-6 lg:pb-8">
                  <div>
                    <h2 class="text-lg font-medium leading-6">Overview</h2>
                    <p class="mt-1 text-sm">Modify the name and core details of the workspace.</p>
                  </div>
                  <div class="mt-8">
                    <label for="name" class="block text-sm font-medium text-snow-storm">Name</label>
                    <div class="relative mt-1 rounded-md shadow-sm">
                      <input
                          @change="setWorkspaceName"
                          type="text"
                          name="name"
                          id="name"
                          class="block w-full rounded-md bg-polar-night-4 pr-10 focus:outline-none sm:text-sm"
                          :value="modifiedWorkspace.name"
                          aria-invalid="true"
                          aria-describedby="name-error"/>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Workspace scope -->
              <div :class="{ hidden: 'scope' !== openTab }">
                <div class="py-6 px-4 sm:p-6 lg:pb-8">
                  <div>
                    <h2 class="text-lg font-medium leading-6">Scope</h2>
                    <p class="mt-1 text-sm">Change the workspace scope to laser focus on your target.</p>
                  </div>

                  <div class="mt-8">
                    <ScopeEditor :scope="modifiedWorkspace.scope" @save="setScope"/>
                  </div>
                </div>
              </div>

              <!-- Workspace interception scope -->
              <div :class="{ hidden: 'interception' !== openTab }">
                <div class="py-6 px-4 sm:p-6 lg:pb-8">
                  <div>
                    <h2 class="text-lg font-medium leading-6">Interception</h2>
                    <p class="mt-1 text-sm">Specify which requests you'd like to intercept and modify before
                      sending. By default, no requests will be intercepted.</p>
                  </div>

                  <div class="mt-8">
                    <ScopeEditor :scope="modifiedWorkspace.interception_scope" @save="setInterceptionScope"/>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="divide-y divide-gray-200 pt-6 text-right">
            <div class="px-4 sm:px-6 pb-4">
              <div>
                <button
                    @click="cancel"
                    type="button"
                    class="rounded border border-transparent bg-aurora-1 px-2.5 py-1.5 text-xs font-medium text-snow-storm-3 shadow-sm hover:bg-aurora-5 focus:outline-none">
                  Cancel
                </button>
                <button
                    @click="saveWorkspace"
                    type="button"
                    class="ml-2 rounded border border-transparent bg-aurora-4 px-2.5 py-1.5 text-xs font-medium text-snow-storm-3 shadow-sm hover:bg-aurora-5 focus:outline-none">
                  Save Changes
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
