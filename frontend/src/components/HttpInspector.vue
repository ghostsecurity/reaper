<script lang="ts">
import {defineComponent, PropType} from 'vue'
import {HttpRequest} from '../lib/Http.js';
import HttpRequestView from "./HttpRequestView.vue";
import HttpResponseView from "./HttpResponseView.vue";

export default /*#__PURE__*/ defineComponent({
  components: {HttpRequestView, HttpResponseView},
  props: {
    request: {type: Object as PropType<HttpRequest>, required: true},
  },
  data() {
    return {
      tab: '',
    }
  },
  watch: {
    request: function () {
      if (this.request.Response == null) {
        this.tab = 'request'
      }
    }
  },
})
</script>

<template>
  <div class="d-flex flex-column fill-height">
  <v-tabs
      v-model="tab"
      bg-color="background"
      show-arrows
      ref="strip"
      style="flex: 0 1 auto"
  >
    <v-tab value="request">
      Request
    </v-tab>
    <v-tab v-if="request.Response" value="response">
      Response
    </v-tab>
  </v-tabs>
  <v-window v-model="tab"  style="flex: 1 1 auto" class="fill-height">
    <v-window-item value="request" class="fill-height">
      <HttpRequestView :readonly="true" :request="request" />
    </v-window-item>
    <v-window-item value="response"  class="fill-height">
      <HttpResponseView v-if="request.Response" :readonly="true" :response="request.Response" />
    </v-window-item>
  </v-window>
  </div>

</template>

<style scoped>
</style>

