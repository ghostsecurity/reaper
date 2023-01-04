<!--
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

</template>-->

<script lang="ts">
import {
  InformationCircleIcon,
  ExclamationCircleIcon,
  PaintBrushIcon, ServerStackIcon,
  ShieldCheckIcon,
} from '@heroicons/vue/20/solid'
import {defineComponent} from "vue";
import {Disclosure, Switch, SwitchDescription, SwitchGroup, SwitchLabel} from '@headlessui/vue'
import Settings from "../lib/Settings";
import {EventsEmit, EventsOn} from "../../wailsjs/runtime";
import setDarkMode from "../lib/theme";

export default /*#__PURE__*/ defineComponent({
  components: {
    Switch, SwitchDescription, SwitchGroup, SwitchLabel,
    Disclosure, ExclamationCircleIcon
  },
  props: {
    settings: {type: Settings, required: true},
    onSave: {type: Function, required: true},
    onCancel: {type: Function, required: true},
  },
  beforeMount() {

  },
  data() {
    return {
      openTab: "display",
      tabs: [
        {name: 'Display', icon: PaintBrushIcon, id: "display"},
        {name: 'Certificates', icon: ShieldCheckIcon, id: "certs"},
        {name: 'Proxy', icon: ServerStackIcon, id: "proxy"},
        {name: 'About', icon: InformationCircleIcon, id: "about"},
      ],
      subNavigation: [],
      modifiedSettings: this.settings,
    }
  },
  methods: {
    setDarkMode(darkMode: boolean) {
      this.modifiedSettings.DarkMode = darkMode
    },
    saveSettings() {
      this.onSave(this.modifiedSettings)
    },
    cancel() {
      this.onCancel()
    },
    toggleTab: function (tabId: string) {
      this.openTab = tabId
    },
    exportCA() {
      EventsEmit("OnExportCA")
    },
    setProxyPort(event: any) {
      let port = parseInt(event.target.value)
      this.modifiedSettings.ProxyPort = port
    },
  }
})
</script>

