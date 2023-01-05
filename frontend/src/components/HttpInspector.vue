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
      currentTab: 'request'
    }
  },
  watch: {
    request: function () {
      if (this.request.Response == null) {
        this.currentTab = 'request'
      }
    },
  },
  methods: {
    selectTab: function (tab: string) {
      this.currentTab = tab
    }
  },
})
</script>

<template>
  <div class="smart">
    <div>
      <div class="border-b border:snow-storm-3 dark:border-polar-night-4">
        <nav class="-mb-px flex space-x-8" aria-label="Tabs">
          <a @click="selectTab('request')"
             :class="['request' == currentTab ?
             'border-frost text-frost' :
             'border-transparent text-polar-night-4 hover:text-frost-4 hover:border-frost-4',
              'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm']" :aria-current="'request' == currentTab ? 'page' : undefined">
            Request
          </a>
          <a v-if="request.Response" @click="selectTab('response')"
             :class="['response' == currentTab ?
             'border-frost text-frost' :
             'border-transparent text-polar-night-4 hover:text-frost-4 hover:border-frost-4',
              'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm']" :aria-current="'response' == currentTab ? 'page' : undefined">
            Response
          </a>
        </nav>
      </div>
    </div>
    <div class="pt-4 h-full">
      <div :class="{'hidden': currentTab != 'request', 'h-full': true}">
        <HttpRequestView :request="request" :readonly="true" />
      </div>
      <div :class="{'hidden': currentTab != 'response', 'h-full': true}">
        <HttpResponseView v-if="request.Response" :response="request.Response" :readonly="true" />
      </div>
    </div>
  </div>
</template>

<style scoped>
a{
  cursor: pointer;
}
.smart{
  height: calc(100% - 4rem);
}
</style>

