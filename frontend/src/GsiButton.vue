<script>
import config from './config.js';

export default {
  mounted() {
    if (document.readyState == 'complete')
      this.init();
    else
      window.addEventListener('load', () => this.init());
  },
  methods: {
    init() {
      window.google.accounts.id.initialize({
        client_id: config.google.clientId,
        callback: this.onSignIn,
      });
      window.google.accounts.id.renderButton(document.getElementById('gsi-button'), {
        theme: 'outline',
        size: 'large',
      });
    },
    onSignIn(response) {
      this.$emit('on-signin', response);
    },
  },
};
</script>

<template>
  <div id="gsi-button"></div>
</template>
