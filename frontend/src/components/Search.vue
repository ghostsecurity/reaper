<script lang="ts">
import {defineComponent, PropType} from "vue";
import {Criteria} from "../lib/Criteria";


export default /*#__PURE__*/ defineComponent({
  props: {
    query: {
      type: String,
      required: true,
    },
    onSearch: {type: Function as PropType<(raw: Criteria) => void>, required: false},
  },
  watch: {
    query: {
      handler: function () {
        this.liveQuery = this.query
      },
      immediate: true,
    },
    liveQuery(value) {
      let crit = new Criteria(value);
      if (this.onSearch !== undefined) {
        this.onSearch(crit)
      }
    }
  },
  data: function () {
    return {
      liveQuery: this.query,
    }
  },
  methods: {
    onChange: function (event: Event) {
      let query = (event.target as HTMLInputElement).value;
      let crit = new Criteria(query);
      if (this.onSearch !== undefined) {
        this.onSearch(crit)
      }
    },
  }
})
</script>

<script lang="ts" setup>
import {MagnifyingGlassIcon} from "@heroicons/vue/20/solid";
</script>

<template>
  <div class="relative w-full text-gray-400 focus-within:text-gray-200 inline-block">
      <span class="absolute p-1">
        <button type="submit" class="p-1 focus:outline-none focus:shadow-outline">
          <MagnifyingGlassIcon class="w-6 h-6"/>
        </button>
      </span>
    <input v-model="liveQuery" type="search"
           autocomplete="off" autocapitalize="off" spellcheck="false"
           class="w-full py-2 text-sm text-snow-storm-1 bg-polar-night-2 focus:bg-polar-night-4 focus:text-snow-storm-3 focus:outline-none rounded-md pl-10"
           placeholder="Search...">
  </div>
</template>

<style scoped>
</style>