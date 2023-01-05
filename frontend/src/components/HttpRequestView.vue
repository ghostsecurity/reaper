

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { HttpRequest, HttpResponse } from "../lib/Http.js";
import Code from "./Code.vue";
import HttpStatus from "./HttpStatus.vue";

export default /*#__PURE__*/ defineComponent({
  components: { HttpStatus, Code },
  props: {
    request: { type: Object as PropType<HttpRequest>, required: true },
    readonly: { type: Boolean, required: true },
    onchange: { type: Function as PropType<(raw: string) => void>, required: false },
  },
})
</script>

<template>
  <!-- TODO add tabs for headers, tags etc. -->
  <div class="h-full min-h-full w-full">
    <Code :code="request.Raw" :onchange="onchange" :readonly="readonly" class="h-full" />
  </div>
  <!--
    <div class="text-left overflow-y-auto">
        <p class="text-h5 text--primary">
          Summary
        </p>
        <table  density="compact" style="max-width: 100% !important;">
          <tbody>
            <tr>
              <th>Status</th>
              <td><HttpStatus :code="request.Response?.StatusCode" /></td>
            </tr>
            <tr>
              <th>Method</th>
              <td>{{request.Method}}</td>
            </tr>
            <tr>
              <th>Host</th>
              <td>{{request.Host}}</td>
            </tr>
            <tr>
              <th>Path</th>
              <td>{{request.Path}}</td>
            </tr>
            <tr>
              <th>Reaper ID</th>
              <td>{{request.ID}}</td>
            </tr>
          </tbody>
        </table>
      <p class="text-h5 text--primary mt-5">
        Headers
      </p>
      <table  density="compact" style="max-width: 100% !important;">
        <tbody>
        <template v-for="(values, key) in request.Headers">
          <tr v-for="value in values" :key="key + ':' + value">
            <th>{{key}}</th>
            <td>{{value}}</td>
          </tr>
        </template>
        </tbody>
      </table>
      <p class="text-h5 text--primary mt-5">
        Query Parameters
      </p>
      <table  density="compact" style="max-width: 100% !important;">
        <tbody>
        <template v-for="(values, key) in request.Query">
          <tr v-for="value in values" :key="key + ':' + value">
            <th>{{key}}</th>
            <td>{{value}}</td>
          </tr>
        </template>
        </tbody>
      </table>
    </div>
    -->
</template>

<style scoped>
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