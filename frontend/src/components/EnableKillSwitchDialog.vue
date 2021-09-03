<script>
import {browsers} from '../consts.js';

import * as ksObjectsPb from '../api_proto/kill_switch_objects_pb.js';
import * as ksPb from '../api_proto/kill_switch_pb.js';
import * as commonPb from '../api_proto/common_pb.js';

import KillSwitch from './KillSwitch.vue';

let browserOptions = {...browsers};
delete browserOptions['0'];
let defaultBrowsers = Object.keys(browserOptions);

export default {
  components: {
    KillSwitch,
  },
  props: {
    modelValue: Boolean,
  },
  emits: [
    'killSwitchAdded',
  ],
  data() {
    return {
      secondOpen: false,
      featuresList: [],
      browserOptions,
      feature: "",
      minVersion: "",
      maxVersion: "",
      browsers: [...defaultBrowsers],
    };
  },
  mounted() {
    this.loadData();
  },
  computed: {
    killSwitch() {
      let ks = new ksObjectsPb.KillSwitch();
      for (let f of this.featuresList) {
        if (f.getId() == this.feature) {
          ks.setFeature(f);
          break;
        }
      }
      ks.setMinVersion(this.minVersion);
      ks.setMaxVersion(this.maxVersion);
      ks.setBrowsersList(this.browsers);
      ks.setActive(true);
      return ks;
    },
  },
  methods: {
    loadData() {
      let request = new ksPb.ListFeaturesRequest();
      request.setWithDeprecatedFeatures(true);

      this.$store.state.client.listFeatures(
          request, {authorization: this.$store.state.jwtToken})
          .then(res => {
            this.$data.featuresList = res.getFeaturesList();
          })
          .catch(err => console.error(err));
    },
    onCancel() {
      this.secondOpen = false;
      this.feature = "";
      this.minVersion = "";
      this.maxVersion = "";
      this.browsers = [...defaultBrowsers];
    },
    onNext() {
      this.secondOpen = true;
    },
    onPrevious() {
      this.secondOpen = false;
      this.$emit('update:modelValue', true);
    },
    onConfirm() {
      let request = new ksPb.EnableKillSwitchRequest();
      request.setKillSwitch(this.killSwitch.cloneMessage());

      this.$store.state.client.enableKillSwitch(
          request, {authorization: this.$store.state.jwtToken})
          .then(res => {
            this.$emit('killSwitchAdded');
            this.onCancel();
          })
          .catch(err => console.error(err));
    },
  },
};
</script>

<template>
  <mcw-dialog
      v-model="modelValue"
      @update:modelValue="$emit('update:modelValue', $event)"
      escape-key-action="close"
      scrim-click-action="close"
      :auto-stack-buttons="true">
    <mcw-dialog-title>Enable Kill Switch</mcw-dialog-title>
    <mcw-dialog-content>
      <div>
        <p>
          <mcw-select ref="featureSelect" v-model="feature" label="Feature to disable" required>
            <mcw-list-item data-value="" role="option"></mcw-list-item>
            <mcw-list-item v-for="feature in featuresList" :data-value="feature.getId()" role="option">
              {{ feature.getCodename() }}
            </mcw-list-item>
          </mcw-select>
        </p>
        <p>
          <span class="helper-text">Disable the feature on extension versions greater or equal to:</span><br>
          <mcw-textfield v-model="minVersion" label="Minimum version" helptext="(optional)" helptext-persistent fullwidth />
        </p>
        <p>
          <span class="helper-text">Disable the feature on extension versions less or equal to:</span><br>
          <mcw-textfield v-model="maxVersion" label="Maximum version" helptext="(optional)" helptext-persistent fullwidth />
        </p>
        <p>
          <span class="helper-text">Disable in the following browsers:</span>
          <template v-for="(b, i) in browserOptions"><br><mcw-checkbox :value="i" v-model="browsers" :label="b" /></template>
        </p>
      </div>
    </mcw-dialog-content>
    <mcw-dialog-footer>
      <mcw-dialog-button @click="onCancel" action="dismiss">Cancel</mcw-dialog-button>
      <mcw-dialog-button @click="onNext" action="accept">Next</mcw-dialog-button>
    </mcw-dialog-footer>
  </mcw-dialog>
  <mcw-dialog
      v-model="secondOpen"
      escape-key-action="close"
      scrim-click-action="close"
      :auto-stack-buttons="true">
    <mcw-dialog-title>Enable Kill Switch?</mcw-dialog-title>
    <mcw-dialog-content>
      <div>
        <p>This will disable the feature globally for all the extension users, if they use one of the versions and browsers specified.</p>
        <p>Please take a look at this preview and confirm that this is indeed what you want to do:</p>
        <kill-switch :kill-switch="killSwitch" />
      </div>
    </mcw-dialog-content>
    <mcw-dialog-footer>
      <mcw-dialog-button @click="onPrevious">Previous</mcw-dialog-button>
      <mcw-dialog-button class="confirm-bad" @click="onConfirm" action="accept">Confirm</mcw-dialog-button>
    </mcw-dialog-footer>
  </mcw-dialog>
</template>

<style scoped lang="scss">
@use "@material/theme/color-palette" as palette;
@use "@material/button";

.helper-text {
  font-size: 14px;
  color: palette.$grey-700;
}

.confirm-bad {
  @include button.ink-color(palette.$red-500);
}
</style>
