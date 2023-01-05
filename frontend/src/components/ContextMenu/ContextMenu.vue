
<template>
  <div
      class="context-menu"
      ref="popper"
      v-show="isVisible"
      tabindex="-1"
      @contextmenu.capture.prevent>
    <ul>
      <slot :contextData="contextData" />
    </ul>
  </div>
</template>

<script lang="ts">
import Popper from 'popper.js';
import {defineComponent} from "vue";

export default /*#__PURE__*/ defineComponent({
  props: {
    boundariesElement: {
      type: String,
      default: 'body',
    },
  },
  components: {
    Popper,
  },
  data() {
    return {
      opened: false,
      contextData: "something",
      clean: () => {},
    };
  },
  computed: {
    isVisible(): boolean{
      return this.opened;
    },
  },
  methods: {
    open(evt: MouseEvent, contextData: any) {
      this.opened = true;
      this.contextData = contextData;

      this.clean()

      let popper = new Popper(this.referenceObject(evt), this.$refs.popper as Element, {
        placement: 'right-start',
        modifiers: {
          preventOverflow: {
            boundariesElement: document.querySelector(this.boundariesElement) as Element,
          },
        },
      });
      this.clean = popper.destroy;

      // Recalculate position
      this.$nextTick(() => {
        popper.scheduleUpdate();
      });

    },
    close() {
      this.opened = false;
      this.contextData = "";
    },
    referenceObject(evt: MouseEvent): Popper.ReferenceObject {
      const left = evt.clientX;
      const top = evt.clientY;
      const right = left + 1;
      const bottom = top + 1;
      const clientWidth = 1;
      const clientHeight = 1;

      function getBoundingClientRect() {
        return {
          left,
          top,
          right,
          bottom,
        };
      }

      return {
        getBoundingClientRect,
        clientWidth,
        clientHeight,
      } as Popper.ReferenceObject;
    },
  },
  beforeUnmount() {
    this.clean()
  },
});

</script>

<style scoped>

.context-menu {
  position: fixed;
  z-index: 999;
  overflow: hidden;
  background: #FFF;
  border-radius: 4px;
  box-shadow: 0 1px 4px 0 #eee;
}

.context-menu:focus {
   outline: none;
 }

.context-menu ul {
  padding:0px;
  margin:0px;
}



</style>
