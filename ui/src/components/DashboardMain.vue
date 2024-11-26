<template>
  <div class="hidden h-full md:flex">
    <div class="flex w-64 flex-col">
      <div v-if="projectsEnabled"
           class="p-2">
        <ProjectSwitcher :is-collapsed="false"
                         :projects="projects" />
      </div>
      <div class="flex items-center gap-2 p-3.5 font-bold">
        <GhostLogo class="size-6 text-primary" />
        Reaper
      </div>
      <Separator />
      <SideNav :is-collapsed="false" />
      <div class="mt-auto">
        <Separator />
        <div class="mx-4 flex items-center justify-between gap-4">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger>
                <div class="my-3 cursor-pointer text-xs font-medium"
                     @click="config.toggleThemeDark()">
                  <Icon icon="lucide:moon"
                        class="size-4" />
                </div>
              </TooltipTrigger>
              <TooltipContent>
                <p>Toggle Dark Mode</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger>
                <div class="flex-none rounded-full"
                     :class="wsConnected ? 'bg-green-400/20 p-1 text-green-400' : 'bg-red-500/20 p-1 text-red-500'">
                  <div class="size-2 rounded-full bg-current"></div>
                </div>
              </TooltipTrigger>
              <TooltipContent>
                <p>Live Stream Connected</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          <div>
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger @click="sessionStore.signOut">
                  <Icon icon="lucide:door-open"
                        class="size-4" />
                </TooltipTrigger>
                <TooltipContent>
                  <p>{{ loggedIn ? 'Log out' : 'Log in' }}</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </div>
        </div>
      </div>
    </div>
    <router-view />
  </div>
</template>

<script lang="ts" setup>
import { Icon } from '@iconify/vue'
import ProjectSwitcher from '@/components/ProjectSwitcher.vue'
import SideNav from '@/components/SideNav.vue'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from '@/components/ui/tooltip'
import Separator from '@/components/ui/separator/Separator.vue'
import GhostLogo from '@/components/brand/GhostLogo.vue'

import { useConfigStore } from '@/stores/config'
import { useSessionStore } from '@/stores/session'
import { storeToRefs } from 'pinia'

defineProps<{
  wsConnected: boolean
}>()

const config = useConfigStore()
const sessionStore = useSessionStore()
const { loggedIn } = storeToRefs(sessionStore)
const projectsEnabled = false
const projects = [{ value: 'ghostbank', label: 'Ghostbank', icon: 'lucide:ghost' }]
</script>
