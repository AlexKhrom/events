import Vue from 'vue'
import App from './App.vue'
import store from './store'
import Vue2TouchEvents from 'vue2-touch-events'
import vuetify from './plugins/vuetify';
import moment from "moment";
import vueRouter from 'vue-router'
import router from './router'
import VueMoment from 'vue-moment'

Vue.use(vueRouter)
Vue.use(Vue2TouchEvents)
Vue.use(moment)
Vue.prototype.$moment = moment;
// moment.locale('ru')

import VueTheMask from 'vue-the-mask'
// import Statistics from "@/components/Statistics";
Vue.use(VueTheMask)

Vue.config.productionTip = false

// const NotFound = { template: '<p>Страница не найдена</p>' }
// const Home = { template: App }
// const Stat = { template:  Statistics}
//
// const routes = {
//   '/': Home,
//   '/statistics': Stat
// }
//
// new Vue({
//   el: '#app',
//   data: {
//     currentRoute: window.location.pathname
//   },
//   computed: {
//     ViewComponent () {
//       return routes[this.currentRoute] || NotFound
//     }
//   },
//   render (h) { return h(this.ViewComponent) }
// })

new Vue({
  store,
  vuetify,
  render: h => h(App),
  router,
  VueMoment,
}).$mount('#app')
