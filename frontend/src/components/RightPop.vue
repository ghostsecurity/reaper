

<script lang="ts" setup>
import { TransitionRoot, TransitionChild, Dialog, DialogPanel } from "@headlessui/vue";
</script>

<script lang="ts">
import { defineComponent, PropType } from "vue";
export default /*#__PURE__*/ defineComponent({
  name: "RightPop",
  props: {
    show: {
      type: Boolean,
      required: true
    },
    onRequestClose: {
      type: Function as PropType<() => void>,
      required: true
    },
  },
  data() {
    return {
      open: true,
    }
  },
  methods: {
    close() {
      this.onRequestClose()
    }
  },
})
</script>

<template>
  <TransitionRoot as="template" :show="show">
    <Dialog as="div" class="relative z-10" @close="close">
      <div class="fixed inset-0" />
      <div class="fixed inset-0 overflow-hidden">
        <div class="absolute inset-0 overflow-hidden">
          <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-full pl-10 sm:pl-16">
            <TransitionChild as="template" enter="transform transition ease-in-out duration-200"
              enter-from="translate-x-full" enter-to="translate-x-0"
              leave="transform transition ease-in-out duration-200" leave-from="translate-x-0"
              leave-to="translate-x-full">
              <DialogPanel class="pointer-events-auto w-screen max-w-4xl">
                <div class="flex h-full flex-col overflow-y-hidden bg-snow-storm dark:bg-polar-night shadow-xl">
                  <div class="relative flex-1 px-4 sm:px-6">
                    <slot />
                  </div>
                </div>
              </DialogPanel>
            </TransitionChild>
          </div>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<style scoped>

</style>