<script lang="ts">
import {defineComponent, PropType} from 'vue'
import {HttpRequest, HttpResponse} from '../lib/Http.js';
import HttpStatus from "./HttpStatus.vue";
import {Criteria}  from "../lib/Criteria";

export default /*#__PURE__*/ defineComponent({
  components: {HttpStatus},
  // type inference enabled
  props: {
    requests: {type: Array<HttpRequest>, required: true},
    onSelect: {type: Function},
    selected: {type: Number},
    criteria: {type: Object as PropType<Criteria>, required: true},
    proxyAddress: {type: String, required: true},
  },
  methods: {
    filterRequests: function(requests: Array<HttpRequest>): Array<HttpRequest> {
      return requests.filter((request) => {
        return this.criteria.Match(request)
      })
    },
    selectRequest(request: HttpRequest | null) {
      if (this.onSelect !== undefined) {
        this.onSelect(request)
      }
    },
    classFromMethod(method: string) {
      switch(method) {
        case 'GET':
          return 'bg-frost-3'
        case 'POST':
          return 'bg-frost-1'
        case 'PUT':
          return 'bg-frost-4'
        case 'HEAD':
          return 'bg-frost-2'
        case 'DELETE':
          return 'bg-aurora-4'
        case 'CONNECT':
          return 'bg-aurora-3'
        case 'OPTIONS':
          return 'bg-aurora-2'
        case 'TRACE':
          return 'bg-aurora-1'
        case 'PATCH':
          return 'bg-aurora-4'
        default:
          return ''
      }
      if (method === 'GET') {
        return 'bg-frost-4'
      } else if (method === 'POST') {
        return 'bg-frost-3'
      } else if (method === 'PUT') {
        return 'bg-frost-2'
      } else if (method === 'DELETE') {
        return 'bg-frost-1'
      }
      return 'bg-aurora-1'
    },
    classFromStatus(status: number | undefined) {
      if (status === undefined) {
        return 'bg-polar-night-1'
      }
      if (status >= 200 && status < 300) {
        return 'bg-aurora-4'
      } else if (status >= 300 && status < 400) {
        return 'bg-aurora-3'
      } else if (status >= 400 && status < 500) {
        return 'bg-aurora-2'
      } else if (status >= 500 && status < 600) {
        return 'bg-aurora-1'
      }
      return ''
    },
  },
})
</script>

<script lang="ts" setup>
import {RocketLaunchIcon, MagnifyingGlassCircleIcon} from "@heroicons/vue/20/solid";
</script>

<template>
  <div v-if="requests.length === 0">
    <div class="text-center pt-8 pl-8">
      <RocketLaunchIcon class="mx-auto h-12 w-12"/>
      <h3 class="mt-2 text-sm font-medium">All Systems Go!</h3>
      <p class="mt-1 text-sm text-snow-storm-1">Reaper is ready to receive requests at {{proxyAddress}}</p>
    </div>
  </div>
  <div v-else-if="filterRequests(requests).length === 0">
    <div class="text-center pt-8 pl-8">
      <MagnifyingGlassCircleIcon class="mx-auto h-12 w-12"/>
      <h3 class="mt-2 text-sm font-medium">No Results</h3>
      <p class="mt-1 text-sm text-snow-storm-1">No requests match your search criteria</p>
    </div>
  </div>
  <div v-else class="sm:rounded-md bg-snow-storm dark:bg-polar-night-1a h-full" >
    <ul role="list" class="divide-y divide-polar-night-3">
      <li v-for="request in filterRequests(requests)"  :key="request.ID">
        <a @click="selectRequest(request)" :class="'block  relative px-4 ' + (request.ID == selected ? 'bg-polar-night-3' : 'hover:bg-gray-50 dark:hover:bg-polar-night-2')">
          <div :class="'left ending '+classFromMethod(request.Method)">{{ request.Method }}</div>
          <div :class="'right ending '+classFromStatus(request.Response?.StatusCode)">{{request.Response ? request.Response.StatusCode :'&nbsp;'}}</div>
          <div class="px-4 py-4 sm:px-6">
            <div class="flex items-center justify-between">
              <p class="truncate text-sm font-medium text-frost">{{ request.Path }}</p>
              <div class="ml-2 flex flex-shrink-0">
                <p class=" px-2 text-xs font-semibold leading-5">
                  text/html
                </p>
              </div>
            </div>
            <div class="mt-2 sm:flex sm:justify-between">
              <div class="sm:flex">
                <p class="flex items-center text-sm text-frost-3">
                  {{ request.Host }}
                </p>
                <p class="mt-2 flex items-center text-sm text-frost-3 sm:mt-0 sm:ml-6">
                  {{ request.QueryString }}
                </p>
              </div>
              <div class="mt-2 flex items-center text-sm text-frost-3 sm:mt-0">
                <p>
                  <span v-if="request.Path.indexOf('a') > -1" class="mr-1 bg-frost-1 text-polar-night-1 inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium">Auth</span>
                  <span :class="classFromMethod(request.Method) + ' text-polar-night-1 inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium'">Something</span>
                </p>
              </div>
            </div>
          </div>
        </a>
      </li>
    </ul>
  </div>
</template>

<style scoped>
li a{
 cursor: pointer;
 border-radius: 6px;
}
.ending {
  position: absolute;
  writing-mode:tb-rl;
  white-space:nowrap;
  display:block;
  bottom:0px;
  height:100%;
  border-radius: 0 6px 6px 0;
}
.ending.left {
  left:0px;
  transform:rotate(180deg);
}
.ending.right {
  right:0px;
}
</style>

