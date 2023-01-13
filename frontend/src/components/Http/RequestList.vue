<script lang="ts" setup>
import { PropType } from 'vue'
import { RocketLaunchIcon, MagnifyingGlassCircleIcon, StarIcon } from '@heroicons/vue/20/solid'
import { StarIcon as EmptyStarIcon } from '@heroicons/vue/24/outline'
import { HttpRequest, MethodClass, StatusClass } from '../../lib/Http'
import { Criteria } from '../../lib/Criteria/Criteria'
import RequestListItemSummary from './RequestItemSummary.vue'

const props = defineProps({
  requests: { type: Array as PropType<HttpRequest[]>, required: true },
  selected: { type: String },
  criteria: { type: Object as PropType<Criteria>, required: true },
  emptyTitle: { type: String, required: false, default: 'All Systems Go!' },
  emptyMessage: { type: String, required: false, default: 'Reaper is ready to receive requests!' },
  emptyIcon: { type: Object, required: false, default: RocketLaunchIcon },
  savedRequestIds: { type: Array as PropType<string[]>, required: false, default: () => [] },
})

const emit = defineEmits(['save-request', 'unsave-request', 'select'])

function filterRequests(requests: Array<HttpRequest>): Array<HttpRequest> {
  return requests.filter((request) => props.criteria.Match(request))
}

function selectRequest(request: HttpRequest | null): void {
  emit('select', request)
}

function isSaved(id: string) {
  return props.savedRequestIds.includes(id)
}
function saveRequest(req: HttpRequest, groupID: string) {
  emit('save-request', req, groupID)
}
function unsaveRequest(req: HttpRequest) {
  emit('unsave-request', req)
}
</script>

<template>
  <div class="sensible-height overflow-y-auto ">
    <div v-if="requests.length === 0">
      <div class="text-center pt-8 pl-8">
        <component :is="emptyIcon" class="mx-auto h-12 w-12" />
        <h3 class="mt-2 text-sm font-medium">{{ emptyTitle }}</h3>
        <p class="mt-1 text-sm text-snow-storm-1">{{ emptyMessage }}</p>
      </div>
    </div>
    <div v-else-if="filterRequests(requests).length === 0">
      <div class="text-center pt-8 pl-8">
        <MagnifyingGlassCircleIcon class="mx-auto h-12 w-12" />
        <h3 class="mt-2 text-sm font-medium">No Results</h3>
        <p class="mt-1 text-sm text-snow-storm-1">No requests match your search criteria</p>
      </div>
    </div>
    <div v-else class="sm:rounded-md bg-snow-storm dark:bg-polar-night-1a h-full">
      <ul role="list" class="divide-y divide-polar-night-3">
        <li v-for="request in filterRequests(requests)" :key="request.ID">
          <a @click="selectRequest(request)" :class="[
            'block  relative px-4 ',
            request.ID == selected ?
              'bg-polar-night-3' :
              'hover:bg-gray-50 dark:hover:bg-polar-night-2'
          ]">
            <div :class="'left ending ' + MethodClass(request)">{{ request.Method }}</div>
            <div :class="'right ending ' + StatusClass(request)">
              {{ request.Response ? request.Response.StatusCode : '&nbsp;' }}
            </div>
            <div class="px-4 py-4 sm:px-6">
              <div class="flex">
                <div class="flex-0 pl-0 pr-4 m-auto">
                  <a v-if="isSaved(request.ID)" class="cursor-pointer group" @click.stop="unsaveRequest(request)">
                    <StarIcon class="h-5 w-5 text-aurora-3 group-hover:text-gray-400" />
                  </a>
                  <a v-else class="cursor-pointer group" @click.stop="saveRequest(request, '')">
                    <EmptyStarIcon class="h-5 w-5 text-gray-400 group-hover:text-aurora-3" />
                  </a>
                </div>
                <div class="flex-1">
                  <RequestListItemSummary :request="request" />
                </div>
              </div>
            </div>
          </a>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.sensible-height {
  max-height: calc(100vh - 8rem);
}

li a {
  cursor: pointer;
  border-radius: 6px;
}

.ending {
  position: absolute;
  writing-mode: tb-rl;
  white-space: nowrap;
  display: block;
  bottom: 0px;
  height: 100%;
  border-radius: 0 6px 6px 0;
}

.ending.left {
  left: 0px;
  transform: rotate(180deg);
}

.ending.right {
  right: 0px;
}
</style>
