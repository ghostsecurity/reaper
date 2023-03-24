<script lang="ts" setup>
import { PropType } from 'vue'
import { RocketLaunchIcon, MagnifyingGlassCircleIcon, StarIcon } from '@heroicons/vue/20/solid'
import { StarIcon as EmptyStarIcon } from '@heroicons/vue/24/outline'
import { HttpRequest, MethodClass, StatusClass } from '../../lib/Http'
import { Criteria } from '../../lib/Criteria/Criteria'
import RequestItemSummary from './RequestItemSummary.vue'

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

function getActions(request: HttpRequest): Map<string, string> {
  const actions = new Map<string, string>([
    ['something', 'Something'],
  ])
  if (isSaved(request.ID)) {
    actions.set('unsave', 'Unsave')
  } else {
    actions.set('save', 'Save')
  }
  return actions
}

function filterRequests(requests: Array<HttpRequest>): Array<HttpRequest> {
  return requests.filter(request => props.criteria.Match(request))
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

function actionRequest(action: string, r: HttpRequest) {
  switch (action) {
    case 'save':
      saveRequest(r, '')
      break
    case 'unsave':
      unsaveRequest(r)
      break
    default:
      break
  }
}
</script>

<template>
  <div class="overflow-y-auto h-full max-h-full max-w-full">
    <div v-if="requests.length === 0">
      <div class="pt-8 pl-8 text-center text-frost-3">
        <component :is="emptyIcon" class="mx-auto h-12 w-12" />
        <h3 class="mt-2 text-sm font-medium">{{ emptyTitle }}</h3>
        <p class="mt-1 text-sm">{{ emptyMessage }}</p>
      </div>
    </div>
    <div v-else-if="filterRequests(requests).length === 0">
      <div class="pt-8 pl-8 text-center text-frost-3">
        <MagnifyingGlassCircleIcon class="mx-auto h-12 w-12" />
        <h3 class="mt-2 text-sm font-bold">No Results</h3>
        <p class="mt-1 text-sm">No requests match your search criteria</p>
      </div>
    </div>
    <div v-else>
      <ul role="list" class="space-y-1">
        <li class="bg-snow-storm-2 dark:bg-polar-night-1a" v-for="request in filterRequests(requests)" :key="request.ID">
          <a @click="selectRequest(request)" :class="[
            'relative  block px-4 ',
            request.ID == selected
              ? 'bg-snow-storm-1 dark:bg-polar-night-3'
              : 'hover:bg-snow-storm-1 dark:hover:bg-polar-night-2',
          ]">
            <div :class="
              'left ending text-xs font-semibold text-snow-storm dark:text-polar-night ' + MethodClass(request)
            ">
              {{ request.Method }}
            </div>
            <div :class="
              'right ending text-xs font-semibold text-snow-storm dark:text-polar-night ' + StatusClass(request)
            ">
              {{ request.Response ? request.Response.StatusCode : '&nbsp;' }}
            </div>
            <div class="px-2 py-1 sm:px-4 sm:py-2">
              <div class="flex">
                <div class="flex-0 m-auto pl-0 pr-4">
                  <a v-if="isSaved(request.ID)" class="group cursor-pointer" @click.stop="unsaveRequest(request)">
                    <StarIcon class="h-5 w-5 text-aurora-3 group-hover:text-gray-400" />
                  </a>
                  <a v-else class="group cursor-pointer" @click.stop="saveRequest(request, '')">
                    <EmptyStarIcon class="h-5 w-5 text-gray-400 group-hover:text-aurora-3" />
                  </a>
                </div>
                <div class="flex-1">
                  <RequestItemSummary :request="request" :actions="getActions(request)"
                    @action="actionRequest($event, request)" />
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
