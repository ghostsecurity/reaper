<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { EventsEmit, EventsOn } from '../../wailsjs/runtime' // eslint-disable-line import/no-unresolved
import { HttpRequest } from '../lib/Http'
import HttpRequestView from './Http/HttpRequest.vue'

const enabled = ref(false)
const request = ref<HttpRequest | null>(null)

onMounted(() => {
  EventsOn('InterceptedRequest', (req: HttpRequest) => {
    request.value = req
  })
})

function onEnabledChange() {
  if (!enabled.value) {
    request.value = null
  }
  EventsEmit('InterceptionEnabledChange', enabled.value)
}

function forwardRequest() {
  EventsEmit('InterceptedRequestChange', request.value)
  request.value = null
}

function dropRequest() {
  EventsEmit('InterceptedRequestDrop', request.value)
  request.value = null
}

function onChange(raw: string) {
  if (request.value !== null) {
    // TODO
    // request.value.Raw = raw
  }
}
</script>

<template>
  <v-card class="text-start">
    <v-card-text>
      <v-switch v-model="enabled" color="primary" label="Intercept Requests" @change="onEnabledChange" hide-details
        inset inline></v-switch>
      <div v-if="request !== null" class="d-flex justify-left align-baseline" style="gap: 1rem">
        <v-btn color="primary" @click="forwardRequest">Forward</v-btn>
        <v-btn color="error" @click="dropRequest">Drop</v-btn>
      </div>
    </v-card-text>
  </v-card>
  <HttpRequestView v-if="request !== null" :request="request" :readonly="false" @change="onChange" />
</template>

<style scoped>

</style>
