<template>
  <li class="context-menu-item">
    <slot/>
  </li>
</template>

<script lang="ts">
export default {}
</script>

<style scoped>
.context-menu-item {
  display: block;
  position: relative;
  padding: 2px 2px;
  background: #000;
  border-radius: 0;
  color: #FFF;
  text-decoration: none;
  font-size: 13px;
  width: 100%;
  text-align: left;
  cursor: pointer;
  padding: 8px;
}

.context-menu-item:hover, .context-menu-item:focus {
  background: #fff;
  color: #000;
  outline: none;
}
</style>

