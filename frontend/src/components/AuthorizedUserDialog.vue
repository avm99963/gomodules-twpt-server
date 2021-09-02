<script>
import {accessLevels} from '../consts.js';
import * as ksObjectsPb from '../api_proto/kill_switch_objects_pb.js';
import * as ksPb from '../api_proto/kill_switch_pb.js';

export default {
  props: {
    modelValue: Boolean, // If the dialog is open
    isUpdate: Boolean,
    user: {
      type: Object,
      default: function() {
        return new ksObjectsPb.KillSwitchAuthorizedUser()
      },
    },
  },
  emits: [
    'update:modelValue',
    'userAdded',
    'userUpdated',
  ],
  data() {
    return {
      accessLevels,
      userM: this.user.cloneMessage(), // User mutable
    };
  },
  watch: {
    user(newUser, oldUser) {
      this.$data.userM = this.user.cloneMessage();
    }
  },
  computed: {
    accessLevel: {
      get() {
        return this.$data.userM.getAccessLevel().toString();
      },
      set(level) {
        this.$data.userM.setAccessLevel(parseInt(level));
      },
    },
    googleUid: {
      get() {
        return this.$data.userM.getGoogleUid();
      },
      set(uid) {
        this.$data.userM.setGoogleUid(uid);
      },
    },
    email: {
      get() {
        return this.$data.userM.getEmail();
      },
      set(email) {
        this.$data.userM.setEmail(email);
      },
    },
  },
  methods: {
    onSubmit() {
      if (this.isUpdate)
        this.doUpdate();
      else
        this.doAdd();
    },
    doUpdate() {
      let request = new ksPb.UpdateAuthorizedUserRequest();
      request.setUserId(this.userM.getId());
      request.setUser(this.userM);

      this.$store.state.client.updateAuthorizedUser(
        request, {authorization: this.$store.state.jwtToken})
          .then(res => {
            this.$emit('userUpdated');
          })
          .catch(err => console.error(err));
    },
    doAdd() {
      let request = new ksPb.AddAuthorizedUserRequest();
      request.setUser(this.userM);

      this.$store.state.client.addAuthorizedUser(
        request, {authorization: this.$store.state.jwtToken})
          .then(res => {
            this.$emit('userAdded');
            this.$data.userM = this.user.cloneMessage();
          })
          .catch(err => console.error(err));
    },
    onCancel() {
      this.$data.userM = this.user.cloneMessage();
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
    <mcw-dialog-title>{{ isUpdate ? 'Update user ' + userM.getId() : 'Add new user' }}</mcw-dialog-title>
    <mcw-dialog-content>
      <div>
        <p><mcw-select v-model="accessLevel" label="Access level" required>
          <mcw-list-item v-for="(al, i) in accessLevels" :data-value="i" role="option">
            {{ al }}
          </mcw-list-item>
        </mcw-select></p>
        <p><mcw-textfield v-model="googleUid" label="Google UID" /></p>
        <p><mcw-textfield v-model="email" type="email" label="E-mail address" /></p>
      </div>
    </mcw-dialog-content>
    <mcw-dialog-footer>
      <mcw-dialog-button @click="onCancel" action="dismiss">Cancel</mcw-dialog-button>
      <mcw-dialog-button @click="onSubmit" action="accept">{{ isUpdate ? 'Update' : 'Add' }}</mcw-dialog-button>
    </mcw-dialog-footer>
  </mcw-dialog>
</template>
