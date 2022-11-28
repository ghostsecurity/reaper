
<script lang="ts" setup>
import TrafficLog from './TrafficLog.vue'
import TrafficStats from './TrafficStats.vue'
import Interceptor from './Interceptor.vue'
import RequestAndResponseInspector from './HttpInspector.vue';
</script>

<script lang="ts">

import {HttpRequest} from '../lib/Http';
import {defineComponent, PropType} from "vue";

export default /*#__PURE__*/ defineComponent({
  data() {
    return {
      tab: null,
      req: null as HttpRequest | null,
    }
  },
  props: {
    //height: {type: Number, required: true},
    requests: {type: Array as PropType<Array<HttpRequest>>, required: true},
  },
  mounted() {
   // this.calculateHeight()
  },
  methods: {
    // calculateHeight () {
    //   console.log(this.height)
    //   let availableHeight = this.height - (this.$refs['strip'] as any).$el.clientHeight
    //   if (availableHeight < 0) {
    //     availableHeight = 0
    //   }
    //   if (this.req !== null) {
    //     this.topHeight = availableHeight / 2
    //     this.bottomHeight = availableHeight / 2
    //   } else {
    //     this.bottomHeight = 100
    //     this.topHeight = availableHeight - this.bottomHeight
    //     if (this.topHeight < 0) {
    //       this.topHeight = 0
    //     }
    //   }
    // },
    examineRequest(request: HttpRequest) {
      this.req = request
     // this.calculateHeight()
    }
  },
  watch: {
    height: function (val, oldVal) {
      //this.calculateHeight()
    }
  }
})
</script>

<template>
  <v-card class="d-flex flex-column fill-height">
    <v-tabs
        v-model="tab"
        ref="strip"
        style="flex: 0 1 auto;"
    >
      <v-tab value="request-log">
        Request Log
      </v-tab>
      <v-tab value="intercept">
        Intercept
      </v-tab>
      <v-tab value="scope">
        Scope
      </v-tab>
    </v-tabs>
      <v-window v-model="tab" style="flex: 1 1 auto;" class="fill-height bg-surface">
        <v-window-item value="request-log" class="fill-height">
          <div class="d-flex flex-column fill-height">
          <TrafficLog :requests="requests" :onSelect="examineRequest" v-bind:style="{flex: '1 1 auto', 'max-height': req === null ? '100%' : '50%', overflow: 'auto'}" /><!--v-bind:style="{flex: req == null ? '80%' : '50%'}"/>-->
          <RequestAndResponseInspector v-if="req" :request="req" style="flex: 1 0 50%; max-height:50%" />
          <TrafficStats v-if="req===null && false" :requests="requests" style="flex: 1 0 10%;max-height:100px" />
          </div>
        </v-window-item>
        <v-window-item value="intercept">
          <Interceptor />
        </v-window-item>
        <v-window-item value="options">
          options here
        </v-window-item>
      </v-window>
  </v-card>
</template>

<style scoped>
</style>