<template>
  <div>
    <main class="relative text-left">
      <div class="mx-auto max-w-screen-[100%] px-4 pb-6 sm:px-6 lg:px-8 lg:pb-16">
        <div class="overflow-hidden rounded-lg bg-snow-storm dark:bg-polar-night shadow text-polar-night dark:text-snow-storm">
          <div class="lg:grid lg:grid-cols-12 lg:divide-y-0 lg:divide-x divide-snow-storm-3 dark:divide-polar-night-3">
            <aside class="py-6 lg:col-span-3">
              <nav class="space-y-1">
                <a @click="toggleTab(tab.id)" v-for="tab in tabs" :key="tab.name" :class="[
                    tab.id === openTab ?
                    'bg-polar-night-4 border-frost-3' :
                    'border-transparent hover:bg-polar-night-3',
                   'group border-l-4 px-3 py-2 flex items-center text-sm font-medium']"
                   :aria-current="tab.id === openTab ? 'page' : undefined">
                  <component :is="tab.icon" :class="['flex-shrink-0 -ml-1 mr-3 h-6 w-6']" aria-hidden="true" />
                  <span class="truncate">{{ tab.name }}</span>
                </a>
              </nav>
            </aside>

            <form class="lg:col-span-9" action="#" method="POST">

              <!-- Display settings -->
              <div :class="{'hidden': 'display' !== openTab}">
                <div class="py-6 px-4 sm:p-6 lg:pb-8">
                <div>
                  <h2 class="text-lg font-medium leading-6">Display</h2>
                  <p class="mt-1 text-sm">Modify the appearance of the application.</p>
                </div>
                  <ul role="list" class="mt-2 divide-y divide-gray-200">
                    <SwitchGroup as="li" class="flex items-center justify-between py-4">
                      <div class="flex flex-col">
                        <SwitchLabel as="p" class="text-sm font-medium" passive>Dark Mode</SwitchLabel>
                        <SwitchDescription class="text-sm">Unruin your eyes.</SwitchDescription>
                      </div>
                      <Switch @update:modelValue="setDarkMode" v-model="modifiedSettings.DarkMode" :class="[modifiedSettings.DarkMode ? 'bg-frost' : 'bg-aurora-1', 'relative ml-4 inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out']">
                        <span aria-hidden="true" :class="[modifiedSettings.DarkMode ? 'translate-x-5' : 'translate-x-0', 'inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out']" />
                      </Switch>
                    </SwitchGroup>
                  </ul>
                </div>
              </div>

              <!-- Proxy settings -->
              <div :class="{'hidden': 'proxy' !== openTab}">
                <div class="py-6 px-4 sm:p-6 lg:pb-8">
                  <div>
                    <h2 class="text-lg font-medium leading-6 ">Proxy</h2>
                    <p class="mt-1 text-sm ">Customise proxy settings to suit your workflow.</p>
                  </div>

                  <div class="mt-8">
                    <label for="port" class="block text-sm font-medium text-snow-storm">Port</label>
                    <div class="relative mt-1 rounded-md shadow-sm">
                      <input @change="setProxyPort" type="text" name="port" id="port" class="block w-full rounded-md bg-polar-night-4 pr-10 focus:outline-none sm:text-sm" :value="modifiedSettings.ProxyPort" aria-invalid="true" aria-describedby="port-error" />
                    </div>
                  </div>

                </div>
              </div>

              <!-- Certificate settings -->
              <div :class="{'hidden': 'certs' !== openTab}">
                <div class="py-6 px-4 sm:p-6 lg:pb-8">
                  <div>
                    <h2 class="text-lg font-medium leading-6 ">Certificates</h2>
                    <p class="mt-1 text-sm ">Manage certificates and configure your browser(s).</p>
                  </div>
                  <div class="mt-2 divide-y divide-gray-200">
                    <div class="flex flex-col">
                      <button @click="exportCA" type="button" class="inline-flex items-center rounded border border-transparent bg-frost px-2.5 py-1.5 text-xs font-medium text-snow-storm-3 shadow-sm hover:bg-frost-2 focus:outline-none">Export CA</button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- About -->
              <div :class="{'hidden': 'about' !== openTab}">
                <div class="py-6 px-4 sm:p-6 lg:pb-8">
                  <div>
                    <h2 class="text-lg font-medium leading-6 ">About</h2>
                    <p class="mt-1 text-sm ">More information about Reaper.</p>
                  </div>

                  <div class="mt-6 flex flex-col lg:flex-row">
                    <div class="flex-grow space-y-6">
                      <div>

                      </div>
                    </div>
                  </div>
                </div>

                <div class="divide-y divide-gray-200 pt-6">
                  <div class="px-4 sm:px-6">
                    <div>
                      <h2 class="text-lg font-medium leading-6 ">Privacy</h2>
                      <p class="mt-1 text-sm ">Ornare eu a volutpat eget vulputate. Fringilla commodo amet.</p>
                    </div>
                  </div>
                </div>
              </div>

              <div class="divide-y divide-gray-200 pt-6 text-right">
                <div class="px-4 sm:px-6 pb-4">
                  <div>
                    <button @click="saveSettings" type="button" class="inline-flex items-center rounded border border-transparent bg-aurora-4 px-2.5 py-1.5 text-xs font-medium text-snow-storm-3 shadow-sm hover:bg-aurora-5 focus:outline-none">Save Changes</button>
                    <button @click="cancel" type="button" class="ml-2 inline-flex items-center rounded border border-transparent bg-aurora-1 px-2.5 py-1.5 text-xs font-medium text-snow-storm-3 shadow-sm hover:bg-aurora-5 focus:outline-none">Cancel</button>
                  </div>
                </div>
              </div>

            </form>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
a {
  cursor: pointer;
}
</style>