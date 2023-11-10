<script lang="ts" setup>
import { PropType } from 'vue'
import { EllipsisVerticalIcon } from '@heroicons/vue/20/solid'
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'
import { Criteria } from '../../lib/Criteria/Criteria'
import { HttpRequest } from '../../lib/api/packaging'

defineProps({
  name: { type: String, required: false, default: '' },
  request: { type: Object as PropType<HttpRequest>, required: true },
  showTags: { type: Boolean, required: false, default: true },
  showResponse: { type: Boolean, required: false, default: true },
  actions: { type: Object as PropType<Map<string, string>>, required: false, default: () => new Map<string, string>() },
})

const emit = defineEmits(['rename', 'action', 'criteria-change'])

function classForTag(tag: string): string {
  const tags = ['bg-aurora-3', 'bg-aurora-4', 'bg-aurora-5', 'bg-frost-1', 'bg-frost-2', 'bg-frost-3', 'bg-frost-4']
  let total = 0
  const chars = [...tag]
  chars.forEach(c => {
    total += c.charCodeAt(0)
  })
  return tags[total % tags.length]
}

function humanSize(size: number): string {
  return Intl.NumberFormat('en', {
    notation: 'compact',
    style: 'unit',
    unit: 'byte',
    maximumSignificantDigits: 2,
    unitDisplay: 'narrow',
  })
    .format(size)
    .replace(/([a-zA-z]+)/, ' $1')
}

function searchTag(tag: string) {
  const q = `tag equals '${tag}'`
  onSearch(new Criteria(q))
}

function onSearch(crit: Criteria) {
  emit('criteria-change', crit)
}
</script>

<template>
  <div class="flex items-center justify-between">
    <div class="relative flex-1">
      <div v-if="name !== ''" class="absolute left-0 top-3 w-full text-center text-sm">
        <a class="rounded-md bg-polar-night-3 px-3 py-1" style="pointer-events: all"
           @click.prevent.stop="emit('rename')">
          {{ name }}
        </a>
      </div>
      <div class="flex items-center justify-between">
        <div class="max-w-4xl flex-1 truncate text-left text-sm font-medium text-frost-4 dark:text-frost">
          {{ request.path }}
          <span class="text-frost-3" v-if="request.query_string !== ''">?{{ request.query_string }}</span>
        </div>
        <div class="flex-0 ml-2 flex text-right">
          <p class="px-2 text-xs font-semibold leading-5">
            <span v-if="showResponse && request.response" class="text-polar-night-1a dark:text-snow-storm-1">
              {{ humanSize(request.response.body_size) }}
            </span>
          </p>
        </div>
      </div>
      <div class="sm:flex sm:justify-between">
        <div class="flex-1">
          <p class="flex items-center text-sm text-frost-3 dark:text-frost-3">
            {{ request.host }}
          </p>
        </div>
        <div v-if="showTags" class="flex-0 mt-2 text-right text-sm text-pink-500 sm:mt-0">
          <a @click="searchTag(tag)" v-for="tag in request.tags" :key="tag"
             :class="['ml-1 rounded-full px-2.5 py-0.5 text-xs font-medium text-polar-night-1', classForTag(tag)]">
            {{ tag }}
          </a>
          <span v-if="showResponse && request.response">
            <span @click="searchTag(tag)" v-for="tag in request.response.tags" :key="tag"
                  :class="['ml-1 rounded-full px-2.5 py-0.5 text-xs font-medium text-polar-night-1', classForTag(tag)]">
              {{ tag }}
            </span>
          </span>
        </div>
      </div>
    </div>
    <div class="flex-0 relative" style="pointer-events: all" @click.prevent.stop>
      <Menu v-if="actions.size > 0" as="div" class="relative inline-block text-left">
        <div>
          <MenuButton
              class="inline-flex w-full justify-center rounded-md p-2 text-sm font-medium text-gray-700 shadow-sm dark:text-snow-storm-1">
            <EllipsisVerticalIcon class="h-4 w-4" aria-hidden="true"/>
          </MenuButton>
        </div>

        <transition enter-active-class="transition ease-out duration-100"
                    enter-from-class="transform opacity-0 scale-95"
                    enter-to-class="transform opacity-100 scale-100" leave-active-class="transition ease-in duration-75"
                    leave-from-class="transform opacity-100 scale-100" leave-to-class="transform opacity-0 scale-95">
          <MenuItems
              class="w-35 right absolute right-0 z-10 mt-2 rounded-md bg-white shadow-lg dark:bg-gray-700 dark:text-snow-storm-1">
            <div class="py-1">
              <MenuItem v-slot="{ active }" v-for="[action, name] in actions" :key="action">
                <a @click="emit('action', action)" :class="[
                active
                  ? 'bg-gray-100 text-gray-900 dark:bg-gray-600 dark:text-snow-storm-1'
                  : 'text-gray-700 dark:text-snow-storm-1',
                'block cursor-pointer px-4 py-2 text-sm',
              ]">
                  {{ name }}
                </a>
              </MenuItem>
            </div>
          </MenuItems>
        </transition>
      </Menu>

    </div>

  </div>
</template>
