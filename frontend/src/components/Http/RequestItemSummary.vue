<script lang="ts" setup>
import { PropType } from 'vue'
import { HttpRequest } from '../../lib/Http'

defineProps({
  name: { type: String, required: false, default: '' },
  request: { type: Object as PropType<HttpRequest>, required: true },
  showTags: { type: Boolean, required: false, default: true },
})

function classForTag(tag: string): string {
  const tags = ['bg-aurora-3', 'bg-aurora-4', 'bg-aurora-5', 'bg-frost-1', 'bg-frost-2', 'bg-frost-3', 'bg-frost-4']
  let total = 0
  const chars = [...tag]
  chars.forEach(c => {
    total += c.charCodeAt(0)
  })
  return tags[total % tags.length]
}

function humanSize(size: number): string {
  return Intl.NumberFormat('en', {
    notation: 'compact',
    style: 'unit',
    unit: 'byte',
    maximumSignificantDigits: 2,
    unitDisplay: 'narrow',
  })
    .format(size)
    .replace(/([a-zA-z]+)/, ' $1')
}
</script>

<template>
  <div class="relative">
    <div v-if="name !== ''" class="absolute left-0 top-3 w-full text-center text-sm">
      <span class="rounded-md bg-polar-night-3 px-3 py-1">{{ name }}</span>
    </div>
    <div class="flex items-center justify-between">
      <p class="flex-1 truncate text-left text-sm font-medium text-frost-4 dark:text-frost">
        {{ request.Path }}
        <span class="max-w-4xl truncate text-frost-3" v-if="request.QueryString !== ''">
          ?{{ request.QueryString }}
        </span>
      </p>
      <div class="flex-0 ml-2 flex text-right">
        <p class="px-2 text-xs font-semibold leading-5">
          <span v-if="request.Response" class="text-polar-night-1a dark:text-snow-storm-1">
            {{ humanSize(request.Response.BodySize) }}
          </span>
        </p>
      </div>
    </div>
    <div class="mt-2 sm:flex sm:justify-between">
      <div class="flex-1">
        <p class="flex items-center text-sm text-frost-3 dark:text-frost-3">
          {{ request.Host }}
        </p>
      </div>
      <div v-if="showTags" class="flex-0 mt-2 text-right text-sm text-pink-500 sm:mt-0">
        <span
          v-for="tag in request.Tags"
          :key="tag"
          :class="['ml-1 rounded-full px-2.5 py-0.5 text-xs font-medium text-polar-night-1', classForTag(tag)]">
          {{ tag }}
        </span>
        <span v-if="request.Response">
          <span
            v-for="tag in request.Response.Tags"
            :key="tag"
            :class="['ml-1 rounded-full px-2.5 py-0.5 text-xs font-medium text-polar-night-1', classForTag(tag)]">
            {{ tag }}
          </span>
        </span>
      </div>
    </div>
  </div>
</template>
