<template>
  <div class="flex h-screen flex-col bg-background">
    <main v-if="!loading"
          class="container flex-1 overflow-auto">
      <Dashboard v-if="loggedIn"
                 :wsConnected="wsConnected" />
      <Session v-if="!loggedIn" />
      <Toaster />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, type Ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { Toaster } from '@/components/ui/toast'
import Session from '@/components/SessionMain.vue'
import Dashboard from '@/components/DashboardMain.vue'
import { useSessionStore } from '@/stores/session'
import { useExploreStore } from '@/stores/explore'
import { useScanStore } from '@/stores/scan'
import { useEndpointStore } from '@/stores/endpoint'
import { useAgentStore } from '@/stores/agent'
import { useConfigStore } from '@/stores/config'

useConfigStore()
const sessionStore = useSessionStore()
const exploreStore = useExploreStore()
const scanStore = useScanStore()
const endpointStore = useEndpointStore()
const agentStore = useAgentStore()
const loggedIn = computed(() => sessionStore.loggedIn)
const loading = computed(() => sessionStore.loading)

const MAX_RECONNECT_ATTEMPTS = 50
const RECONNECT_DELAY_MS = 2000
const HEARTBEAT_INTERVAL_MS = 5000
const wsConnected: Ref<boolean> = ref(false)
const wsStreamUrl: string = getWebSocketUrl()

let ws: WebSocket | null = null
let heartbeatInterval: NodeJS.Timeout | null = null
let reconnectAttempts = 0

/**
 * Get the WebSocket URL.
 *
 * If the environment is not production (i.e. static build),
 * use the VITE_WS_URL environment variable. Otherwise,
 * derive the WebSocket URL from the current URL.
 */
function getWebSocketUrl() {
  const { protocol, hostname, port } = window.location
  if (!import.meta.env.PROD) {
    return import.meta.env.VITE_WS_URL
  }
  const wsProtocol = protocol.replace('http', 'ws')
  const wsUrl = `${wsProtocol}//${hostname}${port ? `:${port}` : ''}/ws`
  return wsUrl
}

function connectWebSocket() {
  ws = new WebSocket(wsStreamUrl)

  ws.onopen = () => {
    wsConnected.value = true;
    reconnectAttempts = 0 // reset the reconnect attempts on a successful connection
    startHeartbeat()
  }

  ws.onclose = (event) => {
    wsConnected.value = false
    if (!event.wasClean && reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
      reconnect()
    }
    stopHeartbeat()
  }

  ws.onerror = (error) => {
    console.error("[ws] WebSocket error:", error)
    ws?.close()
  }

  // Handle incoming messages
  ws.onmessage = (e) => {
    const payload = e.data

    if (payload === "pong") {
      return
    }

    let data

    try {
      data = JSON.parse(payload)
    }
    catch (error) {
      console.error("[ws] Error parsing message:", error)
      return
    }

    // Handle incoming messages
    switch (data.type) {
      case "debug":
        console.log("log:", data)
        break
      case "explore_host":
        console.info("explore_host:", data)
        exploreStore.addHost(data)
        break
      case "explore.response":
        console.info("explore.response:", data)
        exploreStore.addEndpoint(data.host, data)
        break
      case "scan.domain":
        console.info("scan.domain:", data)
        break
      case "scan.domain.sync":
        console.info("scan.domain.sync:", data)
        scanStore.syncDomain(data)
        break
      case "scan_host":
        console.info("scan_host:", data)
        break
      case "attack.result":
        console.info("attack.result:", data)
        endpointStore.addResult(data)
        break
      case "attack.result.clear":
        console.info("attack.result.clear:", data)
        endpointStore.emptyResults()
        break
      case "attack.complete":
        console.info("attack.complete:", data)
        endpointStore.attackCompleted()
        break
      case "navigation.follow":
        console.info("navigation.follow:", data)
        sessionStore.navigationFollow(data.to)
        break
      case "agent.session.message":
        console.info("agent.session.message:", data)
        agentStore.appendMessageToSession(data.session_id, data)
        break
      default:
        console.log("unknown message type:", data)
        break
    }
  }
}

function reconnect() {
  setTimeout(() => {
    console.info("[ws] Attempting to reconnect...")
    reconnectAttempts++
    connectWebSocket()
  }, generateReconnectDelay(reconnectAttempts))
}

function generateReconnectDelay(attempts: number) {
  const delay = RECONNECT_DELAY_MS + attempts * 100
  return delay
}

function startHeartbeat() {
  heartbeatInterval = setInterval(() => {
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send("ping")
    }
  }, HEARTBEAT_INTERVAL_MS)
}

function stopHeartbeat() {
  if (heartbeatInterval) {
    clearInterval(heartbeatInterval)
    heartbeatInterval = null
  }
}

// watch loggedIn status
watch(loggedIn, (newVal) => {
  if (newVal) {
    connectWebSocket()
  }
})

onMounted(() => {
  sessionStore.status()
  if (sessionStore.loggedIn) {
    connectWebSocket()
  }
})

onUnmounted(() => {
  ws?.close()
})

</script>
