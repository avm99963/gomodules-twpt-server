<script>
import Page from './utils/Page.vue';

export default {
  props: [
    'icon',
  ],
  components: {
    Page,
  },
};
</script>

<template>
  <page mini>
    <div class="layout">
      <mcw-material-icon class="main-icon" :icon="icon" />
      <span class="label"><slot></slot></span>
    </div>
  </page>
</template>

<style lang="scss" scoped>
@use "@material/theme/color-palette" as palette;

.layout {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 16px 0;
}

.main-icon {
  font-size: 75px;
  color: palette.$grey-800;
}

.label {
  margin-top: 14px;
  text-align: center;
}
</style>
