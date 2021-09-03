<script>
import {mapGetters} from 'vuex';

import NotAuthorized from './NotAuthorized.vue';
import Page from './utils/Page.vue';
import AuthorizedUserDialog from '../components/AuthorizedUserDialog.vue';

//import * as grpcWeb from "grpc-web";
import * as ksObjectsPb from '../api_proto/kill_switch_objects_pb.js';
import * as ksPb from '../api_proto/kill_switch_pb.js';
import {accessLevels} from '../consts.js';

export default {
  data() {
    let emptyUser = new ksObjectsPb.KillSwitchAuthorizedUser();

    return {
      users: [],
      currentUpdateUser: emptyUser.cloneMessage(), // Current user being updated
      currentDeleteUser: emptyUser.cloneMessage(), // Current user being confirmed deletion
      addDialogOpen: false,
      updateDialogOpen: false,
      deleteDialogOpen: false,
      accessLevels,
    };
  },
  components: {
    NotAuthorized,
    Page,
    AuthorizedUserDialog,
  },
  mounted() {
    this.loadData();
  },
  methods: {
    loadData() {
      if (!this.isSignedIn) return;

      let request = new ksPb.ListAuthorizedUsersRequest();

      this.$store.state.client.listAuthorizedUsers(
          request, {authorization: this.$store.state.jwtToken})
          .then(res => {
            this.$data.users = res.getUsersList();
          })
          .catch(err => console.error(err));
    },
    showAddForm(user) {
      this.$data.addDialogOpen = true;
    },
    showUpdateForm(user) {
      this.currentUpdateUser = user;
      this.$data.updateDialogOpen = true;
    },
    showDeleteDialog(user) {
      this.$data.currentDeleteUser = user;
      this.$data.deleteDialogOpen = true;
    },
    deleteUser(user) {
      let request = new ksPb.DeleteAuthorizedUserRequest();
      request.setUserId(user.getId());

      this.$store.state.client.deleteAuthorizedUser(
          request, {authorization: this.$store.state.jwtToken})
          .then(res => {
            this.loadData();
          })
          .catch(err => console.error(err));
    },
  },
  computed: {
    ...mapGetters([
      'isSignedIn',
    ]),
  },
};
</script>

<template>
  <template v-if="isSignedIn">
    <div class="container">
      <mcw-data-table>
        <table class="mdc-data-table__table" aria-label="Authorized users">
          <thead>
            <tr class="mdc-data-table__header-row">
              <th
                  class="mdc-data-table__header-cell"
                  role="columnheader"
                  scope="col">
                ID
              </th>
              <th
                  class="mdc-data-table__header-cell"
                  role="columnheader"
                  scope="col">
                Google UID
              </th>
              <th
                  class="mdc-data-table__header-cell"
                  role="columnheader"
                  scope="col">
                E-mail address
              </th>
              <th
                  class="mdc-data-table__header-cell"
                  role="columnheader"
                  scope="col">
                Scope
              </th>
              <th
                  class="mdc-data-table__header-cell"
                  role="columnheader"
                  scope="col">
              </th>
            </tr>
          </thead>
          <tbody class="mdc-data-table__content">
            <tr class="mdc-data-table__row" v-for="user in users">
              <td class="mdc-data-table__cell mdc-data-table__cell--numeric">{{ user.getId() }}</td>
              <td class="mdc-data-table__cell">{{ user.getGoogleUid() || "-" }}</td>
              <td class="mdc-data-table__cell">{{ user.getEmail() || "-" }}</td>
              <td class="mdc-data-table__cell">{{ accessLevels[user.getAccessLevel()] }}</td>
              <td class="mdc-data-table__cell">
                <mcw-icon-button @click="(showUpdateForm(user))"><mcw-material-icon icon="edit" /></mcw-icon-button>
                <mcw-icon-button @click="(showDeleteDialog(user))"><mcw-material-icon icon="delete" /></mcw-icon-button>
              </td>
            </tr>
          </tbody>
        </table>
      </mcw-data-table>
    </div>
    <mcw-fab @click="showAddForm()" class="add-user" icon="add" />
    <mcw-dialog
        v-model="deleteDialogOpen"
        escape-key-action="close"
        scrim-click-action="close"
        aria-labelledby="delete-title"
        aria-describedby="delete-content"
        :auto-stack-buttons="true">
      <mcw-dialog-title id="delete-title">Delete authorized user?</mcw-dialog-title>
      <mcw-dialog-content id="delete-content">
        <div>User {{ currentDeleteUser.getId() }} will no longer have access to the TW Power Tools dashboard.</div>
      </mcw-dialog-content>
      <mcw-dialog-footer>
        <mcw-dialog-button action="dismiss">Cancel</mcw-dialog-button>
        <mcw-dialog-button @click="deleteUser(currentDeleteUser)" action="accept">Delete</mcw-dialog-button>
      </mcw-dialog-footer>
    </mcw-dialog>
    <!-- Add user dialog -->
    <authorized-user-dialog @user-added="loadData()" v-model="addDialogOpen" />
    <!-- Update user dialog -->
    <authorized-user-dialog @user-updated="loadData()" v-model="updateDialogOpen" is-update :user="currentUpdateUser" />
  </template>
  <template v-else>
    <not-authorized>
    </not-authorized>
  </template>
</template>

<style scoped>
.container {
  margin-top: 16px;
  display: flex;
  justify-content: space-evenly;
}

.add-user {
  position: fixed;
  right: 16px;
  bottom: 16px;
}
</style>
