import Vue from 'vue'
import Vuex from 'vuex'
import * as files from './modules/files'
import * as info from './modules/info'

Vue.use(Vuex)

export default new Vuex.Store({
    modules: {
        files,
        info
    },
    state: {
    },
    mutations: {
    },
    actions: {
    }
})
