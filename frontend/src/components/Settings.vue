
<script lang="ts">
import {defineComponent} from "vue";
import Settings from '../lib/Settings'
import {EventsEmit, BrowserOpenURL} from "../../wailsjs/runtime";

interface themeInfo {
  id: string
  name: string
}

export default /*#__PURE__*/ defineComponent({
  data() {
    let items: Array<themeInfo> = [{id:"ghost", name:"Ghost"}, {id: "dark", name:"Dark"}, {id: "light", name: "Light"}]
    return {
      tab: null,
      selectedTheme: { id: 'ghost', name: 'Ghost'},
      selectedProxyPort: 0,
      selectedProxyHost: '',
      items: items,
    }
  },
  props: {
    settings: {type: Settings, required: true},
    onSave: {type: Function, required: true},
  },
  beforeMount() {
    // update theme selection to actual theme in use
    this.selectedTheme = this.items.find(item => item.id === this.settings.Theme) as themeInfo
    this.selectedProxyPort = this.settings.ProxyPort
    this.selectedProxyHost = this.settings.ProxyHost
  },
  methods: {
    saveTheme(theme: string) {
      let newSettings = this.settings
      newSettings.Theme = theme
      this.saveSettings(newSettings)
    },
    saveProxyPort(portStr: string) {
      let port = parseInt(portStr)
      let newSettings = this.settings
      newSettings.ProxyPort = port
      this.saveSettings(newSettings)
    },
    saveProxyHost(host: string) {
      let newSettings = this.settings
      newSettings.ProxyHost = host
      this.saveSettings(newSettings)
    },
    saveSettings(settings: Settings) {
      this.onSave(settings)
    },
    validatePort: (portStr: string) => {
      let port = parseInt(portStr)
      if (port <= 0 || port > 65535) {
        return 'Port must be between 0 and 65535'
      }
      return true
    },
    exportCA() {
      EventsEmit("OnExportCA")
    },
    openGithub() {
      BrowserOpenURL('https://github.com/ghostsecurity/reaper')
    },
    openWebsite() {
      BrowserOpenURL('https://ghost.security')
    },
  },
})
</script>

<template>
  <v-card class="d-flex flex-column fill-height">
    <v-tabs
        v-model="tab"
        show-arrows
    >
      <v-tab value="display">
        Display
      </v-tab>
      <v-tab value="certificates">
        Certificates
      </v-tab>
      <v-tab value="proxy">
        Proxy
      </v-tab>
      <v-tab value="about">
        About
      </v-tab>
    </v-tabs>

    <v-card-text>
      <v-window v-model="tab" >
        <v-window-item value="display">
        <v-card class="text-start">
            <v-card-title class="text-h5 pa-2">Appearance</v-card-title>
          <v-card-text>
            <v-select
                v-model="selectedTheme"
                :items="items"
                item-title="name"
                item-value="id"
                label="Theme"
                @update:modelValue="saveTheme"
            ></v-select>
          </v-card-text>
        </v-card>
        </v-window-item>
        <v-window-item value="certificates">
          <v-btn @click="exportCA">Export CA</v-btn>
          <v-btn>Regenerate CA</v-btn>
        </v-window-item>
        <v-window-item value="proxy">
          <v-text-field
              v-model="selectedProxyHost"
              label="Proxy Host"
              @update:modelValue="saveProxyHost"></v-text-field>
          <v-text-field
              v-model="selectedProxyPort"
              label="Proxy Port"
              :rules="[validatePort]"
              @update:modelValue="saveProxyPort"></v-text-field>
        </v-window-item>
        <v-window-item value="about">
          <div  class="d-flex justify-center">
          <v-card width="50vw" variant="tonal" class="mt-10" max-width="500">
            <v-card-item>
              <v-card-text>
                <img src="../../src/assets/images/logo.png" alt="logo" style="width:100px">
                <v-card-title>Reaper</v-card-title>
                <v-card-subtitle>Version 0.0.0</v-card-subtitle>
                <v-card-text>Built by Ghost Security.</v-card-text>
              </v-card-text>
            </v-card-item>
            <v-card-actions>
              <v-btn prepend-icon="mdi:mdi-github" @click="openGithub">Github</v-btn>
              <v-btn prepend-icon="mdi:mdi-earth" @click="openWebsite">Website</v-btn>
            </v-card-actions>
          </v-card>
          </div>
        </v-window-item>
      </v-window>
    </v-card-text>
  </v-card>

</template>
