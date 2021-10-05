import Vue from 'vue'
import {BootstrapVue, BootstrapVueIcons, NavbarPlugin, SidebarPlugin} from 'bootstrap-vue'
import 'bootstrap'
import 'bootstrap/dist/css/bootstrap.min.css'
import App from './App.vue'
import { router } from './router'
import VueCookies from 'vue-cookies'
import VueJwtDecode from 'vue-jwt-decode'
import VueApexCharts from 'vue-apexcharts'
import VueAuthHref from 'vue-auth-href'
import ProgressBar from 'vue-simple-progress'


Vue.use(BootstrapVue)
Vue.use(BootstrapVueIcons)
Vue.use(NavbarPlugin)
Vue.use(SidebarPlugin)
Vue.use(VueCookies)
Vue.use(VueJwtDecode)
Vue.use(VueApexCharts)
Vue.use(VueAuthHref)

Vue.component('progressbar', ProgressBar)
Vue.component('apexchart', VueApexCharts)
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
