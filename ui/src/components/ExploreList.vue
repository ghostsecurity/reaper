<template>
  <ScrollArea class="flex h-screen">
    <div class="flex flex-1 flex-col gap-2 p-4 pt-0">
      <div>
        <ul class="space-y-0.5 text-sm">
          <li v-for="(host, i) in hosts"
              :key="host.name + i">
            <div class="group flex items-center justify-between rounded-md bg-secondary px-2 py-1">
              <span class="flex cursor-pointer items-center gap-2 text-xs font-semibold"
                    @click="searchValue = host.name">{{ host.name }}
                <ScanEyeIcon class="hidden size-4 group-hover:block" />
              </span>
              <span v-if="host.endpoints?.length > 0"
                    class="text-xs font-semibold text-muted-foreground">{{ host.endpoints.length }}</span>
            </div>
            <ul v-if="host.endpoints"
                class="pl-2 pt-1 text-xs">
              <li v-for="(endpoint, j) in host.endpoints"
                  @click="selectedEndpoint = endpoint.id"
                  :key="endpoint.path + j">
                <div class="flex items-center">
                  <svg width="18"
                       height="26"
                       viewBox="0 0 18 44"
                       fill="none"
                       stroke="currentColor"
                       xmlns="http://www.w3.org/2000/svg"
                       class="-top-2 left-4 flex-shrink-0 text-foreground">
                    <path d="M1 -4 V18"
                          stroke-width="2" />
                    <path v-if="j != host.endpoints.length - 1"
                          d="M1 -4 V48"
                          stroke-width="2" />
                    <path d="M1 14V17.5C1 20.2614 3.23858 22.5 6 22.5H15"
                          stroke-width="2" />
                  </svg>
                  <div class="my-1 flex w-full min-w-0 cursor-pointer rounded-sm pr-1 hover:bg-secondary">
                    <div class="w-10 flex-shrink-0">
                      <RequestMethod :code="endpoint.status">
                        {{ endpoint.method }}
                      </RequestMethod>
                    </div>
                    <div class="flex w-full min-w-0 justify-between">
                      <div class="ml-1 truncate text-foreground/80">{{ endpoint.path }}</div>
                      <RequestMethod :code="endpoint.status">
                        {{ endpoint.status }}
                      </RequestMethod>
                    </div>
                  </div>
                </div>
              </li>
            </ul>
          </li>
        </ul>
      </div>
    </div>
  </ScrollArea>
</template>

<script lang="ts" setup>
import { ScrollArea } from '@/components/ui/scroll-area'
import { ScanEyeIcon } from 'lucide-vue-next';
import RequestMethod from '@/components/shared/RequestMethod.vue'
import { type Host } from '@/stores/explore'

interface ExploreListProps {
  hosts: Host[]
}

defineProps<ExploreListProps>()
const searchValue = defineModel('searchValue')

const selectedEndpoint = defineModel<number>('selectedEndpoint', { required: false })
</script>
