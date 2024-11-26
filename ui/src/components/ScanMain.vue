<template>
  <TooltipProvider :delay-duration="0">
    <ResizablePanelGroup id="resize-panel-group-1"
                         direction="horizontal"
                         class="h-full items-stretch">
      <Separator orientation="vertical" />
      <ResizablePanel id="resize-panel-2"
                      :default-size="defaultLayout[1]"
                      :min-size="20">
        <div class="flex items-center justify-between px-4 py-2.5">
          <h1 class="text-xl font-bold">
            Scan
          </h1>
          <Dialog v-model:open="isModalOpen">
            <DialogTrigger as-child>
              <Button variant="outline"
                      size="sm">
                <CircleFadingPlus class="mr-2 w-4" />Add Domain
              </Button>
            </DialogTrigger>
            <DialogContent class="sm:max-w-[425px]">
              <DialogHeader>
                <DialogTitle>New domain</DialogTitle>
                <DialogDescription>
                  Add a new domain.
                </DialogDescription>
              </DialogHeader>
              <div class="py-0">
                <Input id="name"
                       v-model="domainName"
                       placeholder="example.com"
                       @keyup.enter.prevent="handleAddDomain"
                       class="w-full" />
                <div class="h-4 text-xs font-medium text-destructive">{{ errors }}</div>
              </div>
              <div class="items-top flex gap-x-2">
                <Checkbox id="auto-scan"
                          name="auto_scan"
                          :checked="autoScan"
                          @update:checked="handleAutoScanChange" />
                <div class="grid gap-1.5 leading-none">
                  <label for="auto-scan"
                         class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                    Auto-scan
                  </label>
                  <p class="text-sm text-muted-foreground">
                    Immediately scan this domain for hosts and subdomains.
                  </p>
                </div>
              </div>
              <DialogFooter>
                <Button type="submit"
                        @click.prevent="handleAddDomain"
                        :disabled="isSubmitting || !domainName">
                  {{ isSubmitting ? 'Adding...' : `Add ${autoScan ? 'and scan' : ''}` }}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
        <Separator />
        <HelperHint hint-key="scan.main">
          Add a domain to scan for live hosts and subdomains.
        </HelperHint>
        <div class="bg-background/95 p-4 backdrop-blur supports-[backdrop-filter]:bg-background/60">
          <div class="relative flex items-center gap-2">
            <Search class="absolute left-2 top-2.5 size-4 text-muted-foreground" />
            <Input v-model="searchValue"
                   placeholder="filter..."
                   class="pl-8" />
          </div>
        </div>
        <div clas="px-4">
          <ScanList :filter="searchValue" />
        </div>
      </ResizablePanel>
      <ResizableHandle id="resize-handle-2"
                       with-handle />
      <ResizablePanel id="resize-panel-3"
                      :default-size="defaultLayout[2]">
        <ScanDisplay :domain="domain" />
      </ResizablePanel>
    </ResizablePanelGroup>
  </TooltipProvider>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import ScanList from './ScanList.vue'
import ScanDisplay from './ScanDisplay.vue'
import { Search } from 'lucide-vue-next'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Separator } from '@/components/ui/separator'
import { Input } from '@/components/ui/input'
import { TooltipProvider } from '@/components/ui/tooltip'
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from '@/components/ui/resizable'
import { CircleFadingPlus } from 'lucide-vue-next'
import HelperHint from './HelperHint.vue'

import type { Domain } from '@/stores/scan'
import { useScanStore } from '@/stores/scan'

const defaultLayout = ref([20, 30, 70])

const scanStore = useScanStore()
const errors = computed(() => scanStore.errors)
const searchValue = ref('')
const domainName = ref('')
const autoScan = ref(true)
const isModalOpen = ref(false)
const isSubmitting = ref(false)

const handleAutoScanChange = () => {
  autoScan.value = !autoScan.value
}

const handleAddDomain = () => {
  if (!domainName.value) return

  isSubmitting.value = true
  scanStore.createDomain({
    name: domainName.value,
    auto_scan: autoScan.value,
  })
    .then(() => {
      // close the modal and reset the form
      isModalOpen.value = false
      domainName.value = ''
      autoScan.value = false
    })
    .finally(() => {
      isSubmitting.value = false
    })
}

const domain: Domain = {
  name: 'example.com',
  status: 'pending',
  auto_scan: true,
  host_count: 34,
  last_scanned_at: new Date(),
}

onMounted(() => {
  scanStore.errors = ''
})
</script>
