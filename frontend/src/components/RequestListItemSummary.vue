<script lang="ts">
import {defineComponent, PropType} from 'vue'
import {HttpRequest} from '../lib/Http.js';

export default /*#__PURE__*/ defineComponent({
  props: {
    name: {
      type: String,
      required: false,
      default: ""
    },
    request: {
      type: Object as PropType<HttpRequest>,
      required: true
    },
    showTags: {
      type: Boolean,
      required: false,
      default: true
    }
  },
  methods: {
    classForTag(tag: string): string {
      const tags = [
        'bg-aurora-3',
        'bg-aurora-4',
        'bg-aurora-5',
        'bg-frost-1',
        'bg-frost-2',
        'bg-frost-3',
        'bg-frost-4',
      ]
      let total = 0
      for (const c of tag) {
        total += c.charCodeAt(0)
      }
      return tags[total % tags.length]
    },
    humanSize(size: number): string {
      if (size < 1024) {
        return size + " B"
      }
      if (size < 1024 * 1024) {
        return (size / 1024).toFixed(2) + " KB"
      }
      if (size < 1024 * 1024 * 1024) {
        return (size / 1024 / 1024).toFixed(2) + " MB"
      }
      return (size / 1024 / 1024 / 1024).toFixed(2) + " GB"
    }
  }
})
</script>

<template>
  <div class="relative">
    <div v-if="name !== ''" class="absolute left-0 text-center top-3 w-full text-sm">
      <span class="rounded-md bg-polar-night-3 px-3 py-1">{{ name }}</span>
    </div>
    <div class="flex items-center justify-between">
      <p class="truncate text-left  text-sm font-medium text-frost flex-1">
        {{ request.Path }}<span class="text-frost-3 max-w-4xl truncate" v-if="request.QueryString !== ''">?{{
          request.QueryString
        }}</span>
      </p>
      <div class="ml-2 flex flex-0 text-right">
        <p class=" px-2 text-xs font-semibold leading-5">
          <span v-if="request.Response" class="text-snow-storm-1">{{ humanSize(request.Response.BodySize) }}</span>
        </p>
      </div>
    </div>
    <div class="mt-2 sm:flex sm:justify-between">
      <div class="flex-1">
        <p class="flex items-center text-sm text-frost-3">
          {{ request.Host }}
        </p>
      </div>
      <div v-if="showTags" class="mt-2 flex-0 text-sm text-frost-3 sm:mt-0 text-right">
      <span v-for="tag in request.Tags"
            :class="['text-polar-night-1 rounded-full px-2.5 py-0.5 text-xs font-medium ml-1', classForTag(tag)]">{{
          tag
        }}</span>
        <span v-if="request.Response" v-for="tag in request.Response.Tags"
              :class="['text-polar-night-1 rounded-full px-2.5 py-0.5 text-xs font-medium ml-1', classForTag(tag)]">{{
            tag
          }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>