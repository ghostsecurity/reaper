<script lang="ts" setup>
import { InformationCircleIcon, PaintBrushIcon, ServerStackIcon, ShieldCheckIcon } from '@heroicons/vue/20/solid'
import { Switch, SwitchDescription, SwitchGroup, SwitchLabel } from '@headlessui/vue'
import { reactive, ref } from 'vue'

import { EventsEmit } from '../../wailsjs/runtime' // eslint-disable-line import/no-unresolved
import Settings from '../lib/Settings'

import ButtonConfirm from './Shared/ButtonConfirm.vue'
import ButtonCancel from './Shared/ButtonCancel.vue'
import ButtonNetrual from './Shared/ButtonNeutral.vue'

const props = defineProps({
  settings: { type: Settings, required: true },
})

const tabs = [
  { name: 'Display', icon: PaintBrushIcon, id: 'display' },
  { name: 'Certificates', icon: ShieldCheckIcon, id: 'certs' },
  { name: 'Proxy', icon: ServerStackIcon, id: 'proxy' },
  { name: 'About', icon: InformationCircleIcon, id: 'about' },
]
const openTab = ref('display')
const modifiedSettings = reactive(props.settings)
const emit = defineEmits(['save', 'cancel'])

function setDarkMode(darkMode: boolean) {
  modifiedSettings.DarkMode = darkMode
  toggleDarkMode()
}

/**
 * TODO: remove this temp fix for dark mode toggle -jml
 */
function toggleDarkMode() {
  // toggle class="dark" on top level html element
  let isDarkMode = document.documentElement.classList.toggle('dark')

  window.localStorage.darkMode = isDarkMode
}

function saveSettings() {
  emit('save', modifiedSettings)
}

function cancel() {
  emit('cancel')
}

function toggleTab(tabId: string) {
  openTab.value = tabId
}

function exportCA() {
  EventsEmit('CAExport')
}

function setProxyPort(event: Event) {
  const port = parseInt((event.target as HTMLInputElement).value, 10)
  modifiedSettings.ProxyPort = port
}
</script>

<template>
  <div>
    <main class="relative text-left">
      <div class="mx-auto max-w-2xl px-4 pb-6 sm:px-6 lg:px-8 lg:pb-16">
        <div
          class="divide-y divide-snow-storm-3 overflow-hidden rounded-lg bg-snow-storm text-polar-night shadow dark:divide-polar-night-3 dark:bg-polar-night dark:text-snow-storm">
          <div class="divide-snow-storm-3 dark:divide-polar-night-3 lg:grid lg:grid-cols-12 lg:divide-y-0 lg:divide-x">
            <aside class="py-6 lg:col-span-3">
              <nav class="space-y-1">
                <a
                  @click="toggleTab(tab.id)"
                  v-for="tab in tabs"
                  :key="tab.name"
                  :class="[
                    tab.id === openTab
                      ? 'border-frost-3 bg-polar-night-4'
                      : 'border-transparent hover:bg-polar-night-3',
                    'group flex cursor-pointer items-center border-l-4 px-3 py-2 text-sm font-medium',
                  ]"
                  :aria-current="tab.id === openTab ? 'page' : undefined">
                  <component :is="tab.icon" :class="['-ml-1 mr-3 h-6 w-6 flex-shrink-0']" aria-hidden="true" />
                  <span class="truncate">{{ tab.name }}</span>
                </a>
              </nav>
            </aside>

            <form class="lg:col-span-9" method="POST">
              <!-- Display settings -->
              <div :class="{ hidden: 'display' !== openTab }">
                <div class="py-6 px-4 sm:p-6 lg:pb-8">
                  <div>
                    <h2 class="text-lg font-bold leading-6">Display</h2>
                    <p class="mt-1 text-sm">Modify the appearance of the application.</p>
                  </div>
                  <ul role="list" class="mt-2 divide-y divide-gray-200">
                    <SwitchGroup as="li" class="flex items-center justify-between py-4">
                      <div class="flex flex-col">
                        <SwitchLabel as="p" class="text-sm font-bold">
                          {{ modifiedSettings.DarkMode ? `Dark` : `Light` }} Mode
                        </SwitchLabel>
                        <SwitchDescription class="text-xs">
                          {{ modifiedSettings.DarkMode ? `Unruin` : `Ruin` }} your eyes
                        </SwitchDescription>
                      </div>
                      <Switch
                        @update:modelValue="setDarkMode"
                        v-model="modifiedSettings.DarkMode"
                        :class="[
                          modifiedSettings.DarkMode ? 'bg-aurora-4' : 'bg-polar-night-3',
                          'relative ml-4 inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full',
                          'border-2 border-transparent transition-colors duration-200 ease-in-out',
                        ]">
                        <span
                          aria-hidden="true"
                          :class="[
                            modifiedSettings.DarkMode ? 'translate-x-5' : 'translate-x-0',
                            'inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0',
                            'transition duration-200 ease-in-out',
                          ]" />
                      </Switch>
                    </SwitchGroup>
                  </ul>
                </div>
              </div>

              <!-- Proxy settings -->
              <div :class="{ hidden: 'proxy' !== openTab }">
                <div class="py-6 px-4 sm:p-6 lg:pb-8">
                  <div>
                    <h2 class="text-lg font-medium leading-6">Proxy</h2>
                    <p class="mt-1 text-sm">Customise proxy settings to suit your workflow.</p>
                  </div>

                  <div class="mt-8">
                    <label for="port" class="block text-sm font-medium text-snow-storm">Port</label>
                    <div class="relative mt-1 rounded-md shadow-sm">
                      <input
                        @change="setProxyPort"
                        type="text"
                        name="port"
                        id="port"
                        class="block w-20 rounded-md bg-polar-night-4 focus:border-frost-4 focus:outline-none focus:ring-frost-4 sm:text-sm"
                        :value="modifiedSettings.ProxyPort"
                        aria-invalid="true"
                        aria-describedby="port-error"
                        maxlength="5" />
                    </div>
                  </div>
                </div>
              </div>

              <!-- Certificate settings -->
              <div :class="{ hidden: 'certs' !== openTab }">
                <div class="py-6 px-4 sm:p-6 lg:pb-8">
                  <div>
                    <h2 class="text-lg font-bold leading-6">Certificates</h2>
                    <p class="mt-1 text-sm">Manage certificates and configure your browser(s).</p>
                  </div>
                  <div class="mt-2">
                    <ButtonNetrual @click="exportCA">Export CA Certificate</ButtonNetrual>
                  </div>
                </div>
              </div>

              <!-- About -->
              <div :class="{ hidden: 'about' !== openTab }">
                <div class="py-6 px-4 sm:p-6 lg:pb-8">
                  <div>
                    <h2 class="text-lg font-bold leading-6">About</h2>
                    <p class="mt-1 text-sm">More information about Reaper.</p>
                  </div>

                  <div class="mt-6 flex flex-col lg:flex-row">
                    <div class="flex-grow space-y-6">
                      <div></div>
                    </div>
                  </div>
                </div>

                <div class="divide-y divide-gray-200 pt-6">
                  <div class="px-4 sm:px-6">
                    <div>
                      <h2 class="text-lg font-bold leading-6">Privacy</h2>
                      <p class="mt-1 text-sm">Ornare eu a volutpat eget vulputate. Fringilla commodo amet.</p>
                    </div>
                  </div>
                </div>
              </div>
            </form>
          </div>
          <div class="p-4 text-right">
            <div>
              <ButtonConfirm @click="saveSettings">Save Changes</ButtonConfirm>
              <ButtonCancel @click="cancel">Cancel</ButtonCancel>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
