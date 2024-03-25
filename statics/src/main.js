import Vue from 'vue'
import {BootstrapVue, IconsPlugin} from 'bootstrap-vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import 'bootstrap/dist/js/bootstrap.js'

import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {library} from "@fortawesome/fontawesome-svg-core";
import {fas} from '@fortawesome/free-solid-svg-icons'
import {far} from '@fortawesome/free-regular-svg-icons'
import {fab} from '@fortawesome/free-brands-svg-icons'
import '@fortawesome/fontawesome-free/css/all.css'
import '@/assets/_shared.scss'

import App from "@/App.vue";
import VueRouter from 'vue-router'
import routes from '@/router'

library.add(fas);
library.add(far);
library.add(fab);

Vue.component('font-awesome-icon', FontAwesomeIcon);

Vue.use(VueRouter);

const router = new VueRouter({
    routes: routes,
    mode: 'history'
})

Vue.use(BootstrapVue)
Vue.use(IconsPlugin)

Vue.config.productionTip = false

new Vue({
    router,
    render: h => h(App)
}).$mount('#app');
