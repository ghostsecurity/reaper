<template>
  <div :data-collapsed="isCollapsed"
       class="group flex flex-col gap-4 py-2 data-[collapsed=true]:py-2">
    <nav class="grid gap-1 px-2 group-[[data-collapsed=true]]:justify-center group-[[data-collapsed=true]]:px-2">
      <template v-for="(link, index) of links">
        <Tooltip v-if="isCollapsed"
                 :key="`1-${index}`"
                 :delay-duration="0">
          <TooltipTrigger as-child>
            <router-link :to="link.href ?? '#'"
                         :class="cn(
                          buttonVariants({ variant: isActiveRoute(link.href) ? 'default' : 'ghost', size: 'icon' }),
                          'h-9 w-9',
                          isActiveRoute(link.href) && 'dark:bg-muted dark:text-muted-foreground dark:hover:bg-muted dark:hover:text-white',
                        )">
              <Icon :icon="link.icon"
                    class="size-4" />
              <span class="sr-only">{{ link.title }}</span>
            </router-link>
          </TooltipTrigger>
          <TooltipContent side="right"
                          class="flex items-center gap-4">
            {{ link.title }}
            <span v-if="link.label"
                  class="ml-auto text-muted-foreground">
              {{ link.label }}
            </span>
          </TooltipContent>
        </Tooltip>

        <router-link v-else
                     :key="`2-${index}`"
                     :to="link.href ?? '#'"
                     :class="cn(
                      buttonVariants({ variant: isActiveRoute(link.href) ? 'default' : 'ghost', size: 'sm' }),
                      isActiveRoute(link.href) && 'dark:bg-muted dark:text-white dark:hover:bg-muted dark:hover:text-white',
                      'justify-start',
                    )">
          <Icon :icon="link.icon"
                class="mr-2 size-4" />
          {{ link.title }}
          <span v-if="showShortcuts"
                class="py-0.25 ml-2 rounded-sm border border-muted-foreground/25 px-1 text-2xs text-muted-foreground">{{
                  link.shortcut }}</span>
          <span v-if="link.label"
                :class="cn(
                  'ml-auto',
                  isActiveRoute(link.href) && 'text-background dark:text-white',
                )">
            {{ link.label }}
          </span>
        </router-link>
      </template>
      <div class="m-4 text-xs text-muted-foreground">It's lonely here, invite some ghouls to collaborate with.</div>
    </nav>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { Icon } from '@iconify/vue'
import { useRoute } from 'vue-router'
import { cn } from '@/lib/utils'
import { buttonVariants } from '@/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'

export interface LinkProp {
  title: string
  label?: string
  icon: string
  href?: string
  shortcut?: string
}

// interface NavProps {
//   isCollapsed: boolean
//   links: LinkProp[]
// }

// defineProps<NavProps>()
const route = useRoute()
const isCollapsed = ref(false)
const showShortcuts = ref(false)
const isActiveRoute = (href: string | undefined) => {
  if (!href) return false
  return route.path === href
}

const links: LinkProp[] = [
  {
    title: 'Reaper Admin',
    label: '972',
    icon: 'lucide:circle-user-round',
  },
  {
    title: 'Guest 1',
    label: '972',
    icon: 'lucide:user',
  },
  {
    title: 'Guest 2',
    label: '123',
    icon: 'lucide:user',
  },
]
</script>