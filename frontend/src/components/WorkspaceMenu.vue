

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { workspace } from "../../wailsjs/go/models";

export default /*#__PURE__*/ defineComponent({
  props: {
    ws: {
      type: Object as PropType<workspace.Workspace>,
      required: true
    },
    onWorkspaceConfig: {
      type: Function as PropType<() => void>,
      required: true,
    },
  },
  emits: ['switchWorkspace'],
  data() {
    return {
      active: false,
    }
  },
  methods: {
    switchWorkspace() {
      this.$emit('switchWorkspace')
    }
  }
})
</script>

<script lang="ts" setup>
import { ChevronDownIcon, Cog6ToothIcon, ArrowsRightLeftIcon } from "@heroicons/vue/20/solid";
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'
</script>

<template>
  <Menu as="div" class="relative inline-block text-left">
    <div>
      <MenuButton
        class="overflow-y-clip max-h-10 inline-flex w-full justify-center rounded-md border border-2 border-frost-2 bg-frost-2 px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:ring-2 hover:ring-snow-storm-1">
        <i>{{ ws.name }}</i>
        <ChevronDownIcon class="-mr-1 ml-2 h-5 w-5" aria-hidden="true" />
      </MenuButton>
    </div>

    <transition enter-active-class="transition ease-out duration-100" enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100" leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100" leave-to-class="transform opacity-0 scale-95">
      <MenuItems
        class="absolute right-0 z-10 mt-2 w-56 origin-top-right divide-y divide-gray-100 rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
        <div class="py-1">
          <MenuItem v-slot="{ active }">
          <a @click="onWorkspaceConfig"
            :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'group flex items-center px-4 py-2 text-sm']">
            <Cog6ToothIcon class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500" aria-hidden="true" />
            Settings
          </a>
          </MenuItem>
        </div>
        <div class="py-1">
          <MenuItem v-slot="{ active }">
          <a @click="switchWorkspace"
            :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'group flex items-center px-4 py-2 text-sm']">
            <ArrowsRightLeftIcon class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500" aria-hidden="true" />
            Switch workspace...
          </a>
          </MenuItem>
        </div>
      </MenuItems>
    </transition>
  </Menu>
</template>

<style scoped>
a {
  cursor: pointer;
}
</style>