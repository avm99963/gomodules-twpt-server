import './theme.scss';

import {createApp} from 'vue';
import * as VueRouter from 'vue-router';

import App from './App.vue';
import AuthorizedUsers from './pages/AuthorizedUsers.vue';
import Home from './pages/Home.vue';
import KillSwitches from './pages/KillSwitches.vue';
import {store} from './store/index.js';
import VueMaterialAdapter from './vma.js';

const routes = [
  {
    path: '/',
    component: Home,
    meta: {title: 'Home'},
  },
  {
    path: '/kill-switches',
    component: KillSwitches,
    meta: {title: 'Kill Switches'},
  },
  {
    path: '/authorized-users',
    component: AuthorizedUsers,
    meta: {title: 'Authorized Users'},
  },
];

const router =
    VueRouter.createRouter({history: VueRouter.createWebHashHistory(), routes});

const app = createApp(App);
app.use(store);
app.use(router);
app.use(VueMaterialAdapter);
app.mount('app');
