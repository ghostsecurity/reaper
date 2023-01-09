<script lang="ts">
import {
  PencilSquareIcon,
   ViewfinderCircleIcon,
} from '@heroicons/vue/20/solid'
import {defineComponent} from "vue";

import {workspace} from "../../wailsjs/go/models";
import {PropType} from 'vue'

export default /*#__PURE__*/ defineComponent({
  props: {
    scope: {
      type: Object as PropType<workspace.Scope>,
      required: true,
    },
    allowSimpleView: {
      type: Boolean,
      required: false,
      default: false,
    },
  },
  data() {
    return {
      simpleScope: new workspace.Scope({
        include: [],
        exclude: [],
      }),
      advancedScope: this.scope,
      simpleDomains: "",
      includeSubdomains: false,
      showSimple: this.allowSimpleView,
    }
  },
  emits: ['save'],
  methods: {
    saveInclude (include: workspace.Rule[]) {
      this.advancedScope.include = include
      this.$emit('save', this.advancedScope)
    },
    saveExclude (exclude: workspace.Rule[]) {
      this.advancedScope.exclude = exclude
      this.$emit('save', this.advancedScope)
    },
    saveSimple() {
      this.simpleScope.exclude = []
      this.simpleScope.include = []
      let patterns: string[] = []
      this.simpleDomains.split(",").forEach((domain) => {
        domain = domain.trim()
        if (domain === "") {
          return
        }
        patterns.push(escapeForRegExp(domain))
      })
      let pattern = patterns.join("|")
      if (this.includeSubdomains) {
        pattern = `([^\.]+\\.)?${pattern}$`
      }else{
        pattern = `^${pattern}$`
      }
      this.simpleScope.include.push(new workspace.Rule({
        host: pattern,
      }))
      this.$emit('save', this.simpleScope)
    },
  }
})

function escapeForRegExp(target: string) {
  return target.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

</script>

<script lang="ts" setup>
import RulesEditor  from "./RulesEditor.vue";
</script>

<template>
  <div>
    <div v-if="allowSimpleView" class="min-h-16 h-16 max-h-16 px-2">
      <div class="border-b dark:border-polar-night-4">
        <nav class="-mb-px flex space-x-8" aria-label="Tabs">
          <a @click="showSimple=true" :class="[showSimple ? 'border-frost text-frost' : 'border-transparent text-gray-400 hover:text-gray-200 hover:border-gray-500', 'cursor-pointer group inline-flex items-center py-4 px-1 border-b-2 font-medium text-sm']" :aria-current="showSimple ? 'page' : undefined">
            <span>Simple</span>
          </a>
          <a @click="showSimple=false" :class="[!showSimple ? 'border-frost text-frost' : 'border-transparent text-gray-400 hover:text-gray-200 hover:border-gray-500', 'cursor-pointer group inline-flex items-center py-4 px-1 border-b-2 font-medium text-sm']" :aria-current="!showSimple ? 'page' : undefined">
            <span>Advanced</span>
          </a>
        </nav>
      </div>
    </div>
    <div v-if="showSimple">
      <div class="mt-8">
        <label for="domain" class="block text-sm font-medium text-snow-storm">In Scope Domains <span class="text-gray-400">(comma separated, leave blank to allow all domains)</span></label>
        <div class="relative mt-1 rounded-md shadow-sm">
          <input v-model="simpleDomains" type="text" name="domains" id="domains" class="block w-full rounded-md bg-polar-night-4 pr-10 focus:outline-none sm:text-sm" aria-invalid="true" aria-describedby="domains-error" />
        </div>
      </div>
      <div class="mt-8">
        <div class="relative flex items-start">
          <div class="flex h-5 items-center">
            <input v-model="includeSubdomains" @change="saveSimple" id="includeSubdomains" aria-describedby="includeSubdomains-description" name="includeSubdomains" type="checkbox" class="h-4 w-4 rounded border-gray-300 text-frost-1 focus:ring-indigo-500" />
          </div>
          <div class="ml-3 text-sm">
            <label for="includeSubdomains" class="font-medium text-snow-storm-1">Include Subdomains</label>
            <p id="includeSubdomains-description" class="text-gray-400">Include all subdomains of the above domain(s) in scope.</p>
          </div>
        </div>
      </div>
    </div>
    <div v-else>

      <p class="text-sm text-gray-400">A request will be included in scope if ANY <strong><i>include</i></strong> rule matches. A request will be excluded if ANY <strong><i>exclude</i></strong> rule matches. <strong><i>Exclude</i></strong> rules override <strong><i>include</i></strong> rules.</p>

      <h2 class="mt-4 text-lg font-medium leading-6">Include Rules</h2>
      <RulesEditor :rules="advancedScope.include" @save="saveInclude"/>

      <h2 class="text-lg font-medium leading-6">Exclude Rules</h2>
      <RulesEditor :rules="advancedScope.exclude" @save="saveExclude"/>
    </div>
  </div>
</template>

<style scoped>
</style>