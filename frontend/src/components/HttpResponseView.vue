

<script lang="ts">
import {defineComponent, PropType} from "vue";
import {HttpRequest, HttpResponse} from "../lib/Http.js";
import Code from "./Code.vue";
import HttpStatus from "./HttpStatus.vue";

export default /*#__PURE__*/ defineComponent({
  components: {Code, HttpStatus},
  props: {
    response: {type: Object as PropType<HttpResponse>, required: true},
    readonly: {type: Boolean, required: true},
    onchange: {type: Function as PropType<(raw: string) => void>, required: false},
  },
  data() {
    return {
      tab: '',
    }
  },
})
</script>

<template>
  <div class="d-flex fill-height">
    <div class="v-col-8 fill-height">
      <Code v-if="response !== undefined && response !== null"  :onchange="onchange" :readonly="readonly" :code="response.Raw" class="fill-height"/>
    </div>
    <div class="v-col-4 text-left overflow-y-auto">
      <p class="text-h5 text--primary">
        Summary
      </p>
      <v-table  density="compact" style="max-width: 100% !important;">
        <tbody>
        <tr>
          <th>Status</th>
          <td><HttpStatus :code="response.StatusCode" /></td>
        </tr>
        <tr>
          <th>Reaper ID</th>
          <td>{{response.ID}}</td>
        </tr>
        </tbody>
      </v-table>
      <p class="text-h5 text--primary mt-5">
        Headers
      </p>
      <v-table  density="compact" style="max-width: 100% !important;">
        <tbody>
        <template v-for="(values, key) in response.Headers">
          <tr v-for="value in values" :key="key + ':' + value">
            <th>{{key}}</th>
            <td>{{value}}</td>
          </tr>
        </template>
        </tbody>
      </v-table>
    </div>
  </div>
</template>

<style lang="scss" scoped>
tbody {
  max-width: 100% !important;
}
th {
  font-weight: normal !important;
  max-width: 15vw !important;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
td {
  max-width: 15vw !important;
  text-align: right;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;

}
td input {
  width: 100%;
}
</style>