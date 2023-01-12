<script lang="ts">
import {defineComponent} from "vue";

import {workspace} from "../../wailsjs/go/models";
import {PropType} from 'vue'

export default /*#__PURE__*/ defineComponent({
  props: {
    rules: {
      type: Array as PropType<workspace.Rule[]>,
      required: true,
    },
  },
  data() {
    return {
      modifiedRules: this.rules,
      id: this.rules.length,
      hasExisting: this.rules.length > 0,
    }
  },
  emits: ['save'],
  methods: {
    saveRule(rule: workspace.Rule) {
      let found = false;
      this.modifiedRules.forEach((r, i) => {
        if (r.id === rule.id) {
          this.modifiedRules[i] = rule
          found = true
        }
      })
      if (!found) {
        rule.id = this.id
        this.id++
        this.modifiedRules.push(rule)
      }
      this.$emit('save', this.modifiedRules)
    },
    cancelRule(rule: workspace.Rule, saved: boolean) {
      if (saved) {
        return
      }
      this.removeRule(rule)
    },
    removeRule(rule: workspace.Rule) {
      this.modifiedRules = this.modifiedRules.filter((r) => {
        return r.id !== rule.id
      })
      this.$emit('save', this.modifiedRules)
    },
    addRule() {
      this.modifiedRules.push(new workspace.Rule({
        id: this.id,
        ports: [],
        protocol: "",
        host: "",
        path: "",
      }))
      this.id++
    },
  }
})
</script>

<script lang="ts" setup>
import RuleEditor from "./RuleEditor.vue";
</script>

<template>
  <div class="mt-2 mb-6 border border-polar-night-4 rounded-md p-2">
    <div v-if="modifiedRules.length===0" class="border-b border-polar-night-4 pb-2">
      <p class="text-gray-500 italic">No rules defined</p>
    </div>
    <ul v-else>
      <li v-for="rule in modifiedRules">
        <RuleEditor :rule="rule" @save="saveRule" @cancel="cancelRule" @remove="removeRule" :key="rule.id"
                    :saved="hasExisting"/>
      </li>
    </ul>
    <button @click="addRule"
            class="mt-2 inline-flex items-center rounded-md border border-transparent bg-aurora-4 px-3 py-2 text-sm font-medium leading-4 text-white shadow-sm hover:bg-aurora-5 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
      Add Rule
    </button>
  </div>
</template>

<style scoped>
</style>