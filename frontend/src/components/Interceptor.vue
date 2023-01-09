

<script lang="ts">
import {HttpRequest} from '../lib/Http';
import {defineComponent, PropType} from "vue";
import {EventsEmit, EventsOn} from "../../wailsjs/runtime";
import HttpRequestView from "./HttpRequestView.vue";

export default /*#__PURE__*/ defineComponent({
  components: {HttpRequestView},
  data() {
    return {
      enabled: false,
      request: null as HttpRequest | null,
    }
  },
  mounted() {
    EventsOn('InterceptedRequest', (request: HttpRequest) => {
      this.request = request
    })
  },
  methods: {
    onEnabledChange() {
      if(!this.enabled) {
        this.request = null
      }
      EventsEmit("InterceptionEnabledChange", this.enabled)
    },
    forwardRequest(){
      EventsEmit("InterceptedRequestChange", this.request)
      this.request = null
    },
    dropRequest(){
      EventsEmit("InterceptedRequestDrop", this.request)
      this.request = null
    },
    onChange(raw: string) {
      if(this.request !== null) {
        this.request.Raw = raw
      }
    },
  }
})
</script>

<template>
  <v-card class="text-start">
    <v-card-text>
       <v-switch
          v-model="enabled"
          color="primary"
          label="Intercept Requests"
          @change="onEnabledChange"
          hide-details
          inset
          inline
      ></v-switch>
      <div class="d-flex justify-left align-baseline" style="gap: 1rem">
        <v-btn color="primary" @click="forwardRequest" v-if="request !== null">Forward</v-btn>
        <v-btn color="error" @click="dropRequest" v-if="request !== null">Drop</v-btn>
      </div>
    </v-card-text>
  </v-card>
  <HttpRequestView v-if="request !== null" :request="request" :readonly="false" :onchange="onChange"/>
</template>

<style scoped>

</style>