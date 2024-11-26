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
            <span v-if="link.title == 'Explore'"
                  class="ml-auto text-muted-foreground">
              {{ exploreStore.numEndpoints > 0 ? exploreStore.numEndpoints : '' }}
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
          <div v-if="link.title == 'Scan'"
               class="-ml-2 rounded-full p-1"
               :class="scanActive ? 'bg-orange-500' : ''" />
          <div v-if="link.title == 'Explore'"
               class="-ml-2 mr-1 rounded-full p-1"
               :class="proxy.enabled ? 'bg-green-400' : ''" />
          <div v-else-if="link.title == 'Crawl'"
               class="-ml-2 mr-1 rounded-full p-1" />
          <div v-else-if="link.title == 'Collaborate'"
               class="-ml-2 mr-1 rounded-full p-1"
               :class="tunnel.enabled ? 'bg-green-400' : ''" />
          <div v-else
               class="-ml-2 mr-1 rounded-full p-1" />
          <Icon :icon="link.icon"
                class="mr-2 size-4" />
          {{ link.title }}
          <span v-if="link.title == 'Explore'"
                :class="cn(
                  'ml-auto',
                  isActiveRoute(link.href) && 'text-background dark:text-white',
                )">
            {{ exploreStore.numEndpoints > 0 ? utils.customNumberFormat(exploreStore.numEndpoints) : '' }}
          </span>
        </router-link>
      </template>
    </nav>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useRoute } from 'vue-router'
import { cn } from '@/lib/utils'
import { buttonVariants } from '@/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { useScanStore } from '@/stores/scan'
import { useExploreStore } from '@/stores/explore'
import { useCollabStore } from '@/stores/collab'
import { useUtilStore } from '@/utils'


export interface LinkProp {
  title: string
  count?: number
  icon: string
  href?: string
}

const utils = useUtilStore()
const exploreStore = useExploreStore()
const proxy = computed(() => exploreStore.proxy)
const collabStore = useCollabStore()
const tunnel = computed(() => collabStore.tunnel)
const scanStore = useScanStore()
const route = useRoute()
const isCollapsed = ref(false)

const scanActive = computed(() => {
  // true if any domain is 'pending', 'scanning', or 'probing'
  return scanStore.domains.some(domain => domain.status === 'pending' || domain.status === 'scanning' || domain.status === 'probing')
})

// TODO: refactor to use dynamic classes
// https://router.vuejs.org/guide/essentials/active-links.html
const isActiveRoute = (href: string | undefined) => {
  if (!href) return false
  return route.path === href
}

const links: LinkProp[] = [
  {
    title: 'Scan',
    icon: 'lucide:binoculars',
    href: '/scan',
  },
  {
    title: 'Explore',
    icon: 'lucide:earth',
    href: '/explore',
  },
  // {
  //   title: 'Crawl',
  //   label: '279',
  //   icon: 'lucide:bug-play',
  //   href: '/crawl',
  // },
  {
    title: 'Replay',
    icon: 'lucide:replace-all',
    href: '/replay',
  },
  {
    title: 'Tests',
    icon: 'lucide:pocket-knife',
    href: '/test',
  },
  // {
  //   title: 'Collaborate',
  //   label: '',
  //   icon: 'lucide:users',
  //   href: '/collaborate'
  // },
  {
    title: 'AI Agent',
    icon: 'lucide:brain-circuit',
    href: '/agent',
  },
  // {
  //   title: 'Logs',
  //   label: '',
  //   icon: 'lucide:scroll-text',
  //   href: '/logs',
  // },
  {
    title: 'Reports',
    icon: 'lucide:scroll-text',
    href: '/reports',
  },
  // {
  //   title: 'Settings',
  //   icon: 'lucide:settings',
  //   href: '/settings',
  // },
]

onMounted(() => {
  exploreStore.proxyStatus()
  collabStore.tunnelStatus()
})
</script>
