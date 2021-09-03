<script>
import {mapGetters} from 'vuex';

import * as ksPb from '../api_proto/kill_switch_pb.js';
import * as ksObjectsPb from '../api_proto/kill_switch_objects_pb.js';

import EnableKillSwitchDialog from '../components/EnableKillSwitchDialog.vue';
import KillSwitch from '../components/KillSwitch.vue';
import MiniMessage from './MiniMessage.vue';

let emptyKillSwitch = new ksObjectsPb.KillSwitch();

export default {
  data() {
    return {
      enableKillSwitchDialog: false,
      disableDialogOpen: false,
      currentDisableKillSwitch: emptyKillSwitch.cloneMessage(),
      killSwitches: [],
    };
  },
  components: {
    EnableKillSwitchDialog,
    KillSwitch,
    MiniMessage,
  },
  computed: {
    ...mapGetters([
      'isSignedIn',
    ]),
    activeKillSwitches() {
      return this.killSwitches.filter(ks => ks.getActive());
    },
  },
  mounted() {
    this.loadData();
  },
  methods: {
    loadData() {
      let request = new ksPb.GetKillSwitchOverviewRequest();

      this.$store.state.client.getKillSwitchOverview(request)
          .then(res => {
            this.killSwitches = res.getKillSwitchesList();
          })
          .catch(err => console.error(err));
    },
    showDisableKillSwitchDialog(killSwitch) {
      this.currentDisableKillSwitch = killSwitch;
      this.disableDialogOpen = true;
    },
    disableKillSwitch() {
      let request = new ksPb.DisableKillSwitchRequest();
      request.setKillSwitchId(this.currentDisableKillSwitch.getId());

      this.$store.state.client.disableKillSwitch(
          request, {authorization: this.$store.state.jwtToken})
          .then(res => {
            this.loadData();
            this.disableDialogOpen = false;
            this.currentDisableKillSwitch = emptyKillSwitch.cloneMessage();
          })
          .catch(err => console.error(err));
    },
  },
};
</script>

<template>
  <mini-message icon="toggle_off" v-if="activeKillSwitches?.length == 0">
    There aren't any kill switches enabled currently.<br>
    Everything's working normally.
  </mini-message>
  <template v-else>
    <div class="kill-switch-container">
      <kill-switch
          class="kill-switch"
          :kill-switch="killSwitch"
          :can-disable="isSignedIn"
          @wantsToDisableKillSwitch="showDisableKillSwitchDialog($event)"
          v-for="killSwitch in activeKillSwitches" />
    </div>
  </template>

  <template v-if="isSignedIn">
    <mcw-fab @click="enableKillSwitchDialog = true" class="enable-kill-switch" icon="add" />

    <mcw-dialog
        v-model="disableDialogOpen"
        escape-key-action="close"
        scrim-click-action="close"
        aria-labelledby="disable-title"
        aria-describedby="disable-content"
        :auto-stack-buttons="true">
      <mcw-dialog-title id="disable-title">Disable kill switch?</mcw-dialog-title>
      <mcw-dialog-content id="disable-content">
        <div>Feature <span class="feature-codename">{{ currentDisableKillSwitch.getFeature()?.getCodename() }}</span> will be progressively enabled to all the users of the extension.</div>
      </mcw-dialog-content>
      <mcw-dialog-footer>
        <mcw-dialog-button action="dismiss">Cancel</mcw-dialog-button>
        <mcw-dialog-button @click="disableKillSwitch()" action="accept">Confirm</mcw-dialog-button>
      </mcw-dialog-footer>
    </mcw-dialog>

    <enable-kill-switch-dialog @kill-switch-added="loadData" v-model="enableKillSwitchDialog" />
  </template>
</template>

<style scoped>
.kill-switch-container {
  font-family: 'Roboto', 'Helvetica', 'Arial', sans-serif;
  max-width: 500px;
  margin: 16px auto;
}

.kill-switch {
  margin-bottom: 16px;
}

.enable-kill-switch {
  position: fixed;
  right: 16px;
  bottom: 16px;
}

.feature-codename {
  font-family: 'Roboto Mono', monospace;
}
</style>
