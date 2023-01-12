<script lang="ts">
import {defineComponent} from "vue";

import {workspace} from "../../wailsjs/go/models";
import {PropType} from 'vue'

export default /*#__PURE__*/ defineComponent({
  props: {
    rule: {
      type: Object as PropType<workspace.Rule>,
      required: true,
    },
    saved: {
      type: Boolean,
      required: false,
      default: false,
    }
  },
  emits: ['save', 'cancel', 'remove'],
  data() {
    return {
      modifiedRule: Object.assign({}, this.rule),
      editing: !this.saved,
      savedLocally: this.saved,
    }
  },
  methods: {
    startEdit() {
      this.editing = true
    },
    save() {
      this.$emit('save', Object.assign({}, this.modifiedRule))
      this.editing = false
      this.savedLocally = true
    },
    cancel() {
      this.editing = false
      this.modifiedRule = Object.assign({}, this.rule)
      this.$emit('cancel', Object.assign({}, this.rule), this.savedLocally)
    },
    remove() {
      this.$emit('remove', Object.assign({}, this.modifiedRule))
    },
    portsString(): string {
      return this.modifiedRule.ports.join(",")
    },
    changePorts(e: Event) {
      let raw = (e.target as HTMLInputElement).value;
      if (raw === "") {
        this.modifiedRule.ports = []
        return
      }
      this.modifiedRule.ports = raw.split(",").map((port) => {
        return parseInt(port.trim())
      }).filter((port) => {
        return !isNaN(port)
      })
    },
  }
})
</script>

<script lang="ts" setup>
import {CheckIcon, XMarkIcon, PencilIcon, TrashIcon} from "@heroicons/vue/20/solid";
</script>

<template>
  <div>
    <div v-if="editing" class="border border-frost-4 border-dashed p-3">
      <div>
        <div class="space-y-2 sm:space-y-2">
          <div class="sm:grid sm:grid-cols-2 sm:items-start">
            <label for="first-name" class="block text-sm font-medium sm:mt-px sm:pt-2">Host <span class="text-gray-400">(regular expression, leave blank for any)</span></label>
            <div class="mt-1 sm:col-span-1 sm:mt-0">
              <input autocomplete="off" autocapitalize="off" spellcheck="false" v-model="modifiedRule.host" type="text"
                     name="first-name" id="first-name"
                     class="block w-full max-w-lg rounded-md shadow-sm bg-polar-night-4 focus:border-frost-1 focus:ring-frost-1 text-sm"/>
            </div>
          </div>

          <div class="sm:grid sm:grid-cols-2 sm:items-start">
            <label for="last-name" class="block text-sm font-medium sm:mt-px sm:pt-2">Path <span class="text-gray-400">(regular expression, leave blank for any)</span></label>
            <div class="mt-1 sm:col-span-1 sm:mt-0">
              <input autocomplete="off" autocapitalize="off" spellcheck="false" v-model="modifiedRule.path" type="text"
                     name="last-name" id="last-name"
                     class="block w-full max-w-lg rounded-md shadow-sm bg-polar-night-4 focus:border-frost-1 focus:ring-frost-1 text-sm"/>
            </div>
          </div>

          <div class="sm:grid sm:grid-cols-2 sm:items-start">
            <label for="email" class="block text-sm font-medium sm:mt-px sm:pt-2">Ports <span class="text-gray-400">(comma separated, leave blank for any)</span></label>
            <div class="mt-1 sm:col-span-1 sm:mt-0">
              <input autocomplete="off" autocapitalize="off" spellcheck="false" id="email" name="email" type="text"
                     :value="portsString()" @change="changePorts"
                     class="block w-full max-w-lg rounded-md shadow-sm bg-polar-night-4 focus:border-frost-1 focus:ring-frost-1 text-sm"/>
            </div>
          </div>

          <div class="sm:grid sm:grid-cols-2 sm:items-start">
            <label for="username" class="block text-sm font-medium sm:mt-px sm:pt-2">Protocol</label>
            <div class="mt-1 sm:col-span-1 sm:mt-0">
              <select v-model="modifiedRule.protocol" id="location" name="location"
                      class="block w-full max-w-lg rounded-md shadow-sm bg-polar-night-4 focus:border-frost-1 focus:ring-frost-1 text-sm">
                <option value="">any</option>
                <option value="http">http://</option>
                <option value="https">https://</option>
              </select>
            </div>
          </div>

        </div>
      </div>
      <div class="text-right mt-6">
        <button type="button" @click="save"
                class="inline-flex items-center rounded-md border border-transparent bg-aurora-4 p-1 text-sm font-medium leading-4 text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
          <CheckIcon class="h-4 w-4" aria-hidden="true"/>
        </button>
        <button type="button" @click="cancel"
                class="ml-1 inline-flex items-center rounded-md border border-transparent bg-aurora-1 p-1 text-sm font-medium leading-4 text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
          <XMarkIcon class="h-4 w-4" aria-hidden="true"/>
        </button>
      </div>
    </div>
    <div v-else class="border-b p-1 border-polar-night-4">
      <div class="grid grid-cols-12 my-2">
        <div class="col-span-1 truncate">
          Host
        </div>
        <div class="col-span-7 truncate">
          <code v-if="rule.host!==''"
                class="ml-2 rounded-md bg-polar-night-4 p-1 border border-frost-4">{{ rule.host }}</code>
          <span v-else class="ml-2 text-gray-400 italic">any</span>
        </div>
        <div class="col-span-1 truncate">
          Ports
        </div>
        <div class="col-span-2 truncate">
          <span v-if="rule.ports.length>0" class="ml-2 text-snow-storm-1">{{ rule.ports.join(", ") }}</span>
          <span v-else class="ml-2 text-gray-400 italic">any</span>
        </div>
        <div class="col-span-1 truncate text-right">
          <button type="button" @click="startEdit"
                  class="inline-flex items-center rounded-md border border-transparent bg-frost-3 p-1 text-sm font-medium leading-4 text-white shadow-sm hover:bg-aurora-5 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
            <PencilIcon class="h-4 w-4" aria-hidden="true"/>
          </button>
        </div>
      </div>
      <div class="grid grid-cols-12 my-2">
        <div class="truncate col-span-1 ">
          Path
        </div>
        <div class="truncate col-span-7">
          <code v-if="rule.path!==''"
                class="ml-2 rounded-md bg-polar-night-4 p-1 border border-frost-4">{{ rule.path }}</code>
          <span v-else class="ml-2 text-gray-400 italic">any</span>
        </div>
        <div class="col-span-1 truncate">
          Protocol
        </div>
        <div class="col-span-2 truncate">
          <span v-if="rule.protocol!==''" class="ml-2 text-gray-400">{{ rule.protocol }}</span>
          <span v-else class="ml-2 text-gray-400 italic">any</span>
        </div>
        <div class="col-span-1 truncate text-right">
          <button type="button" @click="remove"
                  class="inline-flex items-center rounded-md border border-transparent bg-aurora-1 p-1 text-sm font-medium leading-4 text-white shadow-sm hover:bg-aurora-5 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
            <TrashIcon class="h-4 w-4" aria-hidden="true"/>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
input {
  text-transform: none !important;
}
</style>