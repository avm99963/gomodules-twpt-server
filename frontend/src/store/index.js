import {createStore} from 'vuex';

import {KillSwitchServicePromiseClient} from '../api_proto/kill_switch_grpc_web_pb.js';

export const store = createStore({
  state() {
    return {
      jwtToken: localStorage.getItem('jwtToken'),
      client: null,
    };
  },
  mutations: {
    setJwtToken(state, token) {
      if (token == null)
        localStorage.removeItem('jwtToken');
      else
        localStorage.jwtToken = token;

      state.jwtToken = token;
    },
  },
  getters: {
    getJwtToken(state) {
      return state.jwtToken;
    },
    isSignedIn(state) {
      return state.jwtToken != null;
    },
  },
  actions: {
    connectClient(store, host) {
      // We enable the dev tools in case they are useful sometime in the future.
      const enableDevTools = window.__GRPCWEB_DEVTOOLS__ || (() => {});
      store.state.client = new KillSwitchServicePromiseClient(host, null, null);
      enableDevTools([
        store.state.client,
      ]);
    },
  },
});
