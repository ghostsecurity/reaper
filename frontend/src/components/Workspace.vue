<script lang="ts">
import {
  InformationCircleIcon,
  PaintBrushIcon, ServerStackIcon,
  ShieldCheckIcon,
} from '@heroicons/vue/20/solid'
import {defineComponent} from "vue";

import {Workspace} from "../lib/Workspace";
import {EventsEmit, EventsOn} from "../../wailsjs/runtime";
import {PropType} from 'vue'

export default /*#__PURE__*/ defineComponent({
  props: {
    workspace: {type: Object as PropType<Workspace>, required: true},
    onSave: {type: Function, required: true},
    onCancel: {type: Function, required: true},
  },
  data() {
    return {
      openTab: "overview",
      tabs: [
        {name: 'Overview', icon: PaintBrushIcon, id: "overview"},
        {name: 'Scope', icon: ShieldCheckIcon, id: "scope"},
      ],
      subNavigation: [],
      modifiedWorkspace: this.workspace,
    }
  },
  methods: {
    saveWorkspace() {
      this.onSave(this.modifiedWorkspace)
    },
    setWorkspaceName(event: any) {
      this.modifiedWorkspace.Name = event.target.value
    },
    cancel() {
      this.onCancel()
    },
    toggleTab: function (tabId: string) {
      this.openTab = tabId
    },
  }
})
</script>

<script lang="ts" setup>
import {Switch, SwitchDescription, SwitchGroup, SwitchLabel} from '@headlessui/vue'
</script>

<template>
  <div>
    <main class="relative text-left">
      <div class="mx-auto max-w-screen-[100%] px-4 pb-6 sm:px-6 lg:px-8 lg:pb-16">
        <div class="overflow-hidden rounded-lg bg-snow-storm dark:bg-polar-night shadow text-polar-night dark:text-snow-storm">
          <div class="lg:grid lg:grid-cols-12 lg:divide-y-0 lg:divide-x divide-snow-storm-3 dark:divide-polar-night-3">
            <aside class="py-6 lg:col-span-3">
              <nav class="space-y-1">
                <a @click="toggleTab(tab.id)" v-for="tab in tabs" :key="tab.name" :class="[
                    tab.id === openTab ?
                    'bg-polar-night-4 border-frost-3' :
                    'border-transparent hover:bg-polar-night-3',
                   'group border-l-4 px-3 py-2 flex items-center text-sm font-medium']"
                   :aria-current="tab.id === openTab ? 'page' : undefined">
                  <component :is="tab.icon" :class="['flex-shrink-0 -ml-1 mr-3 h-6 w-6']" aria-hidden="true" />
                  <span class="truncate">{{ tab.name }}</span>
                </a>
              </nav>
            </aside>

            <form class="lg:col-span-9" action="#" method="POST">

              <!-- Workspace overview -->
              <div :class="{'hidden': 'overview' !== openTab}">
                <div class="py-6 px-4 sm:p-6 lg:pb-8">
                <div>
                  <h2 class="text-lg font-medium leading-6">Overview</h2>
                  <p class="mt-1 text-sm">Modify the name and core details of the workspace.</p>
                </div>
                <div class="mt-8">
                  <label for="name" class="block text-sm font-medium text-snow-storm">Name</label>
                  <div class="relative mt-1 rounded-md shadow-sm">
                    <input @change="setWorkspaceName" type="text" name="name" id="name" class="block w-full rounded-md bg-polar-night-4 pr-10 focus:outline-none sm:text-sm" :value="modifiedWorkspace.Name" aria-invalid="true" aria-describedby="name-error" />
                  </div>
                </div>
                </div>
              </div>

              <!-- Workspace scope -->
              <div :class="{'hidden': 'scope' !== openTab}">
                <div class="py-6 px-4 sm:p-6 lg:pb-8">
                  <div>
                    <h2 class="text-lg font-medium leading-6 ">Scope</h2>
                    <p class="mt-1 text-sm ">Change the workspace scope to laser focus on your target.</p>
                  </div>

                  <div class="mt-8">
                    TODO: component for modifying scope goes here...
                  </div>

                </div>
              </div>

              <div class="divide-y divide-gray-200 pt-6 text-right">
                <div class="px-4 sm:px-6 pb-4">
                  <div>
                    <button @click="saveWorkspace" type="button" class="inline-flex items-center rounded border border-transparent bg-aurora-4 px-2.5 py-1.5 text-xs font-medium text-snow-storm-3 shadow-sm hover:bg-aurora-5 focus:outline-none">Save Changes</button>
                    <button @click="cancel" type="button" class="ml-2 inline-flex items-center rounded border border-transparent bg-aurora-1 px-2.5 py-1.5 text-xs font-medium text-snow-storm-3 shadow-sm hover:bg-aurora-5 focus:outline-none">Cancel</button>
                  </div>
                </div>
              </div>

            </form>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
a {
  cursor: pointer;
}
</style>