<template>
  <div class="flex h-full flex-col">
    <div class="flex items-center p-2">
      <div class="flex items-center gap-2">
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost"
                    size="icon">
              <Trash2 class="size-4" />
              <span class="sr-only">Delete endpoint</span>
            </Button>
          </TooltipTrigger>
          <TooltipContent>Delete endpoint</TooltipContent>
        </Tooltip>
      </div>
    </div>
    <Separator />
    <div class="p-8 text-center text-sm font-medium text-muted-foreground">
      To explore hosts and endpoints, start browsing using the Reaper proxy.
    </div>
    <div class="p-4 text-center text-xs">
      The default proxy URL is
      <Badge as="child"
             variant="outline"
             class="text-muted-foreground">
        <a href="http://localhost:8080">http://localhost:8080</a>
      </Badge>
    </div>
    <div>
      {{ endpoint }}
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Trash2 } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { type Endpoint } from '@/stores/explore'

interface ExploreDisplayProps {
  endpoint: Endpoint | undefined
}
defineProps<ExploreDisplayProps>()

</script>