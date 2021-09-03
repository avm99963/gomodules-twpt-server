<script>
import {browsers} from '../consts.js';

export default {
  props: {
    killSwitch: Object,
    canDisable: Boolean,
  },
  emits: [
    'wantsToDisableKillSwitch',
  ],
  computed: {
    versionsText() {
      if (this.killSwitch?.getMinVersion() == '' && this.killSwitch?.getMaxVersion() == '')
        return 'All';

      return 'From ' + (this.killSwitch?.getMinVersion() || '...') + ' to '+ (this.killSwitch?.getMaxVersion() || '...');
    },
    browsersText() {
      return this.killSwitch?.getBrowsersList()?.map(b => browsers[parseInt(b)]).join(', ') ?? 'undefined';
    },
  },
};
</script>

<template>
  <div class="main">
    <mcw-button @click="$emit('wantsToDisableKillSwitch', killSwitch)" class="disable-btn" v-if="canDisable && killSwitch?.getActive()">
      Disable
    </mcw-button>
    <div class="feature-name">{{ killSwitch?.getFeature()?.getCodename() }}</div>
    <div class="status"><span class="status--label">Status:</span> <span class="status-text" :class="{'status-text--active': killSwitch?.getActive()}">{{ killSwitch?.getActive() ? 'active (feature is force disabled)' : 'no longer active' }}</span></div>
    <div class="details-container">
      <div class="details">
        <div class="details--title">Versions affected</div>
        <div class="details--content">{{ versionsText }}</div>
      </div>
      <div class="details">
        <div class="details--title">Browsers affected</div>
        <div class="details--content">{{ browsersText }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use "@material/theme/color-palette" as palette;

.main {
  position: relative;
  width: calc(100% - 32px);
  max-width: 500px;
  padding: 16px;
  border: solid 1px palette.$grey-500;
  border-radius: 4px;
}

.disable-btn {
  position: absolute;
  top: 16px;
  right: 16px;
}

.feature-name {
  color: palette.$grey-900;
  font-family: 'Roboto Mono', monospace;
  font-weight: 500;
  font-size: 20px;
  margin-bottom: 8px;
}

.status {
  color: palette.$grey-700;
}

.status--label {
  font-weight: 500;
}

.status-text.status-text--active {
  color: palette.$red-700;
}

.details-container {
  margin-top: 16px;
  display: flex;
  flex-direction: row;
  justify-content: space-around;
}

.details {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.details--title {
  font-weight: 500;
  font-size: 17px;
  text-align: center;
}

.details--content {
  font-size: 15px;
  text-align: center;
}
</style>
