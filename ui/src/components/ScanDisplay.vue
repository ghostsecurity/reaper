<template>
  <div class="flex h-full flex-col">
    <div class="flex items-center p-2">
      <div class="flex items-center gap-2">
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost"
                    size="icon"
                    :disabled="!domain">
              <Binoculars class="size-4" />
              <span class="sr-only">Re-scan domain</span>
            </Button>
          </TooltipTrigger>
          <TooltipContent>Re-scan domain</TooltipContent>
        </Tooltip>
        <Separator orientation="vertical"
                   class="mx-2 h-6" />
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost"
                    size="icon"
                    :disabled="!domain"
                    @click="handleDeleteDomain">
              <Trash2 class="size-4" />
              <span class="sr-only">Delete domain</span>
            </Button>
          </TooltipTrigger>
          <TooltipContent>Delete domain</TooltipContent>
        </Tooltip>
      </div>
    </div>
    <Separator />
    <div v-if="domain"
         class="flex flex-1 flex-col">
      <div class="flex items-start p-4">
        <div class="flex items-start gap-4 text-sm">
          <Avatar>
            <AvatarFallback>
              {{ domainAvatarName }}
            </AvatarFallback>
          </Avatar>
          <div class="grid gap-1">
            <div class="font-semibold">
              {{ domain.name }}
            </div>
            <div class="line-clamp-1 text-xs">
              <span v-if="tab === 'live'">
                {{ utils.customNumberFormat(filteredHosts.length || 0) }} of
              </span>
              <span>
                {{ utils.customNumberFormat(domain.host_count || 0) }} hosts
              </span>
            </div>
            <Tabs default-value="live"
                  v-model="tab"
                  @update:model-value="handleTabChange">
              <div class="flex items-center pt-2">
                <TabsList class="">
                  <TabsTrigger value="live"
                               class="text-xs text-zinc-600 dark:text-zinc-200">
                    Live
                  </TabsTrigger>
                  <TabsTrigger value="all"
                               class="text-xs text-zinc-600 dark:text-zinc-200">
                    All
                  </TabsTrigger>
                </TabsList>
              </div>
            </Tabs>
          </div>
        </div>
        <div v-if="domain.last_scanned_at"
             class="ml-auto text-xs text-muted-foreground">
          last scanned {{ formatDistanceToNow(new Date(domain.last_scanned_at), { addSuffix: true }) }}
        </div>
      </div>
      <Separator />
      <div class="h-screen space-y-2 overflow-y-auto bg-muted/50 p-2 pb-10 text-xs">
        <div v-for="host in filteredHosts"
             :key="host.id"
             class="group flex justify-between rounded-md bg-background p-1 shadow-sm">
          <div class="flex items-center gap-2">
            <div class="m-2 flex flex-col items-center rounded-sm p-1"
                 :class="host.status === 'live' ? ' text-primary' : 'text-secondary'">
              <Hexagon class="size-4" />
              <span class="h-4 py-1 text-2xs font-semibold text-muted-foreground group-hover:text-primary/80">{{
                host.status_code && host.status_code > 0
                  ? host.status_code : '-' }}</span>
            </div>
            <div class="flex space-x-4">
              <div class="w-64">
                <div class="truncate font-medium">{{ host.name }}</div>
                <div class="lowercase text-muted-foreground">{{ host.content_type }}</div>
              </div>
              <div class="ml-8 min-w-24">
                <div class="font-medium">server</div>
                <div class="lowercase text-muted-foreground">{{ host.webserver && host.webserver.length > 0 ?
                  host.webserver : '-' }}</div>
              </div>
              <div class="ml-8 min-w-24">
                <div class="font-medium">source</div>
                <div class="lowercase text-muted-foreground">{{ host.source }}</div>
              </div>
              <div class="ml-8 min-w-24">
                <div class="font-medium">{{ host.cdn_type || 'cloud' }}</div>
                <div class="lowercase text-muted-foreground">{{ host.cdn_name || '-' }}</div>
              </div>
              <div class="ml-8 min-w-24">
                <div class="font-medium">tech</div>
                <div class="lowercase text-muted-foreground">{{ host.tech && host.tech.length > 0 ? host.tech : '-' }}
                </div>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <div class="text-xs text-muted-foreground group-hover:text-foreground">
              {{ formatDistanceToNow(new Date(host.updated_at), { addSuffix: true }) }}
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else
         class="p-8 text-center text-sm font-medium text-muted-foreground">
      No domain selected.
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue'
import { formatDistanceToNow } from 'date-fns'
import { Binoculars, Hexagon, Trash2 } from 'lucide-vue-next'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import {
  Tabs,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs'

import { useUtilStore } from '@/utils'
import { useScanStore } from '@/stores/scan'

const utils = useUtilStore()
const scanStore = useScanStore()
const domain = computed(() => scanStore.selectedDomain)
const hosts = computed(() => scanStore.hosts)

const domainAvatarName = computed(() => {
  return domain.value?.name.substring(0, 2).toUpperCase()
})

const tab = ref('live')

const handleTabChange = (e: string | number) => {
  tab.value = e as string
}

const handleDeleteDomain = () => {
  if (domain.value) {
    scanStore.deleteDomain(domain.value)
  }
}

const filteredHosts = computed(() => {
  if (tab.value === 'live') {
    return hosts.value.filter((host) => host.status === 'live')
  }
  return hosts.value
})

watch(domain, () => {
  if (domain.value) {
    scanStore.getHosts(domain.value)
  }
})
</script>
