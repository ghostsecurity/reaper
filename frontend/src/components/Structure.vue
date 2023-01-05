<script lang="ts" setup>
import { ChevronDownIcon, ChevronRightIcon, FolderIcon, DocumentIcon, CodeBracketSquareIcon, PhotoIcon } from '@heroicons/vue/20/solid'
import Structure from "./Structure.vue";
import ContextMenuItem from "./ContextMenu/ContextMenuItem.vue";
</script>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import StructureNode from "../lib/StructureNode";
import ContextMenu from "./ContextMenu/ContextMenu.vue";

export default /*#__PURE__*/ defineComponent({
  components: {
    ContextMenu,
  },
  props: {
    nodes: {
      type: Array as PropType<Array<StructureNode>>,
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
    onSelect: {
      type: Function as PropType<(parts: Array<string>) => void>,
      required: false,
    },
  },
  data: () => ({
    visible: new Map<string, boolean>(),
    lastShrink: 0,
  }),
  watch: {
    shrinkIndex: {
      handler: function () {
        if (this.shrinkIndex <= this.lastShrink) {
          return
        }
        this.lastShrink = this.shrinkIndex;
        if (this.hasParent) {
          this.nodes.forEach((node) => {
            this.visible.set(node.Name, false)
          })
        }
      },
      immediate: true,
    },
  },
  methods: {
    toggle(name: string) {
      this.visible.set(name, !this.toggled(name))
      if (!this.toggled(name)) {
        this.lastShrink++
      }
    },
    toggled(name: string) {
      if (this.expanded) {
        return this.visible.get(name) !== false
      }
      return this.visible.get(name) === true
    },
    hasExt(name: string, exts: Array<string>) {
      name = name.toLowerCase()
      for (let i = 0; i < exts.length; i++) {
        if (name.endsWith(exts[i])) {
          return true
        }
      }
      return false
    },
    isPhoto(name: string) {
      return this.hasExt(name, ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.tiff', '.webp', '.svg'])
    },
    isCode(name: string) {
      return this.hasExt(name, ['.js', '.json', '.css', '.html', '.htm'])
    },
    closeMenu() {
      (this.$refs.menu as any).close()
    },
    openMenu(evt: MouseEvent) {
      console.log(this.$refs.menu);
      (this.$refs.menu as any).open(evt, 'something')
    },
    onNodeSelect(node: StructureNode) {
      if (this.onSelect) {
        this.onSelect([node.Name])
      }
    },
    onChildSelect(part: string): (parts: Array<string>) => void {
      return (parts: Array<string>) => {
        if (this.onSelect) {
          this.onSelect([part, ...parts])
        }
      }
    },
  },
})
</script>

<template>
  <ContextMenu ref="menu">
    <template slot-scope="{ contextData }">
      <ContextMenuItem @click.native="closeMenu">
        Action 1
      </ContextMenuItem>
      <ContextMenuItem @click.native="closeMenu">
        Action 2
      </ContextMenuItem>
      <ContextMenuItem @click.native="closeMenu">
        Action 3
      </ContextMenuItem>
      <ContextMenuItem @click.native="closeMenu">
        Action 4
      </ContextMenuItem>
    </template>
  </ContextMenu>
  <div v-if="!hasParent && nodes.length === 0">
    <div class="text-center pt-4 pl-8">
      <FolderIcon class="mx-auto h-12 w-12" />
      <h3 class="mt-2 text-sm font-medium">No requests received</h3>
      <p class="mt-1 text-sm text-snow-storm-1">Configure your browser to use Reaper</p>
    </div>
  </div>
  <ul class="text-xs">
    <li v-for="node in nodes" class="whitespace-nowrap">
      <a @click="toggle(node.Name)" @contextmenu.prevent="openMenu" @dblclick="onNodeSelect(node)">
        <span v-if="node.Children.length === 0" class="w-6 inline-block bg-red h-1" />
        <ChevronDownIcon v-else-if="toggled(node.Name)" class="w-4 inline text-gray-500" />
        <ChevronRightIcon v-else class="w-4 inline text-gray-500" />
        <FolderIcon v-if="node.Children.length > 0" class="text-frost mr-1 w-4 inline" />
        <CodeBracketSquareIcon v-else-if="isCode(node.Name)" class="text-frost-3 mr-1 w-4 inline" />
        <PhotoIcon v-else-if="isPhoto(node.Name)" class="text-frost-3 mr-1 w-4 inline" />
        <DocumentIcon v-else class="text-frost-3 mr-1 w-4 inline" />
      </a>
      <a @click="onNodeSelect(node)" class="hover:bg-polar-night-3">
        {{ node.Name }}
      </a>
      <Structure :on-select="onChildSelect(node.Name)" :key="node.Name" v-if="toggled(node.Name)" :nodes="node.Children"
        :expanded="expanded" :hasParent="true" :shrinkIndex="lastShrink" />
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