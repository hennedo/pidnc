import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '../store'
import Home from '../views/Home.vue'
import Settings from '../views/Settings.vue'

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        component: Home
    },
    {
        path: '/settings',
        component: Settings,
        name: 'Settings',
        beforeEnter(routeTo, routeFrom, next) {
            store
                .dispatch('info/getInfo')
                .then(() => {
                    next()
                });
        }
    }
]

const router = new VueRouter({
    mode: 'history',
    base: '/',
    activeClass: 'active',
    routes
})

export default router