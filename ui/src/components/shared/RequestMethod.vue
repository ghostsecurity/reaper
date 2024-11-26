<template>
  <Badge variant="outline"
         class="rounded-sm px-1 py-0"
         :class="badgeColorFromStatus(code)">
    <div class="text-2xs font-semibold">
      <slot />
    </div>
  </Badge>
</template>

<script lang="ts" setup>
import { Badge } from '@/components/ui/badge'

defineProps<{
  code: number
}>()

function badgeColorFromStatus(code: number) {
  if (code < 0)
    return 'bg-gray-50 border-gray-500/10 text-gray-600 dark:bg-gray-400/10 dark:border-gray-400/20 dark:text-gray-400'

  if (code >= 0 && code < 300)
    return 'bg-green-50 border-green-600/20 text-green-700 dark:bg-green-500/10 dark:border-green-500/20 dark:text-green-400'

  if (code >= 300 && code < 400)
    return 'bg-yellow-50 border-yellow-600/20 text-yellow-700 dark:bg-yellow-400/10 dark:border-yellow-400/20 dark:text-yellow-500'

  return 'bg-red-50 border-red-600/20 text-red-700 dark:bg-red-400/10 dark:border-red-400/20 dark:text-red-400'
}
</script>
