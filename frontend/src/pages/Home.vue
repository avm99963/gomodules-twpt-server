<script>
import Page from './utils/Page.vue';

export default {
  components: {
    Page,
  },
};
</script>

<template>
  <page>
    <h1>TW Power Tools server</h1>
    <p>Welcome to the TW Power Tools server dashboard. Currently the server is only used to manage <a href="https://docs.google.com/document/d/1O5YV6_WcxwrUyz-lwHOSTfZ3oyIFWj2EQee0VuKkhaA/edit" rel="noreferrer noopener">kill switches</a>.</p>
    <p>While you have to sign in to enable/disable kill switches, you can see what's their status anonymously. Only some allowlisted users are currently able to do this. If you wish to be added to the allowlist, please <a href="https://iavm.xyz/b/twpowertools/new" rel="noreferrer noopener">fill in a bug</a> explaining why you're requiring access.</p>
  </page>
</template>
