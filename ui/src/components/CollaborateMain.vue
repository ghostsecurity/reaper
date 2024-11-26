<template>
  <TooltipProvider :delay-duration="0">
    <ResizablePanelGroup id="resize-panel-group-1"
                         direction="horizontal"
                         class="h-full items-stretch">
      <Separator orientation="vertical" />
      <ResizablePanel id="resize-panel-2"
                      :default-size="defaultLayout[1]"
                      :min-size="20">
        <div class="flex items-center justify-between px-4 py-3">
          <h1 class="text-xl font-bold">
            Collaborate
          </h1>
          <div class="flex items-center space-x-2">
            <Label for="proxy-status"
                   class="text-xs text-muted-foreground">Tunnel {{ tunnel.enabled ? 'on' : 'off' }}</Label>
            <Switch id="proxy-status"
                    :checked="tunnel.enabled"
                    @click.prevent="handleTunnelStatusToggle" />
          </div>
        </div>

        <Separator />

        <div class="h-screen bg-muted/50 pt-4">
          <Card v-if="tunnel.enabled"
                class="mx-4">
            <CardHeader>
              <CardTitle>Tunnel URL</CardTitle>
              <CardDescription class="text-xs">Share this URL with your team to collaborate.</CardDescription>
            </CardHeader>
            <CardContent>
              <div v-if="tunnel.url"
                   class="flex items-center space-x-2 text-xs text-foreground">
                <div class="cursor-pointer text-muted-foreground hover:text-primary"
                     @click="copyToClipboard">
                  <Files class="size-4" />
                </div>
                <div>{{ tunnel.url }}</div>
              </div>
            </CardContent>
          </Card>
          <Card v-else
                class="mx-4">
            <CardHeader>
              <CardTitle>Create a Tunnel</CardTitle>
              <CardDescription class="text-xs">Create a secure tunnel to collaborate with your team.</CardDescription>
            </CardHeader>
          </Card>
        </div>
      </ResizablePanel>
      <ResizableHandle id="resiz-handle-2"
                       with-handle />
      <ResizablePanel id="resize-panel-3"
                      :default-size="defaultLayout[2]">
        <ExploreDisplay :endpoint="undefined" />
      </ResizablePanel>
    </ResizablePanelGroup>
  </TooltipProvider>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue'

import ExploreDisplay from './ExploreDisplay.vue'
import { Separator } from '@/components/ui/separator'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { TooltipProvider } from '@/components/ui/tooltip'
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from '@/components/ui/resizable'
import { useToast } from '@/components/ui/toast/use-toast'
import { Files } from 'lucide-vue-next'
import { useCollabStore } from '@/stores/collab'

const defaultLayout = ref([20, 30, 70])
const collabStore = useCollabStore()
const { toast } = useToast()


const tunnel = computed(() => collabStore.tunnel)

const copyToClipboard = () => {
  if (tunnel.value.url) {
    navigator.clipboard.writeText(tunnel.value.url)
    toast({
      title: 'URL copied to clipboard',
      description: 'You can now share this URL with your team to collaborate.',
    })
  }
}

const handleTunnelStatusToggle = () => {
  if (collabStore.tunnel.enabled) {
    collabStore.tunnelStop()
  } else {
    collabStore.tunnelStart()
    pingTunnelStatus()
  }
}

const pingTunnelStatus = () => {
  collabStore.tunnelStatus()
  console.log("pinging tunnel status", tunnel.value.url)
  if (tunnel.value.url === "" || tunnel.value.url === undefined) {
    setTimeout(pingTunnelStatus, 1000)
  }
}

onMounted(() => {
  collabStore.tunnelStatus()
  // pingTunnelStatus()
})
</script>
