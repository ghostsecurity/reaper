<script lang="ts" setup>
import {
  ChevronDownIcon,
  ChevronRightIcon,
  FolderIcon,
  DocumentIcon,
  CodeBracketSquareIcon,
  PhotoIcon,
} from '@heroicons/vue/20/solid'
import { PropType, reactive, ref, watch } from 'vue'
import { workspace } from '../../wailsjs/go/models'

const props = defineProps({
  nodes: {
    type: Array as PropType<Array<workspace.StructureNode>>,
    required: true,
  },
  expanded: {
    type: Boolean,
    required: false,
    default: false,
  },
  hasParent: {
    type: Boolean,
    required: false,
    default: false,
  },
  shrinkIndex: {
    type: Number,
    required: false,
    default: 0,
  },
})

const visible = reactive(new Map<string, boolean>())
const lastShrink = ref(0)

const emit = defineEmits(['select'])

watch(() => props.shrinkIndex, (newVal: number) => {
  if (newVal <= lastShrink.value) {
    return
  }
  lastShrink.value = newVal
  if (props.hasParent) {
    props.nodes.forEach((node) => {
      visible.set(node.name, false)
    })
  }
}, { immediate: true })

function toggle(name: string) {
  visible.set(name, !toggled(name))
  if (!toggled(name)) {
    lastShrink.value += 1
  }
}

function toggled(name: string) {
  if (props.expanded) {
    return visible.get(name) !== false
  }
  return visible.get(name) === true
}

function hasExt(name: string, exts: Array<string>) {
  name = name.toLowerCase()
  for (let i = 0; i < exts.length; i += 1) {
    if (name.endsWith(exts[i])) {
      return true
    }
  }
  return false
}

function isPhoto(name: string) {
  return hasExt(name, ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.tiff', '.webp', '.svg'])
}

function isCode(name: string) {
  return hasExt(name, ['.js', '.json', '.css', '.html', '.htm'])
}

function onNodeSelect(node: workspace.StructureNode) {
  emit('select', [node.name])
}

function onChildSelect(part: string): (parts: Array<string>) => void { // eslint-disable-line no-unused-vars
  return (parts: Array<string>) => {
    emit('select', [part, ...parts])
  }
}
</script>

<template>
  <div v-if="!hasParent && nodes.length === 0">
    <div class="text-center pt-4 pl-8">
      <FolderIcon class="mx-auto h-12 w-12" />
      <h3 class="mt-2 text-sm font-medium">No requests received</h3>
      <p class="mt-1 text-sm text-snow-storm-1">Configure your browser to use Reaper</p>
    </div>
  </div>
  <ul>
    <li v-for="node in nodes" class="whitespace-nowrap text-snow-storm-1" :key="node.id">
      <div class="flex items-center">
        <a @click="toggle(node.name)" @dblclick="onNodeSelect(node)">
          <span v-if="node.children.length === 0" class="w-6 inline-block h-1" />
          <ChevronDownIcon v-else-if="toggled(node.name)" class="w-4 inline text-gray-500" />
          <ChevronRightIcon v-else class="w-4 inline text-gray-500" />
          <FolderIcon v-if="node.children.length > 0" class="text-frost mr-1 w-4 inline" />
          <CodeBracketSquareIcon v-else-if="isCode(node.name)" class="text-frost-3 mr-1 w-4 inline" />
          <PhotoIcon v-else-if="isPhoto(node.name)" class="text-frost-3 mr-1 w-4 inline" />
          <DocumentIcon v-else class="text-frost-3 mr-1 w-4 inline" />
        </a>
        <a @click="onNodeSelect(node)" class="hover:bg-polar-night-3">
          {{ node.name }}
        </a>
      </div>
      <TreeStructure @select="onChildSelect(node.name)($event)" :key="node.name" v-if="toggled(node.name)"
        :nodes="node.children" :expanded="expanded" :hasParent="true" :shrinkIndex="lastShrink" />
    </li>
  </ul>
</template>

<style scoped>
a {
  cursor: pointer;
  user-select: none;
}

ul {
  list-style-type: none;
  display: block;
  clear: both;
  padding: 0;
  margin: 0;
}

li {
  list-style: none;
  text-align: left;
}

li>ul>li {
  padding-left: 1rem;
}

li>a.expand {
  display: none;
}

li.shrunk>a.expand {
  display: inline;
}

li.shrunk>a.expanded {
  display: none;
}
</style>
