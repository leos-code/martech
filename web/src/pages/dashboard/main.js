import Vue from 'vue';
import VueCookies from 'vue-cookies';
import 'es6-promise/auto';
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';

import App from './App.vue';
import router from './router';
import store from './store';
import { permission } from '@dashboard/directive';

Vue.use(ElementUI);
Vue.use(VueCookies);
Vue.config.productionTip = false;
Vue.config.baseApi = process.env.VUE_APP_API_DASHBOARD;
Vue.directive('permission', permission);

if (process.env.NODE_ENV === 'development') {
    require('./mock');
}

export default new Vue({
    router,
    store,
    render: (h) => h(App),
}).$mount('#app');
