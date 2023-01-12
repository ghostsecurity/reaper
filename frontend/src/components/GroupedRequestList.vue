<script lang="ts">
import {defineComponent, PropType} from 'vue'
import {HttpRequest} from '../lib/Http.js';
import HttpStatus from "./HttpStatus.vue";
import {Criteria} from "../lib/Criteria";
import {QuestionMarkCircleIcon} from "@heroicons/vue/20/solid";
import {workspace} from "../../wailsjs/go/models";
import Group = workspace.Group;
import Request = workspace.Request;

// https://learnvue.co/tutorials/vue-drag-and-drop

export default /*#__PURE__*/ defineComponent({
  components: {HttpStatus},
  props: {
    groups: {type: Array as PropType<Group[]>, required: true},
    onSelect: {type: Function},
    selected: {type: String},
    criteria: {type: Object as PropType<Criteria>, required: true},
    emptyTitle: {type: String, required: false, default: "Nothing found"},
    emptyMessage: {type: String, required: false, default: `There are no requests/groups yet`},
    emptyIcon: {type: Function, required: false, default: QuestionMarkCircleIcon},
  },
  emits: ['request-group-change', 'request-group-create', 'group-order-change', 'unsave-request', 'duplicate-request',
    'request-group-delete', 'request-group-rename', 'request-rename'],
  data() {
    return {
      dropGroup: "",
      dragGroupNest: 0,
      dropRequest: "",
      dragRequestNest: 0,
      draggingRequest: false, // request or group
      shrunkenGroups: new Set<string>(),
      renamingGroup: "",
      renamingRequest: "",
      creatingGroup: false,
    }
  },
  methods: {
    filterRequests: function (requests: Array<Request>): Array<Request> {
      return requests.filter((request) => {
        return this.criteria.Match(request.inner as HttpRequest)
      })
    },
    selectRequest(request: HttpRequest | null) {
      if (this.onSelect !== undefined) {
        this.onSelect(request)
      }
    },
    startRequestDrag(evt: DragEvent, request: Request) {
      if (!evt.dataTransfer) {
        return
      }
      this.draggingRequest = true
      evt.dataTransfer.dropEffect = 'move'
      evt.dataTransfer.effectAllowed = 'move'
      evt.dataTransfer.setData('request-id', request.id)
    },
    startGroupDrag(evt: DragEvent, group: Group) {
      if (!evt.dataTransfer) {
        return
      }
      this.draggingRequest = false
      evt.dataTransfer.dropEffect = 'move'
      evt.dataTransfer.effectAllowed = 'move'
      evt.dataTransfer.setData('group-id', group.id)
    },
    onDrop(evt: DragEvent, group: Group, next: Request | null) {
      if (!evt.dataTransfer) {
        return
      }
      if (this.draggingRequest) {
        const id = evt.dataTransfer.getData('request-id')
        const request = this.groups.flatMap((g) => g.requests).find((r) => r.id === id)
        if (request === undefined) {
          return
        }
        const nextID = next ? next.id : ""
        this.$emit("request-group-change", request, group.id, nextID)
        this.disableRequestDrag(evt as MouseEvent)
      } else {
        // dragging a group
        const groupID = evt.dataTransfer.getData('group-id')
        const sourceGroup = this.groups.find((g) => g.id === groupID)
        if (sourceGroup === undefined) {
          return
        }
        // move sourceGroup to position of group
        this.$emit("group-order-change", sourceGroup.id, group.id)
        this.dropGroup = ""
        this.dragGroupNest = 0
      }
    },
    matchTarget(el: HTMLElement | null, nodeName: string, className: string): boolean {
      if (el === null) {
        return false
      }
      if (el.nodeName.toLowerCase() !== nodeName.toLowerCase()) {
        return false
      }
      if (className !== '' && !el.classList.contains(className)) {
        return false
      }
      return true
    },
    findItem(evt: MouseEvent, nodeName: string, className: string): HTMLElement | null {
      let target = evt.target as HTMLElement
      let count = 0
      while (!this.matchTarget(target, nodeName, className) && count < 10) {
        if (target.parentElement === null) {
          return null
        }
        target = target.parentElement as HTMLElement
        count++
      }
      if (!this.matchTarget(target, nodeName, className)) {
        return null
      }
      return target
    },
    enableRequestDrag(evt: MouseEvent) {
      let li = this.findItem(evt, 'li', 'li-request')
      if (li === null) {
        return
      }
      li.setAttribute("draggable", "true")
    },
    disableRequestDrag(evt: MouseEvent) {
      let li = this.findItem(evt, 'li', 'li-request')
      if (li === null) {
        return
      }
      li.setAttribute("draggable", "false")
      this.dropGroup = ""
      this.dragGroupNest = 0
      this.dropRequest = ""
      this.dragRequestNest = 0
    },
    enableGroupDrag(evt: MouseEvent) {
      let li = this.findItem(evt, 'li', 'li-group')
      if (li === null) {
        return
      }
      li.setAttribute("draggable", "true")
    },
    disableGroupDrag(evt: MouseEvent) {
      let li = this.findItem(evt, 'li', 'li-group')
      if (li === null) {
        return
      }
      li.setAttribute("draggable", "false")
      this.dropGroup = ""
      this.dragGroupNest = 0
      this.dropRequest = ""
      this.dragRequestNest = 0
    },
    dragGroupEnter(groupName: string) {
      if (this.dropGroup === groupName) {
        this.dragGroupNest++
      } else {
        this.dragGroupNest = 1
      }
      this.dropGroup = groupName
    },
    dragGroupLeave(groupName: string) {
      if (this.dropGroup === groupName) {
        this.dragGroupNest--
        if (this.dragGroupNest <= 0) {
          this.dropGroup = ""
          this.dragGroupNest = 0
        }
      }
    },
    dragRequestEnter(id: string) {
      if (this.dropRequest === id) {
        this.dragRequestNest++
      } else {
        this.dragRequestNest = 1
        this.dropRequest = id
      }
    },
    dragRequestLeave(id: string) {
      if (this.dropRequest === id) {
        this.dragRequestNest--
        if (this.dragRequestNest <= 0) {
          this.dropRequest = ""
          this.dragRequestNest = 0
        }
      }
    },
    expandGroup(group: Group, expand: boolean) {
      if (expand) {
        this.shrunkenGroups.delete(group.id)
      } else {
        this.shrunkenGroups.add(group.id)
      }
    },
    deleteGroup(group: Group) {
      this.$emit("request-group-delete", group.id)
    },
    renameGroup(groupId: string, name: string) {
      this.$emit("request-group-rename", groupId, name)
      this.renamingGroup = ""
    },
    renameRequest(requestId: string, name: string) {
      this.$emit("request-rename", requestId, name)
      this.renamingRequest = ""
    },
    createGroup(name: string) {
      this.$emit('request-group-create', name)
      this.creatingGroup = false
    }
  },
})
</script>

