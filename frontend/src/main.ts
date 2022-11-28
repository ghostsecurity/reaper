import {createApp} from 'vue'
import App from './App.vue'


// Vuetify
import 'vuetify/styles'
import { createVuetify, ThemeDefinition } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { aliases, mdi } from 'vuetify/iconsets/mdi'

import '@mdi/font/css/materialdesignicons.css'

import ghostTheme from './lib/Theme'

const vuetify = createVuetify({
    components,
    directives,
    theme: {
        defaultTheme: 'ghost',
        themes: {
            ghost: ghostTheme,
        },
        variations: {
            colors: ['primary', 'secondary'],
            lighten: 4,
            darken: 4,
        },
    },
    icons: {
        defaultSet: 'mdi',
        aliases,
        sets: {
            mdi,
        }
    },
})

createApp(App).
use(vuetify).
mount('#app')
