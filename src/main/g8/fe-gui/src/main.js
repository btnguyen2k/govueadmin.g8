import 'core-js/stable'
import Vue from 'vue'
import VueI18n from 'vue-i18n'
import App from './App'
import router from './router'
import CoreuiVue from '@coreui/vue'
import { iconsSet as icons } from './assets/icons/icons.js'
import store from './store'
import messages from './messages'

Vue.config.performance = true
Vue.use(CoreuiVue)
Vue.use(VueI18n)
Vue.prototype.$log = console.log.bind(console)

const i18n = new VueI18n({
    locale: 'vi',
    messages
})

new Vue({
    el: '#app',
    router,
    store,
    icons,
    template: '<App/>',
    components: {
        App
    },
    i18n
})
