<script>
export default {
  props: {
    'mini': Boolean,
  },
};
</script>

<template>
  <div class="page-content">
    <div class="main mdc-elevation--z4" :class="{'main--mini': mini}">
      <slot></slot>
    </div>
  </div>
</template>

<style scoped>
.main {
  position: relative;
  display: block;
  width: Calc(100% - 64px);
  max-width: 1024px;
  margin: 20px auto;
  padding: 16px 16px 16px 16px;
  border-radius: 2px;
  background: #FFFFFF;
  font-family: "Roboto", "Helvetica", "Arial", sans-serif;
}

.main.main--mini {
  max-width: 400px;
}
</style>
