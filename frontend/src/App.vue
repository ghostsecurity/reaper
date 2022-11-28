<script lang="ts">
import {defineComponent} from "vue";
import {HttpRequest, HttpResponse} from "./lib/Http";
import Proxy from './components/Proxy.vue'
import SettingsComponent from './components/Settings.vue'
import { EventsOn, EventsEmit } from '../wailsjs/runtime';
import {useTheme} from "vuetify";
import Settings from "./lib/Settings";


export default defineComponent({
  components: {
    Proxy,
    SettingsComponent
  },
  data: () => ({
    tab: null,
    requests: Array<HttpRequest>(),
    settings: new Settings,
  }),
  beforeMount() {
    let comp = this;
    EventsOn("OnSettingsLoad", (data) => {
      console.log("Settings Loaded", data)
      comp.settings = data as Settings;
      this.changeTheme(comp.settings.Theme);
    });
    EventsOn("OnHttpRequest", (data) => {
      comp.requests.push(data)
    });
    EventsOn("OnHttpResponse", (response: HttpResponse) => {
      for (let i = 0; i < comp.requests.length; i++) {
        if (comp.requests[i].ID === response.ID) {
          comp.requests[i].Response = response;
          break;
        }
      }
    });
    EventsEmit("OnAppReady");
  },
  setup () {
    const theme = useTheme()
    return {
      theme,
      changeTheme: (themeId: string) => {
        theme.global.name.value = themeId;
      },
    }
  },
  methods: {
    saveSettings(settings: Settings) {
      EventsEmit("OnSettingsSave", settings)
    },
  }
})
</script>

<template>
  <v-app class="d-flex flex-column fill-height" :style="{background: theme.current.value.colors.background}">
      <v-tabs
          v-model="tab"
          centered
          stacked
          show-arrows
          ref="strip"
          style="flex: 0 1 auto;"
      >
        <v-tab value="proxy">
          <v-icon icon="mdi:mdi-gate" />
          Proxy
        </v-tab>
        <v-tab value="settings">
          <v-icon icon="mdi:mdi-cog" />
          Settings
        </v-tab>
      </v-tabs>
        <v-window v-model="tab"  style="flex: 1 1 auto;" class="fill-height">
          <v-window-item value="proxy" class="fill-height">
            <Proxy :requests="requests"/>
          </v-window-item>
          <v-window-item value="settings" class="fill-height">
            <SettingsComponent :settings="settings" :onSave="saveSettings"/>
          </v-window-item>
        </v-window>
  </v-app>
</template>