<script lang="ts" setup>
import {
  MagnifyingGlassCircleIcon, Bars3Icon, FolderPlusIcon, TrashIcon,
  PencilSquareIcon,
  DocumentDuplicateIcon, DocumentArrowUpIcon, ChevronDownIcon, ChevronRightIcon
} from "@heroicons/vue/20/solid";
import {MethodClass} from "../lib/Http.js";
import {
  FolderIcon,
  FolderOpenIcon,
} from "@heroicons/vue/24/outline";
import InputBox from "./InputBox.vue";
import RequestListItemSummary from "./RequestListItemSummary.vue";
</script>

<template>
  <InputBox v-if="creatingGroup" title="New Group" message="Enter the group name." @cancel="creatingGroup=false"
            @confirm="createGroup($event)"/>
  <InputBox v-if="renamingGroup" title="Rename group" message="Enter the new group name." @cancel="renamingGroup=''"
            @confirm="renameGroup(renamingGroup, $event)"/>
  <InputBox v-else-if="renamingRequest" title="Rename request" message="Enter the new request name."
            @cancel="renamingRequest=''"
            @confirm="renameRequest(renamingRequest, $event)"/>
  <div class="h-10 max-h-10 w-full flex text-left">
    <div class="flex-1">
      <!-- shrink/expand buttons here? -->
    </div>
    <div class="flex-0">
      <button type="button"
              class="inline-flex items-center rounded-md border border-transparent bg-frost-3 px-4 py-2 text-sm font-medium text-polar-night-1 shadow-sm hover:bg-aurora-5 focus:outline-none focus:ring-2 focus:ring-aurora-5 focus:ring-offset-2">
        <DocumentArrowUpIcon class="-ml-1 mr-2 h-5 w-5" aria-hidden="true"/>
        <span class="hidden sm:inline">New Request</span>
      </button>
      <button @click.stop="creatingGroup=true"
              type="button"
              class="ml-1 inline-flex items-center rounded-md border border-transparent bg-frost-3 px-4 py-2 text-sm font-medium text-polar-night-1 shadow-sm hover:bg-aurora-5 focus:outline-none focus:ring-2 focus:ring-aurora-5 focus:ring-offset-2">
        <FolderPlusIcon class="-ml-1 mr-2 h-5 w-5" aria-hidden="true"/>
        <span class="hidden sm:inline">New Group</span>
      </button>

    </div>
  </div>
  <div class="sensible-height overflow-y-auto pt-2">
    <div v-if="groups.length === 0">
      <div class="text-center pt-8 pl-8">
        <component :is="emptyIcon" class="mx-auto h-12 w-12"/>
        <h3 class="mt-2 text-sm font-medium">{{ emptyTitle }}</h3>
        <p class="mt-1 text-sm text-snow-storm-1">{{ emptyMessage }}</p>
      </div>
    </div>
    <div v-else-if="groups.length === 0">
      <div class="text-center pt-8 pl-8">
        <MagnifyingGlassCircleIcon class="mx-auto h-12 w-12"/>
        <h3 class="mt-2 text-sm font-medium">No Results</h3>
        <p class="mt-1 text-sm text-snow-storm-1">No requests match your search criteria</p>
      </div>
    </div>
    <div v-else class="sm:rounded-md bg-snow-storm dark:bg-polar-night-1a h-full">
      <ul role="list" class="divide-y divide-polar-night-3">
        <li class="li-group pt-4 first:pt-0 bg-polar-night-1" v-for="group in groups"
            :key="group.id" @drop="onDrop($event, group, null)" @dragover.prevent
            @dragenter.prevent="dragGroupEnter(group.id)"
            @dragleave.prevent="dragGroupLeave(group.id)"
            @dragstart.stop="startGroupDrag($event, group)"
        >
          <div
              :class="[(!draggingRequest && group.id === dropGroup ? 'border-t-2 border-aurora-5' : '')]">
            <div :class="['flex',!shrunkenGroups.has(group.id) ? 'border-b border-polar-night-4': '' ]">
              <div
                  :class="['flex-1 drag-handle text-left pb-1 mb-1']"
                  @click.prevent.stop
                  @mousedown.stop="enableGroupDrag"
                  @mouseup.stop="disableGroupDrag"
              >
                <a v-if="!shrunkenGroups.has(group.id)" @click="expandGroup(group, false)">
                  <ChevronDownIcon class="w-5 h-5 inline -mt-1"/>
                  <FolderOpenIcon class="w-5 h-5 inline -mt-1"/>
                </a>
                <a v-else @click="expandGroup(group, true)">
                  <ChevronRightIcon class="w-5 h-5 inline -mt-1"/>
                  <FolderIcon class="w-5 h-5 inline -mt-1"/>
                </a>
                <span class="ml-1">{{ group.name ? group.name : "Untitled" }}</span>
                <span
                    class="text-gray-500 text-sm ml-1"> - {{
                    `${filterRequests(group.requests).length} of ${group.requests.length} ${(group.requests.length === 1 ? 'request' : 'requests')}`
                  }}</span>
              </div>
              <div class="flex-0">
                <a @click.stop="renamingGroup=group.id"
                   class="cursor-pointer text-gray-400 hover:text-frost-2">
                  <PencilSquareIcon class="inline w-4 h-4"/>
                </a>
                <a @click.stop="deleteGroup(group)"
                   class="cursor-pointer text-gray-400 hover:text-aurora-1">
                  <TrashIcon class="inline w-4 h-4"/>
                </a>
              </div>
            </div>
            <ul v-if="!shrunkenGroups.has(group.id)" role="list"
                :class="['divide-y divide-polar-night-3',(draggingRequest && dropRequest === '' && group.id === dropGroup ? 'border-t-2 border-aurora-5' : '')]">
              <div v-if="filterRequests(group.requests).length === 0">
                <div class="text-left py-2 pl-8">
                  <h3 class="mt-2 text-sm font-medium">No Requests</h3>
                  <p class="mt-1 text-sm text-snow-storm-1">There are no requests in this group</p>
                </div>
              </div>
              <li class="li-request bg-snow-storm dark:bg-polar-night-1a" v-else
                  v-for="outer in filterRequests(group.requests)" :key="outer.id"
                  @drop.stop="onDrop($event, group, outer)" @dragover.prevent
                  @dragenter.prevent="dragRequestEnter(outer.id)"
                  @dragleave.prevent="dragRequestLeave(outer.id)"
                  @dragstart.stop="startRequestDrag($event, outer)" @dragend="disableRequestDrag"
              >
                <a :class="[
                  'block relative pl-4',
                  outer.id === selected ? 'bg-polar-night-3' : 'hover:bg-gray-50 dark:hover:bg-polar-night-2',
                  draggingRequest && dropRequest === outer.id ? 'border-b-2 border-aurora-5' : ''
                  ]"
                   @click="selectRequest(outer.inner)">
                  <div :class="['left ending truncate', MethodClass(outer.inner)]">{{ outer.inner.Method }}</div>
                  <div class="pl-4 py-4 sm:pl-6">
                    <div class="flex">
                      <div @click.prevent.stop @mousedown.stop="enableRequestDrag" @mouseup.stop="disableRequestDrag"
                           class="flex-0 pl-0 pr-4 m-auto drag-handle">
                        <Bars3Icon class="w-6 h-6"/>
                      </div>
                      <div class="flex-1">
                        <RequestListItemSummary :request="outer.inner" :name="outer.name" :show-tags="false"/>
                      </div>
                      <div class="flex-0 pl-4 pr-2 pt-2 text-gray-400">
                        <a @click.stop="renamingRequest=outer.id" title="Rename"
                           class="cursor-pointer hover:text-frost-2">
                          <PencilSquareIcon class="inline w-6 h-6"/>
                        </a>
                        <a @click.stop="$emit('duplicate-request', outer)"
                           class="cursor-pointer hover:text-aurora-3">
                          <DocumentDuplicateIcon class="inline w-6 h-6"/>
                        </a>
                        <a @click.stop="$emit('unsave-request', outer)"
                           class="cursor-pointer hover:text-aurora-1">
                          <TrashIcon class="inline w-6 h-6"/>
                        </a>
                      </div>
                    </div>
                  </div>
                </a>
              </li>
            </ul>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>

.sensible-height {
  max-height: calc(100vh - 10.5rem);
}

.li-group * {
  pointer-events: none;
}

.li-group a {
  pointer-events: all;
}

.li-group .drag-handle {
  pointer-events: all;
  cursor: move;
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

</style>

