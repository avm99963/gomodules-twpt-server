<script>
import {mapGetters} from 'vuex';

import config from './config.json5';
import GsiButton from './GsiButton.vue';

export default {
  data() {
    return {
      isDrawerOpen: false,
      mainMenuOpen: false,
    };
  },
  components: {
    GsiButton,
  },
  computed: {
    ...mapGetters([
      'isSignedIn',
    ]),
  },
  methods: {
    onSignIn(response) {
      this.$store.commit('setJwtToken', response.credential);
    },
    onSignOut() {
      this.$store.commit('setJwtToken', null);
    },
    onSelectMainMenu(e) {
      switch (e.item.getAttribute('data-entry')) {
        case 'sign-out':
          this.onSignOut();
          break;

        default:
          console.error('Unknown menu entry.');
      }
    },
  },
  created() {
    this.$store.dispatch('connectClient', config.grpcWebHost);
  },
};
</script>

<template>
  <mcw-drawer
      ref="drawer"
      v-model="isDrawerOpen"
      dismissible
      class="primary-drawer">
    <template #header>
      <div class="mdc-drawer__header"></div>
    </template>

    <mcw-list-item icon="home" to="/" tabindex="0">Home</mcw-list-item>
    <mcw-list-item icon="emergency" to="/kill-switches">Kill switches</mcw-list-item>
    <template v-if="isSignedIn">
      <mcw-list-item icon="person" to="/authorized-users">Authorized users</mcw-list-item>
    </template>
  </mcw-drawer>
  <div ref="app-content" class="mdc-drawer-app-content">
    <mcw-top-app-bar class="main-toolbar">
      <div class="mdc-top-app-bar__row">
        <section class="mdc-top-app-bar__section mdc-top-app-bar__section--align-start">
          <button @click="isDrawerOpen = !isDrawerOpen" class="material-icons mdc-top-app-bar__navigation-icon mdc-icon-button" aria-label="Open navigation menu">menu</button>
          <span class="mdc-top-app-bar__title">{{ $route.meta.title }}</span>
        </section>
        <section class="mdc-top-app-bar__section mdc-top-app-bar__section--align-end" role="toolbar">
          <template v-if="isSignedIn">
            <mcw-menu-anchor>
              <button @click="mainMenuOpen = true" class="material-icons mdc-top-app-bar__action-item mdc-icon-button" aria-label="Options">more_vert</button>
              <mcw-menu v-model="mainMenuOpen" @select="onSelectMainMenu">
                <mcw-list-item data-entry="sign-out">Sign out</mcw-list-item>
              </mcw-menu>
            </mcw-menu-anchor>
          </template>
          <template v-else>
            <gsi-button @on-signin="onSignIn"></gsi-button>
          </template>
        </section>
      </div>
    </mcw-top-app-bar>
    <main class="mdc-top-app-bar--fixed-adjust">
      <router-view />
    </main>
  </div>
</template>
