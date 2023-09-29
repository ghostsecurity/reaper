<script lang="ts" setup>
import { ref, PropType } from 'vue'
import { workspace } from '../../wailsjs/go/models'
import RuleEditor from './RuleEditor.vue'

const props = defineProps({
  rules: {
    type: Array as PropType<workspace.Rule[]>,
    required: true,
  },
})

const modifiedRules = ref(props.rules)
const id = ref(props.rules.length)
const hasExisting = ref(props.rules.length > 0)
const emit = defineEmits(['save'])

function saveRule(rule: workspace.Rule) {
  let found = false
  modifiedRules.value.forEach((r, i) => {
    if (r.id === rule.id) {
      modifiedRules.value[i] = rule
      found = true
    }
  })
  if (!found) {
    rule.id = id.value
    id.value += 1
    modifiedRules.value.push(rule)
  }
  emit('save', modifiedRules.value)
}

function cancelRule(rule: workspace.Rule, saved: boolean) {
  if (saved) {
    return
  }
  removeRule(rule)
}

function removeRule(rule: workspace.Rule) {
  modifiedRules.value = modifiedRules.value.filter(r => r.id !== rule.id)
  emit('save', modifiedRules.value)
}

function addRule() {
  modifiedRules.value.push(
    new workspace.Rule({
      id: id.value,
      ports: [],
      protocol: '',
      host: '',
      path: '',
    }),
  )
  id.value += 1
}
</script>

<template>
  <div class="mb-6 mt-2 rounded-md border border-polar-night-4 p-2">
    <div v-if="modifiedRules.length === 0" class="border-b border-polar-night-4 pb-2">
      <p class="italic text-gray-500">No rules defined</p>
    </div>
    <ul v-else>
      <li v-for="rule in modifiedRules" :key="rule.id">
        <RuleEditor
          :rule="rule"
          :key="rule.id"
          :saved="hasExisting"
          @save="saveRule"
          @cancel="cancelRule"
          @remove="removeRule" />
      </li>
    </ul>
    <button
      @click="addRule"
      class="mt-2 inline-flex items-center rounded-md border border-transparent bg-aurora-4 px-3 py-2 text-sm font-medium leading-4 text-white shadow-sm hover:bg-aurora-5 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
      Add Rule
    </button>
  </div>
</template>

<style scoped></style>
