<script lang="ts">
import { defineComponent } from 'vue'
import {HttpRequest, HttpResponse} from '../lib/Http.js';
import HttpStatus from "./HttpStatus.vue";
//import HttpRequest from '../../wailsjs/go/packaging'

export default /*#__PURE__*/ defineComponent({
  components: {HttpStatus},
  // type inference enabled
  props: {
    requests: {type: Array<HttpRequest>},
    onSelect: {type: Function},
  },
  data() {
    return {
      selected: -1,
    }
  },
  methods: {
    selectRequest(request: HttpRequest | null) {
      if (request === null || this.selected === request.ID) {
        this.selected = -1
        request = null;
      }else {
        this.selected = request.ID
      }
      if (this.onSelect !== undefined) {
        this.onSelect(request)
      }
    }
  },
})
</script>

<template>
  <div class="fill-height">
    <v-table
        density="compact"
        fixed-header
    >
      <thead>
      <tr>
        <th class="text-left">Method</th>
        <th class="text-left">Host</th>
        <th class="text-left">Path</th>
        <th class="text-left">Query</th>
        <th class="text-left">Status</th>
      </tr>
      </thead>
      <tbody>
        <tr v-for="request in requests" @click="selectRequest(request)" v-bind:class="(request.ID === selected)?'selected':''">
          <td>{{request.Method}}</td>
          <td>{{request.Host}}</td>
          <td>{{request.Path}}</td>
          <td>{{request.QueryString}}</td>
          <td><HttpStatus :code="request.Response?.StatusCode"/></td>
        </tr>
      </tbody>
    </v-table>
  </div>
</template>

<style scoped>
td {
  /* keep column sizes sensible */
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 20vw !important;
  text-align: left;
}

.v-theme--dark  tbody tr:hover, .v-theme--ghost tbody tr:hover {
  cursor: pointer;
  background-color: #6B7292;
}

.v-theme--dark tbody tr.selected, .v-theme--ghost tbody tr.selected {
  background-color: #492CFB;
}

.v-theme--light tbody tr:hover {
  cursor: pointer;
  background-color: #F2F4F7;
}

.v-theme--light  tbody tr.selected {
  background-color: #C3C8DF;
}



</style>

