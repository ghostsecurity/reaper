<template>
  <div class="flex h-full flex-col">
    <div class="flex items-center p-2">
      <div class="flex items-center gap-2">
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost"
                    size="icon"
                    :disabled="!report"
                    @click="handleDeleteReport">
              <Trash2 class="size-4" />
              <span class="sr-only">Delete report</span>
            </Button>
          </TooltipTrigger>
          <TooltipContent>Delete report</TooltipContent>
        </Tooltip>
      </div>
    </div>
    <Separator />
    <div v-if="report"
         class="flex flex-1 flex-col">
      <div class="flex items-center p-4">
        <div class="flex items-start gap-4 text-sm">
          <div class="grid gap-1">
            <div class="text-lg font-semibold">
              {{ report.domain }}
            </div>
          </div>
        </div>
        <div v-if="report.created_at"
             class="ml-auto text-xs font-semibold text-muted-foreground">
          {{ format(new Date(report.created_at), "PPpp") }}
        </div>
      </div>
      <Separator />
      <div class="h-screen space-y-2 overflow-y-auto whitespace-pre-wrap bg-muted/50 p-1 text-sm text-foreground/80">
        <div class="rounded-md bg-background p-10">
          <vue-markdown id="report-markdown"
                        :source="report.markdown" />
        </div>
      </div>
    </div>
    <div v-else
         class="p-8 text-center text-sm font-medium text-muted-foreground">
      No report selected.
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Trash2 } from 'lucide-vue-next'
import { format } from 'date-fns'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

import VueMarkdown from 'vue-markdown-render'
import type { Report } from '@/stores/report'
import { useReportStore } from '@/stores/report'

const reportStore = useReportStore()

interface ReportsDisplayProps {
  report: Report | undefined
}

const props = defineProps<ReportsDisplayProps>()

const handleDeleteReport = () => {
  if (props.report) {
    reportStore.deleteReport(props.report)
  }
}
</script>
